package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// --- DATA STRUCTURES ---
type Resume struct {
	Basics     Basics       `json:"basics"`
	Education  []Education  `json:"education"`
	Experience []Experience `json:"experience"`
	Skills     Skills       `json:"skills"`
}

type Basics struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Linkedin string `json:"linkedin"`
	Github   string `json:"github"`
	Website  string `json:"website"`
}

type Education struct {
	School   string `json:"school"`
	Location string `json:"location"`
	Degree   string `json:"degree"`
	Date     string `json:"date"`
}

type Experience struct {
	Company  string   `json:"company"`
	Role     string   `json:"role"`
	Location string   `json:"location"`
	Date     string   `json:"date"`
	Points   []string `json:"points"`
}

type Skills struct {
	Languages  []string `json:"languages"`
	Tools      []string `json:"tools"`
	Frameworks []string `json:"frameworks"`
}

// --- COMMAND SETUP ---
func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the Live Resume Editor",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Route: Serve the HTML Editor
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile("resume.json")
			if err != nil {
				http.Error(w, "Could not read resume.json", 500)
				return
			}

			// We pass the RAW JSON string to the template so JS can use it
			tmpl, err := template.New("editor").Parse(editorHTML)
			if err != nil {
				http.Error(w, "Template error: "+err.Error(), 500)
				return
			}

			tmpl.Execute(w, string(data))
		})

		// 2. Route: Save Changes (API)
		http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "Only POST allowed", 405)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Bad request", 400)
				return
			}

			// Verify it is valid JSON before saving
			var check Resume
			if err := json.Unmarshal(body, &check); err != nil {
				http.Error(w, "Invalid JSON format", 400)
				return
			}

			// Write to disk
			err = os.WriteFile("resume.json", body, 0644)
			if err != nil {
				http.Error(w, "Failed to write file", 500)
				return
			}

			w.WriteHeader(200)
			w.Write([]byte("Saved"))
			fmt.Println("💾 Resume saved successfully.")
		})

		// 3. Start Server
		fmt.Println("✅ Live Editor is running!")
		fmt.Println("🌍 Open: http://localhost:9090")
		http.ListenAndServe(":9090", nil)
	},
}

