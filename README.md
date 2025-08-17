# go-set-mod 🚀  

A simple yet powerful CLI tool to automate Go module initialization based on project structure.  
It ensures correct module naming using a predefined **GitHub namespace** and **relative paths** from the project root.  

---

## 📌 Features  
✅ **Automates `go mod init`** – No need to manually set module names.  
✅ **Project Root Detection** – Finds `mod-name.yaml` to determine the root.  
✅ **Relative Path-Based Module Naming** – Ensures correct module hierarchy.  
✅ **Cross-Platform Support** – Works on macOS, Linux, and Windows.  
✅ **One-Line Installation** – Simple installation using `curl`.  

---

## 🔧 Installation  

You can install `go-set-mod` easily using a one-liner:  

```sh
curl -fsSL https://raw.githubusercontent.com/ananay-nag/go-create-init-module/refs/heads/main/install.sh | bash
```
---
## 📂Project Structure
```
my-project/
│── go.mod
│── go.sum
│── mod-name.yaml  <-- Project root
│── my-module/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── my-sub-module/
|   |   |── mod-name.yaml <-- if you create nested mod-name.yml, will look for this
│   │   ├── go.mod
│   │   ├── file1.go
│   │   ├── my-sub-sub-module/  <-- Running `go-set-mod my-sub-sub-module` here
```
---

## Mod-name.yaml Example

```yaml
    pre-set: "github.com/your-username"
```
---
## 🚀 CLI Usage

## 1️⃣ Initialize a Go module inside a new subdirectory

- Run inside any subdirectory:
```sh
    go-set-mod <module-name>
```
- Example: Running inside my-sub-module/:
```sh
    go-set-mod my-sub-sub-module
```

## 2️⃣ Initialize a Go module in the current directory (-c)
```sh
    go-set-mod -c
```
- Create a default mod-name.yaml if not exist, You need to update pre-set:
- Uses the current directory name for go mod init.
- Does not create a new subdirectory.
```sh
    cd my-sub-module
    go-set-mod -c
```

- Generates:
```sh
    go mod init github.com/your-username/my-module/my-sub-module/my-sub-sub-module
```
---

## 📌 Features Summary
| Feature                                             | Command                     | Behavior                                                                      |
|-----------------------------------------------------|-----------------------------|-------------------------------------------------------------------------------|
| Initialize a module in a new subdirectory           | `go-set-mod my-module`      | Creates a new folder and runs `go mod init` inside it.                        |
| Initialize the current directory as a module        | `go-set-mod -c`             | Runs `go mod init` in current directory, no new folder created.               |
| Customizable module path prefix                     | `mod-name.yaml`             | Uses pre-set value to generate module paths.                                  |
| Create a default mod-name.yaml if not exist         | `mod-name.yaml`             | Create a default mod-name.yaml if not exist                                   |

## 🛠 How It Works
- 1️⃣ Detects Project Root – Searches for config.yaml in parent directories.
- 2️⃣ Computes Relative Path – Finds the path from root to the current directory.
- 3️⃣ Generates Correct Module Name – Uses github.com/your-username/<relative-path>/my-new-module.
- 4️⃣ Runs go mod init – Automatically initializes the Go module.
---
## 📜 License
- This project is licensed under the MIT License. See the LICENSE file for details.