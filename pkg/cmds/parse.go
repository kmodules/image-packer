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
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
)

func NewCmdParseImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "parse",
		Short:                 "Parse image reference",
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ParseImage(args)
		},
	}

	return cmd
}

func ParseImage(args []string) error {
	if len(args) == 0 {
		return errors.New("missing input")
	} else if len(args) > 1 {
		return errors.New("too many inputs")
	}

	img := args[0]
	ref, err := name.ParseReference(img)
	if err != nil {
		return err
	}

	fmt.Println("Image:", ref.String())
	fmt.Println("Full Image Name:", ref.Name())
	fmt.Println("Full Repository Name:", ref.Context().Name())
	fmt.Println("Registry:", ref.Context().RegistryStr())
	fmt.Println("Repository:", ref.Context().RepositoryStr())
	fmt.Println("Identifier:", ref.Identifier())

	return nil
}
