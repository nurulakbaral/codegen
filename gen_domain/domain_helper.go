package domain

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func CreatePath(p ...string) (path string, err error) {
	rootPath, err := os.Getwd()

	if err != nil {
		return "", fmt.Errorf("Error os.Getwd(): %v", err)
	}

	paths := append([]string{rootPath}, p...)

	return filepath.Join(paths...), nil
}

func CheckDir(c Config) (path string, err error) {
	domainPath, err := CreatePath(c.Root)
	completePath := filepath.Join(domainPath, c.Domain)

	if err != nil {
		return completePath, fmt.Errorf("Error CreatePath(): %v", err)
	}

	dirs, err := os.ReadDir(domainPath)
	for _, dir := range dirs {
		if dir.Name() == c.Domain {
			return completePath, fmt.Errorf("Error os.ReadDir(): %v", err)
		}
	}

	return completePath, nil
}

const (
	PermissionMkdirAll = 0o755
)

func GenerateTemplateFiles(data any, input, output string) error {
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, PermissionMkdirAll); err != nil {
		return fmt.Errorf("Error os.MkdirAll(): %v", err)
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Error os.Create(): %v", err)
	}
	defer file.Close()

	tmplBytes, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("Error os.ReadFile(): %v", err)
	}

	tmpl, err := template.New(filepath.Base(input)).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("Error template.New().Parse(): %v", err)
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("Error template.Execute(): %v", err)
	}

	return nil
}

func CreatePairFilePath(name, input, output string) ([][]string, error) {
	var pairFiles [][]string

	si, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("Error os.Stat(): %v", err)
	}

	if !si.IsDir() {
		return nil, fmt.Errorf("Error os.Stat(): %v", err)
	}

	rootBase := filepath.Base(input)

	err = filepath.WalkDir(input, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("Error filepath.WalkDir(): %v", walkErr)
		}

		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(input, path)
		if err != nil {
			return fmt.Errorf("Error filepath.Rel(): %v", err)
		}

		relDir := filepath.Dir(rel)
		base := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))

		var outFileName = ""
		if filepath.Ext(rel) == ".tmpl" {
			outFileName = strings.ReplaceAll(base, "domain", name) + ".go"
		} else {
			outFileName = base + filepath.Ext(rel)
		}

		outPath := filepath.Join(output, rootBase, relDir, outFileName)
		pairFiles = append(pairFiles, []string{path, outPath})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pairFiles, nil
}

func WalkDirs(currPath string) map[string][]string {
	var dirs []string
	var files []string

	filepath.WalkDir(currPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error filepath.WalkDir(): %v", err)
		}

		if currPath == path {
			return nil
		}

		if d.IsDir() {
			dirs = append(dirs, filepath.Base(path))
		} else {
			files = append(files, filepath.Base(path))
		}

		return nil
	})

	return map[string][]string{"dirs": dirs, "files": files}
}
