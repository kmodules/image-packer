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
	"strings"

	"kmodules.xyz/image-packer/pkg/lib"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// NewCmdReplaceImageDigest creates a new cobra command to replace image keys with digests in a YAML file
func NewCmdReplaceImageDigest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-image-digest <input.yaml> <output.yaml>",
		Short: "Replace image tag with image digest in a YAML file",
		Long: `Recursively traverses a YAML document and replaces image tags with
their corresponding digests for "image" and "containerImage" keys.
Tags are stripped so only the digest remains (e.g., nginx:1.21 becomes
nginx@sha256:abc...).

The input and output files can be the same for in-place editing.

Examples:
  image-packer replace-image-digest input.yaml output.yaml
  image-packer replace-image-digest app.yaml app.yaml`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]
			outputFile := args[1]
			data, err := os.ReadFile(inputFile)
			if err != nil {
				return fmt.Errorf("failed to read input file: %w", err)
			}
			var m map[string]any
			if err := yaml.Unmarshal(data, &m); err != nil {
				return fmt.Errorf("failed to unmarshal YAML: %w", err)
			}
			changed, err := replaceImageKeysWithDigest(m)
			if err != nil {
				return err
			}
			out, err := yaml.Marshal(changed)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML: %w", err)
			}
			return os.WriteFile(outputFile, out, 0o644)
		},
	}
	return cmd
}

// replaceImageKeysWithDigest recursively replaces image and containerImage values with their digests
func replaceImageKeysWithDigest(obj any) (any, error) {
	switch v := obj.(type) {
	case map[string]any:
		for key, val := range v {
			switch key {
			case "image", "containerImage":
				newVal, err := replaceImageWithDigest(val)
				if err != nil {
					return nil, err
				}
				v[key] = newVal
			default:
				newVal, err := replaceImageKeysWithDigest(val)
				if err != nil {
					return nil, err
				}
				v[key] = newVal
			}
		}
		return v, nil
	case []any:
		for i, item := range v {
			newItem, err := replaceImageKeysWithDigest(item)
			if err != nil {
				return nil, err
			}
			v[i] = newItem
		}
		return v, nil
	default:
		return v, nil
	}
}

func replaceImageWithDigest(val any) (any, error) {
	img, ok := val.(string)
	if !ok || containsDigest(img) {
		return val, nil
	}
	digest, found, err := lib.ImageDigest(img)
	if err != nil {
		return nil, fmt.Errorf("failed to get digest for %s: %w", img, err)
	}
	if found {
		return fmt.Sprintf("%s@%s", stripTag(img), digest), nil
	}
	return val, nil
}

// containsDigest returns true if the image string contains an '@' (digest)
func containsDigest(img string) bool {
	return strings.Contains(img, "@sha256:")
}

// stripTag removes the tag from an image reference (e.g., "nginx:1.21" -> "nginx")
func stripTag(img string) string {
	if idx := strings.LastIndex(img, ":"); idx != -1 {
		if slashIdx := strings.LastIndex(img, "/"); slashIdx < idx {
			return img[:idx]
		}
	}
	return img
}
