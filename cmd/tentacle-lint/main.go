package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/linter"
	"github.com/daten-krake/tentacle-lint/internal/output"
)

var version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	dir := flag.String("dir", ".", "directory containing yaml files to lint")
	recursive := flag.Bool("recursive", true, "recursively search subdirectories for yaml files")
	strict := flag.Bool("strict", false, "treat warnings as errors")
	format := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tentacle-lint %s\n", version)
		os.Exit(0)
	}

	if *format != "text" && *format != "json" {
		fmt.Fprintf(os.Stderr, "error: invalid format %q, must be 'text' or 'json'\n", *format)
		flag.Usage()
		os.Exit(2)
	}

	cfg := linter.Config{
		Dir:       *dir,
		Recursive: *recursive,
		Strict:    *strict,
	}

	issues, err := linter.Run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	switch strings.ToLower(*format) {
	case "json":
		output.JSON(os.Stdout, issues, *strict)
	default:
		output.Text(os.Stdout, issues, *strict)
	}

	hasErrors := false
	for _, issue := range issues {
		if issue.EffectiveSev(*strict) == linter.Error {
			hasErrors = true
			break
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}
