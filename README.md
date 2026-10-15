# saga-flow — Workflow and Saga Engine

[![Go 1.22](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A small Go library for running state-based workflows. It follows guarded transitions and compensates completed steps in reverse order if a later step fails.

Install with go get github.com/ulrichheringer/saga-flow/saga.
Run the example with go run ./cmd/demo.
