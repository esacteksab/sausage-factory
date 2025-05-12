package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/diff"
)

var conf = koanf.Conf{
	Delim:       ".",
	StrictMerge: true,
}

func main() {
	// Create a new koanf instance
	k := koanf.NewWithConf(conf)

	srcYAML := "../../projectTemplates/go-project-template/.golangci.yaml"
	// Load the destination YAML
	destYAML := ".golangci.yaml"
	// Load the destination YAML
	if err := k.Load(file.Provider(srcYAML), yaml.Parser()); err != nil {
		log.Fatalf("error loading destination config: %v", err)
	}

	srcBytes, err := yaml.Parser().Marshal(k.Raw())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	// Load and merge the source YAML
	if err := k.Load(file.Provider(destYAML), yaml.Parser()); err != nil {
		log.Fatalf("error loading source config: %v", err)
	}

	dstBytes, err := yaml.Parser().Marshal(k.Raw())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	differ(string(srcBytes), string(dstBytes), srcYAML, destYAML)
	// At this point, k contains the merged configuration
	// fmt.Println("Merged configuration:", k)
	pDiffer(string(srcBytes), string(dstBytes), srcYAML, destYAML)

	// Write the merged configuration to a new YAML file
	yamlBytes, err := yaml.Parser().Marshal(k.Raw())
	if err != nil {
		log.Fatalf("error marshaling config: %v", err)
	}

	if err := os.WriteFile(".golangci.yaml", yamlBytes, 0o640); err != nil {
		log.Fatalf("error writing config: %v", err)
	}

	fmt.Println("Successfully wrote merged configuration to .golangci.yaml")
}

func differ(srcFile, dstFile, srcFileName, dstFileName string) {
	src := srcFile
	dst := dstFile
	edits := myers.ComputeEdits("", src, dst)
	unifiedDiff := gotextdiff.ToUnified(srcFileName, dstFileName, src, edits)

	// Capture the string representation
	var buf bytes.Buffer
	fmt.Fprint(&buf, unifiedDiff)

	// Format the output and print
	formattedDiff := formatDiffOutput(buf.String())
	fmt.Println(formattedDiff)
}

func pDiffer(srcFile, dstFile, srcFileName, dstFileName string) {
	// Use a buffer to capture the output instead of writing to os.Stdout
	var buf bytes.Buffer

	// Write the diff to the buffer
	err := diff.Text(
		srcFileName,
		dstFileName,
		srcFile,
		dstFile,
		&buf,
	)
	if err != nil {
		fmt.Printf("Error generating diff: %v\n", err)
		return
	}

	// Get the diff as a string, format it, and print
	diffOutput := buf.String()
	formattedDiff := formatDiffOutput(diffOutput)
	fmt.Println(formattedDiff)
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
