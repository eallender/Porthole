package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Services map[string]Service `yaml:"services"`
	Network  Network            `yaml:"network"`
	Timeout  string             `yaml:"timeout,omitempty"`
}

type Service struct {
	Name    string      `yaml:"name,omitempty"`
	Expects Expectation `yaml:"expects"`
}

type Expectation struct {
	From []Endpoint `yaml:"from,omitempty"`
	To   []Endpoint `yaml:"to,omitempty"`
}

type Endpoint struct {
	Name     string `yaml:"name,omitempty"`
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Method   string `yaml:"method,omitempty"`
}

type Network struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your network",
	Long: `Validate captures and analyzes network traffic to ensure microservices 
communicate according to their defined specifications.

This command monitors Docker network traffic for a specified duration and compares
the observed communication patterns against the expectations defined in your 
configuration file. It validates protocols, endpoints, ports, and message schemas
to ensure your microservices are communicating as designed.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filePath string
		if len(args) > 0 {
			filePath = args[0]
		}
		if filePath == "" {
			var err error
			filePath, err = searchPortholeFile()
			if err != nil {
				slog.Error("", slog.Any("err", err))
				os.Exit(2)
			}
		}

		fileContents, err := getYamlContents(filePath)
		if err != nil {
			slog.Error("", slog.Any("err", err))
			os.Exit(2)
		}

		config, err := parseYamlString(fileContents)
		if err != nil {
			slog.Error("", slog.Any("err", err))
			os.Exit(2)
		}

		slog.Info("Parsed config:\n", slog.Any("config", config))
	},
}

var composeFile bool
var services string

func init() {
	rootCmd.AddCommand(validateCmd)

	// Local flags for validate command only
	validateCmd.Flags().BoolVar(&composeFile, "compose-file", false, "WIP: Provide a compose file to validate")                 //TODO Implement
	validateCmd.Flags().StringVar(&services, "services", "", `WIP: CSV of services you want to validate ("service1,service2")`) //TODO Implement
}

func getYamlContents(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file %s does not exist", filePath)
	}

	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(contents), nil
}

func (c *Config) PopulateServiceNames() {
	for key, service := range c.Services {
		if service.Name == "" {
			service.Name = key
			c.Services[key] = service
		}
	}
}

func parseYamlString(contents string) (*Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(contents), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	config.PopulateServiceNames()

	return &config, nil
}

func searchPortholeFile() (string, error) {
	pattern := `^(?i)(porthole|p?hole|ph)\.ya?ml$`
	regex := regexp.MustCompile(pattern)

	files, err := filepath.Glob("*")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if regex.MatchString(file) {
			return file, nil
		}
	}

	slog.Error("Unable to find valid Porthole file")
	return "", errors.New("unable to find valid Porthole file")
}
