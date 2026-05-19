package domain

import "errors"

type RollbackPolicy struct {
	Auto             bool
	MaxRollbackDepth int
}

func NewRollbackPolicy(auto bool, depth int) (RollbackPolicy, error) {
	if depth < 0 {
		return RollbackPolicy{}, errors.New("rollback depth cannot be negative")
	}

	return RollbackPolicy{
		Auto:             auto,
		MaxRollbackDepth: depth,
	}, nil
}
