package formatter

import (
	"strconv"
	"os"
	"github.com/ike-dai/log-cluster/logcluster"
	"github.com/olekukonko/tablewriter"
)

type TableFormatter struct {
	TableData [][]string
}

func NewTableFormatter(clusters []logcluster.LogCluster) {
	for i, cluster := range clusters {
		for _, log := range cluster.Logs {
			f.TableData = append(f.TableData, []string{strconv.Itoa(i), log})
		}
	}
	return nil
}

func (f *TableFormatter)Output() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"no.", "log data"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(f.TableData)
	table.Render()
}