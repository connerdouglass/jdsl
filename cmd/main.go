package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/connerdouglass/jdsl"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputs    = kingpin.Arg("inputs", "input files to transpile").Required().Strings()
	outputDir = kingpin.Flag("output", "output directory").Short('o').Default(".").String()
	strict    = kingpin.Flag("strict", "strict mode").Bool()
	annotate  = kingpin.Flag("annotate", "enable annotation comments").Bool()
)

func main() {
	kingpin.Parse()

	// Create the context for the program
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	defer cancel()

	// Determine the options
	cwd, _ := os.Getwd()
	opts := jdsl.Options{
		GitRoot:    cwd,
		Inputs:     *inputs,
		OutputRoot: *outputDir,
		Strict:     *strict,
		Annotate:   *annotate,
	}

	// Create the transpiler
	transpiler := jdsl.NewTranspiler()
	if err := transpiler.Transpile(ctx, opts); err != nil {
		log.Fatalf("Transpile error: %s", err)
	}
}
