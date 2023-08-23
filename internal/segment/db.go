package segment

// 正排和倒排相关操作

import (
	"fmt"
	"gsearch/conf"
	"gsearch/internal/storage"
	"gsearch/pkg/utils/log"
)

// CreateNewInvertedIndex 创建新的倒排索引
func CreateNewInvertedIndex(token string, docCount uint64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.DocPositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}

// 读取对应的segment文件下的db
func dbInit(segID SegID, conf *conf.Config) (*storage.InvertedDB, *storage.ForwardDB) {
	if segID < 0 {
		log.Fatalf("dbInit segID:%d < 0", segID)
	}
	term, inverted, forward := GetDBName(conf, segID)
	log.Debugf(
		"index:[termName:%s,invertedName:%s,forwardName:%s]",
		term,
		inverted,
		forward,
	)
	return storage.NewInvertedDB(term, inverted), storage.NewForwardDB(forward)
}

// GetDBName 获取db的路径+名称
func GetDBName(conf *conf.Config, segID SegID) (string, string, string) {
	termName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, TermDBSuffix)
	invertedName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, InvertedDBSuffix)
	forwardName = fmt.Sprintf("%s%d%s", conf.Storage.Path, segID, ForwardDBSuffix)
	return termName, invertedName, forwardName
}
