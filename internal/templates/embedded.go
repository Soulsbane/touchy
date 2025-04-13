package templates

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/touchy/internal/common"
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

func NewEmbeddedTemplates() (*EmbeddedTemplates, error) {
	var templates EmbeddedTemplates
	err := templates.findTemplates(true)

	if err != nil {
		return &templates, err
	}

	return &templates, nil
}

func getEmbeddedData(path string) ([]byte, error) {
	if data, err := embedsDir.ReadFile(path); err != nil {
		return data, fmt.Errorf("%w: %w", common.ErrFailedToReadEmbeddedFile, err)
	} else {
		return data, nil
	}
}

func (g *EmbeddedTemplates) findTemplates(embedded bool) error {
	templatePath := "templates"
	dirs, err := embedsDir.ReadDir("templates")

	if err != nil {
		return fmt.Errorf("failed to read script embeds directory: %w", err)
	}

	for _, languageDir := range dirs {
		if languageDir.IsDir() {
			var templates []os.DirEntry

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
					g.languages = append(g.languages, Languages{languageDir.Name(), config})
				}
			}
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

	if err != nil { // We couldn't read from the embedded file's directory so return an error
		return "", fmt.Errorf("%w: %w", common.ErrTemplateNotFound, err)
	}

	return string(data), nil
}

func (g *EmbeddedTemplates) GetListOfLanguageTemplatesFor(language string) []Languages {
	var values []Languages

	for _, temp := range g.languages {
		if temp.languageName == language {
			values = append(values, temp)
		}
	}

	return values
}

func (g *EmbeddedTemplates) GetListOfAllLanguages() []string {
	var infos []string

	for _, language := range g.languages {
		infos = append(infos, language.languageName)
	}

	return lo.Uniq(infos)
}

func (g *EmbeddedTemplates) GetLanguages() []Languages {
	return g.languages
}

func (g *EmbeddedTemplates) HasLanguage(languageName string) bool {
	return slices.ContainsFunc(g.languages, func(c Languages) bool { return c.languageName == languageName })
}

func (g *EmbeddedTemplates) HasTemplate(languageName string, templateName string) bool {
	return slices.ContainsFunc(g.languages, func(c Languages) bool { return c.languageName == languageName && c.infoFile.GetName() == templateName })
}

func (g *EmbeddedTemplates) GetTemplateIndexFor(languageName string, templateName string) (bool, int) {
	idx := slices.IndexFunc(g.languages, func(c Languages) bool { return c.languageName == languageName && c.infoFile.GetName() == templateName })
	return idx >= 0, idx
}

func (g *EmbeddedTemplates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	if hasTemplate, idx := g.GetTemplateIndexFor(languageName, templateName); hasTemplate {
		if idx >= 0 {
			data, err := g.loadTemplateFile(languageName, templateName)

			if err != nil {
				return "", g.languages[idx].infoFile
			}

			return data, g.languages[idx].infoFile
		}
	}

	return "", infofile.InfoFile{}
}
