package templates

import (
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"os"
	"path"
	"slices"
)

type EmbeddedTemplates struct {
	languages map[string]Language // Map of all languages in the templates directory. Key is the language name.
}

func NewEmbeddedTemplates() *EmbeddedTemplates {
	var templates EmbeddedTemplates

	templates.languages = make(map[string]Language)
	err := templates.findTemplates(true)

	if err != nil {
		panic(err)
	}

	return &templates
}

func (g *EmbeddedTemplates) findTemplates(embedded bool) error {
	templatePath := "templates"
	dirs, err := embedsDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	for _, languageDir := range dirs {
		if languageDir.IsDir() {
			var language Language
			var templates []os.DirEntry

			infoPath := path.Join(templatePath, languageDir.Name(), infofile.DefaultFileName)
			data, err := getFileData(infoPath, embedded)
			language.infoConfig = infofile.Load(languageDir.Name(), infoPath, embedded, data)

			if err != nil {
				fmt.Println(err)
			}

			templates, err = embedsDir.ReadDir(path.Join(templatePath, languageDir.Name()))

			if err != nil {
				fmt.Println("Could not read directory: ", err) // TODO: Handle this better?
			}

			for _, template := range templates {
				if template.IsDir() {
					configPath := path.Join(templatePath, languageDir.Name(), template.Name(), infofile.DefaultFileName)
					templateData, fileReadErr := getFileData(configPath, embedded)

					if fileReadErr != nil {
						fmt.Println(fileReadErr)
					}

					config := infofile.Load(template.Name(), configPath, embedded, templateData)
					config.SetEmbedded(embedded)
					language.templateConfigs = append(language.templateConfigs, config)
				}
			}

			g.languages[languageDir.Name()] = language
		}
	}

	return nil // TODO: Handle errors
}

func (g *EmbeddedTemplates) loadTemplateFile(language string, template string) (string, error) {
	var data []byte
	var templateName string
	var err error

	templateName = path.Join("templates", language, template, template+".template")
	data, err = embedsDir.ReadFile(templateName)

	if err != nil { // We couldn't read from the embedded file or the file in user's config directory so return an error
		return "", ErrTemplateNotFound
	}

	return string(data), nil
}

func (g *EmbeddedTemplates) GetListOfLanguageTemplatesFor(language Language) []infofile.InfoFile {
	return language.templateConfigs
}

func (g *EmbeddedTemplates) GetListOfAllLanguages() map[string]Language {
	return g.languages
}

func (g *EmbeddedTemplates) HasLanguage(languageName string) bool {
	_, found := g.languages[languageName]
	return found
}

func (g *EmbeddedTemplates) HasTemplate(languageName string, templateName string) bool {
	if language, foundLanguage := g.languages[languageName]; foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		return idx >= 0
	}

	return false
}

func (g *EmbeddedTemplates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	language, foundLanguage := g.languages[languageName]

	if foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		if idx >= 0 {
			data, err := g.LoadTemplateFile(languageName, templateName)

			if err != nil {
				return "", language.templateConfigs[idx]
			} else {
				return data, language.templateConfigs[idx]
			}
		}
	}

	return "", infofile.InfoFile{}
}
