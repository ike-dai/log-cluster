package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"github.com/ike-dai/log-cluster/logcluster"
	"github.com/ike-dai/log-cluster/formatter"
)

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
	if interactive {
		for _, cluster := range clusters {
			fmt.Printf("Log cluster:\n%v", cluster.Logs)
			fmt.Printf("Please input the cause of these LOGs\n")
			fmt.Printf("(Send by entering empty line)>>>")
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
			fmt.Printf("Please input the action for  these LOGs\n")
			fmt.Printf("(Send by entering empty line)>>>")
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
		}
	}
	if output == "table" {
		f := formatter.NewTableFormatter(clusters)
		f.Output()

	}
	if output == "json" {
		f := formatter.NewJsonFormatter(clusters)
		f.Output()
	}
	//fmt.Println(clusters)
}
