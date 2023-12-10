package connect

import (
	"context"
	"encoding/json"
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

func (c *Client) Close() {}

func (c *Client) Ping(ctx context.Context) error {

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseurl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

	if _, err := io.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
    }

	return nil
}



func (c *Client) CheckCompatibility(ctx context.Context) error {

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connector-plugins", c.baseurl), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	var plugins []ConnectorPlugin
    if err := json.Unmarshal(body, &plugins); err != nil {
		return err
	}

	var (
		sinkPluginExists   bool = false
		sourcePluginExists bool = false
	)

	for _, plugin := range plugins {
		if plugin.Class == "com.mongodb.kafka.connect.MongoSourceConnector" && plugin.Type == "source" {
			sourcePluginExists = true
		}
		if plugin.Class == "com.mongodb.kafka.connect.MongoSinkConnector" && plugin.Type == "sink" {
			sinkPluginExists = true
		}
	}

	if !sinkPluginExists {
		return fmt.Errorf("sink plugin not found")
	}
	if !sourcePluginExists {
		return fmt.Errorf("source plugin not found")
	}
	return nil
}

