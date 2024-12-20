package repository

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BorrowRepository interface {
	GetAllBorrowLists(ctx context.Context, opts query.QueryOptions) ([]model.Borrow, error)
	GetBorrowByID(ctx context.Context, id int) (*model.Borrow, error)
	GetBorrowByUserName(ctx context.Context, name string) (*model.Borrow, error)
	CreateBorrow(ctx context.Context, b model.Borrow) (int, error)
	UpdateBorrow(ctx context.Context, b model.Borrow) error
	DeleteBorrow(ctx context.Context, id int) error
}

type borrowRepository struct {
	db *sqlx.DB
}

func NewBorrowRepository(db *sqlx.DB) BorrowRepository {
	return &borrowRepository{
		db: db,
	}
}

func (r *borrowRepository) GetAllBorrowLists(ctx context.Context, opts query.QueryOptions) ([]model.Borrow, error) {
	q, args := query.BuildSelectQuery("borrows", opts)
	var borrows []model.Borrow
	err := r.db.SelectContext(ctx, &borrows, q, args...)
	return borrows, err
}

func (r *borrowRepository) GetBorrowByID(ctx context.Context, id int) (*model.Borrow, error) {
	var borrow model.Borrow
	err := r.db.GetContext(ctx, &borrow, "SELECT * FROM borrows WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &borrow, err
}

func (r *borrowRepository) GetBorrowByUserName(ctx context.Context, name string) (*model.Borrow, error) {
	var borrow model.Borrow
	err := r.db.GetContext(ctx, &borrow, "SELECT * FROM borrows WHERE user_name=$1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &borrow, err
}

func (r *borrowRepository) CreateBorrow(ctx context.Context, b model.Borrow) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO borrows (book_id, user_name, borrowed_at) VALUES ($1, $2, $3) RETURNING id",
		b.BookID, b.UserName, b.BorrowedAt,
	).Scan(&id)
	return id, err
}

func (r *borrowRepository) UpdateBorrow(ctx context.Context, b model.Borrow) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE borrows SET book_id=$1, user_name=$2, borrowed_at=$3 WHERE id=$4",
		b.BookID, b.UserName, b.BorrowedAt, b.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *borrowRepository) DeleteBorrow(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM borrows WHERE id=$1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}
