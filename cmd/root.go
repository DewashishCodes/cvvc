package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycelium",
	Short: "Mycelium: The Resume Versioning Network",
	Long: `Mycelium is a professional version control system for career data.
It manages resume iterations via Git-based branching, provides semantic 
analysis, and automates high-fidelity PDF generation.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR]", err)
		os.Exit(1)
	}
}

const BrandASCII = `
   __  __              _ _                      
  |  \/  |_   _  ___ ___| (_)_   _ _ __ ___      
  | |\/| | | | |/ __/ _ \ | | | | | '_ ' _ \     
  | |  | | |_| | (_|  __/ | | |_| | | | | | |    
  |_|  |_|\__, |\___\___|_|_|\__,_|_| |_| |_|    
          |___/                                  
  >>> THE RESUME VERSIONING NETWORK <<<
`

func printBrand() {
	fmt.Println(BrandASCII)
}
