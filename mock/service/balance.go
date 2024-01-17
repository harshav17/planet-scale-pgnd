package service_mock

import (
	"context"

	planetscale "github.com/harshav17/planet_scale"
)

type BalanceService struct {
	GetGroupBalancesFn func(groupID int64) ([]*planetscale.Balance, error)
}

func (s BalanceService) GetGroupBalances(ctx context.Context, groupID int64) ([]*planetscale.Balance, error) {
	return s.GetGroupBalancesFn(groupID)
}
