package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Selahattinn/picus-tcp-message/pkg/server"
	"github.com/Selahattinn/picus-tcp-message/pkg/version"
	"github.com/sirupsen/logrus"
)

var (
	versionFlag = flag.Bool("version", false, "Show version information.")
	debugFlag   = flag.Bool("debug", true, "Show debug information.")
	addrFlag    = flag.String("addr", ":7654", "Listen address")
	logFileFlag = flag.String("log.file", "tcp-message-server.log", "Path to the log file.")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Show version information
	if *versionFlag {
		fmt.Fprintln(os.Stdout, version.Print("tcp-message-server"))
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

	server := server.New(&server.Config{
		Addr: *addrFlag,
	})

	// start server
	go func() {
		logrus.Info("starting server at ", *addrFlag)
		if err := server.Start(); err != nil {
			fmt.Println(err)
			logrus.Fatal(err)
		}

	}()

}
