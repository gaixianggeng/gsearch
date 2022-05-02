package segment

// 正排和倒排相关操作

import (
	"doraemon/conf"
	"doraemon/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

//InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      uint64
	PositionCount uint64 // 查询使用，写入的时候暂时用不到
	TermValues    *storage.TermValue
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, termValue *storage.TermValue) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = termValue.DocCount
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
