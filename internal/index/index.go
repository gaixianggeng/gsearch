package index

import (
	"brain/internal/engine"
	"brain/internal/storage"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Index --
type Index struct {
	*engine.Engine
	IndexCount uint64
}

// AddDocument 添加文档
func (in *Index) AddDocument(doc *storage.Document) error {
	if doc.DocID > 0 && doc.Title != "" {
		err := in.ForwardDB.Add(doc)
		if err != nil {
			return fmt.Errorf("forward doc add err: %v", err)
		}
		err = in.Text2PostingsLists(doc.Title, doc.DocID)
		if err != nil {
			return fmt.Errorf("text2postingslists err: %v", err)
		}
		in.BufCount++
		in.IndexCount++
	}

	log.Debugf("start storage...%v,len:%d", in.PostingsHashBuf, len(in.PostingsHashBuf))

	// 落盘操作 title = ""表示文件读取结束
	if len(in.PostingsHashBuf) > 0 && (in.BufCount > in.BufSize || doc.Title == "") {

		for token, invertedIndex := range in.PostingsHashBuf {

			log.Debugf("token:%s,invertedIndex:%v\n", token, invertedIndex)
			err := in.updatePostings(invertedIndex)
			if err != nil {
				log.Errorf("updatePostings err: %v", err)
				return fmt.Errorf("updatePostings err: %v", err)
			}
		}
		// 更新index count
		if in.IndexCount > 0 {
			err := in.updateCount(in.IndexCount)
			if err != nil {
				return fmt.Errorf("updateCount err: %v", err)
			}
		}

		// 重置
		in.PostingsHashBuf = make(engine.InvertedIndexHash)
		in.BufCount = 0
	}

	return nil

}

func (in *Index) updateCount(num uint64) error {
	count, err := in.ForwardDB.Count()
	if err != nil {
		if err.Error() == engine.ErrCountKeyNotFound {
			count = 0
		} else {
			return fmt.Errorf("updateCount err: %v", err)
		}
	}
	count += num
	return in.ForwardDB.UpdateCount(count)
}

// 落盘
func (in *Index) updatePostings(p *engine.InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}
	// 不需要拉取后merge，直接写入文件
	// // 拉取数据库数据
	// oldPostings, size, err := in.FetchPostings(p.Token)
	// if err != nil {
	// 	return fmt.Errorf("updatePostings fetchPostings err: %v", err)
	// }
	// // merge
	// if size > 0 {
	// 	p.PostingsList = engine.MergePostings(oldPostings, p.PostingsList)
	// 	p.DocsCount += size
	// }
	// 开始写入数据库
	buf, err := engine.EncodePostings(p.PostingsList, p.DocsCount)
	if err != nil {
		return fmt.Errorf("updatePostings encodePostings err: %v", err)
	}
	return in.InvertedDB.DBUpdatePostings(p.Token, buf.Bytes(), p.DocsCount)
}

// Close --
func (in *Index) Close() {
	in.Engine.Close()
}

// NewIndexEngine init
func NewIndexEngine(e *engine.Engine) (*Index, error) {
	if e == nil {
		return nil, fmt.Errorf("NewIndexEngine err: %v", "engine is nil")
	}
	return &Index{e, 0}, nil
}
