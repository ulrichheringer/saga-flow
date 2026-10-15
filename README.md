# saga-flow

A small Go library for running state-based workflows. It follows guarded transitions and compensates completed steps in reverse order if a later step fails.

Install with go get github.com/ulrichheringer/saga-flow/saga.
Run the example with go run ./cmd/demo.
