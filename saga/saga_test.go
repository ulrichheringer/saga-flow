package saga

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestRunCompensatesCompletedTransitionsInReverse(t *testing.T) {
	var order []string
	w := Workflow{
		Start:    "start",
		Terminal: map[string]bool{"done": true},
		Transitions: []Transition{
			{From: "start", To: "reserved", Do: func(context.Context, Data) error { order = append(order, "reserve"); return nil }, Compensate: func(context.Context, Data) error { order = append(order, "release"); return nil }},
			{From: "reserved", To: "charged", Do: func(context.Context, Data) error { order = append(order, "charge"); return nil }, Compensate: func(context.Context, Data) error { order = append(order, "refund"); return nil }},
			{From: "charged", To: "done", Do: func(context.Context, Data) error { return errors.New("email failed") }},
		},
	}
	_, err := w.Run(context.Background(), nil)
	if err == nil || !reflect.DeepEqual(order, []string{"reserve", "charge", "refund", "release"}) {
		t.Fatalf("err=%v order=%v", err, order)
	}
}

func TestGuardSelectsTransition(t *testing.T) {
	w := Workflow{Start: "start", Terminal: map[string]bool{"approved": true}, Transitions: []Transition{
		{From: "start", To: "rejected", Guard: func(Data) bool { return false }},
		{From: "start", To: "approved", Guard: func(d Data) bool { return d["ok"] == true }},
	}}
	got, err := w.Run(context.Background(), Data{"ok": true})
	if err != nil || got.State != "approved" {
		t.Fatalf("result=%+v err=%v", got, err)
	}
}
