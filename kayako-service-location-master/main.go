package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/kayako/service-location/cmd"
	"github.com/kayako/service-location/geo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Initialise the application
func init() {
	cobra.OnInitialize(cmd.Boot)
}

func main() {
	defer handlePanic()
	defer geo.Disconnect()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

// Try to recover from crashes and shut down gracefully
func handlePanic() {
	t := time.Now()
	t, _ = time.Parse(time.RFC3339, t.String())
	fname := "geo-crash-" + t.String()

	if p := recover(); p != nil {
		log.Error(`A fatal error occurred, service will now shut down.
            All the information about crash will be printed here and also 
            stored in the log file mentioned below. Please provide the log
            file along with the bug report.`)

		// make a container to hold 8 KBs
		stack := make([]byte, 8192)
		stack = stack[:runtime.Stack(stack, false)]
		f := "CRASH: %s\n%s"
		trace := fmt.Sprintf(f, p, stack)

		err := ioutil.WriteFile(fname, []byte(trace), 0644)
		if err != nil {
			log.Errorf("Failed to write logs file '%s', %s", fname, err.Error())
		} else {
			log.Info("Logs written to file %s", fname)
		}

		log.Fatalf(trace)
	}
}
