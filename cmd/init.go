package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Mycelium repository",
	Run: func(cmd *cobra.Command, args []string) {
		printBrand() // Show the brand identity

		_, err := git.PlainInit(".", false)
		if err != nil {
			fmt.Println("[ERROR] Failed to initialize Git repository:", err)
			return
		}

		// (Keep your resume content JSON here from previous steps)
		content := `{ "basics": { "name": "Your Name" } }` // Use your full template here

		os.WriteFile("resume.json", []byte(content), 0644)
		os.WriteFile(".gitignore", []byte("*.pdf\nmycelium.exe\n"), 0644)

		fmt.Println("[SUCCESS] Mycelium network initialized in current directory.")
		fmt.Println("[INFO] Edit 'resume.json' or run 'mycelium edit' to begin.")
	},
}
