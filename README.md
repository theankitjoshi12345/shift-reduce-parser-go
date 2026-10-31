# Shift-Reduce Parser in Go

This project implements the CMPS 450 shift/reduce parser for the supplied
expression grammar. It includes the required parse-stack operations and the
optional parse-tree bonus:

- all shift, reduce, accept, and error actions;
- the supplied action and goto tables;
- a step-by-step parser trace;
- parse-tree-stack output during parsing; and
- an indented final parse tree.

## Requirements

Go 1.22 or newer.

## Run

The documented project input is the default:

```bash
go run . 
```

To try the parenthesized example:

```bash
go run . -input='(id)'
```

The `-input` value may contain `id`, `+`, `*`, `(`, and `)` with optional
spaces. Invalid input is reported as ungrammatical.

## Test

```bash
go test ./...
```

The assignment PDFs are kept in [`docs/pdfs/`](docs/pdfs/) and are not needed
to build or run the program.
