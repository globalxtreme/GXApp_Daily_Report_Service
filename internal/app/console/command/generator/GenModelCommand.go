package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"
)

type GenModelCommand struct {
	migrationPath     string
	modelPath         string
	migrationTemplate string
	modelTemplate     string
	reference         string
	migration         bool
}

type modelTemplate struct {
	ModelStruct string
	TableName   string
}

func (c *GenModelCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "gen:model",
		Long: "Model generator command",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()
			c.prepare(cmd, args)
			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
	addCommand.PersistentFlags().Bool("migration", false, "Generate migration file")
}

func (c *GenModelCommand) Handle() {
	if err := c.createFile(c.modelTemplate, modelTemplate{ModelStruct: c.reference, TableName: c.reference}, c.modelPath, c.reference); err != nil {
		log.Println("Error generating model:", err)
	}

	c.generateMigration()
}

func (c *GenModelCommand) generateMigration() {
	if c.migration {
		migrationFileName := fmt.Sprintf("%s_%d", c.reference, time.Now().UnixNano()/1000)
		if err := c.createFile(c.migrationTemplate, migrationTemplate{MigrationStruct: c.reference}, c.migrationPath, migrationFileName); err != nil {
			log.Println("Error generating migration:", err)
			return
		}
	}
}

func (c *GenModelCommand) prepare(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Printf("Please enter migration filename!")
		return
	}

	c.modelTemplate = "./stubs/model.stub"
	c.migrationTemplate = "./stubs/migration.stub"

	c.modelPath = "./internal/pkg/model"
	c.migrationPath = "./internal/app/database/migration"

	c.migration, _ = cmd.Flags().GetBool("migration")
	c.reference = strings.Title(args[0])
}

func (c *GenModelCommand) createFile(templatePath string, data interface{}, targetPath string, filename string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read template: %v", err)
	}

	tmpl, err := template.New("file").Parse(string(content))
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}

	fullPath := fmt.Sprintf("%s/%s.go", targetPath, filename)
	return c.writeFileIfNotExists(fullPath, buf.Bytes())
}

func (c *GenModelCommand) writeFileIfNotExists(filename string, content []byte) error {
	if _, err := os.Stat(filename); err == nil {
		log.Printf("File path already exists!! %s", filename)
		return fmt.Errorf("file already exists")
	}

	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		log.Printf("Could not write file!! %s", err)
		return err
	}

	log.Printf("File %s created successfully.\n", filename)
	return nil
}

func (c *GenModelCommand) camelToSnake(s string) string {
	var result []rune

	for i, r := range s {
		// Jika karakter adalah huruf kapital dan bukan yang pertama, tambahkan underscore
		if unicode.IsUpper(r) {
			// Menambahkan underscore jika bukan karakter pertama
			if i > 0 {
				result = append(result, '_')
			}
			// Menambahkan huruf kecil
			result = append(result, unicode.ToLower(r))
		} else {
			// Menambahkan karakter kecil lainnya tanpa perubahan
			result = append(result, r)
		}
	}

	return string(result) + "s"
}
