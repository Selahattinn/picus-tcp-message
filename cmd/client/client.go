package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Selahattinn/picus-tcp-message/pkg/version"
	"github.com/sirupsen/logrus"
)

var (
	versionFlag    = flag.Bool("version", false, "Show version information.")
	debugFlag      = flag.Bool("debug", false, "Show debug information.")
	logFileFlag    = flag.String("log.file", "tcp-message-client.log", "Path to the log file.")
	serverAddrFlag = flag.String("server", "127.0.0.1:7654", "Server Address")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Show version information
	if *versionFlag {
		fmt.Fprintln(os.Stdout, version.Print("tcp-message-client"))
		os.Exit(0)
	}

	// Log settings
	if *debugFlag {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetReportCaller(false)
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logFile, err := os.OpenFile(*logFileFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.WithError(err).Fatal("Could not open log file")
	}

	logrus.SetOutput(logFile)

}
