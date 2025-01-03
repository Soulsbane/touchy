package templates

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/alecthomas/chroma/quick"
	"golang.org/x/exp/slices"
)

var ErrTemplateNotFound = errors.New("template not found")
var ErrLanguageNotFound = errors.New("language not found")
var ErrFileNameEmpty = errors.New("output filename not specified")
var ErrFailedToCreateFile = errors.New("failed to create file")
var ErrNoUserTemplatesDir = errors.New("no user templates found")
var ErrFailedToReadFile = errors.New("failed to read file")
var ErrFailedToReadEmbeddedFile = errors.New("failed to read embedded file")
var ErrHighlightFailed = errors.New("failed to highlight code")

type Templates interface {
	// CreateFileFromTemplate(languageName string, templateName string, customFileName string) error
	GetListOfAllLanguages() []string
	GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile)
	GetListOfLanguageTemplatesFor(languageName string) []infofile.InfoFile
	HasTemplate(languageName string, templateName string) (bool, int)
	HasLanguage(languageName string) (bool, int)
}

// NOTE: Will no longer be used. Using Languages as main data structure
type Language struct {
	// dirName         string                  // name of the directory under the template's directory.
	templateConfigs []infofile.InfoFile // A list of all the templates in the language directory. The key is the template dir name.
}

type Languages struct {
	languageName string
	infoFile     infofile.InfoFile
}

type TemplateManager struct {
	languages    map[string]Language // Map of all languages in the templates directory. Key is the language name.
	templateList []Templates
}

//func (lang *Language) GetInfoFile() infofile.InfoFile {
//	return lang.infoConfig
//}

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

func New() (*TemplateManager, error, error) {
	var manager TemplateManager

	manager.languages = make(map[string]Language)

	return &manager, nil, nil
}

func (g *TemplateManager) GatherTemplates() {
	embedded := NewEmbeddedTemplates()
	//user := NewUserTemplates()
	g.templateList = append(g.templateList, embedded)

	//maps.Copy(languages, user.GetListOfAllLanguages())
}

func (g *TemplateManager) HasLanguage(languageName string) (bool, []int) {
	indexes := make([]int, 0)

	for _, temp := range g.templateList {
		found, idx := temp.HasLanguage(languageName)

		if found {
			indexes = append(indexes, idx)
		}
	}

	if len(indexes) > 0 {
		return true, indexes
	}

	return false, indexes
}

func (g *TemplateManager) HasTemplate(languageName string, templateName string) (bool, []int) {
	indexes := make([]int, 0)

	for _, temp := range g.templateList {
		found, idx := temp.HasTemplate(languageName, templateName)

		if found {
			indexes = append(indexes, idx)
		}
	}

	if len(indexes) > 0 {
		return true, indexes
	}

	return false, indexes
}

func (g *TemplateManager) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
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

func (g *TemplateManager) loadTemplateFile(language string, template string, info infofile.InfoFile) (string, error) {
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

func (g *TemplateManager) GetListOfLanguageTemplates(language Language) []infofile.InfoFile {
	return language.templateConfigs
}

func (g *TemplateManager) GetListOfAllLanguages() map[string]Language {
	return g.languages
}

func (g *TemplateManager) ShowTemplate(languageName string, templateName string) error {
	for _, temp := range g.templateList {
		foundLanguage, _ := temp.HasLanguage(languageName)

		if foundLanguage {
			// Styles: https://github.com/alecthomas/chroma/tree/master/styles
			sourceCode, info := temp.GetLanguageTemplateFor(languageName, templateName)

			err := quick.Highlight(os.Stdout, sourceCode, info.GetDefaultOutputFileName(), "terminal256", "monokai")

			if err != nil {
				return ErrHighlightFailed
			}
		} else {
			return ErrLanguageNotFound
		}
	}

	return nil
}

// CreateFileFromTemplate Creates a template
func (g *TemplateManager) CreateFileFromTemplate(languageName string, templateName string, customFileName string) error {
	hasLang, _ := g.HasLanguage(languageName)
	hasTemp, _ := g.HasTemplate(languageName, templateName)

	if hasLang {
		if hasTemp {
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
