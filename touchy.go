package main

import (
	"fmt"
	"io/ioutil"
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
	data, err := ioutil.ReadFile("./templates/" + languageName + "/config.toml")

	if err != nil {
		fmt.Println(err)
	}

	config := languageConfig{}
	err = toml.Unmarshal(data, &config)

	if err != nil {
		fmt.Println(err)
	}

	return config
}

func loadTemplate(templateName string) string {
	language := templateName
	template := "default"
	config := languageConfig{}

	if strings.Contains(templateName, ".") {
		var parts = strings.Split(templateName, ".")

		language = parts[0]
		template = parts[1]

		if template == "" {
			template = "default"
		}
	}

	config = loadLanguageConfig(language)

	box := packr.NewBox("./templates")
	data, err := box.FindString(language + "/" + template + "." + config.Extension)

	if err != nil {
		fmt.Println(err)
		return ""
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
