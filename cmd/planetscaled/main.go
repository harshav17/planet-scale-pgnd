package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	planetscale "github.com/harshav17/planet_scale"
	"github.com/harshav17/planet_scale/db"
	"github.com/harshav17/planet_scale/http"
	"github.com/harshav17/planet_scale/service"
	utilities "github.com/harshav17/planet_scale/utilites"
	"github.com/joho/godotenv"
)

func main() {
	// Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	HTTPServer *http.Server
	DB         *db.DB
}

func NewMain() *Main {
	return &Main{}
}

func (m *Main) Run(ctx context.Context) (err error) {
	logger := utilities.GetLogger()
	slog.SetDefault(logger)

	DSN, ok := os.LookupEnv("DSN")
	if !ok {
		slog.Error("DSN not set")
	} else {
		slog.Info(DSN)
	}

	// database
	m.DB = db.NewDB(DSN)
	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	// transaction manager
	tm := db.NewTransactionManager(m.DB)

	// repos
	repos := planetscale.RepoProvider{}
	repos.Product = db.NewProductRepo(m.DB)
	repos.ExpenseGroup = db.NewExpenseGroupRepo(m.DB)
	repos.GroupMember = db.NewGroupMemberRepo(m.DB)
	repos.Expense = db.NewExpenseRepo(m.DB)
	repos.ExpenseParticipant = db.NewExpenseParticipantRepo(m.DB)
	repos.Settlement = db.NewSettlementRepo(m.DB)
	repos.SplitType = db.NewSplitTypeRepo(m.DB)

	// services
	services := planetscale.ServiceProvider{}
	services.Balance = service.NewBalanceService(&repos, tm)

	// controllers
	controllers := planetscale.ControllerProvider{}
	controllers.Product = http.NewProductController(&repos, tm)
	controllers.ExpenseGroup = http.NewExpenseGroupController(&repos, &services, tm)
	controllers.GroupMember = http.NewGroupMemberController(&repos, tm)
	controllers.Expense = http.NewExpenseController(&repos, tm)
	controllers.Settlement = http.NewSettlementController(&repos, tm)
	controllers.SplitType = http.NewSplitTypeController(&repos, tm)

	// start the HTTP server.
	m.HTTPServer = http.NewServer(&controllers)
	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	return nil
}
