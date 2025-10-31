# Contributing to Tagoly

Thank you for considering contributing to Tagoly! We welcome bug reports, feature suggestions, and code contributions. By following these guidelines, you help us maintain the quality and focus of the project.

## 🤝 How to Contribute

We primarily use GitHub for tracking issues and managing code changes.

1.  **Report Bugs:** Use the **Bug Report** template for any unexpected behavior or crashes.
2.  **Suggest Features:** Use the **Feature Request** template to propose new functionality.
3.  **Propose Code Changes:** Submit a Pull Request (PR).

## 💡 Local Development Setup

Tagoly is a CLI tool written in Go.

### 1. Prerequisites

You must have the following installed:
* [Go (version 1.18 or higher)](https://golang.org/dl/)
* [Git](https://git-scm.com/downloads)

### 2. Fork and Clone

```bash
# Fork the meso1007/tagoly repository on GitHub first
git clone [https://github.com/meso1007/tagoly.git](https://github.com/meso1007/tagoly.git)
cd tagoly
```

### 3. Build and Test
You can build and run Tagoly locally for testing purposes:
```bash
# Build the binary in the current directory
go build -o tagoly ./src/main.go

# Run the local binary 
./tagoly
```

## 📝 Committing and Pull Request Guidelines
We enforce the Conventional Commits specification for all commit messages. This helps us generate consistent change logs and determine version releases.

### 1. Commit Message Format
Your commit message must follow this structure:
```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

**Common `<type>` examples:**
| Type | Description |
| :--- | :--- |
| `feat` | A new feature or enhancement. |
| `fix` | A bug fix. |
| `docs` | Documentation only changes (e.g., README updates). |
| `style` | Formatting changes (code refactoring that does not change logic). |
| `refactor` | A code change that neither fixes a bug nor adds a feature. |
| `test` | Adding missing tests or correcting existing tests. |
| `ci` | Changes to CI configuration files and scripts (e.g., GitHub Actions). |

**Common `<scope>` examples:**
Use the part of the codebase affected by the change (e.g., `cli`, `scope-detection`, `config`).

**Example:**

Example:
```
feat(cli): add support for multiple scopes selection

The user can now select several inferred scopes, separated by commas, 
to build a comprehensive commit message.
```

2### 2. Opening a Pull Request (PR)

1.  Create your feature branch: `git checkout -b feature/my-new-feature`
2.  Commit your changes using the Conventional Commit format.
3.  Push your branch: `git push origin feature/my-new-feature`
4.  Open a Pull Request on GitHub.
5.  **Crucially:** Fill out the `PULL_REQUEST_TEMPLATE.md` completely, linking to any relevant issues (e.g., `Closes #123`).

## 👨‍💻 Maintainer Notes (For Core Developers)

* **Triage:** All new Issues are labeled `status: needs-triage`. Triage involves adding a priority label, scope label, and assigning the Issue.
* **Releases:** New version releases must be tagged, and the compiled binaries for all target platforms (macOS, Linux, Windows) must be uploaded to GitHub Releases. The Homebrew Tap and Scoop Bucket must be updated simultaneously.