# JDSL

Ahead-of-time transpiler for JDSL, the world's most powerful programming language. This super-serious project converts your JDSL code into JavaScript at build-time, so you can skip the runtime overhead[^1] of JDSL.

## Inspiration

JDSL, the JSON-based Domain Specific Language, was created by Tom &mdash; a genius &mdash; to solve some of the world's hardest engineering challenges. These challenges included handling user input, and even updating a database. Things all developers struggled with before the advent of JDSL.

Read more about the JDSL language [here](https://thedailywtf.com/articles/the-inner-json-effect).

## Installation

```bash
go install github.com/connerdouglass/jdsl
```

## JDSL Syntax

A JDSL project consists of a series of JSON files. Each file contains a single JSON object, which is a JDSL class. Each class has a name and a series of functions.

Each function is listed as a Git commit hash. The hash is used to look up the contents of a corresponding JS file in the same directory. In that commit, the JS file must contain exactly one function.

To add more functions to a class, simple overwrite the JS file with the new function being added. Commit this change, and take note of the commit hash using `git log`. Then, add the hash to the JSON file.

See an example project in the `example` directory.

## Usage

#### Transpile JDSL to JavaScript

Transpiling JDSL to JavaScript is blazingly fast. Just run the following command. The JavaScript output will be written to the `dist` directory.

```bash
jdsl -o dist example/Customers.json example/Main.json
```

#### Run JDSL code directly

You can even run your JDSL code directly.

```bash
jdsl -r example/Customers.json example/Main.json
```

## Command Line Options

<table>
    <thead>
        <tr>
            <td>Option</td>
            <td>Description</td>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>-r</code>, <code>--run</code></td>
            <td>Run the JDSL code directly</td>
        </tr>
        <tr>
            <td><code>-o</code>, <code>--output</code></td>
            <td>Output directory for transpiled JavaScript</td>
        </tr>
        <tr>
            <td><code>-f</code>, <code>--outfile</code></td>
            <td>Output file for combined transpiled JavaScript</td>
        </tr>
        <tr>
            <td><code>--strict</code></td>
            <td>Enforce build-time checks of your JDSL code</td>
        </tr>
        <tr>
            <td><code>--annotate</code></td>
            <td>Include JDSL header comments in JavaScript output</td>
        </tr>
        <tr>
            <td><code>--verbose</code></td>
            <td>Print transpile logs to stdout</td>
        </tr>
    </tbody>
</table>

## Differences from the original JDSL

- This JDSL transpiler uses Git hashes instead of SVN change numbers. This is due to our technical incompetence.

## Footnotes

[^1]: The runtime overhead of JDSL is not a bug, but a feature. Thanks to Tom, rate limiting is a built-in feature of the original JDSL architecture.
