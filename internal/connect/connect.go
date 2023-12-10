package connect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Goboolean/common/pkg/resolver"
)



type Client struct {
	baseurl  string

	mongourl string
	database string
}

func New(c *resolver.ConfigMap) (*Client, error) {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	flag, _, err := c.GetBoolKeyOptional("APPLY_SSL")
	if err != nil {
		return nil, err
	}

	var scheme string
	if flag {
		scheme = "https"
	} else {
		scheme = "http"
	}

	mongourl, err := c.GetStringKey("MONGODB_CONNECTION_URL")
	if err != nil {
		return nil, err
	}

	database, err := c.GetStringKey("MONGODB_DATABASE")
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
	}

	return &Client{
		baseurl: u.String(),
		mongourl: mongourl,
		database: database,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseurl, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d", resp.StatusCode)
    }
	return nil
}

func (c *Client) Close() {}


