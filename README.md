# go-manifold

Go client API for manifold (unstable)

[Code of Conduct](./.github/CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/go-manifold.svg?label=latest)](https://github.com/manifoldco/go-manifold/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/go-manifold)
[![Travis](https://img.shields.io/travis/manifoldco/go-manifold/master.svg)](https://travis-ci.org/manifoldco/go-manifold)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/go-manifold)](https://goreportcard.com/report/github.com/manifoldco/go-manifold)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

## Overview

This *will* be the Go implementation of a client API for Manifold's
non-provider facing REST services. For now, it is unstable as we build it, and is
used to share code between some of our other packages.

If you are a provider, you'll want to look at
[grafton](https://github.com/manifoldco/grafton) for verifying your
implementation, or [go-signature](https://github.com/manifoldco/go-signature)
for verifying requests have come from Manifold instead.
