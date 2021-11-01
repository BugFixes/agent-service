package agent

import (
  "encoding/json"
  "net/http"

  bugLog "github.com/bugfixes/go-bugfixes/logs"
)

type AgentRequest struct {
	Name      string `json:"name"`
	AccountID string `json:"account_id"`
	AgentID   string `json:"agent_id"`
}

func jsonError(w http.ResponseWriter, msg string, errs error) {
  bugLog.Debugf("jsonError: %+v", errs)

  w.Header().Set("Content-Type", "text/json")
	if err := json.NewEncoder(w).Encode(struct {
		Error string
	}{
		Error: msg,
	}); err != nil {
		bugLog.Debugf("send %s failed: %+v", msg, err)
	}
}

func (ac *AgentClient) CreateAgent(w http.ResponseWriter, r *http.Request) {
	ar := AgentRequest{}

	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonError(w, "Invalid Body", err)
		return
	}

	a := NewBlankAgent(ar.Name, ar.AccountID)
	finalAgent, err := a.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonError(w, "Failed to create agent", err)
		return
	}

  finalAgent.AccountID = ar.AccountID

  if exists := ac.FindAgentByNameAndAccount(finalAgent); exists {
    w.WriteHeader(http.StatusConflict)
    jsonError(w, "Agent already exists", nil)
    return
  }

  if err := ac.saveAgent(finalAgent); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    jsonError(w, "failed to save agent", err)
    return
  }

	if err := json.NewEncoder(w).Encode(finalAgent); err != nil {
		w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "text/json")
		_ = json.NewEncoder(w).Encode(struct {
			AgentID string
			Key     string
			Secret  string
		}{
			AgentID: finalAgent.UUID,
			Key:     finalAgent.Key,
			Secret:  finalAgent.Secret,
		})
	}
}

func (ac *AgentClient) GetAgent(w http.ResponseWriter, r *http.Request) {
  vars := r.URL.Query()
  agentID := vars.Get("agent_uuid")

  if agentID == "" {
    w.WriteHeader(http.StatusBadRequest)
    jsonError(w, "Missing agent_id", nil)
    return
  }

  agent, err := ac.FindAgentByUUID(agentID)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    jsonError(w, "Agent not found", err)
    return
  }

  if err := json.NewEncoder(w).Encode(agent); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    jsonError(w, "Failed to encode agent", err)
    return
  }
}

func (ac *AgentClient) DeleteAgent(w http.ResponseWriter, r *http.Request) {
  vars := r.URL.Query()
  agentID := vars.Get("agent_uuid")

  if agentID == "" {
    w.WriteHeader(http.StatusBadRequest)
    jsonError(w, "Missing agent_id", nil)
    return
  }

  agent, err := ac.FindAgentByUUID(agentID)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    jsonError(w, "Agent not found", err)
    return
  }

  if err := ac.RemoveAgent(&agent); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    jsonError(w, "Failed to delete agent", err)
    return
  }

  w.WriteHeader(http.StatusNoContent)
}
