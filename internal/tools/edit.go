package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	EditToolName = "edit"
	EditToolDescription = `Performs exact string replacements in files.

Usage:
- You must use your Read tool at least once in the conversation before editing. This tool will error if you attempt an edit without reading the file.
- When editing text from Read tool output, ensure you preserve the exact indentation (tabs/spaces) as it appears AFTER the line number prefix. The line number prefix format is: spaces + line number + tab. Everything after that tab is the actual file content to match. Never include any part of the line number prefix in the old_string or new_string.
- ALWAYS prefer editing existing files in the codebase. NEVER write new files unless explicitly required.
- The edit will FAIL if 'old_string' is not unique in the file. Either provide a larger string with more surrounding context to make it unique or use 'replace_all' to change every instance of 'old_string'.
- Use 'replace_all' for replacing and renaming strings across the file. This parameter is useful if you want to rename a variable for instance.`
)

type EditInput struct {
	FilePath   string `json:"file_path" jsonschema:"required" jsonschema_description:"The absolute path to the file to modify"`
	OldString  string `json:"old_string" jsonschema:"required" jsonschema_description:"The text to replace"`
	NewString  string `json:"new_string" jsonschema:"required" jsonschema_description:"The text to replace it with (must be different from old_string)"`
	ReplaceAll bool   `json:"replace_all,omitempty" jsonschema_description:"Replace all occurences of old_string (default false)" jsonschema_default:"false"`
}

type EditOutput struct {
	Success       bool   `json:"success"`
	ReplacedCount int    `json:"replaced_count"`
	Message       string `json:"message"`
}

func NewEditTool() *ToolDefinition[EditInput, EditOutput] {
	return NewToolDefinition(
		EditToolName,
		EditToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input EditInput) (*mcp.CallToolResult, EditOutput, error) {

			content, err := os.ReadFile(input.FilePath)
			if err != nil {
				if os.IsNotExist(err) {
					return nil, EditOutput{}, fmt.Errorf("file does not exist: %s", input.FilePath)
				}
				return nil, EditOutput{}, fmt.Errorf("failed to read file: %w", err)
			}

			fileContent := string(content)

			if !strings.Contains(fileContent, input.OldString) {
				return nil, EditOutput{}, fmt.Errorf("old_string not found in file")
			}

			if input.OldString == input.NewString {
				return nil, EditOutput{}, fmt.Errorf("old_string and new_string must be different")
			}

			var newContent string
			var replacedCount int

			if input.ReplaceAll {

				replacedCount = strings.Count(fileContent, input.OldString)
				newContent = strings.ReplaceAll(fileContent, input.OldString, input.NewString)
			} else {

				count := strings.Count(fileContent, input.OldString)
				if count > 1 {
					return nil, EditOutput{}, fmt.Errorf("old_string appears %d times in the file. Either provide more context to make it unique or use replace_all=true", count)
				}

				newContent = strings.Replace(fileContent, input.OldString, input.NewString, 1)
				replacedCount = 1
			}

			err = os.WriteFile(input.FilePath, []byte(newContent), 0644)
			if err != nil {
				return nil, EditOutput{}, fmt.Errorf("failed to write file: %w", err)
			}

			output := EditOutput{
				Success:       true,
				ReplacedCount: replacedCount,
				Message:       fmt.Sprintf("Successfully replaced %d occurrence(s)", replacedCount),
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: output.Message},
				},
			}, output, nil
		},
	)
}
