package rdbms

import (
	"database/sql"
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

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

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (p *Client) Close() error {
	return p.db.Close()
}

func (p *Client) Ping() error {
	return p.db.Ping()
}