---
name: "🐛 Bug Report"
about: Report an unexpected behavior or a bug in Tagoly.
title: "[BUG] [Concise Summary] A brief description of the issue"
labels: ["type: bug", "status: needs-triage"]
assignees: ''
---

## 🐞 The Issue

What exactly happened when this bug occurred?
(e.g., `tagoly` command crashes on a specific OS, commit scope is not correctly detected for certain files, etc.)

## 🔁 Steps to Reproduce

Please provide clear, step-by-step instructions to reproduce the reported issue.

1. Executed `git add .`.
2. Ran the `tagoly` command.
3. Selected `feat` as the commit type.
4. ...

## 💡 Expected Behavior

What did you expect to happen when following the steps above?

## 🖥️ Environment

Please provide as much detail as possible about your environment.

| Item | Detail |
| :--- | :--- |
| **Tagoly Version** | e.g., `v0.5.0` (If unsure, provide the output of `tagoly --version`) |
| **OS and Version** | e.g., `macOS Sonoma 14.2`, `Ubuntu 22.04` |
| **Git Version** | e.g., `git version 2.40.1` (Provide the output of `git --version`) |
| **Installation Method** | e.g., `Homebrew`, `Scoop`, `Manual Install (Mac M1/M2)` |

## 📝 Related Configuration & Logs

### `.tagolycustom` Content

If you are using a custom configuration file, please paste its content here.

```json
// Paste the content of your .tagolycustom here
{
  // ...
}