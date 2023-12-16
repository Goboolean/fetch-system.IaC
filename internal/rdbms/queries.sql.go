// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: queries.sql

package rdbms

import (
	"context"
)

const deleteAllProducts = `-- name: DeleteAllProducts :exec
DELETE FROM product_meta
`

func (q *Queries) DeleteAllProducts(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllProducts)
	return err
}

const getAllProducts = `-- name: GetAllProducts :many
SELECT id, symbol, locale, market FROM product_meta
`

type GetAllProductsRow struct {
	ID     string
	Symbol string
	Locale string
	Market string
}

func (q *Queries) GetAllProducts(ctx context.Context) ([]GetAllProductsRow, error) {
	rows, err := q.db.Query(ctx, getAllProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllProductsRow
	for rows.Next() {
		var i GetAllProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Symbol,
			&i.Locale,
			&i.Market,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type InsertProductsParams struct {
	ID     string
	Symbol string
	Locale string
	Market string
}