package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/Selahattinn/picus-tcp-message/pkg/server"
	"github.com/Selahattinn/picus-tcp-message/pkg/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	configFileFlag = flag.String("config.file", "config.yml", "Path to the configuration file.")
	versionFlag    = flag.Bool("version", false, "Show version information.")
	debugFlag      = flag.Bool("debug", true, "Show debug information.")
	logFileFlag    = flag.String("log.file", "tcp-message-server.log", "Path to the log file.")
)

func main() {
	// getting flag values
	flag.Parse()

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

	// logrus log file setted
	logrus.SetOutput(logFile)

	// Load configuration file
	data, err := ioutil.ReadFile(*configFileFlag)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load configuration")
	}
	var cfg server.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Could not load configuration")
	}

	// instantiate a server
	s := server.NewServer(&cfg)
	logrus.Info("created new server")

	// run server in separate go routine
	go s.Run()

	// start listening...
	listener, err := net.Listen("tcp", s.Config.ListenAddress)

	if err != nil {
		logrus.WithError(err).Fatal("unable to start server")
	} else {
		logrus.Info("listening to port ", s.Config.ListenAddress)
	}

	defer listener.Close()
	// sometime later.. close listener

	// continuously accept new connections
	for {

		conn, err := listener.Accept()
		if err != nil {
			logrus.WithError(err).Info("failed to accept connection")
			continue
		}

		go s.NewClient(conn)
		logrus.Info("added new client : %s", conn.RemoteAddr().String())
	}

}
