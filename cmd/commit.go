package cmd

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func init() {
	// THIS LINE IS CRITICAL - it connects "commit" to the main tool
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Save current state of the resume",
	Run: func(cmd *cobra.Command, args []string) {
		msg, _ := cmd.Flags().GetString("message")
		if msg == "" {
			fmt.Println("[ERROR] Please provide a commit message: mycelium commit -m 'Updated skills'")
			return
		}

		// Open Repo
		r, err := git.PlainOpen(".")
		if err != nil {
			fmt.Println("Error: Not a mycelium repo. Run 'mycelium init' first.")
			return
		}

		w, _ := r.Worktree()

		// 1. Add resume.json
		_, err = w.Add("resume.json")
		if err != nil {
			fmt.Println("Error staging file:", err)
			return
		}

		// 2. Commit
		commit, err := w.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "mycelium User",
				Email: "user@mycelium.local",
				When:  time.Now(),
			},
		})

		if err != nil {
			fmt.Println("Error committing:", err)
			return
		}

		fmt.Printf("[SUCCESS] Version Saved! [%s] %s\n", commit.String()[:7], msg)
	},
}
