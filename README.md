# Flash

A flashcard application built to explore the Hexagonal Architecture pattern and CI pipelines with GitHub Actions. The app features core requirements of a flashcard application and exposes multiple actors for interacting with the app and storing data.

[![Tests workflow](https://img.shields.io/github/workflow/status/jmcveigh55/flash/Test%20Base?longCache=tru&label=tests&logo=github&logoColor=fff)](https://github.com/jmcveigh55/flash/actions?query=workflow%3ATest%20Base)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmcveigh55/flash)](https://goreportcard.com/report/github.com/jmcveigh55/flash)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Implementations

Currently the Flash application is configured for In-Memory data storage and CLI interaction. Below is the full list of implemented actors:

>### Storage
>
>1. In-Memory
>2. JSON (default)
>
>### Interface
>
>1. CLI (default)
>2. TUI (TODO)

## Resources

Below are a list of resources referenced when creating this application:

1. [Go Structure Examples](https://github.com/katzien/go-structure-examples) by [Kat Zien](https://github.com/katzien)
2. [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3) by [Matías Varela](https://medium.com/@matiasvarela)
