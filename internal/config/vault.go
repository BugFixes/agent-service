package config

import (
  "os"

  bugLog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/hashicorp/vault/api"
)

func GetVaultSecrets(vaultAddress, secretPath string) (map[string]interface{}, error) {
var m = make(map[string]interface{})

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return m, bugLog.Error("token not found")
	}

  cfg := api.DefaultConfig()
  cfg.Address = vaultAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		return m, bugLog.Errorf("client: %+v", err)
	}
	client.SetToken(token)

	data, err := client.Logical().Read(secretPath)
	if err != nil {
		return m, bugLog.Errorf("read: %+v", err)
	}

	return data.Data, nil
}
