package meta

import (
	"encoding/json"
	"fmt"
	"gsearch/conf"
	"gsearch/internal/segment"
	"gsearch/pkg/utils/file"
	"gsearch/pkg/utils/log"
	"os"
)

// ParseProfile 解析数据
func ParseProfile(c *conf.Config) (*Profile, error) {
	metaFile = c.Storage.Path + metaFile
	log.Infof("metaFile:%s", metaFile)
	// 文件不存在表示没有相关数据 第一次创建
	if !file.IsExist(metaFile) {
		log.Debugf("segMetaFile:%s not exist", metaFile)
		f, err := os.Create(metaFile)
		if err != nil {
			return nil, fmt.Errorf("create segmentsGenFile err: %v", err)
		}
		f.Close()
		m := &Profile{
			Version: c.Version,
			path:    metaFile,
			SegMeta: &segment.SegMeta{
				NextSeg:  0,
				SegCount: 0,
				SegInfo:  make(map[segment.SegID]*segment.SegInfo, 0),
			},
			//TODO: 初始化读取正排数据
			IndexCount: 0,
		}
		err = m.writeMeta()
		if err != nil {
			return nil, fmt.Errorf("writeSeg err: %v", err)
		}
		return m, nil
	}

	return readMeta(metaFile)
}

func readMeta(metaFile string) (*Profile, error) {
	metaByte, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, fmt.Errorf("read file err: %v", err)
	}
	h := new(Profile)
	err = json.Unmarshal(metaByte, &h)
	if err != nil {
		return nil, fmt.Errorf("ParseHeader err: %v", err)
	}
	log.Debugf("seg header :%v", h)
	h.path = metaFile
	return h, nil
}
