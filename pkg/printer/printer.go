package printer

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/suminhong/multi-kubectl/pkg/executor"
)

// PrintStream prints results as they come in, prefixed with the context name.
// Used for non-table commands (logs, describe, etc).
func PrintStream(res executor.Result) {
	header := fmt.Sprintf("[%s]", res.Context)
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %v\n", header, res.Error)
		if res.Output != "" {
			fmt.Println(res.Output)
		}
		return
	}

	lines := strings.Split(strings.TrimSpace(res.Output), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("%s %s\n", header, line)
		}
	}
}

// PrintTable aggregates results and prints them in a table format with a CONTEXT column.
// Used for 'get' commands.
func PrintTable(results []executor.Result) {
	// Sort results by Context name for consistent output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Context < results[j].Context
	})

	var validResults []executor.Result
	var header string
	hasResources := false

	// Filter out errors and empty results
	for _, res := range results {
		if res.Error != nil {
			// Print errors to stderr immediately or collect them? 
			// User said "exclude if no resources", but errors are different.
			// Let's print errors to stderr so they don't mess up the table.
			fmt.Fprintf(os.Stderr, "[%s] Error: %v\n", res.Context, res.Error)
			continue
		}
		
		if strings.Contains(res.Output, "No resources found") || strings.TrimSpace(res.Output) == "" {
			continue
		}

		// Extract header from the first valid result
		lines := strings.Split(strings.TrimSpace(res.Output), "\n")
		if len(lines) > 0 {
			if header == "" {
				header = lines[0]
			}
			validResults = append(validResults, res)
			hasResources = true
		}
	}

	if !hasResources {
		fmt.Println("No resources found in any context.")
		return
	}

	// Initialize tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	
	// Print Header
	fmt.Fprintf(w, "CONTEXT\t%s\n", header)

	// Print Rows
	for _, res := range validResults {
		lines := strings.Split(strings.TrimSpace(res.Output), "\n")
		// Skip header (line 0) for each result
		for i := 1; i < len(lines); i++ {
			fmt.Fprintf(w, "%s\t%s\n", res.Context, lines[i])
		}
	}
	w.Flush()
}
