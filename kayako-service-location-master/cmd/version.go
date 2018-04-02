package cmd

import (
	"fmt"
	"os"

	"github.com/kayako/service-location/geo"
	"github.com/spf13/cobra"
)

// versionCmd respresents the serve command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "VCS version of the app",
	RunE:  showVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

// mountServer mounts the HTTP server on provided interface
func showVersion(*cobra.Command, []string) error {
	tmp := "commit: %s"
	tmp += "branch: %s"
	tmp += "build time: %s"

	fmt.Fprintf(os.Stdout, tmp, geo.AppVersion, geo.BuildBranch, geo.BuildTime)
	return nil
}
