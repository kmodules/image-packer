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
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kmodules.xyz/go-containerregistry/name"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func NewCmdGenerateScripts() *cobra.Command {
	var (
		files     []string
		nondistro bool
		insecure  bool
		outDir    string
	)
	cmd := &cobra.Command{
		Use:                   "generate-scripts",
		Short:                 "Generate export/import scripts",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenerateScripts(files, outDir, nondistro, insecure)
		},
	}
	cmd.Flags().StringSliceVar(&files, "src", files, "List of source files (http url or local file)")
	cmd.Flags().BoolVar(&nondistro, "allow-nondistributable-artifacts", nondistro, "Allow pushing non-distributable (foreign) layers")
	cmd.Flags().BoolVar(&insecure, "insecure", insecure, "Allow image references to be fetched without TLS")
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Output directory")

	return cmd
}

func generateImageList(files []string) ([]string, error) {
	var images []string

	for _, file := range files {
		list, err := readImageList(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read image list from %s: %w", file, err)
		}
		images = append(images, list...)
	}
	sort.Strings(images)
	return images, nil
}

func readImageList(file string) ([]string, error) {
	if u, err := url.Parse(file); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(file)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, err
		}
		defer resp.Body.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, resp.Body)
		if err != nil {
			return nil, err
		}
		var images []string
		err = yaml.Unmarshal(buf.Bytes(), &images)
		if err != nil {
			return nil, err
		}
		return images, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var images []string
	err = yaml.Unmarshal(data, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func GenerateScripts(files []string, outdir string, nondistro, insecure bool) error {
	images, err := generateImageList(files)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(`#!/bin/bash

set -x

mkdir -p images

`)
	for _, img := range images {
		// crane pull appscode/cluster-ui:0.4.16 images/cluster-ui.tar
		ref, err := name.ParseReference(img)
		if err != nil {
			return err
		}
		if ref.Tag == "" {
			return fmt.Errorf("image %s has no tag", img)
		}

		buf.WriteString("crane pull")
		if nondistro {
			buf.WriteString(" --allow-nondistributable-artifacts")
		}
		if insecure {
			buf.WriteString(" --insecure")
		}
		buf.WriteString(" ")
		buf.WriteString(img)
		buf.WriteString(" ")
		buf.WriteString("images/" + strings.ReplaceAll(ref.Repository, "/", "-") + "-" + ref.Tag + ".tar")
		buf.WriteRune('\n')
	}

	buf.WriteRune('\n')
	buf.WriteString(`tar -czvf images.tar.gz images
`)
	err = os.WriteFile(filepath.Join(outdir, "export-images.sh"), buf.Bytes(), 0o755)
	if err != nil {
		return err
	}

	buf.Reset()
	buf.WriteString(`#!/bin/bash

set -x

TARBALL=${1:-}
REGISTRY=${2:-}

tar -zxvf $TARBALL

`)
	for _, img := range images {
		// crane push images/cluster-ui.tar $IMAGE_REGISTRY/cluster-ui:0.4.16
		ref, err := name.ParseReference(img)
		if err != nil {
			return err
		}
		if ref.Tag == "" {
			return fmt.Errorf("image %s has no tag", img)
		}

		buf.WriteString("crane push")
		if nondistro {
			buf.WriteString(" --allow-nondistributable-artifacts")
		}
		if insecure {
			buf.WriteString(" --insecure")
		}
		buf.WriteString(" ")
		buf.WriteString("images/" + strings.ReplaceAll(ref.Repository, "/", "-") + "-" + ref.Tag + ".tar")
		buf.WriteString(" ")
		buf.WriteString("$IMAGE_REGISTRY/" + ref.Repository + ":" + ref.Tag)
		buf.WriteRune('\n')
	}
	err = os.WriteFile(filepath.Join(outdir, "import-images.sh"), buf.Bytes(), 0o755)
	if err != nil {
		return err
	}

	buf.Reset()
	buf.WriteString(`#!/bin/bash

set -x

TARBALL=${1:-}
REGISTRY=${2:-}

tar -zxvf $TARBALL

`)
	for _, img := range images {
		// crane push images/cluster-ui.tar $IMAGE_REGISTRY/cluster-ui:0.4.16
		ref, err := name.ParseReference(img)
		if err != nil {
			return err
		}
		if ref.Tag == "" {
			return fmt.Errorf("image %s has no tag", img)
		}

		buf.WriteString("crane cp")
		if nondistro {
			buf.WriteString(" --allow-nondistributable-artifacts")
		}
		if insecure {
			buf.WriteString(" --insecure")
		}
		buf.WriteString(" ")
		buf.WriteString(img)
		buf.WriteString(" ")
		buf.WriteString("$IMAGE_REGISTRY/" + ref.Repository + ":" + ref.Tag)
		buf.WriteRune('\n')
	}
	err = os.WriteFile(filepath.Join(outdir, "copy-images.sh"), buf.Bytes(), 0o755)
	if err != nil {
		return err
	}

	return nil
}
