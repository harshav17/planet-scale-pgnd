package service_mock

import (
	"context"

	planetscale "github.com/harshav17/planet_scale"
)

type ExpenseService struct {
	CreateExpenseFn func(ctx context.Context, expense *planetscale.Expense) error
}

func (s ExpenseService) CreateExpense(ctx context.Context, expense *planetscale.Expense) error {
	return s.CreateExpenseFn(ctx, expense)
}