// --- THE FRONTEND (HTML + CSS + JS) ---
const editorHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>CVVC Live Editor</title>
    <style>
        body { margin: 0; padding: 0; display: flex; height: 100vh; font-family: sans-serif; overflow: hidden; }
        
        /* LEFT PANEL: EDITOR */
        .editor-panel {
            width: 40%;
            background: #1e1e1e;
            color: #d4d4d4;
            display: flex;
            flex-direction: column;
            border-right: 2px solid #333;
        }
        .editor-header { padding: 10px; background: #252526; font-weight: bold; display: flex; justify-content: space-between; align-items: center; }
        .save-btn { background: #0e639c; color: white; border: none; padding: 5px 15px; cursor: pointer; border-radius: 3px; }
        .save-btn:hover { background: #1177bb; }
        
        textarea {
            flex: 1;
            background: #1e1e1e;
            color: #dcdcdc;
            border: none;
            padding: 15px;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 14px;
            resize: none;
            outline: none;
        }

        /* RIGHT PANEL: PREVIEW */
        .preview-panel {
            width: 60%;
            background: #525659;
            padding: 20px;
            overflow-y: auto;
            display: flex;
            justify-content: center;
        }
        .paper {
            background: white;
            width: 210mm;
            min-height: 297mm;
            padding: 40px;
            box-shadow: 0 0 10px rgba(0,0,0,0.5);
            color: #333;
            font-family: "Times New Roman", Times, serif;
            line-height: 1.3;
        }

        /* RESUME STYLES (Jake's Format) */
        h1 { text-align: center; margin-bottom: 5px; text-transform: uppercase; font-size: 24pt; margin-top: 0; }
        .contact { text-align: center; margin-bottom: 20px; font-size: 11pt; }
        .section-title {
            font-size: 14pt;
            font-weight: bold;
            border-bottom: 1px solid black;
            margin-top: 15px;
            margin-bottom: 10px;
            padding-bottom: 2px;
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

    <!-- LEFT: JSON INPUT -->
    <div class="editor-panel">
        <div class="editor-header">
            <span>📝 Resume Data (JSON)</span>
            <button class="save-btn" onclick="saveData()">Save Version</button>
        </div>
        <textarea id="jsonInput" oninput="updatePreview()">{{.}}</textarea>
    </div>

    <!-- RIGHT: LIVE PREVIEW -->
    <div class="preview-panel">
        <div class="paper" id="resumePreview">
            <!-- Content injected by JS -->
        </div>
    </div>

    <script>
        // Initial Render
        const rawData = document.getElementById('jsonInput').value;
        let resumeData = {};
        
        try {
            resumeData = JSON.parse(rawData);
            renderResume(resumeData);
        } catch(e) { console.log("Init error"); }

        // Live Update Function
        function updatePreview() {
            const input = document.getElementById('jsonInput').value;
            try {
                resumeData = JSON.parse(input);
                renderResume(resumeData);
                // Auto-save logic could go here (debounce)
            } catch(e) {
                // Ignore JSON syntax errors while typing
            }
        }

        // Save Function
        async function saveData() {
            const input = document.getElementById('jsonInput').value;
            try {
                // Validate JSON
                JSON.parse(input);
                
                const response = await fetch('/save', {
                    method: 'POST',
                    body: input
                });
                
                if (response.ok) {
                    alert("Saved to disk!");
                } else {
                    alert("Error saving");
                }
            } catch(e) {
                alert("Invalid JSON! Fix errors before saving.");
            }
        }

        // Renderer (Converts JSON -> HTML)
        function renderResume(data) {
            const r = document.getElementById('resumePreview');
            
            // Basics
            let html = "<h1>" + (data.basics.name || "") + "</h1>";
            html += "<div class='contact'>";
            html += [data.basics.phone, data.basics.email, data.basics.linkedin, data.basics.github]
                    .filter(Boolean).join(" | ");
            html += "</div>";

            // Education
            if (data.education && data.education.length > 0) {
                html += "<div class='section-title'>Education</div>";
                data.education.forEach(edu => {
                    html += "<div class='entry'>";
                    html += "<div class='header-row'><span>" + edu.school + "</span><span>" + edu.location + "</span></div>";
                    html += "<div class='sub-row'><span>" + edu.degree + "</span><span>" + edu.date + "</span></div>";
                    html += "</div>";
                });
            }

            // Experience
            if (data.experience && data.experience.length > 0) {
                html += "<div class='section-title'>Experience</div>";
                data.experience.forEach(exp => {
                    html += "<div class='entry'>";
                    html += "<div class='header-row'><span>" + exp.company + "</span><span>" + exp.date + "</span></div>";
                    html += "<div class='sub-row'><span>" + exp.role + "</span><span>" + exp.location + "</span></div>";
                    if (exp.points && exp.points.length > 0) {
                        html += "<ul>";
                        exp.points.forEach(pt => html += "<li>" + pt + "</li>");
                        html += "</ul>";
                    }
                    html += "</div>";
                });
            }

            // Skills
            if (data.skills) {
                html += "<div class='section-title'>Technical Skills</div>";
                html += "<div>";
                if (data.skills.languages) html += "<strong>Languages:</strong> " + data.skills.languages.join(", ") + "<br>";
                if (data.skills.frameworks) html += "<strong>Frameworks:</strong> " + data.skills.frameworks.join(", ") + "<br>";
                if (data.skills.tools) html += "<strong>Tools:</strong> " + data.skills.tools.join(", ");
                html += "</div>";
            }

            r.innerHTML = html;
        }
    </script>
</body>
</html>
`
