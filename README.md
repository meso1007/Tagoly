# Tagoly

Tagoly is a smart CLI tool to assist with Git commits.  
It streamlines commit message creation with features like scope detection, custom tags, and interactive selection.

## 🚀 Key Features

- **Automatic Scope Detection**  
  Detects the scope based on changed file paths automatically.

- **Custom Tag Support**  
  Define your own tags in `.tagolycustom` (e.g., ci, perf).

- **Interactive Commit Generation**  
  Step-by-step selection for commit type, scope, and message.

- **Smart Scope Selection**  
  Automatically selects the most frequent scope if multiple scopes exist.  
  Allows manual selection if needed.

## Installation

### **MacOS**

#### 1. Homebrew
```bash
# Add Tap
brew tap meso1007/tagoly

# Install
brew install meso1007/tagoly/tagoly

# Verify
tagoly --version
tagoly


#### 2. Manual Installation
##### Apple Silicon (M1/M2)
```bash
mv tagoly-darwin-arm64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
##### Intel
```bash
mv tagoly-darwin-amd64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
--------

### **Linux**
#### 1. Homebrew (Linuxbrew)
```bash
brew tap meso1007/tagoly
brew install meso1007/tagoly/tagoly
```
#### 2. Manual Installation
```bash
mv tagoly-linux-amd64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
--------

### **Windows**
#### 1. Scoop
```powershell
# Add bucket
scoop bucket add tagoly https://github.com/meso1007/scoop-tagoly

# Install
scoop install tagoly/tagoly

# Verify
tagoly --version
tagoly

```
#### 2. Manual Installation
```powershell
Move-Item .\tagoly-windows-amd64.exe "C:\Program Files\tagoly\tagoly.exe"
```

--------

## Configuration File .tagolycustom
```json
{
  "customTags": [
    {"key": "ci", "label": "CI/CD changes"},
    {"key": "perf", "label": "Performance improvement"},
    {"key": "test", "label": "Add or update tests"},
    {"key": "hotfix", "label": "Hotfix / urgent fix"}
  ]
}

```

## Usage
```bash
git add .
tagoly
```

## リポジトリ

- Tagoly main: [https://github.com/meso1007/tagoly](https://github.com/meso1007/tagoly)  
- Homebrew Tap: [https://github.com/meso1007/homebrew-tagoly](https://github.com/meso1007/homebrew-tagoly)  
- Scoop bucket: [https://github.com/meso1007/scoop-tagoly](https://github.com/meso1007/scoop-tagoly)

---
