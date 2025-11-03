package connector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DockerContext represents a Docker context configuration
type DockerContext struct {
	Name      string                 `json:"Name"`
	Metadata  map[string]interface{} `json:"Metadata"`
	Endpoints map[string]interface{} `json:"Endpoints"`
}

// DockerEndpoint contains connection details for a Docker daemon
type DockerEndpoint struct {
	Host          string
	SkipTLSVerify bool
	TLSVerify     bool
	CertPath      string
	CAFile        string
	CertFile      string
	KeyFile       string
}

// GetDockerContextsPath returns the path to Docker contexts directory
func GetDockerContextsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".docker", "contexts", "meta")
}

// GetCurrentContext reads the current context from Docker config
func GetCurrentContext() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".docker", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// No config file means use default context
		return "default", nil
	}

	var config struct {
		CurrentContext string `json:"currentContext"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}

	if config.CurrentContext == "" {
		return "default", nil
	}

	return config.CurrentContext, nil
}

// LoadContext loads a Docker context by name
func LoadContext(name string) (*DockerEndpoint, error) {
	if name == "" || name == "default" {
		// Return nil to indicate using environment variables
		return nil, nil
	}

	contextsPath := GetDockerContextsPath()
	if contextsPath == "" {
		return nil, fmt.Errorf("unable to determine contexts path")
	}

	// Look for context metadata
	metaPath := filepath.Join(contextsPath, hashContextName(name), "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("context %q not found: %w", name, err)
	}

	var ctx DockerContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, fmt.Errorf("failed to parse context: %w", err)
	}

	// Extract Docker endpoint
	dockerEndpoint, ok := ctx.Endpoints["docker"]
	if !ok {
		return nil, fmt.Errorf("no docker endpoint in context %q", name)
	}

	endpointMap, ok := dockerEndpoint.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid docker endpoint format in context %q", name)
	}

	endpoint := &DockerEndpoint{}

	// Extract host
	if host, ok := endpointMap["Host"].(string); ok {
		endpoint.Host = host
	}

	// Extract TLS settings
	if skipTLSVerify, ok := endpointMap["SkipTLSVerify"].(bool); ok {
		endpoint.SkipTLSVerify = skipTLSVerify
	}

	// Extract certificate paths if present
	if metadata, ok := ctx.Metadata["docker"].(map[string]interface{}); ok {
		if certPath, ok := metadata["CertPath"].(string); ok {
			endpoint.CertPath = certPath
			endpoint.TLSVerify = true
			// Set full paths to cert files
			endpoint.CAFile = filepath.Join(certPath, "ca.pem")
			endpoint.CertFile = filepath.Join(certPath, "cert.pem")
			endpoint.KeyFile = filepath.Join(certPath, "key.pem")
		}
	}

	return endpoint, nil
}

// hashContextName generates the hash directory name for a context
// Docker uses SHA256 hashing, but for simplicity we'll scan directories
func hashContextName(name string) string {
	contextsPath := GetDockerContextsPath()
	entries, err := os.ReadDir(contextsPath)
	if err != nil {
		return ""
	}

	// Try to find the context by reading each meta.json
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metaPath := filepath.Join(contextsPath, entry.Name(), "meta.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}

		var ctx DockerContext
		if err := json.Unmarshal(data, &ctx); err != nil {
			continue
		}

		if ctx.Name == name {
			return entry.Name()
		}
	}

	return ""
}

// ResolveDockerEndpoint resolves the Docker endpoint based on flags, env, and context
// Priority: hostFlag > contextFlag > DOCKER_CONTEXT env > current context > DOCKER_HOST env > default
func ResolveDockerEndpoint(hostFlag, contextFlag string) (*DockerEndpoint, error) {
	// Priority 1: --host flag
	if hostFlag != "" {
		return &DockerEndpoint{Host: hostFlag}, nil
	}

	// Priority 2: --context flag
	if contextFlag != "" {
		return LoadContext(contextFlag)
	}

	// Priority 3: DOCKER_CONTEXT environment variable
	if envContext := os.Getenv("DOCKER_CONTEXT"); envContext != "" {
		return LoadContext(envContext)
	}

	// Priority 4: Current context from config
	currentContext, err := GetCurrentContext()
	if err == nil && currentContext != "" && currentContext != "default" {
		endpoint, err := LoadContext(currentContext)
		if err == nil && endpoint != nil {
			return endpoint, nil
		}
	}

	// Priority 5: DOCKER_HOST environment variable (handled by NewClientFromEnv)
	// Priority 6: Default (unix:///var/run/docker.sock)
	return nil, nil
}
