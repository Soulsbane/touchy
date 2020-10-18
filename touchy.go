package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/gobuffalo/packr"
	"github.com/pelletier/go-toml"
)

type languageConfig struct {
	Name        string
	Description string
	Extension   string
}

func loadLanguageConfig(languageName string) languageConfig {
	configFileName := "./templates/" + languageName + "/config.toml"
	data, err := ioutil.ReadFile(configFileName)
	config := languageConfig{}

	if err != nil {
		log.Fatal(errors.New("Failed to load config file"))
	}

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(errors.New("Failed to read config file: " + configFileName))
	}

	return config
}

func loadTemplate(name string) string {
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

	return data
}

func buildTemplateList() {
	files, err := ioutil.ReadDir("templates")

	if err != nil {
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func main() {
	var args struct {
		TemplateName string `arg:"positional, required"`
	}

	arg.MustParse(&args)
	fmt.Println(loadTemplate(args.TemplateName))
}
