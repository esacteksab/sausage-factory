// SPDX-License-Identifier: MIT
package internal

import (
	"bufio"
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
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/rogpeppe/go-internal/diff"

	"github.com/esacteksab/annoyed-aardvark/internal/logger"
)

const (
	markerStart = "# annoyed-aardvark start"
	markerStop  = "# annoyed-aardvark stop"
)

// CustomSection represents a block of content to preserve
type CustomSection struct {
	StartLine   int
	EndLine     int
	Content     []string
	MarkerStart string
	MarkerStop  string
}

// CopyFile will copy files from the source path to the destination path.
func CopyFile(sourcePath string, files []string) {
	for _, file := range files {
		logger.Debugf("Copying file from %s: %s", sourcePath, file)
		sourcePathFile := sourcePath + "/" + file
		logger.Debug(sourcePathFile)
		// Open the source file for reading
		sourceFile, err := os.Open(sourcePathFile)
		if err != nil {
			logger.Debug("I'm dying here Jim!")
			logger.Fatal("Error", "error", err.Error())
		}
		defer sourceFile.Close()

		// We need to create the destination directory if it doesn't exist
		destPathDir := filepath.Dir(file)

		// If it doesn't exist, create it
		if !DoesExist(destPathDir) {
			if err := os.MkdirAll(destPathDir, 0o750); err != nil { //nolint:mnd
				logger.Fatal("Error", "error", err.Error())
			}
		}

		// Checking to ensure source exists
		if DoesExist(sourcePathFile) {
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
			} else {
				// If the file does exist locally, we check for custom sections
				sections, err := findCustomSections(file)
				if err != nil {
					logger.Errorf("Error finding custom sections: %v", err)
				}

				if len(sections) > 0 {
					logger.Infof("🔍 Found %d custom sections in %s", len(sections), file)

					// Use smart merge to preserve custom sections
					err = SyncFiles(sourcePathFile, file, false)
					if err != nil {
						logger.Errorf("Error syncing files: %v", err)
					}
				} else {
					// No custom sections, check if files are identical
					same, err := SameFile(sourcePathFile, file)
					if err != nil {
						logger.Debugf("Error comparing files: %v", err)
					}

					// If they're not the same, show diff
					if !same {
						logger.Infof("❌ %s isn't the same...", file)
						diff, err := DiffFile(sourcePathFile, file)
						if err != nil {
							logger.Errorf("%v", err)
						}
						fmt.Print(diff + "\n")

						// Ask if user wants to overwrite
						logger.Info("Would you like to overwrite the local file? (y/N): ")
						var response string
						fmt.Scanln(&response)

						if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
							// Copy the source file over the destination
							destFile, err := os.Create(file)
							if err != nil {
								logger.Error("Error creating destination file", "error", err.Error())
								continue
							}

							// Reopen source file as we might have closed it
							srcFile, err := os.Open(sourcePathFile)
							if err != nil {
								logger.Error("Error opening source file", "error", err.Error())
								destFile.Close()
								continue
							}

							_, err = io.Copy(destFile, srcFile)
							srcFile.Close()
							destFile.Close()

							if err != nil {
								logger.Error("Error copying file", "error", err.Error())
							} else {
								logger.Infof("✅ Copied file: %s...", file)
							}
						} else {
							logger.Info("Skipping file")
						}
					} else {
						logger.Infof("✨ %s file is identical, no changes needed", file)
					}
				}
			}
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
	safeSource := filepath.Clean(source)
	f1, err := os.Open(safeSource)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	h1 := sha256.New()
	if _, err := io.Copy(h1, f1); err != nil {
		return false, err
	}

	safeDest := filepath.Clean(dest)
	f2, err := os.Open(safeDest)
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
	safeSource := filepath.Clean(source)
	src, err := os.ReadFile(safeSource)
	if err != nil {
		logger.Fatal("Error", "error", err.Error())
	}
	safeDest := filepath.Clean(dest)
	dst, err := os.ReadFile(safeDest)
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

	return strings.Join(filteredLines, "\n"), nil
}

