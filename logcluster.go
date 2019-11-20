package main

import (
	"flag"
	"fmt"
	"github.com/ike-dai/log-cluster/logcluster"
)

func main() {
	var logfile string
	var threshold float64
	var limit int
	flag.StringVar(&logfile, "logfile", "./test.log", "Analyze target log")
	flag.Float64Var(&threshold, "threshold", 0.001, "Set cluster threshold")
	flag.IntVar(&limit, "limit", 5, "Set pararell size limit")
	flag.Parse()
	client := logcluster.New(logfile, limit, threshold)
	clusters := client.GetCluster()
	fmt.Println(clusters)
}
