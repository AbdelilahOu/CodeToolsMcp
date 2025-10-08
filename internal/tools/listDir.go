package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ListDirToolName        = "list_dir"
	ListDirToolDescription = `Lists directory contents similar to "ls".

Usage:
- Provide an absolute directory path
- Toggle recursive to walk subdirectories, limit to cap the number of entries, and show_hidden to include dotfiles
- Results include file type, size, permissions, and modification time`
)

type ListDirInput struct {
	Path       string `json:"path" jsonschema:"required" jsonschema_description:"Absolute path to the directory to inspect."`
	Recursive  bool   `json:"recursive,omitempty" jsonschema_description:"Walk subdirectories recursively."`
	ShowHidden bool   `json:"show_hidden,omitempty" jsonschema_description:"Include entries whose names start with a dot."`
	Limit      int    `json:"limit,omitempty" jsonschema_description:"Maximum number of entries to include (0 for unlimited)."`
}

type ListDirEntry struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size_bytes"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
}

type ListDirOutput struct {
	Entries []ListDirEntry `json:"entries"`
}

func NewListDirTool(runner *runners.FileRunner) *ToolDefinition[ListDirInput, ListDirOutput] {
	return NewToolDefinition(
		ListDirToolName,
		ListDirToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input ListDirInput) (*mcp.CallToolResult, ListDirOutput, error) {
			entries, err := runner.ListDir(ctx, runners.ListDirInput{
				Path:       input.Path,
				Recursive:  input.Recursive,
				ShowHidden: input.ShowHidden,
				Limit:      input.Limit,
			})
			if err != nil {
				return nil, ListDirOutput{}, err
			}

			sort.Slice(entries, func(i, j int) bool {
				if entries[i].IsDir == entries[j].IsDir {
					return entries[i].Path < entries[j].Path
				}
				return entries[i].IsDir && !entries[j].IsDir
			})

			summary := formatDirEntries(entries)

			outputEntries := make([]ListDirEntry, len(entries))
			for i, entry := range entries {
				outputEntries[i] = ListDirEntry{
					Path:    entry.Path,
					Name:    entry.Name,
					IsDir:   entry.IsDir,
					Size:    entry.Size,
					Mode:    entry.Mode.String(),
					ModTime: entry.ModTime.Format(time.RFC3339),
				}
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: summary},
				},
			}, ListDirOutput{Entries: outputEntries}, nil
		},
	)
}

func formatDirEntries(entries []runners.DirEntry) string {
	if len(entries) == 0 {
		return "(empty)"
	}

	var builder strings.Builder
	for _, entry := range entries {
		typeLabel := "file"
		if entry.IsDir {
			typeLabel = "dir"
		}
		line := fmt.Sprintf("%-4s %12dB %s %s", typeLabel, entry.Size, entry.ModTime.Format(time.RFC3339), entry.Path)
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}
