package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/daten-krake/tentacle-lint/internal/linter"
	"github.com/daten-krake/tentacle-lint/internal/output"
)

var version = "dev"

type colorMode int

const (
	colorAuto colorMode = iota
	colorAlways
	colorNever
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	dir := flag.String("dir", ".", "directory containing yaml files to lint")
	recursive := flag.Bool("recursive", true, "recursively search subdirectories for yaml files")
	strict := flag.Bool("strict", false, "treat warnings as errors")
	format := flag.String("format", "text", "output format: text or json")
	colorFlag := flag.String("color", "auto", "color output: auto, always, or never")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tentacle-lint %s\n", version)
		os.Exit(0)
	}

	var colors colorMode
	switch strings.ToLower(*colorFlag) {
	case "always":
		colors = colorAlways
	case "never":
		colors = colorNever
	case "auto":
		colors = colorAuto
	default:
		fmt.Fprintf(os.Stderr, "error: invalid color %q, must be 'auto', 'always', or 'never'\n", *colorFlag)
		flag.Usage()
		os.Exit(2)
	}

	formatStr := strings.ToLower(*format)
	if formatStr != "text" && formatStr != "json" {
		errFmt(os.Stderr, colors, "invalid format %q, must be 'text' or 'json'\n", *format)
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
		errFmt(os.Stderr, colors, "%v\n", err)
		os.Exit(2)
	}

	switch formatStr {
	case "json":
		if err := output.JSON(os.Stdout, issues, *strict); err != nil {
			errFmt(os.Stderr, colors, "writing json output: %v\n", err)
			os.Exit(2)
		}
	default:
		opts := output.Options{Color: wantColor(colors, os.Stdout)}
		output.Text(os.Stdout, issues, *strict, opts)
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

func wantColor(mode colorMode, f *os.File) bool {
	switch mode {
	case colorAlways:
		return true
	case colorNever:
		return false
	default:
		fi, _ := f.Stat()
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
}

func errFmt(w io.Writer, colors colorMode, format string, args ...interface{}) {
	useColor := wantColor(colors, os.Stderr)
	prefix := "error: "
	if useColor {
		prefix = output.ColorRed + output.ColorBold + "error:" + output.ColorReset + " "
	}
	fmt.Fprint(w, prefix)
	fmt.Fprintf(w, format, args...)
}
