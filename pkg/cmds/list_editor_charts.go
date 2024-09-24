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
	"fmt"
	"os"
	"path/filepath"

	"kmodules.xyz/resource-metadata/hub/resourceeditors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

func NewCmdListEditorCharts() *cobra.Command {
	var (
		apiGroups []string
		outDir    string
	)
	cmd := &cobra.Command{
		Use:                   "list-editor-charts",
		Short:                 "List all editor charts",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			images, err := ListEditorCharts(sets.New[string](apiGroups...))
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(images)
			if err != nil {
				return err
			}

			filename := filepath.Join(outDir, "editor-charts.yaml")
			err = os.WriteFile(filename, data, 0o644)
			return err
		},
	}

	cmd.Flags().StringSliceVar(&apiGroups, "apiGroup", nil, "API Group to be included in the output")
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Output directory")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "output-dir")

	return cmd
}

func ListEditorCharts(groups sets.Set[string]) ([]string, error) {
	images := sets.New[string]()

	for _, ed := range resourceeditors.List() {
		if !groups.Has(ed.Spec.Resource.Group) {
			continue
		}
		if ed.Spec.UI == nil {
			continue
		}
		if ed.Spec.UI.Options != nil {
			images.Insert(fmt.Sprintf("ghcr.io/appscode-charts/%s:%s", ed.Spec.UI.Options.Name, ed.Spec.UI.Options.Version))
		}
		if ed.Spec.UI.Editor != nil {
			images.Insert(fmt.Sprintf("ghcr.io/appscode-charts/%s:%s", ed.Spec.UI.Editor.Name, ed.Spec.UI.Editor.Version))
		}
		for _, action := range ed.Spec.UI.Actions {
			for _, item := range action.Items {
				if item.Editor != nil {
					images.Insert(fmt.Sprintf("ghcr.io/appscode-charts/%s:%s", item.Editor.Name, item.Editor.Version))
				}
			}
		}
	}

	return sets.List(images), nil
}
