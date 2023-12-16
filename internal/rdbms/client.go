package rdbms

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/Goboolean/common/pkg/resolver"
)


type Client = Queries

func NewDB(c *resolver.ConfigMap) (*Client, error) {
	user, err := c.GetStringKey("USER")
	if err != nil {
		return nil, err
	}

	password, err := c.GetStringKey("PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	database, err := c.GetStringKey("DATABASE")
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err := pgx.Connect(context.Background(), psqlInfo)

	if err != nil {
		return nil, err
	}

	q := New(db)

	return q, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.db.(*pgx.Conn).Ping(ctx)
}

func (c *Client) Close() {
	c.db.(*pgx.Conn).Close(context.Background())
}