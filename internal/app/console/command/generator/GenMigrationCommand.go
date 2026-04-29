package generator

import (
	"bytes"
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

type GenMigrationCommand struct {
	path      string
	filename  string
	reference string
	template  string
}

type migrationTemplate struct {
	MigrationStruct string
}

func (c *GenMigrationCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "gen:migration",
		Long: "Migration generator command",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			c.prepare(cmd, args)
			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
	addCommand.PersistentFlags().String("path", "", "Your custom path")
}

func (c *GenMigrationCommand) Handle() {
	content, err := os.ReadFile(c.template)
	if err != nil {
		log.Panicf("could not read template: %v", err)
	}

	tmpl, err := template.New("migration").Parse(string(content))
	if err != nil {
		log.Panicf("could not parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, migrationTemplate{
		MigrationStruct: c.reference,
	})
	if err != nil {
		log.Panicf("could not execute template: %v", err)
	}

	fullPath := fmt.Sprintf("%s/%s", c.path, c.filename)
	_, err = os.Stat(fullPath)
	if err == nil {
		log.Panicf("file path already exists!! %s", fullPath)
	}

	err = os.WriteFile(fullPath, buf.Bytes(), 0644)
	if err != nil {
		log.Panicf("could not write file!! %s", err)
	}

	fmt.Printf("Model %s created successfully.\n", fullPath)
}

func (c *GenMigrationCommand) prepare(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Panicf("Please enter migration filename!")
	}

	c.template = "./stubs/migration.stub"
	c.path, _ = cmd.Flags().GetString("path")
	if c.path == "" {
		c.path = "./internal/app/database/migration"
	}

	c.reference = fmt.Sprintf("%s_%d", strings.Title(args[0]), time.Now().UnixNano()/1000)
	c.filename = fmt.Sprintf("%s.go", c.reference)
}
