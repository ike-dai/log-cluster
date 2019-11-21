package formatter 

import (
	"fmt"
	"encoding/json"
	"github.com/ike-dai/log-cluster/logcluster"
)

type JsonFormatter struct {
	Clusters []logcluster.LogCluster `json:"clusters"`
}

func NewJsonFormatter(clusters []logcluster.LogCluster) {
	f.Clusters = clusters
}

func (f *JsonFormatter)Output() {
	bytes, err := json.Marshal(f.Clusters)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
}