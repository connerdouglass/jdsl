package transpiler

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
	GitRoot        string
	Inputs         []string
	CombinedOutput string
	OutputRoot     string
	Strict         bool
	Annotate       bool
	Verbose        bool
}

func (o *Options) Log(str string) {
	if o.Verbose {
		fmt.Println(str)
	}
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

	// If we're in single output mode, create that output
	var combinedOutput io.WriteCloser
	if len(opts.CombinedOutput) > 0 {
		// Create the parent directory for the combined output file
		if err := os.MkdirAll(filepath.Dir(opts.CombinedOutput), 0755); err != nil {
			return fmt.Errorf("creating output directory: %s", err)
		}

		// Create the single file we're write all output to
		combinedOutput, err = os.Create(opts.CombinedOutput)
		if err != nil {
			return fmt.Errorf("creating combined output file: %s", err)
		}
		defer combinedOutput.Close()
	} else {
		// Create the output directory
		if err := os.MkdirAll(opts.OutputRoot, 0755); err != nil {
			return fmt.Errorf("creating output directory: %s", err)
		}
	}

	// Process input files
	for _, input := range opts.Inputs {
		// Context cancellation escape route for the largest JDSL compilations
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Process the input file
		err := t.processInput(ctx, repo, &opts, input, combinedOutput)
		if err != nil {
			return fmt.Errorf("processing input: %s", err)
		}
	}
	return nil
}

func (t *transpiler) processInput(ctx context.Context, repo *git.Repository, opts *Options, input string, combinedOutput io.Writer) error {
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

	// Determine if we're writing to a single combined output or not
	var output io.Writer
	if combinedOutput != nil {
		output = combinedOutput
	} else {
		// Create the output file
		outputPath := filepath.Join(opts.OutputRoot, jsPath)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("creating output directory: %s", err)
		}
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("creating output file: %s", err)
		}
		defer outputFile.Close()
		output = outputFile
	}

	// Transpile and write to the output
	return t.transpileFile(ctx, repo, opts, &manifest, jsPath, output)
}

func (t *transpiler) transpileFile(
	ctx context.Context,
	repo *git.Repository,
	opts *Options,
	manifest *Manifest,
	jsPath string,
	output io.Writer,
) error {
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

		if _, err := output.Write([]byte(strings.Join(headerAnnotations, "\n"))); err != nil {
			return fmt.Errorf("writing annotation: %s", err)
		}
	}
	if _, err := output.Write([]byte(fmt.Sprintf("function %s() {}\n", manifest.Class))); err != nil {
		return fmt.Errorf("writing class definition: %s", err)
	}

	// Read the JS file from the specified commit history for each function
	for _, function := range manifest.Functions {
		opts.Log(fmt.Sprintf("Reading %s from commit %s...", jsPath, function))

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
