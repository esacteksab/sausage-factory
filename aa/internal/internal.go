package internal

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rogpeppe/go-internal/diff"

	"github.com/esacteksab/annoyed-aardvark/internal/logger"
)

// CopyFile will copy files from the source path to the destination path.
func CopyFile(sourcePath string, files []string) {
	for _, file := range files {
		logger.Debugf("Copying file from %s: %s", sourcePath, file)
		sourcePath := sourcePath + "/" + file
		logger.Debug(sourcePath)
		// Open the source file for reading
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			logger.Debug("I'm dying here Jim!")
			logger.Fatal("Error", "error", err.Error())
		}
		defer sourceFile.Close()

		// We need to create the destination directory if it doesn't exist
		destPathDir := filepath.Dir(file)

		// If it doesn't exist, create it
		if !DoesExist(destPathDir) {
			if err := os.MkdirAll(destPathDir, os.ModePerm); err != nil {
				logger.Fatal("Error", "error", err.Error())
			}
		}

		// Checking to ensure source exists
		if DoesExist(sourcePath) {
			// If the file doesn't exist locally
			if !DoesExist(file) {
				// Presumably we are copying to the current directory, so this is `file`
				destinationFile, err := os.Create(file)
				if err != nil {
					logger.Debug("Are we here?")
					logger.Fatal("Error", "error", err.Error())
				}
				defer destinationFile.Close()

				// Copy the contents of the source file to the destination file
				logger.Infof("❌ %s does not exist, copying...", file)
				_, err = io.Copy(destinationFile, sourceFile)
				if err != nil {
					logger.Error("Error", "error", err.Error())
				}
				logger.Infof("⭐ Created file: %s...", file)
			}
			// If the file does exist locally, we check to see if they're the same
			same, err := SameFile(sourcePath, file)
			if err != nil {
				logger.Debug("Are we in here?")
				logger.Infof("Error: %s", err)
			}
			// If they're not the same, we say as much and copy from source
			if !same {
				logger.Infof("❌ %s isn't the same...", file)
				diff, err := DiffFile(sourcePath, file)
				if err != nil {
					logger.Errorf("%v", err)
				}
				fmt.Print(diff + "\n")
				// 				// Presumably we are copying to the current directory, so this is `file`
				// 				destinationFile, err := os.Create(file)
				// 				if err != nil {
				// 					logger.Debug("Are we here?")
				// 					logger.Error("Error", "error", err.Error())
				// 				}
				// 				defer destinationFile.Close()
				//
				// 				// Copy the contents of the source file to the destination file
				// 				_, err = io.Copy(destinationFile, sourceFile)
				// 				if err != nil {
				// 					logger.Error("Error", "error", err.Error())
				// 				}
				// 				logger.Infof("✅ Copied file: %s...", file)
			}
			// The file existed and had no drift
			//logger.Infof("✨ %s file exists...", file)
		}
	}
}

func DoesExist(path string) bool {
	// Check if the file or directory exists
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func SameFile(source, dest string) (bool, error) {
	// Get info about both files
	f1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	h1 := sha256.New()
	if _, err := io.Copy(h1, f1); err != nil {
		return false, err
	}

	f2, err := os.Open(dest)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	h2 := sha256.New()
	if _, err := io.Copy(h2, f2); err != nil {
		return false, err
	}

	// Compare the hashes
	return bytes.Equal(h1.Sum(nil), h2.Sum(nil)), nil
}

func DiffFile(source, dest string) (string, error) {
	// Read the contents of the files
	src, err := os.ReadFile(source)
	if err != nil {
		logger.Fatal("Error", "error", err.Error())
	}
	dst, err := os.ReadFile(dest)
	if err != nil {
		logger.Fatal("Error", "error", err.Error())
	}

	// Compute the diff between the two files
	diffs := diff.Diff(string(dst), dst, string(src), src)

	// Format the diff output
	var buf bytes.Buffer
	for _, diff := range diffs {
		buf.WriteString(string(diff))
	}

	// Define styles
	//addStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))   // Green
	//removeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")) // Red
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))   // Green
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

	var filteredLines []string
	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			filteredLines = append(filteredLines, addStyle.Render(line))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			filteredLines = append(filteredLines, removeStyle.Render(line))
		}
	}

	// fmt.Print(strings.Join(filteredLines, "\n"))
	// fmt.Print()
	// return buf.String(), nil
	return strings.Join(filteredLines, "\n"), nil
}
