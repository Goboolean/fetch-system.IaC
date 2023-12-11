package connect

import (
	"bytes"
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

	req.Header.Set("Content-Type", "application/json")

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

	req.Header.Set("Content-Type", "application/json")

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


func (c *Client) CreateConnector(ctx context.Context, topic string) error {

	jsonData, err := json.Marshal(CreateConnectorRequest{
		Name: topic,
		Config: ConnectorConfig{
			ConnecctorClass: "com.mongodb.kafka.connect.MongoSinkConnector",
			Topics:          topic,
			ConnectionUri:   c.mongourl,
			Database:        c.database,
			Collection:      topic,
			KeyConverter:    "org.apache.kafka.connect.storage.StringConverter",
			ValueConverter:  "org.apache.kafka.connect.json.JsonConverter",
			ValueConverterSchemasEnable: "false",
			//KeyIgnore:       "true",
			//InsertMode:      "insert",
			//WritemodelStrategy: "com.mongodb.kafka.connect.sink.writemodel.strategy.ReplaceOneBusinessKeyStrategy",
		},
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/connectors", c.baseurl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusCreated)
	}

	return nil
}


func (c *Client) CheckPluginConfig(ctx context.Context, topic string) error {
	
	jsonData, err := json.Marshal(PluginConfig{
		ConnectorClass: "com.mongodb.kafka.connect.MongoSinkConnector",
		TasksMax:       "1",
		Topics:         topic,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/connector-plugins/%s/config/validate", c.baseurl, "com.mongodb.kafka.connect.MongoSinkConnector"), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	return nil
}



func (c *Client) GetConnectorConfiguration(ctx context.Context, topic string) (ConnectorConfigResponse, error){

	var config ConnectorConfigResponse

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connectors/%s", c.baseurl, topic), nil)
	if err != nil {
		return config, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return config, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return config, fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return config, err
	}

    if err := json.Unmarshal(body, &config); err != nil {
		return config, err
	}
	return config, nil
}


func (c *Client) DeleteConnector(ctx context.Context, topic string) error {

	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/connectors/%s", c.baseurl, topic), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}


func (c *Client) GetConnectors(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connectors", c.baseurl), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

    var connectors []string
    if err := json.Unmarshal(body, &connectors); err != nil {
        return nil, err
    }

	return connectors, nil
}