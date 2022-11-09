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
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
)

func NewCmdGenerateScripts() *cobra.Command {
	var (
		nondistro bool
		insecure  bool
	)
	cmd := &cobra.Command{
		Use:                   "generate-scripts",
		Short:                 "Generate export/import scripts",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateScripts(args, nondistro, insecure)
		},
	}
	cmd.Flags().BoolVar(&nondistro, "allow-nondistributable-artifacts", nondistro, "Allow pushing non-distributable (foreign) layers")
	cmd.Flags().BoolVar(&insecure, "insecure", insecure, "Allow image references to be fetched without TLS")

	return cmd
}

func generateScripts(args []string, nondistro, insecure bool) error {
	dir, manifest, err := readManifest(args)
	if err != nil {
		return err
	}
	outdir, err := os.Getwd()
	if err != nil {
		return err
	}
	return GenerateScripts(dir, manifest, outdir, nondistro, insecure)
}

func GenerateScripts(dir string, manifest bool, outdir string, nondistro, insecure bool) error {
	images, err := ListImages(dir, manifest)
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

		buf.WriteString("crane pull")
		if nondistro {
			buf.WriteString(" --allow-nondistributable-artifacts")
		}
		if insecure {
			buf.WriteString(" --insecure")
		}
		buf.WriteString(" ")
		buf.WriteString(img)

		ref, err := name.ParseReference(img)
		if err != nil {
			return err
		}
		_, bin := path.Split(ref.Context().RepositoryStr())

		buf.WriteString(" ")
		buf.WriteString("images/" + bin + "-" + ref.Identifier() + ".tar")

		buf.WriteRune('\n')
	}

	buf.WriteRune('\n')
	buf.WriteString(`tar -czvf images.tar.gz images
`)
	err = os.WriteFile(filepath.Join(outdir, "export-images.sh"), buf.Bytes(), 0o644)
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
		// crane push images/cluster-ui.tar $REGISTRY/cluster-ui:0.4.16

		buf.WriteString("crane push")
		if nondistro {
			buf.WriteString(" --allow-nondistributable-artifacts")
		}
		if insecure {
			buf.WriteString(" --insecure")
		}

		ref, err := name.ParseReference(img)
		if err != nil {
			return err
		}
		_, bin := path.Split(ref.Context().RepositoryStr())

		buf.WriteString(" ")
		buf.WriteString("images/" + bin + "-" + ref.Identifier() + ".tar")

		buf.WriteString(" ")
		buf.WriteString("$REGISTRY/" + bin + ":" + ref.Identifier())

		buf.WriteRune('\n')
	}
	err = os.WriteFile(filepath.Join(outdir, "import-images.sh"), buf.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}
