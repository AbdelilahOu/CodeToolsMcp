package tools

import (
	"github.com/AbdelilahOu/CodeToolsMcp/internal/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterTools(s *mcp.Server, cfg *config.Config) {

	ripgrepRunner := NewRipgrepRunner()

	NewGrepTool(ripgrepRunner).Register(s)
	NewGlobTool().Register(s)
	NewReadTool().Register(s)
	NewEditTool().Register(s)
	NewWriteTool().Register(s)
}
