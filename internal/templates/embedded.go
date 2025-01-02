package templates

import (
	"embed"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/samber/lo"
)

//go:embed templates
var embedsDir embed.FS

type EmbeddedTemplates struct {
	languages []Languages
}

func NewEmbeddedTemplates() *EmbeddedTemplates {
	var templates EmbeddedTemplates
	templates.languages = make([]Languages, 0)

	err := templates.findTemplates(true)

	if err != nil {
		panic(err)
	}

	return &templates
}

func getEmbeddedData(path string) ([]byte, error) {
	if data, err := embedsDir.ReadFile(path); err != nil {
		return data, fmt.Errorf("%w: %w", ErrFailedToReadEmbeddedFile, err)
	} else {
		return data, nil
	}
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

			//infoPath := path.Join(templatePath, languageDir.Name(), infofile.DefaultFileName)
			//data, err := getEmbeddedData(infoPath)
			//language.infoConfig = infofile.Load(languageDir.Name(), infoPath, embedded, data)

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
					templateData, fileReadErr := getEmbeddedData(configPath)

					if fileReadErr != nil {
						fmt.Println(fileReadErr)
					}

					config := infofile.Load(template.Name(), configPath, embedded, templateData)
					config.SetEmbedded(embedded)
					language.templateConfigs = append(language.templateConfigs, config)
					g.languages = append(g.languages, Languages{languageDir.Name(), config})
				}
			}
		}
	}

	return nil // TODO: Handle errors
}

func (g *EmbeddedTemplates) LoadTemplateFile(language string, template string) (string, error) {
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

func (g *EmbeddedTemplates) GetListOfLanguageTemplatesFor(language string) []infofile.InfoFile {
	return []infofile.InfoFile{} // TODO: Placeholder while I rework
}

func (g *EmbeddedTemplates) GetListOfAllLanguages() []string {
	infos := make([]string, 0)

	for _, language := range g.languages {
		infos = append(infos, language.languageName)
	}

	return lo.Uniq(infos)
}

func (g *EmbeddedTemplates) HasLanguage(languageName string) (bool, int) {
	idx := slices.IndexFunc(g.languages, func(c Languages) bool { return c.languageName == languageName })
	return idx >= 0, idx
}

func (g *EmbeddedTemplates) HasTemplate(languageName string, templateName string) (bool, int) {
	idx := slices.IndexFunc(g.languages, func(c Languages) bool { return c.languageName == languageName && c.infoFile.GetName() == templateName })
	return idx >= 0, idx
}

func (g *EmbeddedTemplates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	hasTemplate, idx := g.HasTemplate(languageName, templateName)

	if hasTemplate {
		if idx >= 0 {
			data, err := g.LoadTemplateFile(languageName, templateName)

			if err != nil {
				return "", g.languages[idx].infoFile
			} else {
				return data, g.languages[idx].infoFile
			}
		}

	}
	return "", infofile.InfoFile{}
}
