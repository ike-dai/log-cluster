package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"encoding/json"
	"github.com/ike-dai/log-cluster/logcluster"
	"github.com/olekukonko/tablewriter"
)

func getTableData(clusters []logcluster.LogCluster) (tableData [][]string){
	for i, cluster := range clusters {
		for _, log := range cluster.Logs {
			tableData = append(tableData, []string{strconv.Itoa(i), log})
		}
	}
	fmt.Println(tableData)
	return tableData
}

func outputTableData(clusters []logcluster.LogCluster) {
	tableData := getTableData(clusters)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"no.", "log data"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(tableData)
	table.Render()
}


type JsonOutput struct {
	Clusters []logcluster.LogCluster `json:"clusters"`
}

func outputJsonData(clusters []logcluster.LogCluster) {
	output := JsonOutput{clusters}
	bytes, _ := json.Marshal(output)
	fmt.Println(string(bytes))
}


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
		outputTableData(clusters)
	}
	if output == "json" {
		outputJsonData(clusters)
	}
	//fmt.Println(clusters)
}
