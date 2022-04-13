package main

import (
	"brain/internal/index"
	"brain/internal/storage"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func run() {
	engine, err := index.NewIndexEngine()
	if err != nil {
		panic(err)
	}
	addDoc(engine)
}

func addDoc(engine *index.Engine) {
	docList := readFile("../../data/source.csv")
	for _, item := range docList {
		log.Debug(item)
		doc, err := doc2Struct(item)
		if err != nil {
			log.Errorf("doc2Struct err: %v", err)
			break
		}
		err = engine.AddDocument(doc)
		if err != nil {
			log.Errorf("AddDocument err: %v", err)
			break
		}
	}
}

func doc2Struct(docStr string) (*storage.Document, error) {

	d := strings.Split(docStr, "\t")

	log.Debugf("doc:%s,len:%d", docStr, len(d))
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
		log.Info("readFile err: %v", "docList is empty")
		return nil
	}
	return docList
}
