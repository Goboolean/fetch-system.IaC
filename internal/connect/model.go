package connect



type ConnectorConfig struct {
	ConnecctorClass string `json:"connector.class"`
	TasksMax        string `json:"tasks.max"`
	Topics          string `json:"topics"`
	ConnectionUrl   string `json:"connection.url"`
	Database        string `json:"database"`
	Collection      string `json:"collection"`
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


type CreateConnectorRequest struct {
	Name   string          `json:"name"`
	Config ConnectorConfig `json:"config"`
}

type ConnectorConfigResponse struct {
	Name string            `json:"name"`
	Config ConnectorConfig `json:"config"`
	Tasks []ConnectorTask  `json:"tasks"`
}