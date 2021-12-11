package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/packr"
)

type languageConfig struct {
	Name            string
	DefaultFileName string
	Description     string
	Extension       string
}

func loadLanguageConfig(languageName string) languageConfig {
	exePath, _ := os.Executable()
	configFileName := filepath.Join(filepath.Dir(exePath), "../../internal/generator/templates/", languageName, "/config.toml")

	data, err := ioutil.ReadFile(configFileName)
	config := languageConfig{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file: " + configFileName))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + configFileName))
	}

	return config
}

func loadTemplate(name string) (string, languageConfig) {
	language := name
	template := "default"
	box := packr.NewBox("./templates")

	if strings.Contains(name, ".") {
		var parts = strings.Split(name, ".")

		language = parts[0]
		template = parts[1]

		if template == "" {
			template = "default"
		}
	}

	config := loadLanguageConfig(language)
	templateName := language + "/" + template + "." + config.Extension
	data, err := box.FindString(templateName)

	if err != nil {
		//log.Fatal(errors.New("That template does not exist: " + config.Name + " => " + template))
		log.Fatal(errors.New("That template does not exist: " + name))
	}

	return data, config
}

func ListTemplates() {
	files, err := ioutil.ReadDir("templates")

	if err != nil {
		fmt.Print(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

// CreateFileFromTemplate Creates a template
func CreateFileFromTemplate(customFileName string, languageName string) {
	var fileName string
	template, config := loadTemplate(languageName)
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
