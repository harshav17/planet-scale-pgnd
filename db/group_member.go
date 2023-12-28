package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	planetscale "github.com/harshav17/planet_scale"
)

type groupMemberRepo struct {
	db *DB
}

func NewGroupMemberRepo(db *DB) *groupMemberRepo {
	return &groupMemberRepo{
		db: db,
	}
}

func (r *groupMemberRepo) Get(tx *sql.Tx, groupID int64, userID string) (*planetscale.GroupMember, error) {
	query := `SELECT group_id, user_id, joined_at FROM group_members WHERE group_id = ? AND user_id = ?`

	var group planetscale.GroupMember
	row := tx.QueryRow(query, groupID, userID)
	err := row.Scan(&group.GroupID, &group.UserID, (*NullTime)(&group.JoinedAt))
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows error specifically if needed
			return nil, fmt.Errorf("no group member found with ID %d", groupID)
		}
		return nil, err
	}
	slog.Info("loaded group member", slog.Int64("id", group.GroupID))

	return &group, nil
}

func (r *groupMemberRepo) Create(tx *sql.Tx, group *planetscale.GroupMember) error {
	query := `INSERT INTO group_members (group_id, user_id) VALUES (?, ?)`

	result, err := tx.Exec(query, group.GroupID, group.UserID)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	slog.Info("created group member", slog.Int64("id", group.GroupID))

	return nil
}

func (r *groupMemberRepo) Delete(tx *sql.Tx, groupID int64, userID string) error {
	query := `DELETE FROM group_members WHERE group_id = ? AND user_id = ?`

	result, err := tx.Exec(query, groupID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no group member found with ID %d", groupID)
	}
	slog.Info("deleted group member", slog.Int64("id", groupID))

	return nil
}
