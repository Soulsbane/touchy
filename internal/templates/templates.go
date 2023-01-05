package templates

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//go:embed templates
var templatesDir embed.FS

const INFO_FILENAME = "info.toml"

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
			defaultConfig := CommonConfig{
				Name:                  languageDir.Name(),
				DefaultOutputFileName: "<Unknown>",
				Description:           "<Unknown>",
			}

			language.templateConfigs = make(map[string]CommonConfig)
			infoPath := filepath.Join("templates", languageDir.Name(), INFO_FILENAME)

			language.infoConfig, err = loadInfoFile(infoPath)

			// If there is no info file, use the default config so it at least shows up in the list command
			if err != nil {
				language.infoConfig = defaultConfig
			}

			templates, err := templatesDir.ReadDir(filepath.Join("templates", languageDir.Name()))

			if err != nil {
				panic(err) // TODO: Handle this better?
			}

			for _, template := range templates {
				if template.IsDir() {
					configPath := filepath.Join("templates", languageDir.Name(), template.Name(), INFO_FILENAME)
					config, err := loadInfoFile(configPath)

					if err != nil {
						language.templateConfigs[template.Name()] = defaultConfig
					} else {
						language.templateConfigs[template.Name()] = config
					}
				}
			}

			g.languages[languageDir.Name()] = language
		}
	}
}

func (g *Templates) GetLanguageTemplateFor(languageName string, templateName string) (string, CommonConfig) {
	for name, language := range g.languages {
		if name == languageName {
			for template, config := range language.templateConfigs {
				if template == templateName {
					return g.loadTemplateFile(languageName, templateName), config
				}
			}
		}
	}

	return "", CommonConfig{}
}

func (g *Templates) loadTemplateFile(language string, template string) string {
	templateName := filepath.Join("templates", language, template, template+".template")
	data, err := templatesDir.ReadFile(templateName)

	if err != nil {
		//log.Fatal(errors.New("That template does not exist: " + config.Name + " => " + template))
		log.Fatal(errors.New("That template does not exist: " + templateName))
	}

	return string(data)

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
func (g *Templates) CreateFileFromTemplate(languageName string, templateName string, customFileName string) {
	var fileName string
	template, config := g.GetLanguageTemplateFor(languageName, templateName)
	currentDir, _ := os.Getwd()

	if customFileName == "DefaultOutputFileName" {
		fileName = config.DefaultOutputFileName
	} else {
		fileName = customFileName
	}

	if fileName == "" {
		fmt.Println("Failed to template file. No file name was provided!")
	} else {
		file, err := os.Create(filepath.Join(currentDir, fileName))

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		file.WriteString(template)
	}
}
