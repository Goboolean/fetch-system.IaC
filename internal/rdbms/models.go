// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package rdbms

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductMetum struct {
	ID          string
	Symbol      string
	Locale      string
	Market      string
	Name        pgtype.Text
	Description pgtype.Text
}