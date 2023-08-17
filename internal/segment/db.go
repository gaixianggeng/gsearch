package segment

// 正排和倒排相关操作

import (
	"doraemon/conf"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string             // 词元
	PostingsList  *PostingsList      // 文档编号的序列
	DocCount      uint64             // 词元关联的文档数量
	PositionCount uint64             // 词元在所有文档中出现的次数 查询使用,用于计算相关性，写入的时候暂时用不到
	TermValues    *storage.TermValue // 存储的doc_count、offset、size
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.PositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}

// 读取对应的segment文件下的db
func dbInit(segID SegID, conf *conf.Config) (*storage.InvertedDB, *storage.ForwardDB) {

	if segID < 0 {
		log.Fatalf("dbInit segID:%d < 0", segID)
	}
	termName, invertedName, forwardName := GetDBName(conf, segID)
	log.Debugf(
		"index:[termName:%s,invertedName:%s,forwardName:%s]",
		termName,
		invertedName,
		forwardName,
	)
	return storage.NewInvertedDB(termName, invertedName), storage.NewForwardDB(forwardName)
}

// GetDBName 获取db的路径+名称
func GetDBName(conf *conf.Config, segID SegID) (string, string, string) {
	termName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, TermDBSuffix)
	invertedName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, InvertedDBSuffix)
	forwardName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, ForwardDBSuffix)
	return termName, invertedName, forwardName
}
