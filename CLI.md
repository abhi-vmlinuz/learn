# learn — CLI Usage Guide

## Getting Started

```bash
# Create a directory for your notes
mkdir ~/learning && cd ~/learning

# Initialize (creates category dirs, config, templates)
learn init

# Check everything is set up
learn doctor
```

## Creating Notes

```bash
# Create a new note (interactive: pick template, enter title, auto-category)
learn new

# Create a note without opening the editor
learn new --no-edit

# Start today's daily journal
learn today
```

`learn new` walks you through:
1. Select a template via fzf (linux, aws, docker, kubernetes, networking, ctf, troubleshooting, challenge)
2. Enter a note title
3. Category is auto-selected when it matches the template name
4. Note opens in your `$EDITOR`

Every note gets YAML frontmatter with auto-injected tags based on category:

```yaml
---
title: "lsof"
date: 2026-06-08
category: linux
created_at: 2026-06-08T11:07:23+05:30
tags: ["linux", "sysadmin", "cli"]
status: active
---
```

## Browsing and Searching

```bash
# Browse all notes via fzf (with bat preview)
learn search

# Full-text search with ripgrep
learn search "open ports"

# Search within a specific category
learn search --category linux "lsof"

# Browse recently edited notes (sorted by modification time)
learn recent

# List all notes by category (tree view)
learn list
```

`learn search` opens a two-step fzf:
1. Select a note from results
2. Choose **read** (opens in glow) or **edit** (opens in $EDITOR)

## Reviewing Notes

```bash
# Review notes older than 7 days (spaced repetition)
learn review

# Review notes older than 14 days
learn review --days 14

# Review notes older than 30 days
learn review --days 30
```

`learn review` randomly selects up to 20 candidates from notes older than the specified number of days. Opens in glow for reading.

## Saving and Syncing

```bash
# Commit all modified notes (interactive prompt for message)
learn commit

# Commit with a message directly
learn commit "linux troubleshooting notes"

# Push to remote
learn push

# Pull from remote
learn pull
```

`learn commit` only stages `.md` files — other files in the repo are ignored.

## Editing

```bash
# Browse and open a note in $EDITOR
learn edit
```

## Exporting

```bash
# Export a note to PDF (interactive selection)
learn export

# Export a specific note
learn export linux/2026-06-08-lsof.md
```

Requires `wkhtmltopdf`. PDF uses a white theme with blue accents, styled code blocks, tables, and tag pills.

```bash
sudo dnf install wkhtmltopdf    # Fedora
sudo apt install wkhtmltopdf    # Debian/Ubuntu
```

## Statistics

```bash
learn stats
```

Output:
```
Total Notes: 182

Categories:
  aws                  37
  ctf                  25
  docker               18
  kubernetes           22
  linux                51
  networking           29
  troubleshooting      25

Current Streak: 12 days
Longest Streak: 31 days

Last Note:
  2026-06-08-lsof.md
```

## Environment Check

```bash
learn doctor
```

Checks: git, fzf, rg, bat, glow, wkhtmltopdf, EDITOR, config file, repository.

## Configuration

Config lives at `~/.config/learn/config.toml`:

```toml
[repo]
  root = "/home/user/learning"
```

Templates live at `~/.config/learn/templates/`. Edit them to customize note structure. Running `learn init` again backs up existing templates to `templates.bak/` before copying defaults.

## File Structure

```
~/learning/
├── aws/
├── challenge/
├── ctf/
├── daily/
│   └── 2026-06-08.md
├── docker/
├── kubernetes/
├── linux/
│   ├── 2026-06-08-lsof.md
│   └── 2026-06-08-lsof.pdf
├── networking/
└── troubleshooting/
```

Notes are plain Markdown. The binary is a convenience layer — your notes remain fully usable without it.

## Dependencies

| Tool          | Required | Used by                    |
|---------------|----------|----------------------------|
| git           | yes      | init, commit, push, pull   |
| fzf           | yes      | new, search, edit, recent, review |
| ripgrep (rg)  | yes      | search (full-text)         |
| bat           | yes      | fzf preview pane           |
| glow          | no       | read mode (search, review, recent) |
| wkhtmltopdf   | no       | export to PDF              |

## Shell Completion

```bash
learn completion install    # auto-detects and installs for your shell
learn completion bash       # output bash completion script
learn completion zsh        # output zsh completion script
learn completion fish       # output fish completion script
```

## All Commands

| Command     | Description                              |
|-------------|------------------------------------------|
| `init`      | Initialize repository structure          |
| `new`       | Create a new note from template          |
| `today`     | Create/open daily journal entry          |
| `edit`      | Open a note in $EDITOR                   |
| `search`    | Full-text search or browse all notes     |
| `recent`    | Browse recently edited notes             |
| `review`    | Spaced repetition review                 |
| `list`      | List all notes by category               |
| `commit`    | Git add + commit modified notes          |
| `push`      | Push to git remote                       |
| `pull`      | Pull from git remote                     |
| `export`    | Export note to PDF                       |
| `stats`     | Show learning statistics                 |
| `doctor`    | Check environment health                 |
| `completion`| Shell completion scripts                 |
