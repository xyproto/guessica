# Guessica [![Build Status](https://travis-ci.com/xyproto/guessica.svg?branch=master)](https://travis-ci.com/xyproto/guessica) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/guessica)](https://goreportcard.com/report/github.com/xyproto/guessica) [![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/guessica/master/LICENSE)

Update a `PKGBUILD` file by guessing the latest version number and finding the latest git tag and hash online.

![logo](img/guessica.svg)

## Installation (development version)

    go get -u github.com/xyproto/guessica/cmd/guessica

## Usage

	guessica PKGBUILD

## Note

The `pkgver` and `source` arrays will be guessed by searching the project webpage as defined by the `url`. For for projects on GitHub, `github.com` may also be visited.

Updating a `PKGBUILD` may or may not work. `guessica` will be doing its best, by guessing. Take a backup of your `PKGBUILD` first, if you need to.

## General info

* Version: 0.0.4
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
