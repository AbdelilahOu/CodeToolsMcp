package tools

import (
	"github.com/AbdelilahOu/CodeToolsMcp/internal/config"
	"github.com/AbdelilahOu/CodeToolsMcp/internal/runners"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {

	ripgrepRunner := runners.NewRipgrepRunner()
	gitRunner := runners.NewGitRunner()
	fileRunner := runners.NewFileRunner()

	NewGrepTool(ripgrepRunner).Register(s)
	NewGlobTool().Register(s)
	NewReadTool().Register(s)
	NewEditTool().Register(s)
	NewWriteTool().Register(s)
	NewGitStatusTool(gitRunner).Register(s)
	NewGitLogTool(gitRunner).Register(s)
	NewGitDiffTool(gitRunner).Register(s)
	NewGitShowTool(gitRunner).Register(s)
	NewGitBranchTool(gitRunner).Register(s)
	NewListDirTool(fileRunner).Register(s)
	NewDeleteTool(fileRunner).Register(s)
	NewRemoveTool(fileRunner).Register(s)
	NewCopyTool(fileRunner).Register(s)
	NewMoveTool(fileRunner).Register(s)
	NewTreeTool(fileRunner).Register(s)
	NewRunTool().Register(s)
}
