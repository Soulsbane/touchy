package templates

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates
var templatesDir embed.FS

type Languages struct {
	languages []Language
}

func New() *Languages {
	return &Languages{}
}

func (g *Languages) Load(name string) (string, Language) {
	language := name
	template := "default"

	if strings.Contains(name, ".") {
		var parts = strings.Split(name, ".")

		language = parts[0]
		template = parts[1]

		if template == "" {
			template = "default"
		}
	}

	configPath := filepath.Join("templates", language, "config.toml")
	config := loadLanguageConfigFile(configPath)
	templateName := filepath.Join("templates", language, template, template+".template")
	data, err := templatesDir.ReadFile(templateName)

	if err != nil {
		//log.Fatal(errors.New("That template does not exist: " + config.Name + " => " + template))
		log.Fatal(errors.New("That template does not exist: " + name))
	}

	return string(data), config
}

func (g *Languages) List(listArg string) {
	languageDirs, err := templatesDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	for _, languageDir := range languageDirs {
		if languageDir.IsDir() {
			fmt.Println("languageDir: " + languageDir.Name())
			templates, err := templatesDir.ReadDir(filepath.Join("templates", languageDir.Name()))

			if err != nil {
				panic(err)
			}

			for _, template := range templates {
				if template.IsDir() {
					fmt.Println("Filename: ", template.Name())
					// TODO: Read info.toml
				}
			}
		}

		fmt.Println()
	}
}

// CreateFileFromTemplate Creates a template
func (g *Languages) CreateFileFromTemplate(customFileName string, languageName string) {
	var fileName string
	template, config := g.Load(languageName)
	currentDir, _ := os.Getwd()

	if customFileName == "DefaultFileName" {
		fileName = config.DefaultFileName
	} else {
		fileName = customFileName
	}

	file, err := os.Create(filepath.Join(currentDir, "/"+fileName+"."+config.Extension))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	file.WriteString(template)
	//fmt.Println(template)
}
