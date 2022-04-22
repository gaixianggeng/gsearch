package main

import (
	"brain/internal/engine"
	"brain/internal/index"
	"brain/internal/storage"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	termDB     = "../data/term.db"
	invertedDB = "../data/inverted.db"
	forwardDB  = "../data/forward.db"
	sourceFile = "../data/source.csv"
)

// 入口
func run() {

	e := engine.NewEngine(termDB, invertedDB, forwardDB)
	index, err := index.NewIndexEngine(e)
	if err != nil {
		panic(err)
	}
	defer index.Close()
	addDoc(index)
}

func addDoc(in *index.Index) {
	docList := readFile(sourceFile)
	for _, item := range docList {
		doc, err := doc2Struct(item)
		if err != nil {
			log.Errorf("doc2Struct err: %v", err)
			doc = new(storage.Document)
		}
		log.Debugf("doc_id:%v,title:%s", doc.DocID, doc.Title)
		err = in.AddDocument(doc)
		if err != nil {
			log.Errorf("AddDocument err: %v", err)
			break
		}
		// // 达到阈值
		if len(in.PostingsHashBuf) > 0 && (in.BufCount > in.BufSize) {
			in.Flush()
		}
	}
	// 读取结束 写入磁盘
	in.Flush()
}

func doc2Struct(docStr string) (*storage.Document, error) {

	d := strings.Split(docStr, "\t")

	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}
	doc := new(storage.Document)

	docID, err := strconv.Atoi(d[0])
	if err != nil {
		return nil, err
	}
	doc.DocID = uint64(docID)
	doc.Title = d[1]
	doc.Body = d[2]
	return doc, nil
}

// 可改用流读取
func readFile(fileName string) []string {
	content, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	docList := strings.Split(string(content), "\n")
	if len(docList) == 0 {
		log.Infof("readFile err: %v", "docList is empty\n")
		return nil
	}
	return docList
}
