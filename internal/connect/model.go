package connect



type ConnectorConfig struct {
	ConnecctorClass             string `json:"connector.class"`
	Topics                      string `json:"topics"`
	ConnectionUri               string `json:"connection.uri"`
	Database                    string `json:"database"`
	Collection                  string `json:"collection"`
	KeyConverter                string `json:"key.converter"`
	ValueConverter              string `json:"value.converter"`
	ValueConverterSchemasEnable string `json:"value.converter.schemas.enable"`
	//TasksMax               string `json:"tasks.max"`
	//KeyIgnore              string `json:"key.ignore"`
	//InsertMode             string `json:"insert.mode"`
	//WritemodelStrategy     string `json:"writemodel.strategy"`
}

type ConnectorTask struct {
	Connector string `json:"connector"`
	Task      int    `json:"task"`
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


type CreateConnectorRequest struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config"`
}

type ConnectorConfigResponse struct {
	Name string            `json:"name"`
	Config ConnectorConfig `json:"config"`
	Tasks []ConnectorTask  `json:"tasks"`
}