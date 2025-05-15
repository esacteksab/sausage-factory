package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/cli/safeexec"
)

var (
	binary     string
	workingDir string
)

func main() {
	os.Setenv("TF_IN_AUTOMATION", "true")
	// Create a context that will be canceled when ctrl+c is pressed
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Call createThing with the cancellable context
	err := CreatePlan(ctx)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Process exited with code: %d\n", exitErr.ExitCode())
		} else if errors.Is(err, context.Canceled) {
			fmt.Println("Operation was canceled (received interrupt signal)")
			// Perform any cleanup needed after cancellation
			return
		} else {
			fmt.Printf("Error occurred: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("Successfully completed the operation")
}

func CreatePlan(ctx context.Context) error {
	binary = "terraform"
	log.Infof("Using %s binary: ", binary)
	execPath, err := safeexec.LookPath(binary)
	log.Infof("execPath is %s: ", execPath)
	if err != nil {
		log.Error(err)
		return err
	}

	workingDir = filepath.Base(".")
	log.Infof("Working dir is %s: ", workingDir)

	// Create a new context that isn't tied to signals
	cmdCtx, cmdCancel := context.WithCancel(context.Background())
	defer cmdCancel()

	// Create the plan command directly
	planCmd := exec.CommandContext(cmdCtx, execPath, "plan", "-out=plan.out", "-no-color")
	planCmd.Dir = workingDir
	planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr

	// Create a channel for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel for command completion
	doneChan := make(chan error, 1)

	log.Info("Starting Terraform plan...")

	// Start the command
	err = planCmd.Start()
	if err != nil {
		log.Errorf("Failed to start plan: %v", err)
		signal.Stop(sigChan)
		return err
	}

	// Store the process ID for signaling
	procID := planCmd.Process.Pid
	log.Infof("Terraform process running with PID: %d", procID)

	// Wait for the command in a goroutine
	go func() {
		doneChan <- planCmd.Wait()
	}()

	// Wait for either completion or signal
	select {
	case <-sigChan:
		log.Info("Interrupt received")

		// Reset signal handling to default
		signal.Reset(os.Interrupt, syscall.SIGTERM)

		// We DON'T need to manually send a signal to Terraform,
		// it should have already received it from the OS

		log.Info("Waiting for Terraform to exit gracefully...")

		// Wait up to 30 seconds for the command to complete naturally
		select {
		case err := <-doneChan:
			log.Info("Terraform exited after signal")
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					log.Infof("Exit code: %d", exitErr.ExitCode())
				}
			}
		case <-time.After(30 * time.Second):
			log.Warn("Terraform didn't exit in time, forcing termination")
			cmdCancel()
		}

		return context.Canceled

	case err := <-doneChan:
		// Command completed normally
		signal.Stop(sigChan)
		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				log.Infof("Command exited with non-zero status: %d", exitErr.ExitCode())
			} else {
				log.Error("Plan error:", err)
			}
			return err
		}

		log.Info("Plan completed successfully")
		return nil
	}
}
