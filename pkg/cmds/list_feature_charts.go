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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"kmodules.xyz/client-go/tools/parser"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
)

func NewCmdListFeatureCharts() *cobra.Command {
	var (
		rootDir string
		outDir  string
	)
	cmd := &cobra.Command{
		Use:                   "list-feature-charts",
		Short:                 "List all feature charts",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			images, err := ListUICharts(rootDir)
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(images)
			if err != nil {
				return err
			}

			filename := filepath.Join(outDir, "feature-charts.yaml")
			err = os.WriteFile(filename, data, 0o644)
			return err
		},
	}

	cmd.Flags().StringVar(&rootDir, "root-dir", "", "Root directory")
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Output directory")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "output-dir")

	return cmd
}

type Skeleton struct {
	Spec struct {
		Chart struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"chart"`
	} `json:"spec"`
}

type ChartInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

func ListUICharts(rootDir string) ([]string, error) {
	sh := shell.NewSession()
	sh.SetDir("/tmp")
	sh.ShowCMD = true

	images := sets.New[string]()
	var out []byte
	var err error

	if rootDir == "" {
		// helm search repo opscenter-features -o json
		out, err = sh.Command("helm", "search", "repo", "opscenter-features", "-o", "json").Output()
		if err != nil {
			return nil, err
		}
		var list []ChartInfo
		err = json.Unmarshal(out, &list)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, errors.New("helm chart opscenter-features not found")
		}
		version := list[0].Version

		out, err = sh.Command("helm", "template", "oci://ghcr.io/appscode-charts/opscenter-features", "--version="+version).Output()
		if err != nil {
			return nil, err
		}
	} else {
		out, err = sh.SetDir(rootDir).Command("helm", "template", "opscenter-features").Output()
		if err != nil {
			return nil, err
		}
	}

	helmout, err := parser.ListResources(out)
	if err != nil {
		panic(err)
	}

	for _, ri := range helmout {
		if ri.Object.GetKind() != "FeatureSet" && ri.Object.GetKind() != "Feature" {
			continue
		}

		chartName, found, err := unstructured.NestedString(ri.Object.UnstructuredContent(), "spec", "chart", "name")
		if err != nil {
			return nil, err
		} else if !found {
			continue
		}
		chartVersion, found, err := unstructured.NestedString(ri.Object.UnstructuredContent(), "spec", "chart", "version")
		if err != nil {
			return nil, err
		} else if !found {
			continue
		}

		images.Insert(fmt.Sprintf("ghcr.io/appscode-charts/%s:%s", chartName, chartVersion))
	}

	return sets.List(images), nil
}
