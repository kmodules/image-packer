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

package lib

import (
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

func CheckImageExists(files []string) error {
	images, err := LoadImageList(files)
	if err != nil {
		return err
	}

	var missing []string
	for _, img := range images {
		_, found, err := ImageDigest(img)
		if err != nil || !found {
			missing = append(missing, img)
			continue
		}
		fmt.Println("✔ " + img)
	}

	if len(missing) > 0 {
		fmt.Println("----------------------------------------")
		fmt.Println("Missing Images:")
		fmt.Println(strings.Join(missing, "\n"))
		return fmt.Errorf("missing %d images", len(missing))
	}

	return nil
}

func LoadImageList(files []string) ([]string, error) {
	result := sets.New[string]()
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		var images []string
		err = yaml.Unmarshal(data, &images)
		if err != nil {
			return nil, err
		}
		result.Insert(images...)
	}
	return sets.List(result), nil
}
