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
	Short: "Open the Full Feature Resume Editor",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile("resume.json")
			if err != nil {
				http.Error(w, "resume.json not found.", 404)
				return
			}
			tmpl, _ := template.New("editor").Parse(editorHTML)
			tmpl.Execute(w, template.HTML(data))
		})

		http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			os.WriteFile("resume.json", body, 0644)
			w.WriteHeader(200)
		})

		fmt.Println("[INFO] mycelium Editor: http://localhost:9090")
		http.ListenAndServe(":9090", nil)
	},
}

const editorHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>mycelium Master Editor</title>
    <style>
        :root { --bg: #f8f9fa; --sidebar: #212529; --primary: #0d6efd; }
        body { margin: 0; display: flex; height: 100vh; font-family: 'Segoe UI', sans-serif; background: var(--bg); overflow: hidden; }
        .nav-sidebar { width: 70px; background: var(--sidebar); display: flex; flex-direction: column; align-items: center; padding-top: 20px; gap: 25px; }
        .nav-item { color: #888; cursor: pointer; font-size: 22px; transition: 0.2s; padding: 10px; border-radius: 8px; }
        .nav-item.active { color: white; background: #343a40; }
        .form-panel { width: 450px; background: white; border-right: 1px solid #dee2e6; display: flex; flex-direction: column; }
        .form-header { padding: 20px; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; align-items: center; }
        .form-content { flex: 1; overflow-y: auto; padding: 20px; }
        label { display: block; font-size: 11px; font-weight: bold; text-transform: uppercase; color: #666; margin-bottom: 5px; margin-top: 15px; }
        input, textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; font-size: 14px; }
        .card { border: 1px solid #eee; padding: 15px; border-radius: 8px; margin-bottom: 20px; background: #fafafa; position: relative; }
        .btn-add { background: #e7f3ff; color: #007bff; border: 1px dashed #007bff; width: 100%; padding: 12px; cursor: pointer; border-radius: 6px; margin-top: 10px; font-weight: bold; }
        .save-btn { background: #198754; color: white; border: none; padding: 8px 20px; border-radius: 4px; cursor: pointer; font-weight: bold; }
        .preview-panel { flex: 1; background: #525659; overflow-y: auto; display: flex; justify-content: center; padding: 40px 0; }
        .paper { background: white; width: 210mm; min-height: 297mm; padding: 50px; box-shadow: 0 0 20px rgba(0,0,0,0.5); font-family: 'Times New Roman', serif; }
        .res-name { text-align: center; font-size: 26pt; border-bottom: 1.5px solid black; padding-bottom: 5px; margin-bottom: 10px; }
        .res-contact { text-align: center; font-size: 11pt; margin-bottom: 10px; }
        .res-sec { font-weight: bold; text-transform: uppercase; border-bottom: 1px solid black; margin-top: 15px; font-size: 12pt; }
        .res-row { display: flex; justify-content: space-between; font-weight: bold; margin-top: 4px; font-size: 11pt; }
        ul { margin: 5px 0; padding-left: 20px; }
        li { font-size: 10.5pt; margin-bottom: 2px; text-align: justify; }
        .del-link { color: #dc3545; font-size: 11px; cursor: pointer; text-decoration: underline; margin-top: 10px; display: inline-block; }
    </style>
</head>
<body>
    <div class="nav-sidebar">
        <div class="nav-item active" onclick="tab('basics', this)">👤</div>
        <div class="nav-item" onclick="tab('education', this)">🎓</div>
        <div class="nav-item" onclick="tab('experience', this)">💼</div>
        <div class="nav-item" onclick="tab('projects', this)">[INFO]</div>
        <div class="nav-item" onclick="tab('skills', this)">🛠️</div>
        <div class="nav-item" onclick="tab('order', this)">🔃</div>
    </div>
    <div class="form-panel">
        <div class="form-header">
            <h3 id="tab-title">Basics</h3>
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
        if (!resume.sectionOrder) resume.sectionOrder = ['education', 'skills', 'experience', 'projects'];
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
                area.innerHTML = '<label>Full Name</label><input id="i-n" value="' + resume.basics.name + '">' +
                                 '<label>Email</label><input id="i-e" value="' + resume.basics.email + '">' +
                                 '<label>Phone</label><input id="i-p" value="' + resume.basics.phone + '">';
                document.getElementById('i-n').oninput = (e) => { resume.basics.name = e.target.value; render(); };
                document.getElementById('i-e').oninput = (e) => { resume.basics.email = e.target.value; render(); };
                document.getElementById('i-p').oninput = (e) => { resume.basics.phone = e.target.value; render(); };

            } else if (currentTab === 'education') {
                resume.education.forEach((edu, i) => {
                    let card = document.createElement('div'); card.className = 'card';
                    card.innerHTML = '<label>School</label><input class="e-s" value="' + edu.school + '">' +
                                     '<label>Degree</label><input class="e-d" value="' + edu.degree + '">' +
                                     '<label>Date</label><input class="e-t" value="' + edu.date + '">' +
                                     '<label>CGPA</label><input class="e-c" value="' + edu.cgpa + '">' +
                                     '<span class="del-link" onclick="remove(\'education\','+i+')">Remove Education</span>';
                    card.querySelector('.e-s').oninput = (e) => { resume.education[i].school = e.target.value; render(); };
                    card.querySelector('.e-d').oninput = (e) => { resume.education[i].degree = e.target.value; render(); };
                    card.querySelector('.e-t').oninput = (e) => { resume.education[i].date = e.target.value; render(); };
                    card.querySelector('.e-c').oninput = (e) => { resume.education[i].cgpa = e.target.value; render(); };
                    area.appendChild(card);
                });
                let b = document.createElement('button'); b.className = 'btn-add'; b.innerText = '+ Add School';
                b.onclick = () => { resume.education.push({school:'', degree:'', date:'', cgpa:''}); renderForm(); render(); };
                area.appendChild(b);

            } else if (currentTab === 'experience') {
                resume.experience.forEach((exp, i) => {
                    let card = document.createElement('div'); card.className = 'card';
                    card.innerHTML = '<label>Company</label><input class="x-c" value="' + exp.company + '">' +
                                     '<label>Role</label><input class="x-r" value="' + exp.role + '">' +
                                     '<label>Date</label><input class="x-d" value="' + exp.date + '">' +
                                     '<label>Points</label><textarea class="x-p" rows="4">' + exp.points.join('\n') + '</textarea>' +
                                     '<span class="del-link" onclick="remove(\'experience\','+i+')">Remove Experience</span>';
                    card.querySelector('.x-c').oninput = (e) => { resume.experience[i].company = e.target.value; render(); };
                    card.querySelector('.x-r').oninput = (e) => { resume.experience[i].role = e.target.value; render(); };
                    card.querySelector('.x-d').oninput = (e) => { resume.experience[i].date = e.target.value; render(); };
                    card.querySelector('.x-p').oninput = (e) => { resume.experience[i].points = e.target.value.split('\n'); render(); };
                    area.appendChild(card);
                });
                let b = document.createElement('button'); b.className = 'btn-add'; b.innerText = '+ Add Experience';
                b.onclick = () => { resume.experience.push({company:'', role:'', date:'', points:[]}); renderForm(); render(); };
                area.appendChild(b);

            } else if (currentTab === 'projects') {
                resume.projects.forEach((prj, i) => {
                    let card = document.createElement('div'); card.className = 'card';
                    card.innerHTML = '<label>Project Name</label><input class="p-n" value="' + prj.name + '">' +
                                     '<label>Tech Stack</label><input class="p-t" value="' + prj.tech + '">' +
                                     '<label>Description Points</label><textarea class="p-p" rows="4">' + prj.points.join('\n') + '</textarea>' +
                                     '<span class="del-link" onclick="remove(\'projects\','+i+')">Remove Project</span>';
                    card.querySelector('.p-n').oninput = (e) => { resume.projects[i].name = e.target.value; render(); };
                    card.querySelector('.p-t').oninput = (e) => { resume.projects[i].tech = e.target.value; render(); };
                    card.querySelector('.p-p').oninput = (e) => { resume.projects[i].points = e.target.value.split('\n'); render(); };
                    area.appendChild(card);
                });
                let b = document.createElement('button'); b.className = 'btn-add'; b.innerText = '+ Add Project';
                b.onclick = () => { resume.projects.push({name:'', tech:'', points:[]}); renderForm(); render(); };
                area.appendChild(b);

            } else if (currentTab === 'skills') {
                for (let k in resume.skills) {
                    let d = document.createElement('div');
                    d.innerHTML = '<label>'+k+'</label><textarea id="sk-'+k+'">'+resume.skills[k]+'</textarea>';
                    d.querySelector('textarea').oninput = (e) => { resume.skills[k] = e.target.value; render(); };
                    area.appendChild(d);
                }
            } else if (currentTab === 'order') {
                resume.sectionOrder.forEach((sec, i) => {
                    let d = document.createElement('div'); d.className = 'card'; d.style.display = 'flex'; d.style.justifyContent = 'space-between';
                    d.innerHTML = '<span>'+sec.toUpperCase()+'</span>' +
                                  '<div><button onclick="moveOrder('+i+', -1)">Up</button> <button onclick="moveOrder('+i+', 1)">Down</button></div>';
                    area.appendChild(d);
                });
            }
        }

        function moveOrder(i, dir) {
            let t = i + dir; if (t < 0 || t >= resume.sectionOrder.length) return;
            let temp = resume.sectionOrder[i]; resume.sectionOrder[i] = resume.sectionOrder[t]; resume.sectionOrder[t] = temp;
            renderForm(); render();
        }

        function remove(arr, i) { resume[arr].splice(i, 1); renderForm(); render(); }

        function render() {
            const paper = document.getElementById('capture-area');
            let h = '<div class="res-name">' + resume.basics.name + '</div>';
            h += '<div class="res-contact">' + resume.basics.phone + ' | ' + resume.basics.email + ' | LinkedIn | Github</div>';

            resume.sectionOrder.forEach(sec => {
                if (sec === 'education' && resume.education.length) {
                    h += '<div class="res-sec">Education</div>';
                    resume.education.forEach(e => { h += '<div class="res-row"><span>'+e.school+'</span><span>'+e.date+'</span></div><div>'+e.degree+' (CGPA: '+e.cgpa+')</div>'; });
                } else if (sec === 'skills') {
                    h += '<div class="res-sec">Technical Skills</div>';
                    for (let k in resume.skills) { h += '<div style="font-size:10.5pt"><strong>'+k+':</strong> '+resume.skills[k]+'</div>'; }
                } else if (sec === 'experience' && resume.experience.length) {
                    h += '<div class="res-sec">Experience</div>';
                    resume.experience.forEach(exp => {
                        h += '<div class="res-row"><span>'+exp.company+'</span><span>'+exp.date+'</span></div><div style="font-style:italic; font-size:10.5pt">'+exp.role+'</div>';
                        h += '<ul>' + exp.points.map(p => p ? '<li>'+p+'</li>' : '').join('') + '</ul>';
                    });
                } else if (sec === 'projects' && resume.projects.length) {
                    h += '<div class="res-sec">Projects</div>';
                    resume.projects.forEach(p => {
                        h += '<div class="res-row"><span>'+p.name+' | <span style="font-weight:normal; font-style:italic;">'+p.tech+'</span></span></div>';
                        h += '<ul>' + p.points.map(pt => pt ? '<li>'+pt+'</li>' : '').join('') + '</ul>';
                    });
                }
            });
            paper.innerHTML = h;
        }

        async function save() {
            await fetch('/save', { method: 'POST', body: JSON.stringify(resume) });
            alert("Saved!");
        }

        renderForm(); render();
    </script>
</body>
</html>
`
