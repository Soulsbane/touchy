package templates

import (
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/samber/lo"
)

type UserTemplates struct {
	languages []Languages
}

func NewUserTemplates() *UserTemplates {
	var templates UserTemplates
	templates.languages = make([]Languages, 0)

	err := templates.findTemplates(true)

	if err != nil {
		panic(err)
	}

	return &templates
}

func getUserData(path string) ([]byte, error) {
	if data, err := os.ReadFile(path); err != nil {
		return data, fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
	} else {
		return data, nil
	}
}

func (g *UserTemplates) findTemplates(embedded bool) error {
	templatePath := pathutils.GetTemplatesDir()
	dirs, err := os.ReadDir(templatePath)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrNoUserTemplatesDir, err)
	}

	for _, languageDir := range dirs {
		if languageDir.IsDir() {
			var templates []os.DirEntry

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
					templateData, fileReadErr := getUserData(configPath)

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

func (g *UserTemplates) LoadTemplateFile(language string, template string) (string, error) {
	var data []byte
	var templateName string
	var err error

	templateName = path.Join("templates", language, template, template+".template")
	data, err = os.ReadFile(templateName)

	if err != nil { // We couldn't read from the embedded file or the file in user's config directory so return an error
		return "", ErrTemplateNotFound
	}

	return string(data), nil
}

func (g *UserTemplates) GetListOfLanguageTemplatesFor(language string) []Languages {
	values := make([]Languages, 0)

	for _, temp := range g.languages {
		if temp.languageName == language {
			values = append(values, temp)
		}
	}

	return values
}

func (g *UserTemplates) GetListOfAllLanguages() []string {
	infos := make([]string, 0)

	for _, language := range g.languages {
		infos = append(infos, language.languageName)
	}

	return lo.Uniq(infos)
}

func (g *UserTemplates) HasLanguage(languageName string) (bool, int) {
	idx := slices.IndexFunc(g.languages, func(c Languages) bool { return c.languageName == languageName })
	return idx >= 0, idx
}

func (g *UserTemplates) HasTemplate(languageName string, templateName string) (bool, int) {
	idx := slices.IndexFunc(g.languages, func(c Languages) bool { return c.languageName == languageName && c.infoFile.GetName() == templateName })
	return idx >= 0, idx
}

func (g *UserTemplates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
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
