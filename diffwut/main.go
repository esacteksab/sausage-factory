package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/rogpeppe/go-internal/diff"
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

func main() {
	sourcePath := "../../projectTemplates/go-project-template/Makefile"
	destPath := "../aa/Makefile"
	fmt.Printf("Scanning file: %s\n", destPath)

	// Debug scan to find all relevant lines
	debugScanFile(destPath)

	// Find the custom sections
	sections, err := findCustomSections(destPath)
	if err != nil {
		log.Fatal(err)
	}

	// Check if we found any sections
	if len(sections) == 0 {
		fmt.Println("No custom sections found")
		return
	}

	// Print information about each section
	for i, section := range sections {
		fmt.Printf("Section %d:\n", i+1)
		fmt.Printf("  Start line: %d\n", section.StartLine)
		fmt.Printf("  End line: %d\n", section.EndLine)
		fmt.Printf("  Start marker: '%s'\n", section.MarkerStart)
		fmt.Printf("  Stop marker: '%s'\n", section.MarkerStop)
		fmt.Printf("  Content (%d lines):\n", len(section.Content))
		for _, line := range section.Content {
			fmt.Printf("    %s\n", line)
		}
		fmt.Println()
	}

	// Then sync the files
	err = SyncFiles(sourcePath, destPath, true) // true for dry-run
	if err != nil {
		log.Fatal("Error syncing files:", err)
	}
}

// debugScanFile prints any lines containing our marker prefix for debugging
func debugScanFile(destPath string) error {
	file, err := os.Open(filepath.Clean(destPath))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	fmt.Println("=== Debug: Marker Lines ===")
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if strings.TrimSpace(line) == markerStart ||
			strings.TrimSpace(line) == markerStop ||
			strings.Contains(line, "# annoyed-aardvark") {
			fmt.Printf("Line %d: '%s'\n", lineNum, line)
		}
	}
	fmt.Println("=========================")

	return scanner.Err()
}

func findCustomSections(destPath string) ([]CustomSection, error) {
	sections := []CustomSection{}

	file, err := os.Open(filepath.Clean(destPath))
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

// SyncFiles synchronizes source file to destination, preserving custom sections
func SyncFiles(sourcePath, destPath string, dryRun bool) error {
	// Find custom sections in destination file
	sections, err := findCustomSections(destPath)
	if err != nil {
		return fmt.Errorf("finding custom sections: %w", err)
	}

	// Read destination file as bytes for diffing
	destBytes, err := os.ReadFile(destPath)
	if err != nil {
		return fmt.Errorf("reading destination file: %w", err)
	}

	// Create merged content
	mergedContent, err := mergeFiles(sourcePath, sections)
	if err != nil {
		return fmt.Errorf("merging files: %w", err)
	}

	// Convert merged content to bytes for diffing
	mergedBytes := []byte(mergedContent)

	// Compare current destination with merged result
	diffs := diff.Diff(string(destBytes), destBytes, string(mergedBytes), mergedBytes)
	if len(diffs) == 0 {
		fmt.Println("Files are already in sync.")
		return nil
	}

	// Format the diff output
	var buf bytes.Buffer
	for _, diffChunk := range diffs {
		buf.WriteString(string(diffChunk)) // Changed from buf.Write(diffChunk)
	}
	diffStr := formatDiffOutput(buf.String())

	fmt.Println("Changes to be applied:")
	fmt.Println(diffStr)

	if dryRun {
		fmt.Println("Dry run - no changes made.")
		return nil
	}

	// Confirm changes
	fmt.Print("Apply these changes? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		// Write the merged content back to destination
		return os.WriteFile(destPath, mergedBytes, 0644)
	}

	fmt.Println("Operation cancelled.")
	return nil
}

// formatDiffOutput formats diff output with colored syntax
func formatDiffOutput(diffText string) string {
	// Define styles
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))   // Green
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red

	var filteredLines []string
	lines := strings.Split(diffText, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			filteredLines = append(filteredLines, addStyle.Render(line))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			filteredLines = append(filteredLines, removeStyle.Render(line))
		}
	}

	return strings.Join(filteredLines, "\n")
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
