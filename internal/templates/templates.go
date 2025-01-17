package templates

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/ui"
	"github.com/alecthomas/chroma/quick"
	"github.com/jedib0t/go-pretty/v6/table"
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
	GetLanguages() []Languages
	GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile)
	GetListOfLanguageTemplatesFor(languageName string) []Languages
	HasTemplate(languageName string, templateName string) (bool, int)
	HasLanguage(languageName string) bool
}

type Languages struct {
	languageName string
	infoFile     infofile.InfoFile
}

type TemplateManager struct {
	templateList []Templates
}

func (lang *Languages) GetInfoFile() infofile.InfoFile {
	return lang.infoFile
}

func New() (*TemplateManager, error, error) {
	var manager TemplateManager
	return &manager, nil, nil
}

func (g *TemplateManager) GatherTemplates() {
	embedded := NewEmbeddedTemplates()
	//user := NewUserTemplates()
	g.templateList = append(g.templateList, embedded)
}

func (g *TemplateManager) hasLanguage(languageName string) bool {
	for _, temp := range g.templateList {
		if temp.HasLanguage(languageName) {
			return true
		}
	}

	return false
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

func (g *TemplateManager) outputTable(languageName string, languages []Languages) {
	if len(languages) > 0 {
		outputTable := ui.CreateNewTableWriter(languageName+" Templates", "Language", "Name", "Description", "Default Output File name")

		for _, info := range languages {
			outputTable.AppendRow(table.Row{info.languageName, info.infoFile.GetName(), info.infoFile.GetDescription(), info.infoFile.GetDefaultOutputFileName()})
		}

		outputTable.Render()
	}
}

func (g *TemplateManager) ListTemplates(languageName string) {
	languages := make([]Languages, 0)

	for _, temp := range g.templateList {
		hasLang := temp.HasLanguage(languageName)

		if hasLang {
			if languages != nil {
				languages = append(languages, temp.GetListOfLanguageTemplatesFor(languageName)...)
			}
		} else if languageName == "all" {
			if languages != nil {
				languages = append(languages, temp.GetLanguages()...)
			}
		} else {
			fmt.Println("That language could not be found! Use 'list all' to see all available languages.")
		}
	}

	g.outputTable(languageName, languages)
}

func (g *TemplateManager) ListLanguages() {
	// TODO: Download language info from the programming language DB and store locally
	outputTable := ui.CreateNewTableWriter("Languages", "Language Name", "Description", "URL")

	for _, temp := range g.templateList {
		languages := temp.GetListOfAllLanguages()

		for _, language := range languages {
			outputTable.AppendRow(table.Row{language, "<no description>", "<no url>"})
		}

	}
	outputTable.Render()
}

func (g *TemplateManager) ShowTemplate(languageName string, templateName string) error {
	for _, temp := range g.templateList {
		foundLanguage := temp.HasLanguage(languageName)

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
	hasLang := g.hasLanguage(languageName)
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
