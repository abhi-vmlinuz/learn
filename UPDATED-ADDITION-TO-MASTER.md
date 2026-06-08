# Specification Updates

## Metadata & Frontmatter

All generated notes must include YAML frontmatter automatically.

Example:

---

title: lsof
date: 2026-06-06
category: linux
created_at: 2026-06-06T14:32:11+05:30
tags: []
status: active
--------------

This metadata should be automatically injected during template rendering and should not require user input.

Rationale:

* Enables future search, filtering, statistics, review, and export functionality.
* Avoids filename parsing where metadata is available.

---

## learn commit

Replace the existing commit workflow.

Current behavior:

* Creates one git commit per note.

New behavior:

learn commit

Workflow:

1. Detect all modified and untracked markdown files.
2. Display summary:

Modified:

* linux/2026-06-06-lsof.md
* linux/2026-06-06-fuser.md

3. Prompt:

Commit message:

4. Execute:

git add .
git commit -m "<user message>"

Example:

learn: linux troubleshooting notes

Rationale:

* Produces cleaner git history.
* Better reflects learning sessions.

---

## learn recent

New command:

learn recent

Workflow:

1. Discover all markdown notes.
2. Sort by modification time descending.
3. Present results via fzf.
4. Use bat preview.
5. Open selected note in $EDITOR.

Purpose:
Quickly reopen recently edited notes.

---

## learn review

New command:

learn review

Workflow:

1. Discover notes older than N days.
2. Randomly select candidates.
3. Present through fzf.
4. Preview with bat.
5. Open selected note in $EDITOR.

Default:

learn review

uses notes older than 7 days.

Optional:

learn review --days 14

Purpose:
Spaced repetition for engineering knowledge retention.

---

## learn stats

New command:

learn stats

Example output:

Total Notes: 182

Categories:
linux            51
aws              37
docker           18
kubernetes       22
networking       29
ctf              25

Current Streak: 12 days
Longest Streak: 31 days

Last Note:
2026-06-06-lsof.md

Metrics should be calculated from metadata and filesystem information.

---

## learn doctor

New command:

learn doctor

Checks:

* git available
* fzf available
* ripgrep available
* bat available
* EDITOR configured
* repository initialized
* config file exists

Example output:

✓ git
✓ fzf
✓ rg
✓ bat
✓ EDITOR
✓ repository

Repository healthy

Purpose:
Troubleshooting and environment validation.

---

## Challenge Template

Promote challenge.md to a first-class bundled template.

Template sections:

# {title}

## Scenario

## Goal

## Initial Hypothesis

## Investigation

## Commands Used

## Findings

## Root Cause

## Resolution

## Lessons Learned

## Related Concepts

## Flashcards

Q:

A:

Rationale:
Challenge-based learning is a primary workflow for DevOps, Linux, and Platform Engineering practice.

---

## Category Selection

Remove "uncategorized" from category choices.

Category selection should be mandatory.

Rationale:
Forces proper organization and improves long-term searchability.

---

## Search Enhancements

learn search

Should support:

--category <name>

and continue using:

ripgrep + fzf + bat preview

No custom indexing system should be implemented.

Prefer Unix tooling over internal search engines.

---

## Design Philosophy Additions

* Notes are operational runbooks, not passive notes.
* Learning should optimize for recall, troubleshooting, and reuse.
* Favor existing Unix tools (git, fzf, rg, bat) over custom implementations.
* Keep the codebase small and composable.
* Every note should remain readable and useful without the learn binary.
