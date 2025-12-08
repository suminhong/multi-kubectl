package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/suminhong/multi-kubectl/pkg/executor"
	"github.com/suminhong/multi-kubectl/pkg/groups"
	"github.com/suminhong/multi-kubectl/pkg/kube"
	"github.com/suminhong/multi-kubectl/pkg/printer"
)

func main() {
	// 1. Parse arguments
	contextFilter, groupFilter, kubectlArgs := parseArgs(os.Args[1:])

	if len(kubectlArgs) == 0 || (len(kubectlArgs) == 1 && kubectlArgs[0] == "help") {
		printHelp()
		os.Exit(0)
	}

	// Resolve group if specified
	if groupFilter != "" {
		groupList, err := groups.LoadGroups()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load groups: %v\n", err)
		} else {
			if contexts, found := groups.FindContexts(groupList, groupFilter); found {
				if contextFilter != "" {
					fmt.Fprintln(os.Stderr, "Warning: --context and -g cannot be used together. Only -g will be applied.")
				}
				contextFilter = contexts
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Group '%s' not found in ~/.kube/mk_config\n", groupFilter)
			}
		}
	}

	// 2. Get kubeconfig and list contexts
	kubeconfig := kube.GetKubeconfigPath()
	contexts, err := kube.ListContexts(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// 3. Filter contexts
	targetContexts := kube.FilterContexts(contexts, contextFilter)
	if len(targetContexts) == 0 {
		fmt.Println("No contexts matched the filter.")
		os.Exit(0)
	}

	fmt.Printf("Running command on %d contexts: %s\n", len(targetContexts), strings.Join(targetContexts, ", "))

	// 4. Execute concurrently
	resultsChan := executor.Execute(targetContexts, kubectlArgs)

	// Determine if we should use table output (only for 'get', 'top', etc.)
	useTableOutput := isTableViewCommand(kubectlArgs)

	if useTableOutput {
		// Check for --no-headers flag
		noHeaders := false
		for _, arg := range kubectlArgs {
			if arg == "--no-headers" {
				noHeaders = true
				break
			}
		}

		// Collect all results for table formatting
		var allResults []executor.Result
		for res := range resultsChan {
			allResults = append(allResults, res)
		}
		printer.PrintTable(allResults, noHeaders)
	} else {
		// Stream results for other commands
		for res := range resultsChan {
			printer.PrintStream(res)
		}
	}
}

func isTableViewCommand(args []string) bool {
	// Commands that typically output a table and should be aggregated.
	tableCommands := map[string]bool{
		"get":           true,
		"top":           true,
		"events":        true,
		"api-resources": true,
	}

	for _, arg := range args {
		// If the argument matches a known table command, return true.
		if tableCommands[arg] {
			return true
		}
	}
	return false
}

func parseArgs(args []string) (string, string, []string) {
	var contextFilter string
	var groupFilter string
	var kubectlArgs []string

	for i := 0; i < len(args); i++ {
		if args[i] == "--context" {
			if i+1 < len(args) {
				contextFilter = args[i+1]
				i++ // Skip value
			}
		} else if strings.HasPrefix(args[i], "--context=") {
			contextFilter = strings.TrimPrefix(args[i], "--context=")
		} else if args[i] == "-g" || args[i] == "--group" {
			if i+1 < len(args) {
				groupFilter = args[i+1]
				i++ // Skip value
			}
		} else if strings.HasPrefix(args[i], "--group=") {
			groupFilter = strings.TrimPrefix(args[i], "--group=")
		} else {
			kubectlArgs = append(kubectlArgs, args[i])
		}
	}
	return contextFilter, groupFilter, kubectlArgs
}

func printHelp() {
	fmt.Println("Usage: mk [kubectl-args] [--context <filter>] [-g <group>]")
	fmt.Println("  (Note: --context and -g cannot be used together)")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --context <filter>   Filter contexts by substring or comma-separated list")
	fmt.Println("  -g, --group <name>   Run on contexts defined in a group (comma-separated list supported)")
	fmt.Println("")
	fmt.Println("Configuration (~/.kube/mk_config):")
	fmt.Println("  - name: dev")
	fmt.Println("    contexts: dev-cluster1, dev-cluster2")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  mk get pods -n devops               # Run on all contexts")
	fmt.Println("  mk get pods --context prod          # Run on contexts containing 'prod'")
	fmt.Println("  mk get pods --context prod,stage    # Run on contexts containing 'prod' or 'stage'")
	fmt.Println("  mk get pods -g dev                  # Run on contexts in 'dev' group")
	fmt.Println("  mk get pods -g dev,aws              # Run on contexts in 'dev' and 'aws' groups")
}
