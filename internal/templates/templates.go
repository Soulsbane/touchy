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
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func New() *TemplateManager {
	var manager TemplateManager
	return &manager
}

func (g *TemplateManager) GatherTemplates() {
	embedded := NewEmbeddedTemplates()
	user := NewUserTemplates()
	g.templateList = append(g.templateList, embedded)
	g.templateList = append(g.templateList, user)
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
	for _, temp := range g.templateList {
		hasTemp, _ := temp.HasTemplate(languageName, templateName)

		if hasTemp {
			return temp.GetLanguageTemplateFor(languageName, templateName)
		}
	}

	return "", infofile.InfoFile{}
}

func (g *TemplateManager) outputTemplateList(languageName string, languages []Languages) {
	headerName := cases.Title(language.English).String(languageName) + " Templates"

	if len(languages) > 0 {
		outputTable := ui.CreateNewTableWriter(headerName, "Language", "Name", "Description", "Output File name")

		for _, info := range languages {
			outputTable.AppendRow(table.Row{
				info.languageName, info.infoFile.GetName(),
				info.infoFile.GetDescription(),
				info.infoFile.GetDefaultOutputFileName(),
			})
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

	g.outputTemplateList(languageName, languages)
}

func (g *TemplateManager) ListLanguages() {
	// TODO: Download language info from the programming language DB and store locally
	var languages []string
	outputTable := ui.CreateNewTableWriter("Languages", "Language Name", "Description", "URL")

	for _, temp := range g.templateList {
		languages = append(languages, temp.GetListOfAllLanguages()...)
	}

	for _, lang := range lo.Uniq(languages) {
		outputTable.AppendRow(table.Row{lang, "<no description>", "<no url>"})
	}

	outputTable.Render()
}

func (g *TemplateManager) ShowTemplate(languageName string, templateName string) error {
	foundLanguage := g.hasLanguage(languageName)

	if foundLanguage {
		// Styles: https://github.com/alecthomas/chroma/tree/master/styles
		sourceCode, info := g.GetLanguageTemplateFor(languageName, templateName)

		err := quick.Highlight(os.Stdout, sourceCode, info.GetDefaultOutputFileName(), "terminal256", "monokai")

		if err != nil {
			return ErrHighlightFailed
		}
	} else {
		return ErrLanguageNotFound
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
