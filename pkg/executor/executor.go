package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

// Result represents the output of a command execution on a specific context.
type Result struct {
	Context string
	Output  string
	Error   error
}

// Execute runs the given kubectl arguments on the specified contexts concurrently.
// It returns a channel that receives the results as they complete.
func Execute(contexts []string, args []string) <-chan Result {
	results := make(chan Result, len(contexts))
	var wg sync.WaitGroup

	for _, ctx := range contexts {
		wg.Add(1)
		go func(ctx string) {
			defer wg.Done()
			output, err := runKubectl(ctx, args)
			results <- Result{
				Context: ctx,
				Output:  output,
				Error:   err,
			}
		}(ctx)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func runKubectl(context string, args []string) (string, error) {
	// Prepend --context flag
	cmdArgs := append([]string{"--context", context}, args...)
	cmd := exec.Command("kubectl", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	if err != nil {
		return output, fmt.Errorf("command failed: %w", err)
	}

	return output, nil
}
