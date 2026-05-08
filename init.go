package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 1. detect project name dari folder
	wd, _ := os.Getwd()
	parts := strings.Split(wd, string(os.PathSeparator))
	projectName := parts[len(parts)-1]

	// 2. ambil module lama
	data, err := os.ReadFile("go.mod")
	if err != nil {
		panic("go.mod not found")
	}

	var oldModule string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			oldModule = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			break
		}
	}

	if oldModule == "" {
		panic("module not found")
	}

	// 3. ambil username dari git
	// cmd := exec.Command("git", "config", "user.name")
	// out, _ := cmd.Output()
	// username := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(string(out)), " ", ""))
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		panic("failed to get git remote url")
	}

	remote := strings.TrimSpace(string(out))

	var username string

	if strings.HasPrefix(remote, "git@github.com:") {
		// SSH
		parts := strings.Split(strings.TrimPrefix(remote, "git@github.com:"), "/")
		username = parts[0]
	} else if strings.HasPrefix(remote, "https://github.com/") {
		// HTTPS
		parts := strings.Split(strings.TrimPrefix(remote, "https://github.com/"), "/")
		username = parts[0]
	} else {
		panic("unsupported git remote format")
	}

	if username == "" {
		panic("failed to extract github username")
	}

	if username == "" {
		panic("git username not set")
	}

	newModule := "github.com/" + username + "/" + projectName

	fmt.Println("🚀 Initializing project")
	fmt.Println("📦", oldModule, "→", newModule)

	// 4. update go.mod
	newGoMod := strings.ReplaceAll(string(data), oldModule, newModule)
	os.WriteFile("go.mod", []byte(newGoMod), 0644)

	// 5. replace import (AST)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".go" || path == "./init.go" {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		changed := false

		for _, imp := range node.Imports {
			if strings.Contains(imp.Path.Value, oldModule) {
				// imp.Path.Value = `` + strings.ReplaceAll(imp.Path.Value, oldModule, newModule) + ``
				newPath := strings.ReplaceAll(strings.Trim(imp.Path.Value, `"`), oldModule, newModule)
				imp.Path.Value = fmt.Sprintf("%q", newPath)
				changed = true
			}
		}

		if changed {
			var buf bytes.Buffer
			printer.Fprint(&buf, fset, node)
			os.WriteFile(path, buf.Bytes(), 0644)
			fmt.Println("✔ updated:", path)
		}

		return nil
	})

	// 6. copy .env
	if _, err := os.Stat(".env.example"); err == nil {
		input, _ := os.ReadFile(".env.example")
		os.WriteFile(".env", input, 0644)
		fmt.Println("🧪 .env created")
	}

	// 7. tidy deps
	fmt.Println("📦 running go mod tidy...")
	exec.Command("go", "mod", "tidy").Run()

	// 8. hapus init.go sendiri
	fmt.Println("🧹 cleaning init.go...")
	os.Remove("init.go")

	fmt.Println("")
	fmt.Println("✅ Project ready!")
	fmt.Println("👉 go run .")
}
