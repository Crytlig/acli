package pkg

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

type ModelType string

const (
	// Only gpt-3.5-turbo-0613 is being supported
	// and that is a chat type model
	ModelTypeChat ModelType = "chat"

	// Other models, such as text-davinci could probably be added for testing
	// ModelTypeCompletion ModelType = "completion"
)

type Model struct {
	Name      string
	MaxTokens int
	Type      ModelType
}

// Add available, however, this is only for pretty print.
// For now, only GPT3Dot5Turbo0613 is supported, hence the hardcoding in query.go
// and lack of a CLI flag
var (
	GPT432K0613          = Model{"gpt-4-32k-0613", 32384, ModelTypeChat}
	GPT40613             = Model{"gpt-4-0613", 16192, ModelTypeChat}
	GPT3Dot5Turbo0613    = Model{"gpt-3.5-turbo-0613", 4096, ModelTypeChat}
	GPT3Dot5Turbo16K0613 = Model{"gpt-3.5-turbo-16k-0613", 16192, ModelTypeChat}

	supportedModels = []Model{
		GPT432K0613,
		GPT40613,
		GPT3Dot5Turbo0613,
		GPT3Dot5Turbo16K0613,
	}
)

func AvailableModels(c *cli.Context) error {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Type", "Max Tokens")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, model := range supportedModels {
		tbl.AddRow(model.Name, model.Type, model.MaxTokens)
	}

	tbl.Print()
	return nil
}
