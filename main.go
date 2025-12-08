package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/suminhong/multi-kubectl/pkg/executor"
	"github.com/suminhong/multi-kubectl/pkg/kube"
	"github.com/suminhong/multi-kubectl/pkg/printer"
)

func main() {
	// 1. Parse arguments
	contextFilter, kubectlArgs := parseArgs(os.Args[1:])

	if len(kubectlArgs) == 0 {
		printHelp()
		os.Exit(0)
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

	// Determine if we should use table output (only for 'get' commands)
	command := getCommand(kubectlArgs)
	isGetCommand := command == "get"

	if isGetCommand {
		// Collect all results for table formatting
		var allResults []executor.Result
		for res := range resultsChan {
			allResults = append(allResults, res)
		}
		printer.PrintTable(allResults)
	} else {
		// Stream results for other commands
		for res := range resultsChan {
			printer.PrintStream(res)
		}
	}
}

func getCommand(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}
	return ""
}

func parseArgs(args []string) (string, []string) {
	var contextFilter string
	var kubectlArgs []string

	for i := 0; i < len(args); i++ {
		if args[i] == "--context" {
			if i+1 < len(args) {
				contextFilter = args[i+1]
				i++ // Skip value
			}
		} else if strings.HasPrefix(args[i], "--context=") {
			contextFilter = strings.TrimPrefix(args[i], "--context=")
		} else {
			kubectlArgs = append(kubectlArgs, args[i])
		}
	}
	return contextFilter, kubectlArgs
}

func printHelp() {
	fmt.Println("Usage: mk [kubectl-args] [--context <filter>]")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  mk get pods -n devops               # Run on all contexts")
	fmt.Println("  mk get pods --context prod          # Run on contexts containing 'prod'")
	fmt.Println("  mk get pods --context prod,stage    # Run on contexts containing 'prod' or 'stage'")
}
