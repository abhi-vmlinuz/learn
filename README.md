# learn

A lightweight, Git-first engineering knowledge base CLI for DevOps, platform engineers, cloud engineers, and Linux system administrators.

Notes are plain Markdown stored in a Git repository. The binary is a convenience layer — your notes remain fully usable and discoverable without it.

## Install

### Build from source

Requires Go 1.21+.

```bash
git clone https://github.com/abhi-vmlinuz/learn.git
cd learn
make build
sudo make install
```

### Dependencies

#### Required

| Tool | Purpose | Install |
|------|---------|---------|
| [git](https://git-scm.com/) | Version control | Pre-installed on most systems |
| [fzf](https://github.com/junegunn/fzf) | Interactive selection | See below |
| [ripgrep](https://github.com/BurntSushi/ripgrep) | Full-text search | See below |
| [bat](https://github.com/sharkdp/bat) | Syntax-highlighted preview | See below |
| [glow](https://github.com/charmbracelet/glow) | Markdown terminal viewer | See below |

#### Optional

| Tool | Purpose | Install |
|------|---------|---------|
| [tdf](https://github.com/justjavac/tdf) | Terminal PDF viewer | `go install github.com/justjavac/tdf@latest` |
| [wkhtmltopdf](https://wkhtmltopdf.org/) | PDF export | See below |

#### Install all dependencies

**Fedora / RHEL:**
```bash
sudo dnf install fzf ripgrep bat glow
# optional: sudo dnf install wkhtmltopdf
```

**Debian / Ubuntu:**
```bash
sudo apt install fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# wkhtmltopdf: sudo apt install wkhtmltopdf
```

**Arch:**
```bash
sudo pacman -S fzf ripgrep bat glow
# optional: sudo pacman -S wkhtmltopdf
```

**openSUSE:**
```bash
sudo zypper install fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# optional: sudo zypper install wkhtmltopdf
```

**Alpine:**
```bash
apk add fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# optional: apk add wkhtmltopdf
```

**macOS (Homebrew):**
```bash
brew install fzf ripgrep bat glow
# optional: brew install wkhtmltopdf
```

**tdf (all platforms, requires Go):**
```bash
go install github.com/justjavac/tdf@latest
```

> **Note:** `bat` may be installed as `batcat` on Debian/Ubuntu. If so, create a symlink:
> ```bash
> sudo ln -s /usr/bin/batcat /usr/local/bin/bat
> ```

After installing, verify everything is set up:

```bash
learn doctor
```

### Shell completion

```bash
learn completion install
```

Auto-detects your shell (bash, zsh, fish) and installs completions to the standard location.

## Quick start

```bash
# Create a directory for your notes
mkdir ~/learning && cd ~/learning

# Initialize (creates default categories)
learn init

# Create your first note
learn new

# Add a custom category — just make a directory
mkdir ~/learning/math

# It shows up automatically in learn new, learn list, learn search

# Start a daily journal
learn today

# Search or browse notes
learn search

# Review old notes (spaced repetition)
learn review

# Commit and push
learn commit
learn push
```

## Commands

| Command | Description |
|---------|-------------|
| `learn init` | Initialize repository structure, config, and templates |
| `learn new` | Create a new note from a template |
| `learn today` | Create or open today's daily journal |
| `learn edit` | Open a note in your editor |
| `learn search [query]` | Full-text search or browse all notes |
| `learn recent` | Browse recently edited notes |
| `learn review` | Spaced repetition review of old notes |
| `learn list` | List all notes by category |
| `learn delete [file]` | Delete a note (with confirmation) |
| `learn move [file]` | Move a note to a different category |
| `learn tag [file]` | Edit tags on an existing note |
| `learn commit [message]` | Stage and commit modified notes |
| `learn push` | Push to git remote |
| `learn pull` | Pull from git remote |
| `learn log` | Show git history of notes |
| `learn export [file]` | Export a note to styled PDF |
| `learn stats` | Show learning statistics and streaks |
| `learn doctor` | Check environment health |
| `learn completion` | Shell completion scripts |

## Usage

See [CLI.md](CLI.md) for detailed usage, examples, and workflows.

## Philosophy

- **Plain Markdown.** Notes are `.md` files. Read them with any editor, viewer, or on GitHub.
- **Git-first.** Version control, sync, and history come free. No proprietary storage.
- **Unix philosophy.** Uses fzf, ripgrep, bat, glow — proven tools, not reinvented wheels.
- **Fast.** No startup lag, no Electron, no runtime dependencies beyond the tools above.
- **Operational runbooks.** Notes are structured for recall, troubleshooting, and reuse — not passive reading.

## Templates

Bundled templates: linux, aws, docker, kubernetes, networking, ctf, troubleshooting, daily, challenge, general.

All notes get YAML frontmatter with auto-generated tags:

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

**Generating templates with AI:** See [TEMPLATE-GUIDE.md](TEMPLATE-GUIDE.md) — a structured guide for AI models (ChatGPT, Claude, Gemini, etc.) to generate custom templates that follow the correct format with frontmatter and placeholders. Feed it to your AI of choice and ask it to create a template for your topic.

## Categories

Default categories are created by `learn init`: aws, linux, docker, kubernetes, networking, ctf, troubleshooting, daily, challenge, general.

**Adding custom categories** — just create a directory:

```bash
mkdir ~/learning/math
mkdir ~/learning/physics
mkdir ~/learning/exam-prep
```

New categories are discovered automatically. They appear in `learn new`, `learn list`, `learn search --category`, and everywhere else. No config changes needed.

The **general** template works for any topic — use it for study notes, exam prep, or anything that doesn't fit a specific category.

**Custom templates:** If you want a category-specific template instead of falling back to general:

```bash
# Copy general as a starting point
cp ~/.config/learn/templates/general.md ~/.config/learn/templates/math.md
# Edit it to add math-specific sections
```

Customize templates at `~/.config/learn/templates/`. Running `learn init` again backs up existing templates before overwriting.

## License

MIT
