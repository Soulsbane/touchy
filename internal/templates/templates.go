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

const CONFIG_FILENAME = "config.toml"

type Templates struct {
	languages map[string]Language
}

func New() *Templates {
	var templates Templates

	templates.languages = make(map[string]Language)
	templates.findTemplates()

	return &templates
}

func (g *Templates) findTemplates() {
	languageDirs, err := templatesDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	for _, languageDir := range languageDirs {
		if languageDir.IsDir() {
			var language Language

			language.templateConfigs = make(map[string]CommonConfig)
			infoPath := filepath.Join("templates", languageDir.Name(), "info.toml")

			language.infoConfig = loadLanguageInfoFile(infoPath)
			templates, err := templatesDir.ReadDir(filepath.Join("templates", languageDir.Name()))

			if err != nil {
				panic(err)
			}

			for _, template := range templates {
				if template.IsDir() {
					configPath := filepath.Join("templates", languageDir.Name(), template.Name(), "config.toml")
					config := loadLanguageConfigFile(configPath)
					language.templateConfigs[template.Name()] = config
				}
			}

			g.languages[languageDir.Name()] = language
		}
	}
}

func (g *Templates) Load(name string) (string, CommonConfig) {
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

	configPath := filepath.Join("templates", language, template, "config.toml")
	config := loadLanguageConfigFile(configPath)
	templateName := filepath.Join("templates", language, template, template+".template")
	data, err := templatesDir.ReadFile(templateName)

	if err != nil {
		//log.Fatal(errors.New("That template does not exist: " + config.Name + " => " + template))
		log.Fatal(errors.New("That template does not exist: " + name))
	}

	return string(data), config
}

func (g *Templates) List(listArg string) {
	for languageName, language := range g.languages {
		fmt.Println("Language: ", languageName)

		for templateName, config := range language.templateConfigs {
			fmt.Println("Template Name: ", templateName)
			fmt.Println(config.Description)
			fmt.Println()
		}
	}
}

// CreateFileFromTemplate Creates a template
func (g *Templates) CreateFileFromTemplate(customFileName string, languageName string) {
	var fileName string
	template, config := g.Load(languageName)
	currentDir, _ := os.Getwd()

	if customFileName == "DefaultFileName" {
		fileName = config.DefaultOutputFileName
	} else {
		fileName = customFileName
	}

	file, err := os.Create(filepath.Join(currentDir, "/"+fileName+"."+config.Extension))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	file.WriteString(template)
}
