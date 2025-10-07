package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GlobToolName = "glob"
	GlobToolDescription = `- Fast file pattern matching tool that works with any codebase size
- Supports glob patterns like "**/*.js" or "src/**/*.ts"
- Returns matching file paths sorted by modification time
- Use this tool when you need to find files by name patterns`
)

type GlobInput struct {
	Pattern string `json:"pattern" jsonschema:"required" jsonschema_description:"The glob pattern to match files against"`
	Path    string `json:"path,omitempty" jsonschema_description:"The directory to search in. If not specified, the current working directory will be used. IMPORTANT: Omit this field to use the default directory. DO NOT enter \"undefined\" or \"null\" - simply omit it for the default behavior. Must be a valid directory path if provided."`
}

type GlobOutput struct {
	Files []string `json:"files"`
}

type fileInfo struct {
	path    string
	modTime time.Time
}

func NewGlobTool() *ToolDefinition[GlobInput, GlobOutput] {
	return NewToolDefinition(
		GlobToolName,
		GlobToolDescription,
		func(ctx context.Context, req *mcp.CallToolRequest, input GlobInput) (*mcp.CallToolResult, GlobOutput, error) {
			searchPath := input.Path
			if searchPath == "" {
				cwd, err := os.Getwd()
				if err != nil {
					return nil, GlobOutput{}, fmt.Errorf("failed to get current directory: %w", err)
				}
				searchPath = cwd
			}

			absPath, err := filepath.Abs(searchPath)
			if err != nil {
				return nil, GlobOutput{}, fmt.Errorf("failed to get absolute path: %w", err)
			}

			fullPattern := filepath.Join(absPath, input.Pattern)

			matches, err := filepath.Glob(fullPattern)
			if err != nil {
				return nil, GlobOutput{}, fmt.Errorf("glob pattern error: %w", err)
			}

			if containsDoubleStar(input.Pattern) {
				matches, err = walkGlob(absPath, input.Pattern)
				if err != nil {
					return nil, GlobOutput{}, err
				}
			}

			fileInfos := make([]fileInfo, 0, len(matches))
			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					continue
				}
				if !info.IsDir() {
					fileInfos = append(fileInfos, fileInfo{
						path:    match,
						modTime: info.ModTime(),
					})
				}
			}

			sort.Slice(fileInfos, func(i, j int) bool {
				return fileInfos[i].modTime.After(fileInfos[j].modTime)
			})

			files := make([]string, len(fileInfos))
			for i, fi := range fileInfos {
				files[i] = fi.path
			}

			output := GlobOutput{Files: files}
			resultText := fmt.Sprintf("Found %d files:\n%s", len(files), formatFileList(files))

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: resultText},
				},
			}, output, nil
		},
	)
}

func containsDoubleStar(pattern string) bool {
	return filepath.ToSlash(pattern) != pattern &&
		(filepath.Clean(pattern) != pattern ||
			filepath.IsAbs(pattern))
}

func walkGlob(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			return err
		}

		if matched && !info.IsDir() {
			matches = append(matches, path)
		}

		if matchPattern(relPath, pattern) && !info.IsDir() {
			matches = append(matches, path)
		}

		return nil
	})
	return matches, err
}

func matchPattern(path, pattern string) bool {

	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	parts := splitPattern(pattern)
	pathParts := splitPath(path)

	return matchParts(pathParts, parts)
}

func splitPattern(pattern string) []string {
	return filepath.SplitList(strings.ReplaceAll(pattern, "/", string(filepath.ListSeparator)))
}

func splitPath(path string) []string {
	return filepath.SplitList(strings.ReplaceAll(path, "/", string(filepath.ListSeparator)))
}

func matchParts(pathParts, patternParts []string) bool {
	pi, ppi := 0, 0

	for pi < len(pathParts) && ppi < len(patternParts) {
		if patternParts[ppi] == "**" {
			if ppi == len(patternParts)-1 {
				return true
			}
			ppi++

			for pi < len(pathParts) {
				if matchParts(pathParts[pi:], patternParts[ppi:]) {
					return true
				}
				pi++
			}
			return false
		}

		matched, _ := filepath.Match(patternParts[ppi], pathParts[pi])
		if !matched {
			return false
		}

		pi++
		ppi++
	}

	return pi == len(pathParts) && ppi == len(patternParts)
}

func formatFileList(files []string) string {
	if len(files) == 0 {
		return "(none)"
	}
	var sb strings.Builder
	for _, f := range files {
		sb.WriteString(f)
		sb.WriteString("\n")
	}
	return sb.String()
}
