package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"github.com/ike-dai/log-cluster/logcluster"
	"github.com/ike-dai/log-cluster/formatter"
	"github.com/olekukonko/tablewriter"
)

func viewLog(logs []string) {
	var viewData [][]string
	for _, log := range logs {
		viewData = append(viewData, []string{log})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.AppendBulk(viewData)
	table.Render()
}

func main() {
	var logfile string
	var threshold float64
	var limit int
	var output string
	var interactive bool
	flag.StringVar(&logfile, "logfile", "./test.log", "Analyze target log")
	flag.Float64Var(&threshold, "threshold", 0.001, "Set cluster threshold")
	flag.IntVar(&limit, "limit", 5, "Set pararell size limit")
	flag.StringVar(&output, "output", "table", "Set output type (table/json)")
	flag.BoolVar(&interactive, "interactive", false, "Select interactive mode. (true/false)")
	flag.Parse()
	client := logcluster.New(logfile, limit, threshold)
	clusters := client.GetCluster()
	outputData := clusters
	if interactive {
		fmt.Printf("Clustered to [ %d ] log clusters. Please set cause and action \n", len(clusters))
		for _, cluster := range clusters {
			fmt.Println("Log cluster:")
			viewLog(cluster.Logs)
			fmt.Printf("Please input the 'CAUSE' of these logs\n")
			fmt.Printf("(Send by entering a empty line) >>> ")
			stdin := bufio.NewScanner(os.Stdin)
			var causeText string
			for stdin.Scan(){
				line := stdin.Text()
				causeText += line
				causeText += "\n"
				if len(line) == 0 {
					break
				}
			}
			cluster.CauseComment = causeText
			fmt.Printf("Please input the 'ACTION' for these logs\n")
			fmt.Printf("(Send by entering a empty line) >>> ")
			var actionText string
			for stdin.Scan(){
				line := stdin.Text()
				actionText += line
				actionText += "\n"
				if len(line) == 0 {
					break
				}
			}
			cluster.ActionPlan = actionText
			outputData = append(outputData, cluster)
		}
	}
	if output == "table" {
		f := formatter.NewTableFormatter(outputData)
		f.Output()

	}
	if output == "json" {
		f := formatter.NewJsonFormatter(outputData)
		f.Output()
	}
}
