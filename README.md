[![Build Status](https://travis-ci.org/Azer0s/Persephone-Go.svg?branch=master)](https://travis-ci.org/Azer0s/Persephone-Go)  [![Go Report Card](https://goreportcard.com/badge/github.com/Azer0s/Persephone-Go)](https://goreportcard.com/report/github.com/Azer0s/Persephone-Go)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Azer0s/Persephone-Go/blob/master/README.md)

# Persephone-Go

A simple Persephone implementation in Go. 
This repository is the, de facto, reference implementation for the Persephone instruction set.
It is not incredibly fast or feature rich but it implements all Persephone language features.
There will be no additional things like JIT compilation or other fancy runtime techniques. This is just a proof-of-concept.
Whenever I add a new feature to Persephone, this implementation will get it first, so it's the most up-to-date.

[Language Spec](https://github.com/Azer0s/Persephone)

## Usage

```bash
$ git clone https://github.com/Azer0s/Persephone-Go.git
$ cd Persephone-Go
$ make
$ make run/examples/fib.psph
```
 
