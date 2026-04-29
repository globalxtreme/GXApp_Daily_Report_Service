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

type GenParserCommand struct {
	path      string
	filename  string
	reference string
	template  string
	model     string
}

type parserTemplate struct {
	ParserStruct string
	Model        string
	HasModel     bool
}

func (c *GenParserCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "gen:parser",
		Long: "Command to generate a parser",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			c.prepare(cmd, args)
			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
	addCommand.PersistentFlags().String("path", "", "Custom path for the parser file")
	addCommand.PersistentFlags().Bool("model", false, "Use model mode (optional flag)")
}

func (c *GenParserCommand) Handle() {
	content, err := os.ReadFile(c.template)
	if err != nil {
		log.Printf("Failed to read template: %v", err)
		return
	}

	tmpl, err := template.New("parser").Parse(string(content))
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return
	}

	var buf bytes.Buffer

	if c.model == "" {
		err = tmpl.Execute(&buf, parserTemplate{
			ParserStruct: c.reference,
			HasModel:     false,
		})
	} else {
		err = tmpl.Execute(&buf, parserTemplate{
			ParserStruct: c.reference,
			Model:        c.model,
			HasModel:     true,
		})
	}

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

	fmt.Printf("Parser %s successfully created.\n", fullPath)
}

func (c *GenParserCommand) prepare(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Printf("Please provide a parser name!")
		return
	}

	c.template = "./stubs/parser.stub"
	c.path, _ = cmd.Flags().GetString("path")
	if c.path == "" {
		c.path = "./internal/pkg/parser"
	}

	title := strings.Title(args[0])
	model, _ := cmd.Flags().GetBool("model")
	if model {
		c.model = title
	}

	c.reference = fmt.Sprintf("%sParser", title)
	c.filename = fmt.Sprintf("%s.go", c.reference)
}
