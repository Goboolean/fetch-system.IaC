package connect



const (
	MongoSinkConnector	= "com.mongodb.kafka.connect.MongoSinkConnector"
	MongoSourceConnector= "com.mongodb.kafka.connect.MongoSourceConnector"
	StringConverter		= "org.apache.kafka.connect.storage.StringConverter"
	JsonConverter		= "org.apache.kafka.connect.json.JsonConverter"
)


type ConnectorConfig struct {
	ConnecctorClass             string `json:"connector.class"`
	Topics                      string `json:"topics"`
	ConnectionUri               string `json:"connection.uri"`
	Database                    string `json:"database"`
	Collection                  string `json:"collection"`
	KeyConverter                string `json:"key.converter"`
	ValueConverter              string `json:"value.converter"`
	ValueConverterSchemasEnable string `json:"value.converter.schemas.enable"`
	RotateIntervalMs            string `json:"rotate.interval.ms"`
}

type ConnectorPlugin struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type PluginConfig struct {
	ConnectorClass string `json:"connector.class"`
	TasksMax       string `json:"tasks.max"`
	Topics         string `json:"topics"`
}

type TaskDetail struct {
	Connector string `json:"connector"`
	Task      int    `json:"task"`
}

type TaskStatus struct {
	State    string `json:"state"`
	ID       int    `json:"id"`
	WorkerId string `json:"worker_id"`
}

type TaskConfig struct {
	TaskClass string `json:"task.class"`
	Topics   string `json:"topics"`
}

type Task struct {
	TaskDetail TaskDetail `json:"id"`
	Config     TaskConfig `json:"config"`
}

type CreateConnectorRequest struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config"`
}