// findCustomSections identifies the protected sections in a file
func findCustomSections(filePath string) ([]CustomSection, error) {
	sections := []CustomSection{}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentSection *CustomSection
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Use exact matching for markers, but allow for leading/trailing whitespace
		if strings.TrimSpace(line) == markerStart {
			// Start a new section
			currentSection = &CustomSection{
				StartLine:   lineNum,
				MarkerStart: line,
				Content:     []string{line},
			}
		} else if currentSection != nil && strings.TrimSpace(line) == markerStop {
			// End the current section
			currentSection.EndLine = lineNum
			currentSection.MarkerStop = line
			currentSection.Content = append(currentSection.Content, line)

			// Add the complete section to our list
			sections = append(sections, *currentSection)
			currentSection = nil
		} else if currentSection != nil {
			// Add line to current section
			currentSection.Content = append(currentSection.Content, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

func SyncFiles(sourcePath, destPath string, dryRun bool) error {
	// Find custom sections in destination file
	sections, err := findCustomSections(destPath)
	if err != nil {
		return fmt.Errorf("finding custom sections: %w", err)
	}

	// Read destination file content
	destContent, err := os.ReadFile(destPath)
	if err != nil {
		return fmt.Errorf("reading destination file: %w", err)
	}

	// Create merged content
	mergedContent, err := mergeFiles(sourcePath, sections)
	if err != nil {
		return fmt.Errorf("merging files: %w", err)
	}

	// Convert content to strings for diffing
	destContentStr := string(destContent)
	mergedContentStr := mergedContent

	// Compute edits and diff
	edits := myers.ComputeEdits("", destContentStr, mergedContentStr)
	unifiedDiff := gotextdiff.ToUnified(sourcePath, destPath, destContentStr, edits)

	// Capture the string representation
	var buf bytes.Buffer
	fmt.Fprint(&buf, unifiedDiff)

	// Format the output and print
	diffStr := formatDiffOutput(buf.String())

	// Check if there are actual changes
	if len(edits) == 0 {
		logger.Info("Files are already in sync.")
		return nil
	}

	logger.Info("Changes to be applied:")
	fmt.Println(diffStr)

	if dryRun {
		logger.Info("Dry run - no changes made.")
		return nil
	}

	// Confirm changes if not already confirmed
	fmt.Print("Apply these changes? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		// Write the merged content back to destination
		logger.Infof("✅ Updating file: %s...", destPath)
		return os.WriteFile(destPath, []byte(mergedContent), 0o640)
	}

	logger.Info("Operation cancelled.")
	return nil
}

// formatDiffOutput formats diff output with colored syntax
func formatDiffOutput(diffText string) string {
	// Define styles
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))   // Green
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

	var formattedLines []string
	lines := strings.Split(diffText, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
			// Added lines
			formattedLines = append(formattedLines, addStyle.Render(line))
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "---") {
			// Removed lines
			formattedLines = append(formattedLines, removeStyle.Render(line))
		} else {
			// Header lines and context lines - keep them as is
			formattedLines = append(formattedLines, line)
		}
	}

	return strings.Join(formattedLines, "\n")
}

// mergeFiles creates a merged version of source preserving provided custom sections
func mergeFiles(sourcePath string, sections []CustomSection) (string, error) {
	// Read source file content
	sourceLines, err := readLines(sourcePath)
	if err != nil {
		return "", fmt.Errorf("reading source file: %w", err)
	}

	// Track which sections we've already used
	usedSections := make(map[int]bool)

	// Create a result slice to build our output
	var result []string

	// Process source file line by line
	insideMarker := false
	for i := 0; i < len(sourceLines); i++ {
		line := sourceLines[i]

		// If we hit a start marker
		if strings.TrimSpace(line) == markerStart && !insideMarker {
			insideMarker = true

			// Look for a matching section in our custom sections
			foundMatch := false
			for sectionIdx, section := range sections {
				// If we haven't used this section yet
				if !usedSections[sectionIdx] {
					// Use the entire custom section instead of source
					result = append(result, section.Content...)
					usedSections[sectionIdx] = true
					foundMatch = true

					// Skip past the end marker in the source
					for j := i; j < len(sourceLines); j++ {
						if strings.TrimSpace(sourceLines[j]) == markerStop {
							i = j // Move loop counter to end marker
							break
						}
					}
					break
				}
			}

			// If no matching section found, use the source section
			if !foundMatch {
				result = append(result, line)
			}
		} else if strings.TrimSpace(line) == markerStop && insideMarker {
			// End of marker section
			insideMarker = false

			// If we didn't add a custom section, add the stop marker
			if result[len(result)-1] != line {
				result = append(result, line)
			}
		} else if !insideMarker {
			// Regular line outside a marker section
			result = append(result, line)
		}
		// Skip lines inside marker sections that weren't replaced
	}

	// Add any unused custom sections to the end
	for i, section := range sections {
		if !usedSections[i] {
			if len(result) > 0 && !strings.HasSuffix(result[len(result)-1], "\n") {
				result = append(result, "") // Add blank line for separation
			}
			result = append(result, section.Content...)
		}
	}

	return strings.Join(result, "\n"), nil
}

// readLines reads file lines with proper line ending handling
func readLines(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Normalize line endings
	text := string(content)
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// Split into lines
	lines := strings.Split(text, "\n")

	// Remove trailing empty line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}
