package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeconfigPath returns the path to the kubeconfig file.
// It checks the KUBECONFIG environment variable, otherwise defaults to ~/.kube/config.
func GetKubeconfigPath() string {
	if path := os.Getenv("KUBECONFIG"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// ListContexts returns a list of all context names from the kubeconfig.
func ListContexts(kubeconfigPath string) ([]string, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	var contexts []string
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}
	return contexts, nil
}

// FilterContexts filters the list of contexts based on the provided filter string.
// The filter string can be a comma-separated list of substrings.
// If the filter is empty, all contexts are returned.
func FilterContexts(contexts []string, filter string) []string {
	if filter == "" {
		return contexts
	}

	filterParts := strings.Split(filter, ",")
	var filtered []string

	for _, ctx := range contexts {
		for _, part := range filterParts {
			if strings.Contains(ctx, strings.TrimSpace(part)) {
				filtered = append(filtered, ctx)
				break 
			}
		}
	}
	return filtered
}
