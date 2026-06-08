LEARN CLI — Master Build Prompt
Project Overview
Learn is a lightweight, Git-first engineering knowledge base CLI tool written in Go using Cobra. It's designed for DevOps, platform engineers, cloud engineers, and Linux system administrators who need fast, searchable access to troubleshooting notes, operational runbooks, and learning documentation—without vendor lock-in or proprietary formats.
Core Philosophy: Notes are plain Markdown stored in a Git repository. The binary is a convenience layer; notes remain usable and discoverable without it. Optimize for speed, simplicity, and long-term knowledge retention.

Design Specifications
Storage Model:

All notes stored as Markdown (.md) files in a user-chosen Git repository.
Templates stored in ~/.config/learn/templates/ (user-discoverable, bundled defaults embedded at build).
Configuration stored in ~/.config/learn/config.toml (tracks repo root).
Directory structure initialized once, managed by learn init:

  learning/
  ├── aws/
  ├── linux/
  ├── docker/
  ├── kubernetes/
  ├── networking/
  ├── ctf/
  ├── troubleshooting/
  └── daily/
File Naming Convention:

Standard notes: YYYY-MM-DD-title.md
Daily journal: daily/YYYY-MM-DD.md

Template System:

Simple placeholder substitution: {title}, {date}, {category}
Users can add custom templates to ~/.config/learn/templates/
learn new discovers and lists all templates via fzf.

Git Integration:

learn init detects .git. If missing, warns user with manual setup instructions (no auto-remote).
learn commit handles: git add <file> && git commit -m "Add: <title>" using filename-derived title.
Auto-commit support deferred to v2 (manual learn commit for now).


Command Reference
1. learn init — Initialize Repository
Workflow:

Verify CWD is writable.
Check for .git:

If present: print "Git repo detected. Continuing..."
If absent: warn user, print instructions: git init && git remote add origin <url>, continue.


Create directory structure (aws/, linux/, docker/, kubernetes/, networking/, ctf/, troubleshooting/, daily/).
Create ~/.config/learn/config.toml:

toml   [repo]
   root = "/absolute/path/to/cwd"

Create ~/.config/learn/templates/ if missing.
Copy bundled template defaults to ~/.config/learn/templates/ (linux.md, aws.md, docker.md, kubernetes.md, networking.md, ctf.md, troubleshooting.md, daily.md).
Print success message with next steps (e.g., "Run learn new to create your first note").

2. learn new — Create Note
Workflow:

Load repo root from config. Fail if not initialized.
List all templates in ~/.config/learn/templates/ via fzf selection.
Prompt: "Note title?"
Prompt: "Category?" with fzf listing all subdirs in repo root (aws/, linux/, etc.) + "uncategorized" option.
Render template:

Substitute {title} → user input.
Substitute {date} → current date (YYYY-MM-DD).
Substitute {category} → selected category.


Write to $LEARN_ROOT/<category>/YYYY-MM-DD-title.md.
Open file in $EDITOR (or skip if --no-edit flag).
Print file path on success.

3. learn today — Daily Journal Entry
Workflow:

Load repo root from config. Fail if not initialized.
Use daily journal template (from ~/.config/learn/templates/daily.md).
Render template with {date} → current date.
Write to $LEARN_ROOT/daily/YYYY-MM-DD.md (fail if file exists; suggest editing existing).
Open in $EDITOR.
Print file path on success.

4. learn search — Full-Text Search
Workflow:

Load repo root from config. Fail if not initialized.
Run: rg --files-with-matches '<query>' $LEARN_ROOT (if no query provided, prompt for it).
Pipe results to fzf with preview:

fzf --preview 'bat --color=always {}'


User selects a result.
Open selected file in $EDITOR.
Print file path on success.

Optional flags:

--category <cat> — filter search to specific category dir.
--limit <n> — limit result count.

5. learn commit — Git Add + Commit
Workflow:

Load repo root from config. Fail if not initialized.
If arg provided (filepath):

git add <filepath>
Extract title from filename (YYYY-MM-DD-title.md → title).
git commit -m "Add: <title>"


If no arg:

Discover all modified/untracked .md files in repo.
For each: git add <file> && git commit -m "Add: <title-from-filename>".


Print summary of committed files.


Code Structure
learn/
├── cmd/
│   ├── root.go              (Cobra root, load config, global flags)
│   ├── init.go              (learn init)
│   ├── new.go               (learn new)
│   ├── today.go             (learn today)
│   ├── search.go            (learn search)
│   └── commit.go            (learn commit)
├── internal/
│   ├── config/
│   │   └── config.go        (load/save ~/.config/learn/config.toml)
│   ├── template/
│   │   └── render.go        (placeholder substitution)
│   ├── git/
│   │   └── git.go           (detect .git, commit, status checks)
│   ├── fzf/
│   │   └── select.go        (fzf wrapper for templates, categories, search results)
│   └── file/
│       └── file.go          (create dirs, write files, parse filenames)
├── templates/               (bundled defaults, embed at build)
│   ├── linux.md
│   ├── aws.md
│   ├── docker.md
│   ├── kubernetes.md
│   ├── networking.md
│   ├── ctf.md
│   ├── troubleshooting.md
│   └── daily.md
├── main.go
└── go.mod

Implementation Notes
Dependencies:

github.com/spf13/cobra — CLI framework.
github.com/BurntSushi/toml — config parsing.
(No external fzf/bat/ripgrep deps; shell out to them via os/exec).

Embedded Templates:

Use Go's //go:embed to bundle templates at build time.
Copy to ~/.config/learn/templates/ on learn init.
User can override by editing those files.

Error Handling:

Fail early with clear messages if config missing, repo not initialized, .md not found.
Print actionable next steps (e.g., "Run learn init to set up your repo").

Design Principles:

Plain Markdown. Git-first. Unix philosophy. Fast startup. Minimal deps. No GUI. No vendor lock-in.
Notes remain usable without the binary.
Optimized for DevOps, Platform Engineering, Cloud Engineering, troubleshooting knowledge retention.


Related Prompts

Scaffold Project & Dependencies — Initialize Go module, set up Cobra structure, define config schema, write stubs for all commands.
Implement Config System — Load/save ~/.config/learn/config.toml, create config dir if missing, validate repo root.
Implement Template System — Embed bundled templates, copy to user config on init, discover user templates, render placeholders.
Implement learn init — Detect .git, create directories, initialize config, print next steps.
Implement learn new — fzf template/category selection, render template, write file, open in editor.
Implement learn today — Load daily template, render, write to daily/, open in editor.
Implement learn search — ripgrep + fzf + bat preview integration.
Implement learn commit — Parse filenames, run git add/commit with auto-message.
Testing & Validation — Unit tests for config, template rendering, filename parsing. Integration tests for commands.
Polish & Optimization — Add global flags (--verbose, --no-edit, etc.), improve error messages, optimize startup time.
