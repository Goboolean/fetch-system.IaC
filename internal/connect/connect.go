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
	"github.com/pkg/errors"
)



type Client struct {
	baseurl  string
	mongouri string
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

	mongouri, err := c.GetStringKey("MONGODB_CONNECTION_URI")
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
		mongouri: mongouri,
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(string(body))
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
		return fmt.Errorf(string(body))
	}

	var plugins []ConnectorPlugin
    if err := json.Unmarshal(body, &plugins); err != nil {
		return err
	}

	for _, plugin := range plugins {
		if plugin.Class == MongoSinkConnector && plugin.Type == "sink" {
			return nil
		}
	}

	return fmt.Errorf("plugin %s not found", MongoSinkConnector)
}



func (c *Client) CreateConnector(ctx context.Context, topic string) error {

	jsonData, err := json.Marshal(CreateConnectorRequest{
		Name: topic,
		Config: ConnectorConfig{
			ConnecctorClass: MongoSinkConnector,
			Topics:          topic,
			ConnectionUri:   c.mongouri,
			Database:        c.database,
			Collection:      topic,
			KeyConverter:    StringConverter,
			ValueConverter:  JsonConverter,
			ValueConverterSchemasEnable: "false",
			RotateIntervalMs: "1000",
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(string(body))
	}

	return nil
}


func (c *Client) CheckPluginConfig(ctx context.Context, topic string) error {
	
	jsonData, err := json.Marshal(PluginConfig{
		ConnectorClass: MongoSinkConnector,
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: received %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	return nil
}


func (c *Client) CheckTaskStatus(ctx context.Context, topic string, taskid int) error {

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connectors/%s/tasks/%d/status", c.baseurl, topic, taskid), nil)
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
		return fmt.Errorf(string(body))
	}

	var status TaskStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return err
	}

	if status.State != "RUNNING" {
		return fmt.Errorf("unexpected status: received %s, expected %s", status.State, "RUNNING")
	}

	return nil
}


func (c *Client) CheckTasksStatus(ctx context.Context, topic string)  error {

	var taskList []Task

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connectors/%s/tasks", c.baseurl, topic), nil)
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
		return fmt.Errorf(string(body))
	}

    if err := json.Unmarshal(body, &taskList); err != nil {
		return err
	}

	fmt.Println(string(body))

	for _, task := range taskList {
		if err := c.CheckTaskStatus(ctx, topic, task.TaskDetail.Task); err != nil {
			return errors.Wrap(err, "failed to call CheckTaskStatus")
		}
	}
	return nil
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf(string(body))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

    var connectors []string
    if err := json.Unmarshal(body, &connectors); err != nil {
        return nil, err
    }

	return connectors, nil
}