package cmd

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	binary      string
	cfgFile     string
	configDir   string
	cwd         string
	homeDir     string
	mdFile      string
	planFile    string
	projectName string
)

type Project struct {
	Project map[string]configParams
}

type configParams struct {
	Binary   string `toml:"binary" comment:"(type: string) 'tofu' or 'terraform'. Must exist on your $PATH."`
	MdFile   string `toml:"mdFile" comment:"(type: string)The name of the Markdown file created by 'gh tp'."`
	PlanFile string `toml:"planFile" comment:"(type: string) The name of the plan file created by 'gh tp'."`
	Verbose  bool   `toml:"verbose" comment:"(type: bool) Enable Verbose Logging. Default is false."`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		homeDir, configDir = getDirectories()
		repo := getRepo()
		dir, cwderr := getCWD()
		scwd := strings.Split(dir, "/")
		cwd = scwd[len(scwd)-1]
		if cwderr != nil {
			log.Fatalf("Error: %s", cwderr)
		}

		form := huh.NewForm(
			huh.NewGroup(

				huh.NewSelect[string]().
					Title("What is the name of your project?").
					Options(
						huh.NewOption("Repository: "+repo, repo).Selected(true),
						huh.NewOption("Project Root: "+cwd, cwd),
					).
					Value(&projectName),

				huh.NewSelect[string]().
					Title("Where would you like to save your .tp.toml config file?").
					Options(
						huh.NewOption("Home Config Directory: "+configDir+"/.tp.toml", configDir).Selected(true),
						huh.NewOption("Home Directory: "+homeDir+"/.tp.toml", homeDir),
						huh.NewOption("Project Root: "+cwd+"/.tp.toml", cwd),
					).
					Value(&cfgFile),

				huh.NewSelect[string]().
					Title("Choose your binary").
					Options(
						huh.NewOption("OpenTofu", "tofu").Selected(true),
						huh.NewOption("Terraform", "terraform"),
					).
					Value(&binary),

				huh.NewInput().
					Title("What do you want to name your plan's output file? ").
					PlaceholderFunc(func() string {
						switch planFile {
						case "":
							return ""
						default:
							return "plan.out"
						}
					}, &planFile).
					Inline(true).
					Value(&planFile),

				huh.NewInput().
					Title("What do you want to name your Markdown file?  ").
					Placeholder("plan.md").
					Value(&mdFile),
			),
		).WithTheme(huh.ThemeBase())

		runerr := form.Run()
		if runerr != nil {
			log.Fatal(runerr)
		}

		conf := make(map[string]interface{})

		conf[projectName] = configParams{
			Binary:   binary,
			PlanFile: planFile,
			MdFile:   mdFile,
			Verbose:  false,
		}

		genConfig(conf)
		//log.Infof("Project name is: %s", projectName)
		//log.Infof("Binary is: %s", binary)
		//log.Infof("Plan out file is: %s", planFile)
		//log.Infof("Markdown file is: %s", mdFile)
		//log.Infof("File saved in: %s", cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}
