package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io" // Used now!
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Generate Dewashish's Professional PDF",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Start temporary server for the PDF engine
		go func() {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				data, _ := os.ReadFile("resume.json")
				var res interface{}
				json.Unmarshal(data, &res)
				tmpl, _ := template.New("p").Parse(pdfHTML)
				tmpl.Execute(w, res)
			})
			http.ListenAndServe(":9091", nil)
		}()

		fmt.Println("⏳ Starting PDF generation...")
		time.Sleep(1 * time.Second)

		// 2. Launch Local Browser
		path, _ := launcher.LookPath()
		if path == "" {
			fmt.Println("[ERROR] Error: Could not find Chrome or Edge.")
			return
		}

		u := launcher.New().Bin(path).Leakless(false).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.MustClose()

		page := browser.MustPage("http://localhost:9091")
		page.MustWaitLoad()

		// 3. Capture PDF Stream
		fmt.Println("[INFO] Rendering PDF...")
		pdfStream, err := page.PDF(&proto.PagePrintToPDF{
			PrintBackground: true,
			PaperWidth:      toPtr(8.27),
			PaperHeight:     toPtr(11.69),
			MarginTop:       toPtr(0),
			MarginBottom:    toPtr(0),
			MarginLeft:      toPtr(0),
			MarginRight:     toPtr(0),
		})

		if err != nil {
			fmt.Println("[ERROR] Error rendering:", err)
			return
		}

		// 4. FIX: Convert Stream to Bytes using 'io'
		pdfBytes, err := io.ReadAll(pdfStream)
		if err != nil {
			fmt.Println("[ERROR] Error reading stream:", err)
			return
		}

		// 5. Save file
		outputName := "Dewashish_Resume.pdf"
		err = os.WriteFile(outputName, pdfBytes, 0644)
		if err != nil {
			fmt.Println("[ERROR] Error saving file:", err)
			return
		}

		fmt.Printf("[SUCCESS] Success! Exported to %s\n", outputName)
	},
}

// Helper for Rod library
func toPtr(f float64) *float64 { return &f }

// --- THE CSS/HTML DESIGN (Matches your PDF) ---
const pdfHTML = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: "Times New Roman", serif; padding: 45px; line-height: 1.15; color: black; }
        .name { text-align: center; font-size: 28pt; margin: 0; }
        .contact { text-align: center; font-size: 11pt; border-bottom: 1.5px solid black; padding-bottom: 6px; margin-bottom: 10px; }
        .section { font-weight: bold; text-transform: uppercase; border-bottom: 1.5px solid black; margin-top: 15px; font-size: 12.5pt; }
        .row { display: flex; justify-content: space-between; font-weight: bold; margin-top: 5px; font-size: 11pt; }
        .sub-row { display: flex; justify-content: space-between; font-style: italic; font-size: 10.5pt; }
        ul { margin: 4px 0; padding-left: 18px; }
        li { margin-bottom: 1.5px; font-size: 10.5pt; text-align: justify; }
    </style>
</head>
<body>
    <div class="name">{{.basics.name}}</div>
    <div class="contact">{{.basics.phone}} | {{.basics.email}} | LinkedIn | Github</div>
    <div class="section">Education</div>
    {{range .education}}<div class="row"><span>{{.school}}</span><span>{{.date}}</span></div><div class="sub-row"><span>{{.degree}}</span><span>(current): {{.cgpa}}</span></div>{{end}}
    <div class="section">Technical Skills</div>
    <div style="font-size: 10.5pt; margin-top: 5px;">{{range $cat, $val := .skills}}<strong>{{$cat}}:</strong> {{$val}}<br>{{end}}</div>
    <div class="section">Experience</div>
    {{range .experience}}<div class="row"><span>{{.company}}</span><span>{{.date}}</span></div><div class="sub-row"><span>{{.role}}</span></div><ul>{{range .points}}<li>{{.}}</li>{{end}}</ul>{{end}}
    <div class="section">Projects</div>
    {{range .projects}}<div class="row"><span>{{.name}} | <span style="font-weight:normal; font-style:italic;">{{.tech}}</span></span></div><ul>{{range .points}}<li>{{.}}</li>{{end}}</ul>{{end}}
</body>
</html>`
