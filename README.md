# learn

A lightweight, Git-first engineering knowledge base CLI for DevOps, platform engineers, cloud engineers, and Linux system administrators.

Notes are plain Markdown stored in a Git repository. The binary is a convenience layer — your notes remain fully usable and discoverable without it.

## Install

### Build from source

Requires Go 1.21+.

```bash
git clone git@github.com:abhi-vmlinuz/learn.git
cd learn
make build
sudo make install
```

### Dependencies

#### Required

| Tool | Purpose | Install |
|------|---------|---------|
| [git](https://git-scm.com/) | Version control | Pre-installed on most systems |
| [fzf](https://github.com/junegunn/fzf) | Interactive selection | `sudo dnf install fzf` |
| [ripgrep](https://github.com/BurntSushi/ripgrep) | Full-text search | `sudo dnf install ripgrep` |
| [bat](https://github.com/sharkdp/bat) | Syntax-highlighted preview | `sudo dnf install bat` |

#### Optional

| Tool | Purpose | Install |
|------|---------|---------|
| [glow](https://github.com/charmbracelet/glow) | Markdown terminal viewer | `sudo dnf install glow` |
| [tdf](https://github.com/justjavac/tdf) | Terminal PDF viewer | `go install github.com/justjavac/tdf@latest` |
| [wkhtmltopdf](https://wkhtmltopdf.org/) | PDF export | `sudo dnf install wkhtmltopdf` |

**Fedora / RHEL:**
```bash
sudo dnf install fzf ripgrep bat glow wkhtmltopdf
```

**Debian / Ubuntu:**
```bash
sudo apt install fzf ripgrep bat glow wkhtmltopdf
```

**Arch:**
```bash
sudo pacman -S fzf ripgrep bat glow wkhtmltopdf
```

**macOS (Homebrew):**
```bash
brew install fzf ripgrep bat glow wkhtmltopdf
```

**tdf (all platforms):**
```bash
go install github.com/justjavac/tdf@latest
```

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

# Initialize
learn init

# Create your first note
learn new

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
| `learn commit [message]` | Stage and commit modified notes |
| `learn push` | Push to git remote |
| `learn pull` | Pull from git remote |
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

Bundled templates: linux, aws, docker, kubernetes, networking, ctf, troubleshooting, daily, challenge.

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

Customize templates at `~/.config/learn/templates/`. Running `learn init` again backs up existing templates before overwriting.

## License

MIT
