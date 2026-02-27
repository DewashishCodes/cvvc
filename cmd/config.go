package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	GeminiKey string `json:"gemini_key"`
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().String("key", "", "Your Gemini API Key")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI API keys",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("[ERROR] Please provide a key: mycelium config --key YOUR_KEY")
			return
		}

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".mycelium_config.json")

		cfg := Config{GeminiKey: key}
		data, _ := json.MarshalIndent(cfg, "", "  ")

		os.WriteFile(configPath, data, 0644)
		fmt.Println("[SUCCESS] API Key saved to", configPath)
	},
}

// Helper to get the key in other files
func getAPIKey() string {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".mycelium_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg.GeminiKey
}
