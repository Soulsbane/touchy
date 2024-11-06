package templates

import (
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"os"
	"path"
)

type UserTemplates struct {
	languages map[string]Language // Map of all languages in the templates directory. Key is the language name.
}

func NewUserTemplates() *UserTemplates {
	var templates UserTemplates

	templates.languages = make(map[string]Language)
	err := templates.findTemplates(false)

	if err != nil {
		fmt.Println(err)
	}

	return &templates
}

func (g *UserTemplates) findTemplates(embedded bool) error {
	templatePath := pathutils.GetTemplatesDir()
	dirs, err := os.ReadDir(templatePath)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrNoUserTemplatesDir, err)
	}

	for _, languageDir := range dirs {
		if languageDir.IsDir() {
			var language Language
			var templates []os.DirEntry

			infoPath := path.Join(templatePath, languageDir.Name(), infofile.DefaultFileName)
			data, err := getFileData(infoPath, false)
			language.infoConfig = infofile.Load(languageDir.Name(), infoPath, embedded, data)

			if err != nil {
				fmt.Println(err)
			}

			templates, err = os.ReadDir(path.Join(templatePath, languageDir.Name()))

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

func (g *UserTemplates) GetListOfAllLanguages() map[string]Language {
	return g.languages
}

func (g *UserTemplates) HasLanguage(languageName string) bool {
	_, found := g.languages[languageName]
	return found
}
