package main

import (
  "fmt"
  "net/http"
  "os"
  "time"

  "github.com/bugfixes/agent_service/agent"
  "github.com/bugfixes/agent_service/config"
  bugLog "github.com/bugfixes/go-bugfixes/logs"
  bugfixes "github.com/bugfixes/go-bugfixes/middleware"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/keloran/go-probe"
)

func main() {
	bugLog.Local().Info("Starting Agent")

	cfg, err := config.BuildConfig()
	if err != nil {
		_ = bugLog.Errorf("buildConfig: %+v", err)
		return
	}

	if err := route(cfg); err != nil {
		_ = bugLog.Errorf("route failed: %+v", err)
		return
	}
}

func route(cfg config.Config) error {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)
	r.Use(bugfixes.BugFixes)

  agentPrefix := ""
  if os.Getenv("DEVELOPMENT") == "true" {
    agentPrefix = "agent"
  }

  r.Route(fmt.Sprintf("/%s", agentPrefix), func(r chi.Router) {
    r.Post("/", agent.NewAgent(cfg).CreateAgent)
    r.Delete("/", agent.NewAgent(cfg).DeleteAgent)
    r.Get("/", agent.NewAgent(cfg).GetAgent)
  })

  r.Get("/probe", probe.HTTP)

	bugLog.Local().Infof("Listening on port: %d\n", cfg.Local.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Local.Port), r); err != nil {
		return bugLog.Errorf("port: %+v", err)
	}

	return nil
}
