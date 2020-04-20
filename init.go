package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	db                  *bolt.DB
	verbose, debug      bool
	repFile, cpuprofile string
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       false,
	})
	log.SetLevel(logrus.InfoLevel)

	d, err := setupDB("words.db")
	if err != nil {
		log.Fatal(err)
	}
	db = d
}

func initFlags() {
	flag.StringVar(&repFile, "report", "report.csv", "Report file path")
	flag.StringVar(&dict, "dict", "kryss.csv", "Kryss words dictionary")
	flag.StringVar(&dict, "d", "kryss.csv", "Kryss words dictionary")
	//d.Parse([]string{"-d"})

	flag.StringVar(&boards, "boards", "boards", "Kryss boards directory")
	flag.StringVar(&boards, "b", "boards", "Kryss boards directory")

	flag.BoolVar(&debug, "debug", false, "Debug program")
	flag.BoolVar(&verbose, "verbose", false, "Print interim messages")

	flag.StringVar(&cpuprofile,	"cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	if debug {
		log.SetReportCaller(true)
		log.SetLevel(logrus.TraceLevel)
	}

	if verbose {
		log.SetLevel(logrus.TraceLevel)
	}
}
