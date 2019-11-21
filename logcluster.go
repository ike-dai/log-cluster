package main

import (
	"flag"
	"fmt"
	"github.com/ike-dai/log-cluster/logcluster"
	"github.com/ike-dai/log-cluster/formatter"
)

func main() {
	var logfile string
	var threshold float64
	var limit int
	var output string
	flag.StringVar(&logfile, "logfile", "./test.log", "Analyze target log")
	flag.Float64Var(&threshold, "threshold", 0.001, "Set cluster threshold")
	flag.IntVar(&limit, "limit", 5, "Set pararell size limit")
	flag.StringVar(&output, "output", "table", "Set output type (table/json)")
	flag.Parse()
	client := logcluster.New(logfile, limit, threshold)
	clusters := client.GetCluster()
	if output == "table" {
		output.NewTableFormatter(clusters)
	}
	if output == "json" {
		output.NewJsonFormatter(clusters)
	}
	output.Output()
	//fmt.Println(clusters)
}
