package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current branch and unsaved changes",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen(".")
		if err != nil {
			fmt.Println("[ERROR] Not a mycelium repo. Run 'mycelium init'")
			return
		}

		// 1. Get current branch (With Safety Check)
		ref, err := r.Head()
		if err != nil {
			fmt.Println("📍 Current Branch: (initial branch)")
		} else {
			fmt.Printf("📍 Current Branch: %s\n", ref.Name().Short())
		}

		// 2. Check for changes
		w, _ := r.Worktree()
		status, _ := w.Status()

		if status.IsClean() {
			fmt.Println("✨ Everything is committed and saved.")
		} else {
			fmt.Println("[WARN]  Uncommitted changes found in resume.json!")
			fmt.Println("👉 Run 'mycelium commit -m \"message\"' to save your first version.")
		}
	},
}
