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
	Short: "Initialize a new Mycelium resume network",
	Run: func(cmd *cobra.Command, args []string) {
		printBrand()

		// 1. Initialize Git
		_, err := git.PlainInit(".", false)
		if err != nil {
			fmt.Println("[ERROR] Failed to initialize repository. Is this folder already a git repo?")
			return
		}

		// 2. Create the Mock Resume (John Doe)
		mockData := `{
  "basics": {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1 555-0199",
    "linkedin": "linkedin.com/in/johndoe",
    "github": "github.com/johndoe"
  },
  "sectionOrder": ["education", "skills", "experience", "projects"],
  "education": [
    {
      "school": "University of Technology",
      "degree": "B.S. in Computer Science",
      "date": "2018 - 2022",
      "cgpa": "3.9/4.0",
      "location": "San Francisco, CA"
    }
  ],
  "skills": {
    "Languages": "Golang, Python, TypeScript, SQL",
    "Cloud": "AWS, Docker, Kubernetes",
    "AI/ML": "PyTorch, Scikit-Learn, OpenAI API"
  },
  "experience": [
    {
      "company": "Tech Solutions Inc.",
      "role": "Software Engineer",
      "date": "2022 - Present",
      "points": [
        "Led development of a high-throughput data pipeline in Go.",
        "Reduced cloud infrastructure costs by 25% through container optimization.",
        "Mentored junior developers on best practices for version control."
      ]
    }
  ],
  "projects": [
    {
      "name": "Distributed Crawler",
      "tech": "Golang, Redis, Docker",
      "points": [
        "Built a concurrent web crawler capable of processing 10k pages/minute.",
        "Implemented Redis-based deduplication logic to prevent redundant crawls."
      ]
    }
  ]
}`

		err = os.WriteFile("resume.json", []byte(mockData), 0644)
		if err != nil {
			fmt.Println("[ERROR] Failed to create resume.json:", err)
			return
		}

		// 3. Create .gitignore
		ignore := "*.pdf\nmycelium.exe\nnode_modules/\n.DS_Store\n"
		os.WriteFile(".gitignore", []byte(ignore), 0644)

		fmt.Println("[SUCCESS] Mycelium network initialized successfully.")
		fmt.Println("[INFO] 'resume.json' has been seeded with a professional template.")
		fmt.Println("[INFO] Run 'mycelium edit' to begin tailoring your profile.")
	},
}
