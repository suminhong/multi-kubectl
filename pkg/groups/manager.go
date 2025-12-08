package groups

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"sigs.k8s.io/yaml"
)

type Group struct {
	Name     string `json:"name"`
	Contexts string `json:"contexts"`
}

// LoadGroups loads the group configuration from ~/.kube/mk_config
func LoadGroups() ([]Group, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".kube", "mk_config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// It's okay if the file doesn't exist, just return empty
		return nil, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var groups []Group
	if err := yaml.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return groups, nil
}

// FindContexts returns the contexts for a given group name (or comma-separated names).
func FindContexts(groups []Group, groupNames string) (string, bool) {
	names := strings.Split(groupNames, ",")
	var allContexts []string
	foundAny := false

	for _, name := range names {
		name = strings.TrimSpace(name)
		for _, g := range groups {
			if g.Name == name {
				if g.Contexts != "" {
					allContexts = append(allContexts, g.Contexts)
				}
				foundAny = true
				break
			}
		}
	}
	
	if !foundAny {
		return "", false
	}
	
	return strings.Join(allContexts, ","), true
}
