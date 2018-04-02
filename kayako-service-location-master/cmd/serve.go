package cmd

import (
	"github.com/kayako/service-location/server"
	"github.com/spf13/cobra"
)

// iface is the interface and port to mount the
// HTTP server
var iface string

// serveCmd respresents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the location service over HTTP",
	Long:  "Start an HTTP server to serve the geo location data in JSON format",
	RunE:  mountServer,
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&iface, "http-addr", ":80", "An interface and port to mount the HTTP server [Required]")
}

// mountServer mounts the HTTP server on provided interface
func mountServer(*cobra.Command, []string) error {
	return server.Mount(iface)
}
