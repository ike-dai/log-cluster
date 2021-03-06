package logcluster

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"bufio"
	"strconv"
	"sync"
	"bytes"
	"gopkg.in/jdkato/prose.v2"
	"github.com/ike-dai/wego/builder"
	"github.com/ike-dai/wego/model/word2vec"
	"github.com/cipepser/goClustering/ward"
	"github.com/cheggaaa/pb/v3"
)

type LogClusterClient struct {
	FileData []byte
	Limit int
	Threshold float64
}

type LogCluster struct {
	MemberCount int `json:"count"`
	Logs []string `json:"log"`
	CauseComment string `json:"cause"`
	ActionPlan string `json:"action"`
}

func New(fileData []byte, limit int, threshold float64) LogClusterClient {
	return LogClusterClient{fileData, limit, threshold}
}

func (c *LogClusterClient) GetCluster() (clusters []LogCluster) {
	logDataSlice := readLog(c.FileData, c.Limit)
	logData := strings.Join(logDataSlice, "\n")
	vectors := calcVector(logData)
	if len(vectors) == 0 {
		return clusters
	}
	matrix := make([][]float64, 0)
	for _, logRow := range logDataSlice {
		v := getLogVector(logRow, vectors)
		matrix = append(matrix, v)
	}
	tree := execClustering(matrix)
	roots := getClusterRootNodesNo(tree, c.Threshold)
	for _, r := range roots {
		cluster := LogCluster{}
		clusterMember := getChildNodes(r, tree)
		cluster.MemberCount = len(clusterMember)
		for _, logno := range clusterMember {
			cluster.Logs = append(cluster.Logs, logDataSlice[logno])
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

func removeDateString(logStr string) string {
	timeReg1 := regexp.MustCompile(`[0-9]{2,4}(-|\/)[0-9]{2}(-|\/)[0-9]{2}`)
	timeReg2 := regexp.MustCompile(`[0-9]{2}:[0-9]{2}:[0-9]{2}`)
	logStr = timeReg1.ReplaceAllString(logStr, "")
	logStr = timeReg2.ReplaceAllString(logStr, "")
	return removeSymbolString(logStr)
}

func removeSymbolString(logStr string) string {
	symbolReg := regexp.MustCompile(`(,|-|;|:|<|>|\?|!|\(|\)|\+|=|\.|\[|\])`)
	return symbolReg.ReplaceAllString(logStr, "")
}

func strListToFloatList(strList []string) (floatList []float64) {
	for _, value := range strList {
		if len(value) == 0 {
			continue
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("ParseFloat error: %v\n", err)
			return floatList
		}
		floatList = append(floatList, floatValue)
	}
	return floatList
}


func pickupImportantWords(rawLogData string) (pickupedLogData string) {
	rawLogData = removeDateString(rawLogData)
	var pickup []string
	doc, err := prose.NewDocument(rawLogData)
	if err != nil {
		log.Fatal(err)
		return pickupedLogData
	}
	for _, tok := range doc.Tokens() {
		if !strings.HasPrefix(tok.Tag, "CD") && !strings.HasPrefix(tok.Tag, "SYM") && !strings.HasPrefix(tok.Tag, "LS") {
			pickup = append(pickup, tok.Text)
		}
	}
	return strings.Join(pickup, " ")
}

func getLogLineCount(fileData []byte) (count int) {
	scanner := bufio.NewScanner(bytes.NewReader(fileData))
	count = 0
	for scanner.Scan() {
		count += 1
	}
	return count
}

func readLog(fileData []byte, limit int) (logData []string) {
	fmt.Printf("### Start read log & Morphological Analysis ###\n")

	lineCount := getLogLineCount(fileData)
	bar := pb.StartNew(lineCount)
	// ch := make(chan string, lineCount)
	scanner := bufio.NewScanner(bytes.NewReader(fileData))
	wg := new(sync.WaitGroup) //並行処理のため、WaitGroupを使ってloopを待つように。
	semaphore := make(chan struct{}, limit) //同時並行処理件数の制御用セマフォ
	for scanner.Scan() {
		wg.Add(1) //goroutineに入る前にインクリメントしてgoroutineが終わればデクリメントされるように。最終的に全部が終わった時点で処理が抜けるようにする。
		semaphore <- struct{}{}
		line := scanner.Text()
		go func() {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			importantWords := pickupImportantWords(string(line))
			if len(importantWords) > 0 {
				logData = append(logData, importantWords)
			}
			bar.Increment()
		}()
	}
	wg.Wait()
	bar.Finish()
	fmt.Printf("___ Finish read log & Morphological Analysis ___\n")
	return logData
}

func calcVector(source string) (wordVector map[string][]float64) {
	fmt.Printf("### Start word2vec Analysis ###\n")

	b := builder.NewWord2vecBuilder()

	b.Dimension(10).
		Window(5).
		Model(word2vec.CBOW).
		Optimizer(word2vec.NEGATIVE_SAMPLING).
		NegativeSampleSize(5)

	m, err := b.Build()
	if err != nil {
		// Failed to build word2vec.
	}
	input := strings.NewReader(source)
	err = m.Train(input)
	if err != nil {
		fmt.Println(err)
		return wordVector
	}
	wordVector, err = m.Get()
	if err != nil {
		fmt.Println(err)
		return wordVector
	}
	fmt.Printf("___ Finish word2vec Analysis ___\n")
	return wordVector
}

func getLogVector(logData string, wordVec map[string][]float64) []float64 {
	// 1行のログに含まれる単語のベクトルを平均して1行のログの特徴ベクトルを返す
	splittedLog := strings.Split(logData, " ")
	wordCount := len(splittedLog)
	dimention := 10
	sumVector := make([]float64, dimention)
	vector := make([]float64, dimention)
	for _, word := range splittedLog {
		vec := wordVec[word]
		if len(vec) == 0 {
			return vector
		}
		for i, _ := range(make([]int, dimention)) {
			sumVector[i] += vec[i]
		}
	}
	for i, _ := range(make([]int, dimention)) {
		vector[i] = sumVector[i] / float64(wordCount)
	}
	return vector
}

func execClustering(matrix [][]float64) ward.Tree {
	tree := ward.Ward(matrix)
	return tree
}

func getChildNodes(parentNodeNo int, tree ward.Tree) (childs []int) {
	left := tree[parentNodeNo].Left
	right := tree[parentNodeNo].Right
	//fmt.Printf(">>>>> Parent: %d, Left: %d, Right: %d\n", parentNodeNo, left, right)
	if left != -1 {
		childs = append(childs, getChildNodes(left, tree)...)
	}
	if right != -1 {
		childs = append(childs, getChildNodes(right, tree)...)
	}
	if left == -1 && right == -1 {
		childs = append(childs, parentNodeNo)
	}
	return childs
}

func includes(i int, list []int) bool {
	for _, n := range list {
		if i == n {
			return true
		}
	}
	return false
}

func getClusterRootNodesNo(tree ward.Tree, threshold float64) (roots []int) {
	// treeの中からしきい値以上のdistanceのクラスタノードを抽出し、その中でも最下層のレイヤに該当するクラスタノードの番号一覧を取り出す。
	var removedNodes []int
	var parentNodes ward.Tree
	for i, node := range tree {
		if node.GetDist() < threshold {
			removedNodes = append(removedNodes, i)
		} else {
			parentNodes = append(parentNodes, node)
		}
	}

	for _, parent := range parentNodes {
		if includes(parent.Left, removedNodes) {
			roots = append(roots, parent.Left)
		}
		if includes(parent.Right, removedNodes) {
			roots = append(roots, parent.Right)
		}
	}
	return roots
}

