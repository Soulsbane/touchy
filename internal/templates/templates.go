package templates

import (
	"fmt"
	"os"
	"path"

	"github.com/Soulsbane/touchy/internal/common"

	"github.com/Soulsbane/touchy/internal/infofile"
	"github.com/Soulsbane/touchy/internal/pathutils"
	"github.com/Soulsbane/touchy/internal/ui"
	"github.com/alecthomas/chroma/quick"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Templates interface {
	GetListOfAllLanguages() []string
	GetLanguages() []Languages
	GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile)
	GetListOfLanguageTemplatesFor(languageName string) []Languages
	HasTemplate(languageName string, templateName string) bool
	HasLanguage(languageName string) bool
	GetTemplateIndexFor(languageName string, templateName string) (bool, int)
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
	embedded, embedErr := NewEmbeddedTemplates()
	user, userErr := NewUserTemplates()

	g.templateList = append(g.templateList, embedded)
	g.templateList = append(g.templateList, user)

	g.handleTemplateError(embedErr, userErr)
}

func (g *TemplateManager) handleTemplateError(embedErr error, userErr error) {
	if embedErr != nil {
		fmt.Println(embedErr)
	}

	if userErr != nil {
		fmt.Println(userErr)
	}
}

func (g *TemplateManager) HasLanguage(languageName string) bool {
	for _, temp := range g.templateList {
		if temp.HasLanguage(languageName) {
			return true
		}
	}

	return false
}

func (g *TemplateManager) HasTemplate(languageName string, templateName string) bool {
	for _, temp := range g.templateList {
		if temp.HasTemplate(languageName, templateName) {
			return true
		}
	}

	return false
}

func (g *TemplateManager) GetLanguageTemplateFor(languageName string, templateName string) (string, infofile.InfoFile) {
	for _, temp := range g.templateList {
		hasTemp := temp.HasTemplate(languageName, templateName)

		if hasTemp {
			return temp.GetLanguageTemplateFor(languageName, templateName)
		}
	}

	return "", infofile.InfoFile{}
}

func (g *TemplateManager) outputTemplateList(languageName string, languages []Languages) {
	headerName := cases.Title(language.English).String(languageName) + " Templates"

	if len(languages) > 0 {
		var outputTable table.Writer

		if languageName == "all" {
			outputTable = ui.CreateNewTableWriter(headerName, "Name", "Language", "Description", "Output File name")
		} else {
			outputTable = ui.CreateNewTableWriter(headerName, "Name", "Description", "Output File name")
		}

		for _, info := range languages {
			var outputRow table.Row

			if languageName == "all" {
				outputRow = table.Row{
					info.infoFile.GetName(),
					info.languageName, info.infoFile.GetDescription(),
					info.infoFile.GetDefaultOutputFileName(),
				}
			} else {
				outputRow = table.Row{info.infoFile.GetName(),
					info.infoFile.GetDescription(),
					info.infoFile.GetDefaultOutputFileName(),
				}
			}

			outputTable.AppendRow(outputRow)
		}

		outputTable.Render()
	}
}

func (g *TemplateManager) ListTemplates(languageName string) {
	var languages []Languages

	for _, temp := range g.templateList {
		hasLang := temp.HasLanguage(languageName)

		if hasLang {
			languages = append(languages, temp.GetListOfLanguageTemplatesFor(languageName)...)
		}

		if languageName == "all" {
			languages = append(languages, temp.GetLanguages()...)
		}
	}

	if len(languages) > 0 {
		g.outputTemplateList(languageName, languages)
	} else {
		fmt.Println("No templates found. Use 'list templates' to see available templates.")
	}
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
	foundLanguage := g.HasLanguage(languageName)

	if foundLanguage {
		// Styles: https://github.com/alecthomas/chroma/tree/master/styles
		sourceCode, info := g.GetLanguageTemplateFor(languageName, templateName)

		err := quick.Highlight(os.Stdout, sourceCode, info.GetDefaultOutputFileName(), "terminal256", "monokai")

		if err != nil {
			return common.ErrHighlightFailed
		}
	} else {
		return common.ErrLanguageNotFound
	}

	return nil
}

// CreateFileFromTemplate Creates a template
func (g *TemplateManager) CreateFileFromTemplate(languageName string, templateName string, customFileName string) error {
	hasLang := g.HasLanguage(languageName)
	hasTemp := g.HasTemplate(languageName, templateName)

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
				return common.ErrFileNameEmpty
			} else {
				fullFileName := path.Join(currentDir, pathutils.CleanPath(fileName))

				if err := os.WriteFile(fullFileName, []byte(template), 0600); err != nil {
					return fmt.Errorf("%w %s", common.ErrFailedToCreateFile, fullFileName)
				}
			}

			return nil
		} else {
			return common.ErrTemplateNotFound
		}
	}

	return common.ErrLanguageNotFound

}
