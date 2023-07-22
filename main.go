package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/connerdouglass/jdsl/transpiler"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	run        = kingpin.Flag("run", "run the transpiled code").Short('r').Bool()
	inputs     = kingpin.Arg("inputs", "input files to transpile").Required().Strings()
	outputDir  = kingpin.Flag("output", "output directory").Short('o').Default("dist").String()
	outputFile = kingpin.Flag("outfile", "output file").Short('f').String()
	strict     = kingpin.Flag("strict", "strict mode").Bool()
	annotate   = kingpin.Flag("annotate", "enable annotation comments").Bool()
	verbose    = kingpin.Flag("verbose", "enable verbose logging").Bool()
)

func main() {
	kingpin.Parse()

	// Create the context for the program
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	defer cancel()

	// Determine the options
	cwd, _ := os.Getwd()
	opts := transpiler.Options{
		GitRoot:  cwd,
		Inputs:   *inputs,
		Strict:   *strict,
		Annotate: *annotate,
		Verbose:  *verbose,
	}

	// If we're in run mode, output to a temporary directory
	if *run {
		tempdir, _ := os.MkdirTemp("", "jdsl-run")
		opts.CombinedOutput = filepath.Join(tempdir, "main.js")
	} else if outputFile != nil && len(*outputFile) > 0 {
		opts.CombinedOutput = *outputFile
	} else {
		opts.OutputRoot = *outputDir
	}

	// Create the transpiler
	transpiler := transpiler.NewTranspiler()
	if err := transpiler.Transpile(ctx, opts); err != nil {
		log.Fatalf("Transpile error: %s", err)
	}

	// In run mode, trigger node on the output file
	if *run {
		cmd := exec.CommandContext(ctx, "node", opts.CombinedOutput)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Run error: %s", err)
		}
	}
}
