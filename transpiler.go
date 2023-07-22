package jdsl

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Options struct {
	GitRoot    string
	Inputs     []string
	OutputRoot string
	Strict     bool
	Annotate   bool
}

type Transpiler interface {
	Transpile(ctx context.Context, opts Options) error
}

func NewTranspiler() Transpiler {
	return &transpiler{}
}

type transpiler struct{}

func (t *transpiler) Transpile(ctx context.Context, opts Options) error {
	// Open the git repository
	repo, err := git.PlainOpen(opts.GitRoot)
	if err != nil {
		return fmt.Errorf("opening repo: %s", err)
	}

	// Create the output directory
	if err := os.MkdirAll(opts.OutputRoot, 0755); err != nil {
		return fmt.Errorf("creating output directory: %s", err)
	}

	// Process input files
	for _, input := range opts.Inputs {
		err := t.processInput(ctx, repo, &opts, input)
		if err != nil {
			return fmt.Errorf("processing input: %s", err)
		}
	}
	return nil
}

func (t *transpiler) processInput(ctx context.Context, repo *git.Repository, opts *Options, input string) error {
	// Trim the extension from the input filename
	input = strings.TrimSuffix(input, filepath.Ext(input))
	manifestPath := input + ".json"
	jsPath := input + ".js"

	// Read the manifest file
	var manifest Manifest
	if err := manifest.ReadFile(manifestPath); err != nil {
		return fmt.Errorf("reading manifest: %s", err)
	}

	// In strict mode, enforce the file name in the manifest
	if opts.Strict && manifest.File != filepath.Base(manifestPath) {
		return fmt.Errorf("file name in manifest (%s) does not match: %s", manifest.File, filepath.Base(manifestPath))
	}

	// Create the output file
	outputPath := filepath.Join(opts.OutputRoot, jsPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %s", err)
	}
	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %s", err)
	}
	defer output.Close()

	// Start with the standard JDSL class definition.
	// ES3 of course.
	if opts.Annotate {
		headerAnnotations := []string{
			"// JDSL 1.0",
			"// File: " + manifest.File,
			"// Class: " + manifest.Class,
			"// Author: " + manifest.Author,
			"// Purpose: " + manifest.Purpose,
			"",
		}
		if _, err := output.WriteString(strings.Join(headerAnnotations, "\n")); err != nil {
			return fmt.Errorf("writing annotation: %s", err)
		}
	}
	if _, err := output.WriteString(fmt.Sprintf("function %s() {}\n", manifest.Class)); err != nil {
		return fmt.Errorf("writing class definition: %s", err)
	}

	// Read the JS file from the specified commit history for each function
	for _, function := range manifest.Functions {
		fmt.Printf("Reading %s from commit %s...\n", jsPath, function)

		commit, err := repo.CommitObject(plumbing.NewHash(function))
		if err != nil {
			return fmt.Errorf("getting commit %s: %s", function, err)
		}

		// Read the file from the specicified commit hash
		commitFile, err := commit.File(jsPath)
		if err != nil {
			return fmt.Errorf("getting file %s: %s", jsPath, err)
		}
		reader, err := commitFile.Reader()
		if err != nil {
			return fmt.Errorf("getting reader for %s: %s", jsPath, err)
		}

		// Copy the contents of the file to the output
		if _, err := io.Copy(output, reader); err != nil {
			reader.Close()
			return fmt.Errorf("copying file contents: %s", err)
		}
		reader.Close()
	}
	return nil
}
