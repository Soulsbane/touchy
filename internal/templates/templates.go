package templates

import (
	"embed"
	"errors"
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/alecthomas/chroma/quick"
	"golang.org/x/exp/slices"
	"io/fs"
	"os"
	"path"
)

//go:embed templates
var embedsDir embed.FS
var ErrTemplateNotFound = errors.New("template not found")
var ErrLanguageNotFound = errors.New("language not found")
var ErrFileNameEmpty = errors.New("output filename not specified")
var ErrFailedToCreateFile = errors.New("failed to create file")
var ErrNoUserTemplatesDir = errors.New("no user templates found")
var ErrFailedToReadFile = errors.New("failed to read file")
var ErrFailedToReadEmbeddedFile = errors.New("failed to read embedded file")
var ErrHighlightFailed = errors.New("failed to highlight code")

type Templates2 interface {
	CreateFileFromTemplate(languageName string, templateName string, customFileName string) error
	GetListOfAllLanguages() []string
	GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile)
	GetListOfLanguageTemplatesFor(languageName string) []infofile.InfoFile
	HasTemplate(languageName string, templateName string) bool
	HasLanguage(languageName string) bool
	ShowTemplate(languageName string, templateName string) error
}

type Language struct {
	// dirName         string                  // name of the directory under the template's directory.
	infoConfig      infofile.InfoFile   // Each language has a config file in its root directory call config.toml
	templateConfigs []infofile.InfoFile // A list of all the templates in the language directory. The key is the template dir name.
}

type Templates struct {
	languages map[string]Language // Map of all languages in the templates directory. Key is the language name.
}

func (lang *Language) GetInfoFile() infofile.InfoFile {
	return lang.infoConfig
}

func (lang *Language) GetTemplatesInfoFiles() []infofile.InfoFile {
	return lang.templateConfigs
}

func getFileData(path string, embedded bool) ([]byte, error) {
	if embedded {
		if data, err := embedsDir.ReadFile(path); err != nil {
			return data, fmt.Errorf("%w: %w", ErrFailedToReadEmbeddedFile, err)
		} else {
			return data, nil
		}
	} else {
		if data, err := os.ReadFile(path); err != nil {
			return data, fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
		} else {
			return data, nil
		}
	}
}

func New() (*Templates, error, error) {
	var templates Templates

	templates.languages = make(map[string]Language)
	userTemplatesErr := templates.findUserTemplates()
	embeddedTemplatesErr := templates.findEmbeddedTemplates()

	return &templates, userTemplatesErr, embeddedTemplatesErr
}

func (g *Templates) findUserTemplates() error {
	dirs, err := os.ReadDir(pathutils.GetTemplatesDir())

	if err != nil {
		return fmt.Errorf("%w: %w", ErrNoUserTemplatesDir, err)
	}

	return g.findTemplates(dirs, false)
}

func (g *Templates) findEmbeddedTemplates() error {
	dirs, err := embedsDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	return g.findTemplates(dirs, true)
}

func (g *Templates) findTemplates(dirs []fs.DirEntry, embedded bool) error {
	var templatePath string

	if embedded {
		templatePath = "templates"
	} else {
		templatePath = pathutils.GetTemplatesDir()
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

			if embedded {
				templates, err = embedsDir.ReadDir(path.Join(templatePath, languageDir.Name()))
			} else {
				templates, err = os.ReadDir(path.Join(templatePath, languageDir.Name()))
			}

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

func (g *Templates) HasLanguage(languageName string) bool {
	_, found := g.languages[languageName]
	return found
}

func (g *Templates) HasTemplate(languageName string, templateName string) bool {
	if language, foundLanguage := g.languages[languageName]; foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		return idx >= 0
	}

	return false
}

func (g *Templates) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	language, foundLanguage := g.languages[languageName]

	if foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		if idx >= 0 {
			info := language.templateConfigs[idx]
			data, err := g.loadTemplateFile(languageName, templateName, info)

			if err != nil {
				return "", language.templateConfigs[idx]
			} else {
				return data, language.templateConfigs[idx]
			}
		}
	}

	return "", infofile.InfoFile{}
}

func (g *Templates) loadTemplateFile(language string, template string, info infofile.InfoFile) (string, error) {
	var data []byte
	var templateName string
	var err error

	if info.IsEmbedded() {
		templateName = path.Join("templates", language, template, template+".template")
		data, err = embedsDir.ReadFile(templateName)

	} else {
		templateName = path.Join(pathutils.GetTemplatesDir(), language, template, template+".template")
		data, err = os.ReadFile(templateName)
	}

	if err != nil { // We couldn't read from the embedded file or the file in user's config directory so return an error
		return "", ErrTemplateNotFound
	}

	return string(data), nil
}

func (g *Templates) GetListOfLanguageTemplates(language Language) []infofile.InfoFile {
	return language.templateConfigs
}

func (g *Templates) GetListOfAllLanguages() map[string]Language {
	return g.languages
}

func (g *Templates) ShowTemplate(languageName string, templateName string) error {
	if language, languageFound := g.languages[languageName]; languageFound {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		if idx >= 0 {
			sourceCode, err := g.loadTemplateFile(languageName, templateName, language.templateConfigs[idx])

			if err != nil {
				return ErrTemplateNotFound
			}

			// Formatters: terminal, terminal8, terminal16, terminal256, terminal16m
			// Styles: https://github.com/alecthomas/chroma/tree/master/styles
			err = quick.Highlight(os.Stdout, sourceCode, language.templateConfigs[idx].GetDefaultOutputFileName(), "terminal256", "monokai")

			if err != nil {
				return ErrHighlightFailed
			}

		} else {
			return ErrTemplateNotFound
		}
	} else {
		return ErrLanguageNotFound
	}

	return nil
}

// CreateFileFromTemplate Creates a template
func (g *Templates) CreateFileFromTemplate(languageName string, templateName string, customFileName string) error {
	if g.HasLanguage(languageName) {
		if g.HasTemplate(languageName, templateName) {
			var fileName string

			template, config := g.GetLanguageTemplateFor(languageName, templateName)
			currentDir, _ := os.Getwd()

			if customFileName == "DefaultOutputFileName" {
				fileName = config.GetDefaultOutputFileName()
			} else {
				fileName = customFileName
			}

			if fileName == "" {
				return ErrFileNameEmpty
			} else {
				fullFileName := path.Join(currentDir, pathutils.CleanPath(fileName))

				if err := os.WriteFile(fullFileName, []byte(template), 0600); err != nil {
					return fmt.Errorf("%w %s", ErrFailedToCreateFile, fullFileName)
				}
			}

			return nil
		} else {
			return ErrTemplateNotFound
		}
	} else {
		return ErrLanguageNotFound
	}

}
