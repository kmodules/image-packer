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
	"os"
	"path/filepath"

	"kmodules.xyz/image-packer/pkg/lib"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewCmdListImages() *cobra.Command {
	var (
		rootDir string
		outDir  string
	)
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List all Docker images in a directory",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			images, err := lib.ListImages(rootDir)
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(images)
			if err != nil {
				return err
			}

			filename := filepath.Join(outDir, "imagelist.yaml")
			err = os.WriteFile(filename, data, 0o644)
			return err
		},
	}

	cmd.Flags().StringVar(&rootDir, "root-dir", "", "Root directory")
	cmd.Flags().StringVar(&rootDir, "output-dir", "", "Output directory")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "root-dir")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "output-dir")

	return cmd
}
