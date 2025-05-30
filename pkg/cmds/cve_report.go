/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"kmodules.xyz/image-packer/pkg/lib"

	"github.com/spf13/cobra"
	//"kubedb.dev/installer/cmd/lib"
	"github.com/olekukonko/tablewriter"
	shell "gomodules.xyz/go-sh"
	"k8s.io/klog/v2"
	"kubeops.dev/scanner/apis/trivy"
)

func NewCmdGenerateCVEReport() *cobra.Command {
	var (
		files  []string
		outDir string
	)
	cmd := &cobra.Command{
		Use:                   "generate-cve-report",
		Short:                 "Generate Image Vulnerability report using Trivy",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			reports, err := GatherReports(files)
			if err != nil {
				return err
			}
			data := GenerateMarkdownReport(reports)

			readmeFile := filepath.Join(outDir, "README.md")
			return os.WriteFile(readmeFile, data, 0o644)
		},
	}
	cmd.Flags().StringSliceVar(&files, "src", files, "List of source files (http url or local file)")
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Output directory")

	return cmd
}

type CVEReport struct {
	Ref      string
	Digest   string
	OS       string
	Critical Stats
	High     Stats
	Medium   Stats
	Low      Stats
	Unknown  Stats
}

func (r *CVEReport) MarkAsMissing() {
	r.Critical.OS = -1
	r.Critical.Other = -1

	r.High.OS = -1
	r.High.Other = -1

	r.Medium.OS = -1
	r.Medium.Other = -1

	r.Low.OS = -1
	r.Low.Other = -1

	r.Unknown.OS = -1
	r.Unknown.Other = -1
}

type Stats struct {
	OS    int
	Other int
}

func (s Stats) PrettyPrint() string {
	b, a := "0", "0"
	if s.OS >= 0 {
		b = strconv.Itoa(s.OS)
	}
	if s.Other >= 0 {
		a = strconv.Itoa(s.Other)
	}
	if s.OS > 0 {
		return fmt.Sprintf("**%s**, %s", b, a)
	}
	return fmt.Sprintf("%s, %s", b, a)
}

func (s Stats) String() string {
	b, a := "0", "0"
	if s.OS >= 0 {
		b = strconv.Itoa(s.OS)
	}
	if s.Other >= 0 {
		a = strconv.Itoa(s.Other)
	}
	return fmt.Sprintf("%s, %s", b, a)
}

func (s Stats) Zero() bool {
	return s.OS+s.Other == 0
}

func (r CVEReport) NoCVE() bool {
	return r.Critical.Zero() &&
		r.High.Zero() &&
		r.Medium.Zero() &&
		r.Low.Zero() &&
		r.Unknown.Zero()
}

func (r CVEReport) Headers() []string {
	return []string{
		"Image Ref",
		"OS",
		"Critical<br>(os, other)",
		"High<br>(os, other)",
		"Medium<br>(os, other)",
		"Low<br>(os, other)",
		"Unknown<br>(os, other)",
	}
}

func (r CVEReport) Strings() []string {
	ref := r.Ref
	if r.Digest != "" {
		ref += "<br>" + r.Digest
	}
	return []string{
		ref,
		r.OS,
		r.Critical.PrettyPrint(),
		r.High.PrettyPrint(),
		r.Medium.String(),
		r.Low.String(),
		r.Unknown.String(),
	}
}

// "Class": "os-pkgs",
func GatherReports(files []string) ([]CVEReport, error) {
	images, err := GenerateImageList(files, false)
	if err != nil {
		return nil, err
	}

	sh := lib.NewShell()

	reports := make([]CVEReport, 0, len(images))
	for _, ref := range images {
		cveReport, err := gatherReport(sh, ref)
		if err != nil {
			klog.Warningln(err.Error())
			continue
		}
		reports = append(reports, cveReport)
	}

	return reports, nil
}

func gatherReport(sh *shell.Session, ref string) (CVEReport, error) {
	cveReport := CVEReport{
		Ref: ref,
	}
	if digest, found, err := lib.ImageDigest(ref); err != nil {
		cveReport.MarkAsMissing()
		return cveReport, err
	} else if found {
		cveReport.Digest = digest
		report, err := lib.Scan(sh, ref)
		if err != nil {
			cveReport.MarkAsMissing()
			return cveReport, err
		}
		setReport(report, &cveReport)
	}
	return cveReport, nil
}

func setReport(report *trivy.SingleReport, result *CVEReport) {
	result.OS = fmt.Sprintf("%s %s", report.Metadata.Os.Family, report.Metadata.Os.Name)

	for _, rpt := range report.Results {
		for _, tv := range rpt.Vulnerabilities {
			switch tv.Severity {
			case "CRITICAL":
				if rpt.Class == "os-pkgs" {
					result.Critical.OS += 1
				} else {
					result.Critical.Other += 1
				}
			case "HIGH":
				if rpt.Class == "os-pkgs" {
					result.High.OS += 1
				} else {
					result.High.Other += 1
				}
			case "MEDIUM":
				if rpt.Class == "os-pkgs" {
					result.Medium.OS += 1
				} else {
					result.Medium.Other += 1
				}
			case "LOW":
				if rpt.Class == "os-pkgs" {
					result.Low.OS += 1
				} else {
					result.Low.Other += 1
				}
			case "UNKNOWN":
				if rpt.Class == "os-pkgs" {
					result.Unknown.OS += 1
				} else {
					result.Unknown.Other += 1
				}
			}
		}
	}
}

func GenerateMarkdownReport(reports []CVEReport) []byte {
	var buf bytes.Buffer
	buf.WriteString("# CVE Report")
	buf.WriteRune('\n')
	buf.Write(generateMarkdownTable(reports))

	return buf.Bytes()
}

func generateMarkdownTable(reports []CVEReport) []byte {
	var tr CVEReport

	data := make([][]string, 0, len(reports))
	for _, r := range reports {
		data = append(data, r.Strings())
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(tr.Headers())
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return buf.Bytes()
}
