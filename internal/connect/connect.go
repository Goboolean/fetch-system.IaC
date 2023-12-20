package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
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



func (c *Client) CreateSingleTopicConnector(ctx context.Context, name string, tasks int, conf ConnectorTopicConfig) error {

	jsonData, err := json.Marshal(CreateConnectorRequest{
		Name: name,
		Config: ConnectorConfig{
			ConnecctorClass: MongoSinkConnector,
			Topics:          conf.Topic,
			ConnectionUri:   c.mongouri,
			Database:        c.database,
			Collection:      conf.Collection,
			KeyConverter:    StringConverter,
			ValueConverter:  JsonConverter,
			ValueConverterSchemasEnable: "false",
			RotateIntervalMs: strconv.Itoa(conf.RotateIntervalMs),
			//DocumentIdStrategy: DocumentIdStrategy,
			MaxTasks: strconv.Itoa(tasks),
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


func (c *Client) CreateBulkTopicConnector(ctx context.Context, name string, tasks int, configs []ConnectorTopicConfig) error {
	config := make(map[string]string)

	config["connector.class"] = MongoSinkConnector
	config["connection.uri"] = c.mongouri
	config["database"] = c.database
	config["key.converter"] = StringConverter
	config["value.converter"] = JsonConverter
	config["value.converter.schemas.enable"] = "false"
	config["rotate.interval.ms"] = "1000000"
	config["document.id.strategy"] = "com.mongodb.kafka.connect.sink.processor.id.strategy.ProvidedInKeyStrategy"
	config["max.tasks"] = strconv.Itoa(tasks)

	topicList := make([]string, len(configs))
	for i, conf := range configs {
		topicList[i] = conf.Topic
	}
	config["topics"] = strings.Join(topicList, ",")

	for _, topicConfig := range configs {
		config[fmt.Sprintf("topic.override.%s.collection", topicConfig.Topic)] = topicConfig.Collection
		config[fmt.Sprintf("topic.override.%s.rotate.interval.ms", topicConfig.Topic)] = strconv.Itoa(topicConfig.RotateIntervalMs)
	}

	jsonData, err := json.Marshal(CreateBulkConnectorRequest{
		Name: name,
		Config: config,
	})
	if err != nil {
		return err
	}

	deadline, exists := ctx.Deadline()
	if !exists {
		deadline = time.Now().Add(time.Hour)
	}
	client := &http.Client{Timeout: time.Until(deadline)}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/connectors", c.baseurl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(string(body))
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


func (c *Client) CheckTasksStatus(ctx context.Context, name string) (int, error) {

	time.Sleep(1 * time.Second)

	var taskList []Task

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/connectors/%s/tasks", c.baseurl, name), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(string(body))
	}

    if err := json.Unmarshal(body, &taskList); err != nil {
		return 0, err
	}

	for _, task := range taskList {

		var queryErr error
		var retry = 0

		for {
			if err := c.CheckTaskStatus(ctx, name, task.TaskDetail.Task); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return len(taskList), errors.Join(err, fmt.Errorf("retry %d times", retry), queryErr)
				} else {
					queryErr = err
					retry++
					time.Sleep(1 * time.Second)
				}
			} else {
				break
			}
		}
	}
	return len(taskList), nil
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


func (c *Client) DeleteAllConnectors(ctx context.Context) error {

	connectors, err := c.GetConnectors(ctx)
	if err != nil {
		return err
	}

	for _, connector := range connectors {
		if err := c.DeleteConnector(ctx, connector); err != nil {
			return err
		}
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