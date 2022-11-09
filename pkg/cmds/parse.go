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
