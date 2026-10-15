package saga

import (
	"context"
	"errors"
	"fmt"
)

type Data map[string]any

type Guard func(Data) bool
type Action func(context.Context, Data) error

type Transition struct {
	From       string
	To         string
	Guard      Guard
	Do         Action
	Compensate Action
}

type Workflow struct {
	Start       string
	Terminal    map[string]bool
	Transitions []Transition
	MaxSteps    int
}

type Result struct {
	State     string
	Completed []string
}

var ErrNoTransition = errors.New("no eligible transition")

func (w Workflow) Run(ctx context.Context, data Data) (Result, error) {
	if w.Start == "" {
		return Result{}, errors.New("workflow start state is required")
	}
	if data == nil {
		data = Data{}
	}
	limit := w.MaxSteps
	if limit <= 0 {
		limit = 100
	}
	state := w.Start
	var completed []Transition
	result := Result{State: state}
	for range limit {
		if w.Terminal[state] {
			result.State = state
			return result, nil
		}
		var selected *Transition
		for i := range w.Transitions {
			t := &w.Transitions[i]
			if t.From == state && (t.Guard == nil || t.Guard(data)) {
				selected = t
				break
			}
		}
		if selected == nil {
			err := fmt.Errorf("%w from state %q", ErrNoTransition, state)
			return result, errors.Join(err, compensate(ctx, data, completed))
		}
		if selected.Do != nil {
			if err := selected.Do(ctx, data); err != nil {
				cause := fmt.Errorf("transition %q -> %q: %w", selected.From, selected.To, err)
				return result, errors.Join(cause, compensate(ctx, data, completed))
			}
		}
		completed = append(completed, *selected)
		result.Completed = append(result.Completed, selected.To)
		state = selected.To
		result.State = state
	}
	err := fmt.Errorf("workflow exceeded %d transitions", limit)
	return result, errors.Join(err, compensate(ctx, data, completed))
}

func compensate(ctx context.Context, data Data, completed []Transition) error {
	var errs []error
	for i := len(completed) - 1; i >= 0; i-- {
		t := completed[i]
		if t.Compensate != nil {
			if err := t.Compensate(ctx, data); err != nil {
				errs = append(errs, fmt.Errorf("compensating %q -> %q: %w", t.From, t.To, err))
			}
		}
	}
	return errors.Join(errs...)
}
