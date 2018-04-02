package cmd

import (
	"errors"

	"github.com/kayako/service-location/geo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dbNames are the full path to database files
type dbNames []string

// String return all database files path in a comma separated string
func (n dbNames) String() string {
	s := ""
	for _, v := range n {
		s += ", " + v
	}

	return s
}

// Set populates the dbNames
func (n *dbNames) Set(s string) error {
	*n = append(*n, s)
	return nil
}

// Type returns the flag type for dbNames
func (n *dbNames) Type() string {
	return "string"
}

var (
	// Path to the database with City data
	databases dbNames

	// Minimum level to logs to emit
	logLevel string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "location",
	Short:             "IP Intelligence Discovery",
	Long:              "Look up the geo location data for any IP address",
	PersistentPreRunE: initializeDB,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "error", "Minimum level of logs to print [Optional]")
	RootCmd.PersistentFlags().Var(&databases, "db", "Full path to the database file")
}

// Execute adds all child commands to the root command and
// sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

// Boot the command configuration and set up flags
func Boot() {
	setLogLevel()
}

// initializeDB attempts to connect with Geo database
func initializeDB(*cobra.Command, []string) error {
	if len(databases) == 0 {
		return errors.New("No database file specified")
	}

	return geo.Connect(databases...)
}

// Parse and set the log level
func setLogLevel() {
	if logLevel == "" {
		log.Info("no log level specified, only errors will be logged")
		logLevel = "error"
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Info("invalid log level %s, only errors will be logged", logLevel)
		log.Info("valid log levels are 'debug', 'info', 'warn' and 'error'")
		level = log.ErrorLevel
	}

	log.SetLevel(level)
}
