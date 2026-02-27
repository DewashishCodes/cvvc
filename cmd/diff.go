package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// Fully mapped structs for the diff engine
type DiffResume struct {
	Basics     DiffBasics       `json:"basics"`
	Education  []DiffEducation  `json:"education"`
	Experience []DiffExperience `json:"experience"`
	Skills     DiffSkills       `json:"skills"`
}

type DiffBasics struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type DiffEducation struct {
	School string `json:"school"`
	Degree string `json:"degree"`
}

type DiffExperience struct {
	Company  string   `json:"company"`
	Role     string   `json:"role"`
	Location string   `json:"location"`
	Date     string   `json:"date"`
	Points   []string `json:"points"`
}

type DiffSkills struct {
	Languages []string `json:"languages"`
	Tools     []string `json:"tools"`
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show all changes in your resume data",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Read Current from disk
		data, _ := os.ReadFile("resume.json")
		var current DiffResume
		json.Unmarshal(data, &current)

		// 2. Read Previous from Git
		r, _ := git.PlainOpen(".")
		ref, err := r.Head()
		if err != nil {
			fmt.Println("[ERROR] No commit history found. Commit once first.")
			return
		}
		commit, _ := r.CommitObject(ref.Hash())
		tree, _ := commit.Tree()
		file, _ := tree.File("resume.json")
		prevData, _ := file.Contents()
		var prev DiffResume
		json.Unmarshal([]byte(prevData), &prev)

		fmt.Printf("🔍 Diffing current changes against: %s\n", ref.Hash().String()[:7])
		fmt.Println("------------------------------------------------------------")
		hasChanges := false

		// --- CHECK BASICS ---
		if prev.Basics.Name != current.Basics.Name {
			fmt.Printf("[BASICS] Name: %s -> %s\n", prev.Basics.Name, current.Basics.Name)
			hasChanges = true
		}
		if prev.Basics.Phone != current.Basics.Phone {
			fmt.Printf("[BASICS] Phone: %s -> %s\n", prev.Basics.Phone, current.Basics.Phone)
			hasChanges = true
		}

		// --- CHECK EXPERIENCE (Comprehensive) ---
		if len(prev.Experience) != len(current.Experience) {
			fmt.Printf("[EXP] Number of jobs changed: %d -> %d\n", len(prev.Experience), len(current.Experience))
			hasChanges = true
		} else {
			for i := range current.Experience {
				p, c := prev.Experience[i], current.Experience[i]
				if p.Company != c.Company {
					fmt.Printf("[EXP] Company: %s -> %s\n", p.Company, c.Company)
					hasChanges = true
				}
				if p.Role != c.Role {
					fmt.Printf("[EXP] Role at %s: %s -> %s\n", c.Company, p.Role, c.Role)
					hasChanges = true
				}
				if !reflect.DeepEqual(p.Points, c.Points) {
					fmt.Printf("[EXP] Bullet points updated for %s\n", c.Company)
					hasChanges = true
				}
			}
		}

		// --- CHECK EDUCATION ---
		if len(prev.Education) != len(current.Education) {
			fmt.Printf("[EDU] Schools changed.\n")
			hasChanges = true
		} else {
			for i := range current.Education {
				if prev.Education[i].School != current.Education[i].School {
					fmt.Printf("[EDU] School: %s -> %s\n", prev.Education[i].School, current.Education[i].School)
					hasChanges = true
				}
			}
		}

		// --- CHECK SKILLS ---
		if !reflect.DeepEqual(prev.Skills.Languages, current.Skills.Languages) {
			fmt.Println("[SKILLS] Languages list was modified.")
			hasChanges = true
		}

		if !hasChanges {
			fmt.Println("✨ No changes found. Your JSON on disk matches the Git history.")
		}
	},
}
