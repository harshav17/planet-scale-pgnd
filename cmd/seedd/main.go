package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	planetscale "github.com/harshav17/planet_scale"
	"github.com/harshav17/planet_scale/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	DSN, ok := os.LookupEnv("DSN")
	if !ok {
		slog.Error("DSN not set")
	} else {
		slog.Info(DSN)
	}

	// database
	dbNew := db.NewDB(DSN)
	if err := dbNew.Open(); err != nil {
		log.Fatal("cannot open db: %w", err)
	}

	// transaction manager
	tm := db.NewTransactionManager(dbNew)

	// repositories
	expenseGroupRepo := db.NewExpenseGroupRepo(dbNew)
	gmRepo := db.NewGroupMemberRepo(dbNew)
	userRepo := db.NewUserRepo(dbNew)
	expenseRepo := db.NewExpenseRepo(dbNew)

	createFakeDataFunc := func(tx *sql.Tx) error {
		// create fakeUser user
		fakeUser := &planetscale.User{
			UserID: "test_user_id",
			Name:   gofakeit.Name(),
			Email:  gofakeit.Email(),
		}
		if err := userRepo.Create(tx, fakeUser); err != nil {
			return fmt.Errorf("cannot create user: %w", err)
		}

		// create fakeUser user 2
		fakeUser2 := &planetscale.User{
			UserID: gofakeit.UUID(),
			Name:   gofakeit.Name(),
			Email:  gofakeit.Email(),
		}
		if err := userRepo.Create(tx, fakeUser2); err != nil {
			return fmt.Errorf("cannot create user: %w", err)
		}

		// create fake expense group
		fakeExpenseGroup := &planetscale.ExpenseGroup{
			GroupName: gofakeit.Name(),
			CreateBy:  fakeUser.UserID,
		}
		if err := expenseGroupRepo.Create(tx, fakeExpenseGroup); err != nil {
			return fmt.Errorf("cannot create expense group: %w", err)
		}

		// create fake group member
		fakeGroupMember := &planetscale.GroupMember{
			GroupID: fakeExpenseGroup.ExpenseGroupID,
			UserID:  fakeUser.UserID,
		}
		if err := gmRepo.Create(tx, fakeGroupMember); err != nil {
			return fmt.Errorf("cannot create group member: %w", err)
		}

		// create fake group member 2
		fakeGroupMember2 := &planetscale.GroupMember{
			GroupID: fakeExpenseGroup.ExpenseGroupID,
			UserID:  fakeUser2.UserID,
		}
		if err := gmRepo.Create(tx, fakeGroupMember2); err != nil {
			return fmt.Errorf("cannot create group member: %w", err)
		}

		// create fake expense
		fakeExpense := &planetscale.Expense{
			GroupID:     &fakeExpenseGroup.ExpenseGroupID,
			PaidBy:      fakeUser.UserID,
			Amount:      gofakeit.Price(0, 1000),
			Description: gofakeit.ProductDescription(),
			Timestamp:   gofakeit.Date(),
			SplitTypeID: 1, // TODO load from db into cache
			CreatedBy:   fakeUser.UserID,
			UpdatedBy:   fakeUser.UserID,
		}
		if err := expenseRepo.Create(tx, fakeExpense); err != nil {
			return fmt.Errorf("cannot create expense: %w", err)
		}

		// create fake expense 2
		fakeExpense2 := &planetscale.Expense{
			GroupID:     &fakeExpenseGroup.ExpenseGroupID,
			PaidBy:      fakeUser2.UserID,
			Amount:      gofakeit.Price(0, 1000),
			Description: gofakeit.ProductDescription(),
			Timestamp:   gofakeit.Date(),
			SplitTypeID: 1,
			CreatedBy:   fakeUser2.UserID,
			UpdatedBy:   fakeUser2.UserID,
		}
		if err := expenseRepo.Create(tx, fakeExpense2); err != nil {
			return fmt.Errorf("cannot create expense: %w", err)
		}

		return nil
	}

	err = tm.ExecuteInTx(context.Background(), createFakeDataFunc)
	if err != nil {
		log.Fatal("cannot execute in tx: %w", err)
	}
}
