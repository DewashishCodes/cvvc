package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	Short: "Compile resume.json into a PDF",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⏳ Starting PDF generation...")

		// 1. Start a temporary local server to serve the HTML
		go startExportServer()
		time.Sleep(1 * time.Second) // Give server time to start

		// 2. Launch Local Browser (Edge or Chrome)
		fmt.Println("🚀 Launching browser engine...")

		path, _ := launcher.LookPath() // Finds Chrome/Edge on Windows
		if path == "" {
			fmt.Println("❌ Error: Could not find Chrome or Edge installed on this computer.")
			return
		}

		// Add .Leakless(false) to stop Windows Defender from panicking
		u := launcher.New().Bin(path).Leakless(false).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.MustClose()

		// 3. Open the page
		page := browser.MustPage("http://localhost:9091")
		page.MustWaitLoad()

		// 4. PDF Settings
		fmt.Println("🖨️  Rendering PDF...")

		pdfStream, err := page.PDF(&proto.PagePrintToPDF{
			PaperWidth:      toPtr(8.27),  // A4 Width
			PaperHeight:     toPtr(11.69), // A4 Height
			MarginTop:       toPtr(0.0),
			MarginBottom:    toPtr(0.0),
			MarginLeft:      toPtr(0.0),
			MarginRight:     toPtr(0.0),
			PrintBackground: true,
		})

		if err != nil {
			fmt.Println("Error generating PDF:", err)
			return
		}

		// 5. Save to file
		pdfBytes, err := io.ReadAll(pdfStream)
		if err != nil {
			fmt.Println("Error reading PDF stream:", err)
			return
		}

		outputName := "resume.pdf"
		err = os.WriteFile(outputName, pdfBytes, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}

		fmt.Printf("✅ Success! Resume exported to '%s'\n", outputName)
	},
}

// --- HELPER: SERVER LOGIC ---
func startExportServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("resume.json")
		var resume Resume
		json.Unmarshal(data, &resume)
		tmpl, _ := template.New("export").Parse(exportHTML)
		tmpl.Execute(w, resume)
	})
	http.ListenAndServe(":9091", nil)
}

// --- HELPER: POINTER CONVERTER ---
func toPtr(v float64) *float64 {
	return &v
}

// --- JAKE'S HTML ---
const exportHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Resume Export</title>
    <style>
        @page { size: A4; margin: 0; }
        body {
            font-family: "Times New Roman", Times, serif;
            width: 210mm;
            min-height: 297mm;
            margin: 0 auto;
            padding: 40px;
            box-sizing: border-box;
            color: #333;
            line-height: 1.3;
        }
        h1 { text-align: center; margin-bottom: 5px; text-transform: uppercase; font-size: 24pt; margin-top: 0; }
        .contact { text-align: center; margin-bottom: 20px; font-size: 11pt; }
        .section-title {
            font-size: 14pt;
            font-weight: bold;
            border-bottom: 1px solid black;
            margin-top: 15px;
            margin-bottom: 10px;
            padding-bottom: 2px;s
            text-transform: uppercase;
        }
        .entry { margin-bottom: 10px; }
        .header-row { display: flex; justify-content: space-between; font-weight: bold; }
        .sub-row { display: flex; justify-content: space-between; font-style: italic; }
        ul { margin: 5px 0; padding-left: 20px; }
        li { margin-bottom: 2px; }
    </style>
</head>
<body>
    <h1>{{.Basics.Name}}</h1>
    <div class="contact">
        {{.Basics.Phone}} | {{.Basics.Email}} | {{.Basics.Linkedin}} | {{.Basics.Github}}
    </div>

    <div class="section-title">Education</div>
    {{range .Education}}
    <div class="entry">
        <div class="header-row">
            <span>{{.School}}</span>
            <span>{{.Location}}</span>
        </div>
        <div class="sub-row">
            <span>{{.Degree}}</span>
            <span>{{.Date}}</span>
        </div>
    </div>
    {{end}}

    <div class="section-title">Experience</div>
    {{range .Experience}}
    <div class="entry">
        <div class="header-row">
            <span>{{.Company}}</span>
            <span>{{.Date}}</span>
        </div>
        <div class="sub-row">
            <span>{{.Role}}</span>
            <span>{{.Location}}</span>
        </div>
        <ul>
            {{range .Points}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    </div>
    {{end}}

    <div class="section-title">Technical Skills</div>
    <div>
        <strong>Languages:</strong> {{range .Skills.Languages}}{{.}}, {{end}}<br>
        <strong>Frameworks:</strong> {{range .Skills.Frameworks}}{{.}}, {{end}}<br>
        <strong>Tools:</strong> {{range .Skills.Tools}}{{.}}, {{end}}
    </div>
</body>
</html>
`
