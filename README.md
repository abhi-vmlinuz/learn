# learn

A lightweight, Git-first engineering knowledge base CLI for DevOps, platform engineers, cloud engineers, and Linux system administrators.

Notes are plain Markdown stored in a Git repository. The binary is a convenience layer — your notes remain fully usable and discoverable without it.

## Install

### Build from source

Requires Go 1.21+.

```bash
git clone git@github.com:abhi-vmlinuz/learn.git
cd learn
make
sudo make install
```

### PDF export (optional)

Requires Python 3 and weasyprint:

```bash
pip install -r requirements.txt
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
| [weasyprint](https://weasyprint.org/) | PDF export | `pip install -r requirements.txt` |

#### Install all dependencies

**Fedora / RHEL:**
```bash
sudo dnf install fzf ripgrep bat glow
# optional: pip install -r requirements.txt
```

**Debian / Ubuntu:**
```bash
sudo apt install fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# weasyprint: pip install weasyprint
```

**Arch:**
```bash
sudo pacman -S fzf ripgrep bat glow
# optional: pip install -r requirements.txt
```

**openSUSE:**
```bash
sudo zypper install fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# optional: pip install -r requirements.txt
```

**Alpine:**
```bash
apk add fzf ripgrep bat
# glow: https://github.com/charmbracelet/glow/releases
# optional: pip install -r requirements.txt
```

**macOS (Homebrew):**
```bash
brew install fzf ripgrep bat glow
# optional: pip install -r requirements.txt
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
| `learn export` | Export a note to styled PDF |
| `learn stats` | Show learning statistics and streaks |
| `learn doctor` | Check environment health |
| `learn completion` | Shell completion scripts |

## Usage

See [CLI.md](docs/CLI.md) for detailed usage, examples, and workflows.

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

## AI Skills

Two AI skills are included in `docs/` — load them into your AI model (ChatGPT, Claude, Gemini, etc.) as system prompts, custom instructions, or skills.

| Skill | File | What It Does |
|-------|------|-------------|
| Template Creator | [docs/TEMPLATE-GUIDE.md](docs/TEMPLATE-GUIDE.md) | Generates custom template files for new categories |
| Note Generator | [docs/NOTE-GENERATOR.md](docs/NOTE-GENERATOR.md) | Asks you questions and generates a ready-to-save note |

**How to use:**
1. Open your AI model's settings (system prompt, custom instructions, or skills)
2. Paste the contents of the skill file
3. Start a conversation — the AI will follow the skill's workflow

The **Note Generator** asks one question at a time based on the category (commands used, symptoms, root cause, etc.) and outputs a complete `.md` file with frontmatter that you save directly to your learn repo.

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
