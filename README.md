## C.V.V.C. (Curriculum Vitae Version Control)

C.V.V.C. is a command-line interface for version controlling resumes and CVs. It wraps Git functionality using the Go library **go-git**, enabling structured version control for resume files, including binary outputs such as PDFs.

The goal is simple: eliminate multiple scattered resume versions and manage your CV with proper branching, history, and diffs—just like source code.

---

## Features

### Core Commands

* `cvvc init`
  Initializes a resume template (currently based on Jake’s Resume Template from Overleaf).

* `cvvc edit`
  Launches a local editor at `localhost:9090` to modify, add, remove, or reorder sections and entries.

* `cvvc status`
  Displays the current branch and commit status.

* `cvvc commit -m "message"`
  Commits the current state of the resume with a message, storing it in permanent history.

* `cvvc list`
  Lists all previous commits for the resume.

* `cvvc branch create <branch-name>`
  Creates a new branch from the current branch.

* `cvvc switch <branch-name>`
  Switches between branches.

* `cvvc export`
  Exports the current resume to PDF format.

* `cvvc diff`
  Shows changes made compared to the previous commit.

---

## Roadmap

* **Rebase**
  Move the current branch to a previous commit.

* **Selective Diffing**
  Compare any two chosen resume versions.

* **Multiple Templates**
  Support for additional resume templates selectable during initialization.

* **Template Transformation**
  Convert a resume from one template format to another.

---

## Contributing

Contributions are welcome.
Please open a pull request for improvements or submit an issue if you encounter any problems.

You can also reach out via email or LinkedIn.

---

## Release

C.V.V.C. is currently under active development.
Version 1.0 is planned for release in March 2026.

If you find this project useful, consider starring the repository and following for future updates.

Built by Dewashish (with love ofcc).
