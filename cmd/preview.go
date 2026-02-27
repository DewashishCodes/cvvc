package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func init() {
	rootCmd.AddCommand(reviewCmd)
	// Add the role flag
	reviewCmd.Flags().StringP("role", "r", "General Software Engineer", "The specific job role you are targeting")
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Get an AI recruiter to critique your resume for a specific role",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := getAPIKey()
		if apiKey == "" {
			fmt.Println("[ERROR] API Key not found. Run: mycelium config --key YOUR_KEY")
			return
		}

		// Get the role from the flag
		targetRole, _ := cmd.Flags().GetString("role")

		fmt.Printf("[AI] AI Recruiter is analyzing your resume for the role: [%s]...\n", targetRole)

		// 1. Read Resume
		resumeData, err := os.ReadFile("resume.json")
		if err != nil {
			fmt.Println("[ERROR] Error: resume.json not found.")
			return
		}

		// 2. Setup Gemini
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			fmt.Println("[ERROR] Error connecting to AI:", err)
			return
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-2.5-flash")

		// 3. The Specialized Prompt
		prompt := fmt.Sprintf(`
			Act as a Senior Technical Recruiter who specializes in hiring for %s roles.
			I am going to provide you with a resume in JSON format.
			
			Critique this resume ONLY through the lens of a %s position.
			
			Please provide:
			1. A Role-Match Score (0-100).
			2. Missing Keywords: What technologies or skills typical for a %s are missing?
			3. Bullet Point Strengthening: Pick 2 bullet points and show how to make them more impressive for this specific role.
			4. "Red Flags": Anything that would make a recruiter for this role hesitate.
			
			RESUME DATA:
			%s
		`, targetRole, targetRole, targetRole, string(resumeData))

		// 4. Generate
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			fmt.Println("[ERROR] AI Error:", err)
			return
		}

		// 5. Print Result
		fmt.Printf("\n--- AI RECRUITER FEEDBACK FOR: %s ---\n", targetRole)
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
		fmt.Println("\n------------------------------------------------")
	},
}
