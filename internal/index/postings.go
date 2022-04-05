package index

import (
	"brain/internal/query"
	"brain/internal/storage"
	"bytes"
	"fmt"
)

// mergePostings merge two postings list
// https://leetcode-cn.com/problems/he-bing-liang-ge-pai-xu-de-lian-biao-lcof/
func mergePostings(pa, pb *PostingsList) *PostingsList {
	ret := new(PostingsList)
	p := new(PostingsList)
	p = nil
	for pa != nil || pb != nil {

		temp := new(PostingsList)
		if pb == nil || (pa != nil && pa.DocID <= pb.DocID) {
			temp = pa
			pa = pa.next
		} else if pa == nil || (pb != nil && pa.DocID > pb.DocID) {
			temp = pb
			pb = pb.next
		} else {
			break
		}
		temp.next = nil

		if p == nil {
			ret.next = temp
		} else {
			p.next = temp
		}

		p = temp
	}

	return ret.next
}

// mergeInvertedIndex 合并两个倒排索引
func mergeInvertedIndex(base, toBeAdded InvertedIndexHash) {
	for tokenID, index := range base {
		if toBeAddedIndex, ok := (toBeAdded)[tokenID]; ok {
			index.postingList = mergePostings(index.postingList, toBeAddedIndex.postingList)
			index.docsCount += toBeAddedIndex.docsCount
			delete(toBeAdded, tokenID)
		}
	}
	for tokenID, index := range toBeAdded {
		(base)[tokenID] = index
	}

}

// 解码
func decodePostings() {

}

// 编码
// bytes.Buffer
func encodePostings(postings *PostingsList, docCount int64) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})

	return buf

	// static int encode_postings_none(
	// const postings_list *postings,
	// const int postings_len,
	// buffer *postings_e) {
	//     const postings_list *p;

	//     LL_FOREACH(postings, p) {
	//         int *pos = NULL;
	//         append_buffer(postings_e, (void *)&p->document_id, sizeof(int));
	//         append_buffer(postings_e, (void *)&p->positions_count, sizeof(int));
	//         while ((pos = (int *)utarray_next(p->positions, pos))) {
	//             append_buffer(postings_e, (void *)pos, sizeof(int));
	//         }
	//     }
	//     return 0;
	// }

}

func fetchPostings(tokenID int64) (*PostingsList, int64, error) {

	return nil, 0, nil
}

func updatePostings(p *InvertedIndexValue) error {
	if p == nil {
		return fmt.Errorf("updatePostings p is nil")
	}
	// 拉取数据库数据
	oldPostings, size, err := fetchPostings(p.TokenID)
	if err != nil {
		return fmt.Errorf("updatePostings fetchPostings err: %v", err)
	}
	// merge
	if size > 0 {
		p.postingList = mergePostings(oldPostings, p.postingList)
		p.docsCount += size
	}
	// 开始写入数据库
	buf := encodePostings(p.postingList, p.docsCount)

	return storage.DBUpdatePostings(p.TokenID, p.docsCount, buf, int64(buf.Len()))
}

// /**
//  * 将内存上（小倒排索引中）的倒排列表与存储器上的倒排列表合并后存储到数据库中
//  * @param[in] env 存储着应用程序运行环境的结构体
//  * @param[in] p 含有倒排列表的倒排索引中的索引项
//  */
// void update_postings(const wiser_env *env, inverted_index_value *p) {
//     int old_postings_len;
//     postings_list *old_postings;

//     if (!fetch_postings(env, p->token_id, &old_postings, &old_postings_len)) {
//         buffer *buf;
//         if (old_postings_len) {
//             p->postings_list = merge_postings(old_postings, p->postings_list);
//             p->docs_count += old_postings_len;
//         }
//         if ((buf = alloc_buffer())) {
//             encode_postings(env, p->postings_list, p->docs_count, buf);
// 				// #define BUFFER_PTR(b) ((b)->head)              /* 返回指向缓冲区开头的指针 */
//             db_update_postings(env, p->token_id, p->docs_count, BUFFER_PTR(buf), BUFFER_SIZE(buf));
//             free_buffer(buf);
//         }
//     } else {
//         print_error("cannot fetch old postings list of token(%d) for update.", p->token_id);
//     }
// }

// int text_to_postings_lists(wiser_env* env, const int document_id, const UTF32Char* text, const unsigned int text_len,
// const int n, inverted_index_hash** postings) {

// text2PostingsLists --
func (e *Engine) text2PostingsLists(docID int64, text []byte) error {
	tokens, err := query.Ngram(string(text), e.N)
	if err != nil {
		return fmt.Errorf("text2PostingsLists Ngram err: %v", err)
	}
	bufInvertedHash := make(InvertedIndexHash)

	for _, token := range tokens {
		err := e.token2PostingsLists(bufInvertedHash, token.Token, token.Position, docID)
		if err != nil {
			return fmt.Errorf("text2PostingsLists token2PostingsLists err: %v", err)
		}
	}

	if e.postingsHashBuf != nil && len(e.postingsHashBuf) > 0 {
		mergeInvertedIndex(e.postingsHashBuf, bufInvertedHash)
	} else {
		e.postingsHashBuf = make(InvertedIndexHash)
		e.postingsHashBuf = bufInvertedHash
	}
	return nil

}

func (e *Engine) token2PostingsLists(
	bufInvertHash InvertedIndexHash, token []byte,
	position int64, docID int64) error {

	bufInvert := new(InvertedIndexValue)

	// 查询的是整个索引库 不是临时库
	// doc_id用来标识写入数据还是查询数据
	tokenID, docCount, err := e.db.GetTokenID(token, docID)
	if err != nil {
		return fmt.Errorf("token2PostingsLists GetTokenID err: %v", err)
	}

	if len(bufInvertHash) > 0 {
		if item, ok := bufInvertHash[tokenID]; ok {
			bufInvert = item
		}
	}

	pl := new(PostingsList)
	if bufInvert != nil {
		pl = bufInvert.postingList
		// 这里的positioinCount和下面bufInvert的po
		pl.positionCount++
	} else {
		if docID != 0 {
			docCount = 1
		}
		bufInvert = createNewInvertedIndex(tokenID, docCount)
		bufInvertHash[tokenID] = bufInvert
		pl = createNewPostingList(docID)
		bufInvert.postingList = pl
	}
	// 存储位置信息
	pl.positions = append(pl.positions, position)
	bufInvert.positionCount++

	return nil
}
