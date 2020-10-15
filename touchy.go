package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/gobuffalo/packr"
)

func loadTemplate(templateName string) string {
	language := templateName
	template := "default"

	if strings.Contains(templateName, ".") {
		var parts = strings.Split(templateName, ".")

		language = parts[0]
		template = parts[1]

		if template == "" {
			template = "default"
		}
	}

	box := packr.NewBox("./templates")
	// TODO: Better way to handle file extensions. Perhaps a config file in each template folder.
	data, err := box.FindString(language + "/" + template + "." + language)

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
