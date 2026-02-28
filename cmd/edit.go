package cmd

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the Mycelium Live Form Editor",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile("resume.json")
			if err != nil {
				http.Error(w, "[ERROR] resume.json not found. Run 'mycelium init' first.", 404)
				return
			}
			tmpl, _ := template.New("editor").Parse(editorHTML)
			// Pass as template.HTML to preserve JSON quotes
			tmpl.Execute(w, template.HTML(data))
		})

		http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				return
			}
			body, _ := io.ReadAll(r.Body)
			err := os.WriteFile("resume.json", body, 0644)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
		})

		fmt.Println("[INFO] Mycelium Editor started.")
		fmt.Println("[INFO] Local Network Link: http://localhost:9090")
		fmt.Println("[INFO] Press Ctrl+C to disconnect from the network.")
		http.ListenAndServe(":9090", nil)
	},
}

const editorHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Mycelium Editor</title>
    <style>
        :root { --bg: #f4f7f6; --sidebar: #1a1c1e; --border: #e0e0e0; --primary: #007bff; }
        body { margin: 0; display: flex; height: 100vh; font-family: 'Segoe UI', system-ui, sans-serif; background: var(--bg); overflow: hidden; }
        
        /* SIDEBAR NAV */
        .nav-sidebar { width: 70px; background: var(--sidebar); display: flex; flex-direction: column; align-items: center; padding-top: 30px; gap: 30px; }
        .nav-item { color: #5f6368; cursor: pointer; font-size: 24px; transition: 0.3s; padding: 12px; border-radius: 12px; }
        .nav-item.active { color: white; background: #3c4043; }
        .nav-item:hover { color: #fff; }

        /* FORM PANEL */
        .form-panel { width: 480px; background: white; border-right: 1px solid var(--border); display: flex; flex-direction: column; box-shadow: 2px 0 5px rgba(0,0,0,0.05); }
        .form-header { padding: 25px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
        .form-content { flex: 1; overflow-y: auto; padding: 25px; }
        
        label { display: block; font-size: 11px; font-weight: 700; text-transform: uppercase; color: #5f6368; margin-top: 20px; margin-bottom: 8px; letter-spacing: 0.5px; }
        input, textarea { width: 100%; padding: 12px; border: 1px solid #dadce0; border-radius: 6px; box-sizing: border-box; font-size: 14px; transition: border 0.2s; }
        input:focus { border-color: var(--primary); outline: none; }
        
        .card { border: 1px solid #e8eaed; padding: 20px; border-radius: 10px; margin-bottom: 25px; background: #fdfdfd; position: relative; }
        .controls { display: flex; justify-content: flex-end; gap: 10px; margin-top: 15px; border-top: 1px solid #eee; padding-top: 10px; }
        
        .btn-add { background: #fff; color: var(--primary); border: 2px dashed #d2e3fc; width: 100%; padding: 15px; cursor: pointer; border-radius: 8px; font-weight: 600; margin-bottom: 40px; }
        .btn-add:hover { background: #f8faff; border-color: var(--primary); }
        .save-btn { background: #1a73e8; color: white; border: none; padding: 10px 25px; border-radius: 6px; cursor: pointer; font-weight: 700; }
        .save-btn:hover { background: #185abc; }
        .btn-sm { padding: 6px 12px; font-size: 11px; border: 1px solid #dadce0; background: white; cursor: pointer; border-radius: 4px; font-weight: 600; }
        .btn-danger { color: #d93025; border-color: #f5c2c7; }

        /* PREVIEW PANEL */
        .preview-panel { flex: 1; background: #525659; overflow-y: auto; display: flex; justify-content: center; padding: 50px 0; }
        .paper { background: white; width: 210mm; min-height: 297mm; padding: 50px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); font-family: 'Times New Roman', Times, serif; color: black; }

        /* PREVIEW CONTENT */
        .res-name { text-align: center; font-size: 26pt; border-bottom: 1.5px solid black; padding-bottom: 5px; margin-bottom: 10px; }
        .res-contact { text-align: center; font-size: 11pt; margin-bottom: 10px; }
        .res-sec { font-weight: bold; text-transform: uppercase; border-bottom: 1px solid black; margin-top: 18px; font-size: 12.5pt; margin-bottom: 6px; }
        .res-row { display: flex; justify-content: space-between; font-weight: bold; margin-top: 5px; font-size: 11pt; }
        ul { margin: 5px 0; padding-left: 20px; }
        li { font-size: 10.5pt; margin-bottom: 3px; text-align: justify; }
    </style>
</head>
<body>
    <div class="nav-sidebar">
        <div class="nav-item active" onclick="tab('basics', this)" title="Profile">👤</div>
        <div class="nav-item" onclick="tab('education', this)" title="Education">🎓</div>
        <div class="nav-item" onclick="tab('experience', this)" title="Work Experience">💼</div>
        <div class="nav-item" onclick="tab('projects', this)" title="Projects">🚀</div>
        <div class="nav-item" onclick="tab('skills', this)" title="Technical Skills">🛠️</div>
        <div class="nav-item" onclick="tab('order', this)" title="Reorder Sections">🔃</div>
    </div>

    <div class="form-panel">
        <div class="form-header">
            <h3 id="tab-title" style="margin:0">Basics</h3>
            <button class="save-btn" onclick="save()">SAVE</button>
        </div>
        <div class="form-content" id="form-area"></div>
    </div>

    <div class="preview-panel">
        <div id="capture-area" class="paper"></div>
    </div>

    <script id="data-raw" type="text/plain">{{.}}</script>

    <script>
        let resume = JSON.parse(document.getElementById('data-raw').textContent);
        
        // Initial Defaults if missing
        if (!resume.sectionOrder) resume.sectionOrder = ['education', 'skills', 'experience', 'projects'];
        if (!resume.education) resume.education = [];
        if (!resume.experience) resume.experience = [];
        if (!resume.projects) resume.projects = [];
        
        let currentTab = 'basics';

        function tab(t, el) {
            currentTab = t;
            document.querySelectorAll('.nav-item').forEach(i => i.classList.remove('active'));
            el.classList.add('active');
            document.getElementById('tab-title').innerText = t.charAt(0).toUpperCase() + t.slice(1);
            renderForm();
        }

        function renderForm() {
            const area = document.getElementById('form-area');
            area.innerHTML = '';

            if (currentTab === 'basics') {
                area.innerHTML = '<label>Full Name</label><input id="inp-name" value="' + resume.basics.name + '">' +
                                 '<label>Email</label><input id="inp-email" value="' + resume.basics.email + '">' +
                                 '<label>Phone</label><input id="inp-phone" value="' + resume.basics.phone + '">' +
                                 '<label>LinkedIn</label><input id="inp-link" value="' + (resume.basics.linkedin || "") + '">' +
                                 '<label>GitHub</label><input id="inp-git" value="' + (resume.basics.github || "") + '">';
                
                const set = (f,v) => { resume.basics[f] = v; render(); };
                document.getElementById('inp-name').oninput = (e) => set('name', e.target.value);
                document.getElementById('inp-email').oninput = (e) => set('email', e.target.value);
                document.getElementById('inp-phone').oninput = (e) => set('phone', e.target.value);
                document.getElementById('inp-link').oninput = (e) => set('linkedin', e.target.value);
                document.getElementById('inp-git').oninput = (e) => set('github', e.target.value);

            } else if (currentTab === 'education') {
                resume.education.forEach((edu, i) => {
                    let card = createCard('education', i);
                    card.innerHTML += '<label>Institution</label><input value="' + edu.school + '" oninput="resume.education['+i+'].school=this.value;render()">' +
                                     '<label>Degree</label><input value="' + edu.degree + '" oninput="resume.education['+i+'].degree=this.value;render()">' +
                                     '<label>Date Range</label><input value="' + edu.date + '" oninput="resume.education['+i+'].date=this.value;render()">' +
                                     '<label>Score (GPA/CGPA)</label><input value="' + edu.cgpa + '" oninput="resume.education['+i+'].cgpa=this.value;render()">';
                    area.appendChild(card);
                });
                area.appendChild(createAddBtn('education', {school:'', degree:'', date:'', cgpa:''}));

            } else if (currentTab === 'experience') {
                resume.experience.forEach((exp, i) => {
                    let card = createCard('experience', i);
                    card.innerHTML += '<label>Company</label><input value="' + exp.company + '" oninput="resume.experience['+i+'].company=this.value;render()">' +
                                     '<label>Role</label><input value="' + exp.role + '" oninput="resume.experience['+i+'].role=this.value;render()">' +
                                     '<label>Date</label><input value="' + exp.date + '" oninput="resume.experience['+i+'].date=this.value;render()">' +
                                     '<label>Bullet Points (New line for each)</label><textarea rows="5" oninput="resume.experience['+i+'].points=this.value.split(\'\\n\');render()">' + exp.points.join('\n') + '</textarea>';
                    area.appendChild(card);
                });
                area.appendChild(createAddBtn('experience', {company:'', role:'', date:'', points:[]}));

            } else if (currentTab === 'projects') {
                resume.projects.forEach((prj, i) => {
                    let card = createCard('projects', i);
                    card.innerHTML += '<label>Project Name</label><input value="' + prj.name + '" oninput="resume.projects['+i+'].name=this.value;render()">' +
                                     '<label>Technologies</label><input value="' + prj.tech + '" oninput="resume.projects['+i+'].tech=this.value;render()">' +
                                     '<label>Details</label><textarea rows="4" oninput="resume.projects['+i+'].points=this.value.split(\'\\n\');render()">' + prj.points.join('\n') + '</textarea>';
                    area.appendChild(card);
                });
                area.appendChild(createAddBtn('projects', {name:'', tech:'', points:[]}));

            } else if (currentTab === 'skills') {
                for (let k in resume.skills) {
                    let d = document.createElement('div');
                    d.innerHTML = '<label>'+k+'</label><textarea rows="3" oninput="resume.skills[\''+k+'\']=this.value;render()">'+resume.skills[k]+'</textarea>';
                    area.appendChild(d);
                }
            } else if (currentTab === 'order') {
                resume.sectionOrder.forEach((sec, i) => {
                    let d = document.createElement('div'); d.className = 'card'; d.style.display = 'flex'; d.style.justifyContent = 'space-between'; d.style.alignItems = 'center';
                    d.innerHTML = '<span style="font-weight:bold; font-size:12px">'+sec.toUpperCase()+'</span>' +
                                  '<div class="controls"><button class="btn-sm" onclick="move(\'sectionOrder\','+i+',-1)">Up</button> <button class="btn-sm" onclick="move(\'sectionOrder\','+i+',1)">Down</button></div>';
                    area.appendChild(d);
                });
            }
        }

        // --- UI HELPERS ---
        function createCard(arr, i) {
            let d = document.createElement('div'); d.className = 'card';
            let ctrl = document.createElement('div'); ctrl.className = 'controls';
            ctrl.innerHTML = '<button class="btn-sm" onclick="move(\''+arr+'\','+i+',-1)">Up</button>' +
                             '<button class="btn-sm" onclick="move(\''+arr+'\','+i+',1)">Down</button>' +
                             '<button class="btn-sm btn-danger" onclick="remove(\''+arr+'\','+i+')">Delete</button>';
            d.appendChild(ctrl);
            return d;
        }

        function createAddBtn(arr, template) {
            let b = document.createElement('button'); b.className = 'btn-add'; b.innerText = '+ Add Entry';
            b.onclick = () => { resume[arr].push(template); renderForm(); render(); };
            return b;
        }

        function move(arr, i, dir) {
            let t = i + dir; if (t < 0 || t >= resume[arr].length) return;
            [resume[arr][i], resume[arr][t]] = [resume[arr][t], resume[arr][i]];
            renderForm(); render();
        }

        function remove(arr, i) {
            if(confirm('Remove this entry?')) { resume[arr].splice(i, 1); renderForm(); render(); }
        }

        // --- PREVIEW RENDERER ---
        function render() {
            const paper = document.getElementById('capture-area');
            if (!resume.basics) return;
            
            let h = '<div class="res-name">' + (resume.basics.name || "Name") + '</div>';
            h += '<div class="res-contact">' + 
                 [resume.basics.phone, resume.basics.email, resume.basics.linkedin, resume.basics.github].filter(Boolean).join(" | ") + 
                 '</div>';

            (resume.sectionOrder || []).forEach(sec => {
                if (sec === 'education' && resume.education && resume.education.length) {
                    h += '<div class="res-sec">Education</div>';
                    resume.education.forEach(e => { 
                        h += '<div class="res-row"><span>'+(e.school||"")+'</span><span>'+(e.date||"")+'</span></div>';
                        h += '<div>'+(e.degree||"")+' | CGPA: '+(e.cgpa||"")+'</div>'; 
                    });
                } else if (sec === 'skills' && resume.skills) {
                    h += '<div class="res-sec">Technical Skills</div>';
                    for (let k in resume.skills) { h += '<div style="font-size:10.5pt"><strong>'+k+':</strong> '+resume.skills[k]+'</div>'; }
                } else if (sec === 'experience' && resume.experience && resume.experience.length) {
                    h += '<div class="res-sec">Experience</div>';
                    resume.experience.forEach(exp => {
                        h += '<div class="res-row"><span>'+(exp.company||"")+'</span><span>'+(exp.date||"")+'</span></div><div style="font-style:italic; font-size:10.5pt">'+(exp.role||"")+'</div>';
                        h += '<ul>' + (exp.points || []).map(p => p.trim() ? '<li>'+p+'</li>' : '').join('') + '</ul>';
                    });
                } else if (sec === 'projects' && resume.projects && resume.projects.length) {
                    h += '<div class="res-sec">Projects</div>';
                    resume.projects.forEach(p => {
                        h += '<div class="res-row"><span>'+(p.name||"")+' | <span style="font-weight:normal; font-style:italic;">'+(p.tech||"")+'</span></span></div>';
                        h += '<ul>' + (p.points || []).map(pt => pt.trim() ? '<li>'+pt+'</li>' : '').join('') + '</ul>';
                    });
                }
            });
            paper.innerHTML = h;
        }

        async function save() {
            const btn = document.querySelector('.save-btn');
            btn.innerText = 'SAVING...';
            const res = await fetch('/save', { method: 'POST', body: JSON.stringify(resume) });
            if(res.ok) { btn.innerText = 'SAVED!'; setTimeout(() => btn.innerText = 'SAVE', 2000); }
            else { alert('Save Failed'); btn.innerText = 'SAVE'; }
        }

        renderForm(); render();
    </script>
</body>
</html>
`
