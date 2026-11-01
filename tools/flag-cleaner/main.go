package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	cleaner := NewCleaner()

	switch command {
	case "scan":
		scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
		pathFlag := scanCmd.String("path", ".", "Target repository directory to scan")
		flagKey := scanCmd.String("flag", "", "Feature flag key to search for (required)")
		_ = scanCmd.Parse(os.Args[2:])

		if *flagKey == "" {
			fmt.Println("Error: --flag <flag-key> is required")
			os.Exit(1)
		}

		fmt.Printf("🔍 Scanning %q for flag %q...\n", *pathFlag, *flagKey)
		results, err := cleaner.Scan(*pathFlag, *flagKey)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Printf("No references found for flag %q.\n", *flagKey)
			return
		}

		fmt.Printf("Found %d references to %q:\n", len(results), *flagKey)
		for _, r := range results {
			fmt.Printf("  • %s:%d [%s]\n      %s\n", r.File, r.Line, r.FileType, r.ContextLine)
		}

	case "clean":
		cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
		pathFlag := cleanCmd.String("path", ".", "Target repository directory to clean")
		flagKey := cleanCmd.String("flag", "", "Feature flag key to eliminate (required)")
		resolution := cleanCmd.Bool("resolution", true, "Permanent resolution value: true or false")
		write := cleanCmd.Bool("write", false, "Write changes to disk (default false: dry-run mode)")
		_ = cleanCmd.Parse(os.Args[2:])

		if *flagKey == "" {
			fmt.Println("Error: --flag <flag-key> is required")
			os.Exit(1)
		}

		dryRun := !*write
		if dryRun {
			fmt.Println("⚠️  DRY RUN MODE: No files will be modified on disk. Use --write to apply changes.")
		}

		results, err := cleaner.Scan(*pathFlag, *flagKey)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			os.Exit(1)
		}

		filesMap := make(map[string]bool)
		for _, r := range results {
			filesMap[r.File] = true
		}

		totalRemoved := 0
		for file := range filesMap {
			res, err := cleaner.CleanFlag(file, *flagKey, *resolution, dryRun)
			if err != nil {
				fmt.Printf("Failed to clean %s: %v\n", file, err)
				continue
			}
			if res.ReferencesRemoved > 0 {
				totalRemoved += res.ReferencesRemoved
				fmt.Print(res.Diff)
			}
		}

		fmt.Printf("\n✨ Cleaned %d references across %d files for flag %q.\n", totalRemoved, len(filesMap), *flagKey)

	case "help", "--help", "-h":
		printUsage()

	default:
		fmt.Printf("Unknown command %q\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`FlagManagment Dead Flag Cleaner (flag-cleaner)

Usage:
  flag-cleaner <command> [options]

Commands:
  scan    Search codebase for all usages of a flag key
  clean   Safely remove flag checks and simplify dead branches

Examples:
  flag-cleaner scan --path ./src --flag checkout-v2
  flag-cleaner clean --path ./src --flag checkout-v2 --resolution true
  flag-cleaner clean --path ./src --flag checkout-v2 --resolution true --write`)
}
