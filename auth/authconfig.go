// authconfig.go
package auth

import ()

// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth"`
	Email         string `json:"email"`
	ServerAddress string `json:"serveraddress,omitempty"`
}

//AuthConfigs
type AuthConfigs map[string]AuthConfig
