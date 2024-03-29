package templates

import (
	"embed"
	"fmt"
	"github.com/Soulsbane/touchy/internal/infofile"
	"golang.org/x/exp/slices"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Soulsbane/touchy/internal/path"
	"github.com/alecthomas/chroma/quick"
)

//go:embed templates
var embedsDir embed.FS

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
		data, err := embedsDir.ReadFile(path)

		if err != nil {
			return data, err
		} else {
			return data, nil
		}
	} else {
		data, err := os.ReadFile(path)

		if err != nil {
			return data, err
		} else {
			return data, nil
		}
	}
}

func New() *Templates {
	var templates Templates

	templates.languages = make(map[string]Language)
	templates.findUserTemplates()
	templates.findEmbeddedTemplates()

	return &templates
}

func (g *Templates) findUserTemplates() {
	dirs, err := os.ReadDir(path.GetTemplatesDir())

	if err != nil {
		//panic(err)
	}

	g.findTemplates(dirs, false)
}

func (g *Templates) findEmbeddedTemplates() {
	dirs, err := embedsDir.ReadDir("templates")

	if err != nil {
		panic(err)
	}

	g.findTemplates(dirs, true)
}

func (g *Templates) findTemplates(dirs []fs.DirEntry, embedded bool) {
	var templatePath string

	if embedded {
		templatePath = "templates"
	} else {
		templatePath = path.GetTemplatesDir()
	}

	for _, languageDir := range dirs {
		if languageDir.IsDir() {
			var language Language
			var templates []os.DirEntry

			infoPath := filepath.Join(templatePath, languageDir.Name(), infofile.DefaultFileName)
			data, err := getFileData(infoPath, embedded)
			language.infoConfig = infofile.Load(languageDir.Name(), infoPath, embedded, data)

			if err != nil {
				fmt.Println(err)
			}

			if embedded {
				templates, err = embedsDir.ReadDir(filepath.Join(templatePath, languageDir.Name()))
			} else {
				templates, err = os.ReadDir(filepath.Join(templatePath, languageDir.Name()))
			}

			if err != nil {
				fmt.Println("Could not read directory: ", err) // TODO: Handle this better?
			}

			for _, template := range templates {
				if template.IsDir() {
					configPath := filepath.Join(templatePath, languageDir.Name(), template.Name(), infofile.DefaultFileName)
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
}

func (g *Templates) HasLanguage(languageName string) bool {
	_, found := g.languages[languageName]
	return found
}

func (g *Templates) HasTemplate(languageName string, templateName string) bool {
	language, foundLanguage := g.languages[languageName]

	if foundLanguage {
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

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
		idx := slices.IndexFunc(language.templateConfigs, func(c infofile.InfoFile) bool { return c.GetName() == templateName })

		if idx >= 0 {
			info := language.templateConfigs[idx]
			//return g.loadTemplateFile(languageName, templateName, info), language.templateConfigs[idx]
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
		templateName = filepath.Join("templates", language, template, template+".template")
		data, err = embedsDir.ReadFile(templateName)

	} else {
		templateName = filepath.Join(path.GetTemplatesDir(), language, template, template+".template")
		data, err = os.ReadFile(templateName)
	}

	if err != nil { // We couldn't read from the embedded file or the file in user's config directory so return an error
		return "", fmt.Errorf("That template does not exist: %s", templateName)
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
				return fmt.Errorf("That template does not exist: %s", templateName)
			}

			// Formatters: terminal, terminal8, terminal16, terminal256, terminal16m
			// Styles: https://github.com/alecthomas/chroma/tree/master/styles
			err = quick.Highlight(os.Stdout, sourceCode, language.templateConfigs[idx].GetDefaultOutputFileName(), "terminal256", "monokai")

			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("That template does not exist: %s", templateName)
		}
	} else {
		return fmt.Errorf("That language does not exist: %s", languageName)
	}

	return nil
}

// CreateFileFromTemplate Creates a template
func (g *Templates) CreateFileFromTemplate(languageName string, templateName string, customFileName string) error {
	var fileName string
	template, config := g.GetLanguageTemplateFor(languageName, templateName)
	currentDir, _ := os.Getwd()

	if customFileName == "DefaultOutputFileName" {
		fileName = config.GetDefaultOutputFileName()
	} else {
		fileName = customFileName
	}

	if fileName == "" {
		return fmt.Errorf("Failed to load template file. No file name was provided!")
	} else {
		fullFileName := filepath.Join(currentDir, path.CleanPath(fileName))
		file, err := os.Create(fullFileName)

		if err != nil {
			return fmt.Errorf("Failed to create file: %s", fullFileName)
		}

		_, err = file.WriteString(template)

		if err != nil {
			return fmt.Errorf("Failed to write to file: %s", fullFileName)
		}

		err = file.Close()

		if err != nil {
			return fmt.Errorf("Failed to close file: %s", fullFileName)
		}
	}

	return nil
}
