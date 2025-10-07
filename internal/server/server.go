package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/config"
	"github.com/AbdelilahOu/CodeToolsMcp/internal/logger"
	"github.com/AbdelilahOu/CodeToolsMcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPServerConfig struct {
	Version string
	Config  *config.Config
}

func NewMCPServer(cfg MCPServerConfig) (*mcp.Server, error) {
	logCfg := logger.ConfigFromLoggingConfig(cfg.Config.Logging)
	if err := logger.Initialize(logCfg); err != nil {
		fmt.Printf("Warning: Failed to initialize logger: %v\n", err)
	} else {
		logger.Info("Logger initialized successfully", map[string]interface{}{
			"level":       logger.LogLevelString(logCfg.Level),
			"output_file": logCfg.OutputFile,
			"console":     logCfg.Console,
		})
	}

	impl := &mcp.Implementation{Name: "CodeToolsMcp", Version: cfg.Version}
	server := mcp.NewServer(impl, nil)

	logger.Info("Code Tools MCP Server starting", map[string]interface{}{
		"version": cfg.Version,
	})

	tools.RegisterTools(server, cfg.Config)

	return server, nil
}

type StdioServerConfig struct {
	Version string
	Config  *config.Config
}

func RunStdioServer(cfg StdioServerConfig) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if err := logger.Shutdown(); err != nil {
			fmt.Printf("Error shutting down logger: %v\n", err)
		}
	}()

	server, err := NewMCPServer(MCPServerConfig{
		Version: cfg.Version,
		Config:  cfg.Config,
	})

	if err != nil {
		logger.Error("Failed to create MCP server", err)
		return fmt.Errorf("failed to create MCP server: %w", err)
	}

	logger.Info("Code Tools MCP Server started and running", map[string]interface{}{
		"version": cfg.Version,
	})
	fmt.Printf("Code Tools MCP Server running ...\n")

	err = server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		logger.Error("Server stopped with error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}

	return err
}
