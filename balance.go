package planetscale

import (
	"context"
)

type (
	Balance struct {
		ExpenseGroupID int64              `json:"group_id"`
		UserID         string             `json:"user_id"`
		Amount         float64            `json:"amount"`
		BalanceItems   map[string]float64 `json:"balance_items"`
	}

	BalanceService interface {
		GetGroupBalances(ctx context.Context, groupID int64) ([]*Balance, error)
	}
)
