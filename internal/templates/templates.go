package templates

import (
	"embed"
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path/filepath"

	"github.com/Soulsbane/touchy/internal/path"
	"github.com/alecthomas/chroma/quick"
	"github.com/jedib0t/go-pretty/v6/table"
)

//go:embed templates
var embedsDir embed.FS

type Language struct {
	// dirName         string                  // Name of the directory under the template's directory.
	infoConfig      infofile.InfoFile   // Each language has a config file in its root directory call config.toml
	templateConfigs []infofile.InfoFile // A list of all the templates in the language directory. The key is the template dir name.
}

type Templates struct {
	languages map[string]Language // Map of all languages in the templates directory. Key is the language name.
}

func New() *Templates {
	var templates Templates

	templates.languages = make(map[string]Language)
	templates.findTemplates()

	return &templates
}

func (g *Templates) findTemplates() {
	languageDirs, err := embedsDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	for _, languageDir := range languageDirs {
		if languageDir.IsDir() {
			var language Language
			defaultConfig := infofile.InfoFile{
				Name:                  languageDir.Name(),
				DefaultOutputFileName: "<Unknown>",
				Description:           "<Unknown>",
			}

			infoPath := filepath.Join("templates", languageDir.Name(), infofile.DefaultFileName)

			language.infoConfig, err = infofile.Load(infoPath, embedsDir)

			// If there is no info file, use the default config so it at least shows up in the list command
			if err != nil {
				language.infoConfig = defaultConfig
			}

			templates, err := embedsDir.ReadDir(filepath.Join("templates", languageDir.Name()))

			if err != nil {
				panic(err) // TODO: Handle this better?
			}

			for _, template := range templates {
				if template.IsDir() {
					configPath := filepath.Join("templates", languageDir.Name(), template.Name(), infofile.DefaultFileName)
					config, err := infofile.Load(configPath, embedsDir)

					if err != nil {
						language.templateConfigs = append(language.templateConfigs, defaultConfig)
					} else {
						language.templateConfigs = append(language.templateConfigs, config)
					}
				}
			}

			g.languages[languageDir.Name()] = language
		}
	}
}

func (g *Templates) HasLanguage(languageName string) bool {
	_, found := g.languages[languageName]
	return found
}

func (g *Templates) HasTemplate(languageName string, templateName string) bool {
	language, foundLanguage := g.languages[languageName]

	if foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.Name == templateName })

		if idx >= 0 {
			return true
		}

		return false
	}

	return false
}

func (g *Templates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	language, foundLanguage := g.languages[languageName]

	if foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.Name == templateName })

		if idx >= 0 {
			return g.loadTemplateFile(languageName, templateName), language.templateConfigs[idx]
		}
	}

	return "", infofile.InfoFile{}
}

func (g *Templates) loadTemplateFile(language string, template string) string {
	templateName := filepath.Join("templates", language, template, template+".template")
	data, err := embedsDir.ReadFile(templateName)

	if err != nil {
		log.Fatal(errors.New("That template does not exist: " + templateName))
	}

	return string(data)

}

func (g *Templates) List(listArg string) {
	if language, found := g.languages[listArg]; found {
		g.listLanguageTemplates(language)
	} else if listArg == "all" {
		g.listAllLanguages()
	} else {
		fmt.Println("That language could not be found! Use 'list all' to see all available languages.")
	}
}

func (g *Templates) listLanguageTemplates(language Language) {
	outputTable := table.NewWriter()

	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"Template Name", "Description"})

	for templateName, config := range language.templateConfigs {
		outputTable.AppendRow(table.Row{templateName, config.Description})
	}

	outputTable.SetStyle(table.StyleRounded)
	outputTable.Style().Options.SeparateRows = true
	outputTable.Render()
}

func (g *Templates) listAllLanguages() {
	outputTable := table.NewWriter()

	outputTable.SetOutputMirror(os.Stdout)
	outputTable.AppendHeader(table.Row{"Name", "Description", "Default Output File Name"})

	for languageName, language := range g.languages {
		info := language.infoConfig
		outputTable.AppendRow(table.Row{languageName, info.Description, info.DefaultOutputFileName})
	}

	outputTable.SetStyle(table.StyleRounded)
	outputTable.Style().Options.SeparateRows = true
	outputTable.Render()
}

func (g *Templates) ShowTemplate(languageName string, templateName string) {
	if language, languageFound := g.languages[languageName]; languageFound {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.Name == templateName })

		if idx >= 0 {
			sourceCode := g.loadTemplateFile(languageName, templateName)

			// Formatters: terminal, terminal8, terminal16, terminal256, terminal16m
			// Styles: https://github.com/alecthomas/chroma/tree/master/styles
			err := quick.Highlight(os.Stdout, sourceCode, language.templateConfigs[idx].DefaultOutputFileName, "terminal256", "monokai")

			if err != nil {
				fmt.Println(err)
			}

		} else {
			fmt.Println("That template does not exist: ", templateName)
		}
	} else {
		fmt.Println("That language does not exist: ", languageName)
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
		fmt.Println("Failed to load template file. No file name was provided!")
	} else {
		file, err := os.Create(filepath.Join(currentDir, path.CleanPath(fileName)))

		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(template)
		err = file.Close()

		if err != nil {
			// TODO This function should return an error instead of panicking
			log.Fatal("Failed to close file: ", err)
		}
	}
}
