package generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) loadTemplate(name string) (string, Language) {
	language := name
	template := "default"
	box := packr.New("Templates", "./templates")

	if strings.Contains(name, ".") {
		var parts = strings.Split(name, ".")

		language = parts[0]
		template = parts[1]

		if template == "" {
			template = "default"
		}
	}

	config := LoadLanguageConfigFile(language)
	templateName := language + "/" + template + "." + config.Extension
	data, err := box.FindString(templateName)

	if err != nil {
		//log.Fatal(errors.New("That template does not exist: " + config.Name + " => " + template))
		log.Fatal(errors.New("That template does not exist: " + name))
	}

	return data, config
}

func (g *Generator) ListTemplates() {
	box := packr.New("Templates", "./templates")

	box.Walk(func(path string, f file.File) error {
		fmt.Println(path)
		return nil
	})
}

// CreateFileFromTemplate Creates a template
func (g *Generator) CreateFileFromTemplate(customFileName string, languageName string) {
	var fileName string
	template, config := g.loadTemplate(languageName)
	currentDir, _ := os.Getwd()

	if customFileName == "DefaultFileName" {
		fileName = config.DefaultFileName
	} else {
		fileName = customFileName
	}

	file, err := os.Create(filepath.Join(currentDir, "/"+fileName+"."+config.Extension))

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(template)
	//fmt.Println(template)
}
