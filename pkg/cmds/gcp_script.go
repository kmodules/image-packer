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
	"path"
	"path/filepath"
	"strings"

	"kmodules.xyz/go-containerregistry/name"

	"github.com/spf13/cobra"
)

func NewCmdGenerateGCPScript() *cobra.Command {
	var (
		files     []string
		nondistro bool
		insecure  bool
		outDir    string
	)
	cmd := &cobra.Command{
		Use:                   "generate-gcp-script",
		Short:                 "Generate GCP Marketplace image syncer script",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenerateGCPScript(files, outDir, nondistro, insecure)
		},
	}
	cmd.Flags().StringSliceVar(&files, "src", files, "List of source files (http url or local file)")
	cmd.Flags().BoolVar(&nondistro, "allow-nondistributable-artifacts", nondistro, "Allow pushing non-distributable (foreign) layers")
	cmd.Flags().BoolVar(&insecure, "insecure", insecure, "Allow image references to be fetched without TLS")
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Output directory")

	return cmd
}

var gcpImageMap = map[string]string{
	"defaultbackend-amd64":               "ingress-nginx-defaultbackend",
	"fluxcd/helm-controller":             "flux-helm-controller",
	"fluxcd/kustomize-controller":        "flux-kustomize-controller",
	"fluxcd/notification-controller":     "flux-notification-controller",
	"fluxcd/source-controller":           "flux-source-controller",
	"ingress-nginx/controller":           "ingress-nginx-controller",
	"ingress-nginx/kube-webhook-certgen": "ingress-nginx-kube-webhook-certgen",
	"kedacore/http-add-on-interceptor":   "keda-http-add-on-interceptor",
	"kedacore/http-add-on-operator":      "keda-http-add-on-operator",
	"kedacore/http-add-on-scaler":        "keda-http-add-on-scaler",
	"prometheus/node-exporter":           "prometheus-node-exporter",
	"sig-storage/livenessprobe":          "csi-driver-livenessprobe",
}

func GenerateGCPScript(files []string, outdir string, nondistro, insecure bool) error {
	images, err := GenerateImageList(files, true)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(`#!/bin/bash

set -x

if [ -z "${IMAGE_REGISTRY}" ]; then
	echo "IMAGE_REGISTRY is not set"
	exit 1
fi
if [ -z "${TAG}" ]; then
	echo "TAG is not set"
	exit 1
fi

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

		repo := ref.Repository
		if repo == "prometheus-operator/prometheus-operator" {
			continue
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

		if strings.HasPrefix(ref.Repository, "library/") {
			repo = ref.Repository[len("library/"):]
		}
		if v, found := gcpImageMap[repo]; found {
			repo = v
		}
		_, bin := path.Split(repo)

		buf.WriteString("$IMAGE_REGISTRY/" + bin + ":$TAG")
		buf.WriteRune('\n')
	}
	err = os.WriteFile(filepath.Join(outdir, "sync-gcp-mp-images.sh"), buf.Bytes(), 0o755)
	if err != nil {
		return err
	}

	return nil
}
