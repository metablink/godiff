# godiff
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)]()
[![Build Status](https://travis-ci.org/metablink/godiff.svg?branch=master)](https://travis-ci.org/metablink/godiff)
[![Test Coverage](https://api.codeclimate.com/v1/badges/368430c7858f2a9afaac/test_coverage)](https://codeclimate.com/github/metablink/godiff/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/368430c7858f2a9afaac/maintainability)](https://codeclimate.com/github/metablink/godiff/maintainability)

CSV differ written in Go

Currently nonfunctional. Provides some aggregate-level diff data and nothing more. Also currently assumes that the input CSV is pre-sorted based on key field. In the future, the sorting will be done automatically, but first we have to write a Go external sorting function that can handle arbitrarily large files.
