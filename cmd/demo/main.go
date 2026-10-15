package main

import (
	"context"
	"fmt"

	"github.com/ulrichheringer/saga-flow/saga"
)

func main() {
	w := saga.Workflow{Start: "new", Terminal: map[string]bool{"complete": true}, Transitions: []saga.Transition{
		{From: "new", To: "reserved", Do: func(_ context.Context, d saga.Data) error {
			d["reservation"] = "R-17"
			fmt.Println("inventory reserved")
			return nil
		}, Compensate: func(_ context.Context, d saga.Data) error { fmt.Println("inventory released"); return nil }},
		{From: "reserved", To: "complete", Guard: func(d saga.Data) bool { return d["payment_ok"] == true }, Do: func(context.Context, saga.Data) error { fmt.Println("payment accepted"); return nil }},
	}}
	result, err := w.Run(context.Background(), saga.Data{"payment_ok": true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("workflow ended in %s (%v)\n", result.State, result.Completed)
}
