package agent

import (
  "context"

  "github.com/bugfixes/agent_service/config"
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/google/uuid"
)

//go:generate mockery --name=Agents
type Agents interface {
	Create() (*Agent, error)
	Fetch(id string) (*Agent, error)
	Delete(a Agent) error
}

type Credentials struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type Agent struct {
	ID   int `json:"-"`
	UUID string `json:"agent_id"`
	Name string `json:"name"`

	Credentials
	AccountID string `json:"account_id"`
}

type AgentClient struct {
	Config  config.Config
	Context context.Context
}

func NewAgent(c config.Config) *AgentClient {
	return &AgentClient{
		Config:  c,
		Context: context.Background(),
	}
}

func NewBlankAgent(name, accountID string) *Agent {
	return &Agent{
		Name:      name,
		AccountID: accountID,
	}
}

func (a Agent) Create() (*Agent, error) {
	id, err := createID()
	if err != nil {
		return nil, bugLog.Errorf("agent create: %+v", err)
	}
	a.UUID = id

	key, err := createKey()
	if err != nil {
		return nil, bugLog.Errorf("key create: %+v", err)
	}
	a.Key = key

	secret, err := createSecret()
	if err != nil {
		return nil, bugLog.Errorf("secret create: %+v", err)
	}
	a.Secret = secret

	return &a, nil
}

func createID() (string, error) {
	return generateUUID()
}

func createKey() (string, error) {
	return generateUUID()
}

func createSecret() (string, error) {
	return generateUUID()
}

func generateUUID() (string, error) {
	s, err := uuid.NewUUID()
	if err != nil {
		return "", bugLog.Errorf("generateUUID: %+v", err)
	}

	return s.String(), nil
}
