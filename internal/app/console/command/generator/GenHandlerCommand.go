package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
)

type GenHandlerCommand struct {
	path      string
	filename  string
	reference string
	template  string
	resource  bool
}

type handlerTemplate struct {
	HandlerStruct string
	HasResource   bool
}

func (c *GenHandlerCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "gen:handler",
		Long: "Command to generate a handler",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			c.prepare(cmd, args)
			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
	addCommand.PersistentFlags().String("type", "", "Type of the handler (web/mobile)")
	addCommand.PersistentFlags().Bool("resource", false, "Use resource mode (optional flag)")
}

func (c *GenHandlerCommand) Handle() {
	content, err := os.ReadFile(c.template)
	if err != nil {
		log.Printf("Failed to read template: %v", err)
		return
	}

	tmpl, err := template.New("handler").Parse(string(content))
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, handlerTemplate{
		HandlerStruct: c.reference,
		HasResource:   c.resource,
	})

	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		return
	}

	fullPath := fmt.Sprintf("%s/%s", c.path, c.filename)
	_, err = os.Stat(fullPath)
	if err == nil {
		log.Printf("File already exists at path: %s", fullPath)
		return
	}

	err = os.WriteFile(fullPath, buf.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to write file: %s", err)
		return
	}

	fmt.Printf("Handler %s successfully created.\n", fullPath)
}

func (c *GenHandlerCommand) prepare(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Printf("Please provide a handler name!")
		return
	}

	c.template = "./stubs/handler.stub"

	cType, _ := cmd.Flags().GetString("type")
	if cType == "" {
		cType = "web"
	}

	c.path, _ = cmd.Flags().GetString("path")
	if c.path == "" {
		c.path = fmt.Sprintf("./internal/app/api/%s/handler", cType)
	}

	title := strings.Title(args[0])

	c.reference = fmt.Sprintf("%sHandler", title)
	c.filename = fmt.Sprintf("%s.go", c.reference)

	c.resource, _ = cmd.Flags().GetBool("resource")
}
