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
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
)

func NewCmdAceUp() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:                   "ace-up",
		Short:                 "Update ace.yaml",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			aceImageMap := map[string]string{}
			err := LoadLatestImageMap(filepath.Join(dir, "imagelist.yaml"), aceImageMap)
			if err != nil {
				return err
			}

			tagB3, found := aceImageMap["ghcr.io/appscode/b3"]
			if !found {
				return fmt.Errorf("no b3 image found in imagelist.yaml")
			}

			featureMap := map[string]string{}
			err = LoadLatestImageMap(filepath.Join(dir, "feature-charts.yaml"), featureMap)
			if err != nil {
				return err
			}
			files := []string{
				filepath.Join(dir, "imagelist.yaml"),
				fmt.Sprintf("https://github.com/kluster-manager/installer/raw/%s/catalog/imagelist.yaml", tagB3),
			}
			if tag, ok := featureMap["ghcr.io/appscode-charts/kubedb"]; ok {
				files = append(files, fmt.Sprintf("https://github.com/kubedb/installer/raw/%s/catalog/imagelist.yaml", tag))
			}
			if tag, ok := featureMap["ghcr.io/appscode-charts/kubestash"]; ok {
				files = append(files, fmt.Sprintf("https://github.com/kubestash/installer/raw/%s/catalog/imagelist.yaml", tag))
			}
			if tag, ok := featureMap["ghcr.io/appscode-charts/kubevault"]; ok {
				files = append(files, fmt.Sprintf("https://github.com/kubestash/installer/raw/%s/catalog/imagelist.yaml", tag))
			}
			if tag, ok := featureMap["ghcr.io/appscode-charts/capi-catalog"]; ok {
				files = append(files, fmt.Sprintf("https://github.com/kluster-api/installer/raw/%s/catalog/imagelist.yaml", tag))
			}

			var images map[string]string
			if imageList, err := GenerateImageList(files, true); err != nil {
				return err
			} else {
				images = ToImageMap(imageList)
			}

			aceMap, err := LoadImageMap(filepath.Join(dir, "ace.yaml"))
			if err != nil {
				return err
			}

			// update ace with images
			for img, tag := range images {
				tags := aceMap[img]
				switch len(tags) {
				case 0:
					// skip
				case 1:
					aceMap[img] = []string{tag}
				default:
					aceMap[img] = sets.List(sets.New[string](aceMap[img]...).Insert(tag))
				}
			}

			return write(ToImageList2(aceMap), filepath.Join(dir, "ace.yaml"))
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Directory containing ace.yaml and imagelist.yaml files")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "dir")

	return cmd
}

func ToImageList(in map[string]string) []string {
	images := make([]string, 0, len(in))
	for img, tag := range in {
		images = append(images, img+":"+tag)
	}
	return images
}

func ToImageList2(in map[string][]string) []string {
	images := make([]string, 0, len(in))
	for img, tags := range in {
		for _, tag := range tags {
			images = append(images, img+":"+tag)
		}
	}
	sort.Strings(images)
	return images
}

func ToImageMap(in []string) map[string]string {
	images := make(map[string]string, len(in))
	for _, repo := range in {
		if img, tag, ok := strings.Cut(repo, ":"); ok {
			images[img] = tag
		}
	}
	return images
}
