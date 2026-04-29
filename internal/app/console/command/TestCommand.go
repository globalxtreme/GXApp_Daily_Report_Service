package command

import (
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
	"os"
)

type TestCommand struct{}

func (class *TestCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "dev-test",
		Long: "Development Test Command",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			class.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
}

func (class *TestCommand) Handle() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Working Directory:", dir)
}
