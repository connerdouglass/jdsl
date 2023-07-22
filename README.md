# JDSL

Ahead-of-time transpiler for JDSL, the world's most powerful programming language. This super-serious project converts your JDSL code into JavaScript at build-time, so you can skip the runtime overhead of JDSL.

## Inspiration

JDSL, the JSON-based Domain Specific Language, was created by Tom &mdash; a genius &mdash; to solve some of the world's hardest engineering challenges. These challenges included handling user input, and even updating a database. Things all developers struggled with before the advent of JDSL.

## Installation

```bash
go install github.com/connerdouglass/jdsl
```

## Usage

#### Transpile JDSL to JavaScript

Transpiling JDSL to JavaScript is blazingly fast. Just run the following command. The JavaScript output will be written to the `dist` directory.

```bash
jdsl -o dist --annotate --strict example/Customers.json example/Main.json
```

#### Run JDSL code directly

You can even run your JDSL code directly.

```bash
jdsl run example/Customers.json example/Main.json
```

