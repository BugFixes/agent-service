package agent

import (
  "fmt"

  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/jackc/pgx/v4"
)

func (ac AgentClient) getConnection() (*pgx.Conn, error) {
  conn, err := pgx.Connect(ac.Context,
    fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
      ac.Config.RDS.Username,
      ac.Config.RDS.Password,
      ac.Config.RDS.Hostname,
      ac.Config.RDS.Port,
      ac.Config.RDS.Database))
  if err != nil {
    return nil, bugLog.Error(err)
  }

  return conn, nil
}

func (ac *AgentClient) saveAgent(a *Agent) error {
  conn, err := ac.getConnection()
  if err != nil {
    return bugLog.Error(err)
  }

  defer func() {
    if err := conn.Close(ac.Context); err != nil {
      bugLog.Debugf("saveAgent close conn: %+v", err)
    }
  }()

  if _, err := conn.Exec(
    ac.Context,
    `INSERT INTO agents (uuid, key, secret, name, account_id) VALUES ($1, $2, $3, $4, $5)`,
    a.UUID,
    a.Key,
    a.Secret,
    a.Name,
    a.AccountID); err != nil {
    return bugLog.Error(err)
  }

  return nil
}

func (ac *AgentClient) FindAgentByNameAndAccount(a *Agent) bool {
  conn, err := ac.getConnection()
  if err != nil {
    bugLog.Debug(err)
    return true
  }

  defer func() {
    if err := conn.Close(ac.Context); err != nil {
      bugLog.Debugf("FindAgentByNameAndAccount close conn: %+v", err)
    }
  }()

  bugLog.Infof("AgentDebug: %+v", a)

  var foundAgent string
  if err := conn.QueryRow(
    ac.Context,
    `SELECT name FROM agents WHERE name = $1 AND account_id = $2`,
    a.Name,
    a.AccountID).Scan(&foundAgent); err != nil {
    if err == pgx.ErrNoRows {
      return false
    }
    return true
  }

  return false
}

func (ac AgentClient) FindAgentByID(agentID string) (Agent, error) {
  conn, err := ac.getConnection()
  if err != nil {
    bugLog.Debug(err)
    return Agent{}, bugLog.Error(err)
  }

  defer func() {
    if err := conn.Close(ac.Context); err != nil {
      bugLog.Debugf("FindAgentByID close conn: %+v", err)
    }
  }()

  var foundAgent Agent
  if err := conn.QueryRow(
    ac.Context,
    `SELECT uuid, key, secret, name, account_id FROM agents WHERE id = $1`,
    agentID).Scan(
    &foundAgent.UUID,
    &foundAgent.Key,
    &foundAgent.Secret,
    &foundAgent.Name,
    &foundAgent.AccountID); err != nil {
    bugLog.Debug(err)
    return Agent{}, bugLog.Error(err)
  }

  return foundAgent, nil
}

func (ac AgentClient) FindAgentByUUID(agentUUID string) (Agent, error) {
  conn, err := ac.getConnection()
  if err != nil {
    bugLog.Debug(err)
    return Agent{}, bugLog.Error(err)
  }

  defer func() {
    if err := conn.Close(ac.Context); err != nil {
      bugLog.Debugf("FindAgentBYUUID close conn: %+v", err)
    }
  }()

  var foundAgent Agent
  if err := conn.QueryRow(
    ac.Context,
    `SELECT uuid, key, secret, name, account_id FROM agents WHERE uuid = $1`,
    agentUUID).Scan(
    &foundAgent.UUID,
    &foundAgent.Key,
    &foundAgent.Secret,
    &foundAgent.Name,
    &foundAgent.AccountID); err != nil {
    bugLog.Debug(err)
    return Agent{}, bugLog.Error(err)
  }

  return foundAgent, nil
}

func (ac AgentClient) RemoveAgent(a *Agent) error {
  conn, err := ac.getConnection()
  if err != nil {
    return bugLog.Error(err)
  }

  defer func() {
    if err := conn.Close(ac.Context); err != nil {
      bugLog.Debugf("RemoveAgent close conn: %+v", err)
    }
  }()

  if _, err := conn.Exec(
    ac.Context,
    `DELETE FROM agents WHERE uuid = $1 AND account_id = $2`,
    a.UUID, a.AccountID); err != nil {
    return bugLog.Error(err)
  }

  return nil
}
