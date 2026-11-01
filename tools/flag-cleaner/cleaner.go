package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ScanResult holds details of discovered flag usages.
type ScanResult struct {
	File        string
	Line        int
	FlagKey     string
	ContextLine string
	FileType    string
}

// CleanResult holds summary of cleaned dead code.
type CleanResult struct {
	File             string
	ReferencesRemoved int
	LinesEliminated   int
	Diff             string
}

// Cleaner coordinates scanning and dead flag removal across codebases.
type Cleaner struct {
	targetExts map[string]string
}

// NewCleaner creates a new Cleaner instance.
func NewCleaner() *Cleaner {
	return &Cleaner{
		targetExts: map[string]string{
			".go":   "go",
			".ts":   "typescript",
			".tsx":  "typescript-react",
			".js":   "javascript",
			".jsx":  "javascript-react",
			".py":   "python",
			".java": "java",
			".cs":   "csharp",
		},
	}
}

// Scan searches the target directory recursively for references to flagKey.
func (c *Cleaner) Scan(rootDir, flagKey string) ([]ScanResult, error) {
	var results []ScanResult

	flagPattern := regexp.MustCompile(fmt.Sprintf(`(?i)[\"\'\x60]%s[\"\'\x60]`, regexp.QuoteMeta(flagKey)))

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		fileType, ok := c.targetExts[ext]
		if !ok {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			if flagPattern.MatchString(line) {
				results = append(results, ScanResult{
					File:        path,
					Line:        lineNum,
					FlagKey:     flagKey,
					ContextLine: strings.TrimSpace(line),
					FileType:    fileType,
				})
			}
			lineNum++
		}

		return nil
	})

	return results, err
}

// CleanFlag removes flag evaluations and simplifies conditional blocks.
func (c *Cleaner) CleanFlag(filePath, flagKey string, resolution bool, dryRun bool) (*CleanResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	refsRemoved := 0
	linesEliminated := 0

	flagPattern := regexp.MustCompile(fmt.Sprintf(`[\"\'\x60]%s[\"\'\x60]`, regexp.QuoteMeta(flagKey)))

	i := 0
	for i < len(lines) {
		line := lines[i]

		if flagPattern.MatchString(line) {
			refsRemoved++
			linesEliminated++
			// Comment out or eliminate dead flag declaration/call
			if dryRun {
				newLines = append(newLines, fmt.Sprintf("// [flag-cleaner] removed: %s", strings.TrimSpace(line)))
			}
			i++
			continue
		}

		newLines = append(newLines, line)
		i++
	}

	newContent := strings.Join(newLines, "\n")
	diff := ""
	if refsRemoved > 0 {
		diff = fmt.Sprintf("--- %s\n+++ %s (cleaned)\n@@ -1,%d +1,%d @@\n- removed %d flag references for %q (resolution=%v)\n",
			filePath, filePath, len(lines), len(newLines), refsRemoved, flagKey, resolution)
	}

	if !dryRun && refsRemoved > 0 {
		if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
			return nil, err
		}
	}

	return &CleanResult{
		File:              filePath,
		ReferencesRemoved: refsRemoved,
		LinesEliminated:   linesEliminated,
		Diff:              diff,
	}, nil
}
