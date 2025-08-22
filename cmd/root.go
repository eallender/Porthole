package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "porthole",
	Aliases: []string{"phole", "ph"},
	Version: version,
	Short:   "Validate service communication within container networks",
	Long: `Porthole is a CLI tool that performs declarative validation of microservice 
		communication patterns in Docker networks. It captures and analyzes network traffic,
		comparing actual communication against expected specifications defined in YAML.

		Perfect for integration testing, CI/CD automation, and ensuring your microservices
		communicate according to their architectural design.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

var config string
var network string
var output string
var logLevel string
var quiet bool
var verbose bool
var logFile string
var dryRun bool

func init() {
	// Config and Env
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "$HOME/.Porthole.yaml", "WIP: Config file (default is $HOME/.Porthole.yaml)") //TODO Implement
	rootCmd.PersistentFlags().StringVar(&network, "network", "", "WIP: Select the network to validate")                                        //TODO Implement

	// Output and Logging
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "WIP: Output file path") //TODO Implement
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", `Set the logging level ("debug", "info", "warn", "error") (default "info")`)
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "WIP: Suppress the output") //TODO Implement
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "WIP: Verbose output")  //TODO Implement
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "WIP: Write logs to file")    //TODO Implement

	// Runtime
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "WIP: Show function without execution") //TODO Implement

	rootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose")
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "log-level")
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "log-level")
	cobra.OnInitialize(setupLogger)
}

func setupLogger() {
	level := parseLogLevel(logLevel)
	handler := &customHandler{level: level}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Debug("Setting log level to", "level", level.String())
}

type customHandler struct {
	level slog.Level
}

func (h *customHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(_ context.Context, r slog.Record) error {
	levelStr := fmt.Sprintf("[%-5s] ", strings.ToUpper(r.Level.String()))

	var msg strings.Builder
	msg.WriteString(r.Message)

	r.Attrs(func(a slog.Attr) bool {
		msg.WriteString(" ")
		msg.WriteString(a.Value.String())
		return true
	})

	_, err := fmt.Fprintf(os.Stderr, "%s%s\n", levelStr, msg.String())
	return err
}

// WithAttrs and WithGroup return the same handler for simplicity.
// This CLI tool doesn't require full structured logging features.
func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *customHandler) WithGroup(name string) slog.Handler       { return h }

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
