package main

import (
	"fmt"
	"os"

	"github.com/AbdelilahOu/CodeToolsMcp/internal/config"
	"github.com/AbdelilahOu/CodeToolsMcp/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "CodeToolsMcp",
	Short: "Code Tools MCP Server - Claude Code compatible tools",
	Long:  `A Model Context Protocol (MCP) server providing the same powerful code tools that Claude Code uses: Grep (ripgrep), Glob, Read, Edit, Write, Git Status/Log/Diff/Show/Branch, filesystem helpers (list_dir, delete, remove, copy, move, tree), and Run.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path (for logging configuration)")

	stdioCmd := &cobra.Command{
		Use:   "stdio",
		Short: "Run over stdio transport (for local MCP clients)",
		RunE:  runStdioServer,
	}
	rootCmd.AddCommand(stdioCmd)
}

func runStdioServer(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		fmt.Println("Server will start with default configuration.")
		cfg = &config.Config{
			Logging: config.LoggingConfig{
				Level:      "INFO",
				OutputFile: "code-tools-mcp.log",
				MaxSizeMB:  10,
				Console:    true,
			},
		}
	}

	return server.RunStdioServer(server.StdioServerConfig{
		Version: "v0.1.0",
		Config:  cfg,
	})
}
