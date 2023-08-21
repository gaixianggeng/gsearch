package index

import (
	"fmt"
	"gsearch/conf"
	"gsearch/internal/meta"
	"gsearch/internal/storage"
	"gsearch/pkg/utils/log"
	"strconv"
	"strings"
)

// Run 索引写入入口
func Run(meta *meta.Profile, conf *conf.Config) {

	log.Infof("index run...")
	index, err := NewIndexEngine(meta, conf)
	if err != nil {
		panic(err)
	}
	defer index.Close()

	// TODO: 这样调用不合适
	addDoc(index)
	log.Infof("index run end")
}

func addDoc(in *Index) {
	log.Infof("addDoc start")
	docList := readFiles(in.Conf.Source.Files)
	go in.Scheduler.Merge()
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
	}
	// 读取结束 写入磁盘
	in.Flush(true)
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
