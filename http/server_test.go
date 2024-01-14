package http

import (
	"context"
	"database/sql"
	"testing"

	planetscale "github.com/harshav17/planet_scale"
	db_mock "github.com/harshav17/planet_scale/mock/db"
)

type TestServer struct {
	*Server
	repos *planetscale.RepoProvider
}

func MustOpenServer(tb testing.TB) TestServer {
	tb.Helper()

	tm := db_mock.TransactionManager{}
	tm.ExecuteInTxFn = func(ctx context.Context, fn func(*sql.Tx) error) error {
		return fn(nil)
	}

	repos := planetscale.RepoProvider{}

	services := planetscale.ServiceProvider{}

	controllers := planetscale.ControllerProvider{}
	controllers.Product = NewProductController(&repos, &tm)
	controllers.ExpenseGroup = NewExpenseGroupController(&repos, &services, &tm)
	controllers.GroupMember = NewGroupMemberController(&repos, &tm)
	controllers.Expense = NewExpenseController(&repos, &tm)
	controllers.Settlement = NewSettlementController(&repos, &tm)
	controllers.SplitType = NewSplitTypeController(&repos, &tm)

	server := NewServer(&controllers)

	// Begin running test server.
	if err := server.Open(); err != nil {
		tb.Fatal(err)
	}

	return TestServer{
		Server: server,
		repos:  &repos,
	}
}

// MustCloseServer is a test helper function for shutting down the server.
// Fail on error.
func MustCloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}
