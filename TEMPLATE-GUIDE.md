---
name: learn-cli-template-creation
description: Create custom templates for the learn CLI knowledge base tool
category: software-development
---

# Learn CLI — Template Creation Guide

## What This Is

The `learn` CLI is a Git-first knowledge base tool. Templates define the structure of new notes. Users can create custom templates for any category.

## Template Location

Templates live at `~/.config/learn/templates/<name>.md`

Each template file name must match the category directory name:
- `linux.md` -> creates notes in `<repo-root>/linux/`
- `math.md` -> creates notes in `<repo-root>/math/`
- `general.md` -> fallback for categories without a matching template

The repo root is whatever directory the user ran `learn init` in (e.g., `~/notes`, `~/study`, `~/knowledge`). It's stored in `~/.config/learn/config.toml`.

## Template Format

Every template MUST have YAML frontmatter at the top, followed by markdown sections.

```markdown
---
title: "{title}"
date: {date}
category: CATEGORY_NAME
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Section One

## Section Two

## References
```

## Available Placeholders

These are auto-replaced at note creation time:

| Placeholder | Replaced With | Example |
|-------------|---------------|---------|
| `{title}` | User's note title | "Binary Search" |
| `{date}` | Current date | 2026-06-08 |
| `{datetime}` | ISO 8601 timestamp | 2026-06-08T11:07:23+05:30 |
| `{category}` | Category name | "math" |
| `{tags}` | Auto-generated tag list | ["math", "notes"] |

## Rules

1. Frontmatter MUST start with `---` on the first line
2. `title` field MUST use `"{title}"` (with quotes)
3. `tags` field MUST use `{tags}` (no quotes, auto-injected)
4. `category` field should match the template filename
5. `status` defaults to `active`
6. All other frontmatter fields are optional but recommended

## Auto-Tagging

Tags are auto-generated based on category. Known categories have predefined tags:
- linux -> ["linux", "sysadmin", "cli"]
- aws -> ["aws", "cloud", "iaas"]
- docker -> ["docker", "containers", "devops"]
- kubernetes -> ["kubernetes", "k8s", "containers", "devops"]
- networking -> ["networking", "tcp", "sysadmin"]
- ctf -> ["ctf", "security", "challenge"]
- troubleshooting -> ["troubleshooting", "debugging"]
- daily -> ["daily", "journal"]
- challenge -> ["challenge", "learning"]
- general -> ["general", "notes"]

Unknown categories get `["<category>"]` as tags.

## How To Create A Custom Template

1. Create the category directory:
   ```bash
   mkdir ~/learning/<category>
   ```

2. Create the template file:
   ```bash
   # Copy general as a starting point
   cp ~/.config/learn/templates/general.md ~/.config/learn/templates/<category>.md
   ```

3. Edit the template — update frontmatter and add sections

4. Test with `learn new` — select the new category

## Template Examples

### Study / Exam Prep

```markdown
---
title: "{title}"
date: {date}
category: study
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Syllabus Topic


## Key Concepts


## Formulas / Equations


## Practice Problems


## Mistakes to Avoid


## Flashcards

Q:

A:

## References

```

### Programming

```markdown
---
title: "{title}"
date: {date}
category: programming
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Problem Statement


## Approach


## Solution

`` ` ``
// code here
`` ` ``

## Complexity

- Time:
- Space:

## Edge Cases


## Related Problems

```

### Meetings / Standups

```markdown
---
title: "{title}"
date: {date}
category: meetings
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Attendees


## Agenda


## Discussion


## Action Items

- [ ] 

## Decisions


## Follow-up

```

### Research / Papers

```markdown
---
title: "{title}"
date: {date}
category: research
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Paper Info

- **Authors:**
- **Published:**
- **Link:**

## Abstract / Summary


## Key Contributions


## Methodology


## Results


## Strengths


## Weaknesses


## How to Apply


## References

```

### Book Notes

```markdown
---
title: "{title}"
date: {date}
category: books
created_at: {datetime}
tags: {tags}
status: active
---

# {title}

## Book Info

- **Author:**
- **Pages:**
- **Genre:**

## Why I Read This


## Key Takeaways

1. 

## Chapter Notes


## Favorite Quotes

> 

## Rating

/10

## Would I Recommend?

```

## Bundled Templates (Cannot Be Overwritten)

These ship with learn init and are always available:
- linux, aws, docker, kubernetes, networking
- ctf, troubleshooting, daily, challenge, general

Custom templates are user-created and persist across `learn init` runs (init backs up existing templates to `.bak/` before overwriting defaults).

## Troubleshooting

- Template not showing up in `learn new`? Check filename matches category directory name exactly.
- Tags wrong? The `{tags}` placeholder must be exactly `{tags}` — no quotes around it.
- Frontmatter broken? Make sure `---` appears on line 1 and after the last frontmatter field.
