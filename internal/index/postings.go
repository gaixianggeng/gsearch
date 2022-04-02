package index

import (
	"brain/internal/storage"
	"bytes"
	"fmt"
)

// MergePostings merge two postings list
// https://leetcode-cn.com/problems/he-bing-liang-ge-pai-xu-de-lian-biao-lcof/
func MergePostings(pa, pb *PostingsList) *PostingsList {
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

// MergeInvertedIndex 合并两个倒排索引
func MergeInvertedIndex(base, toBeAdded *InvertedIndexHash) {
	for tokenID, index := range *base {
		if toBeAddedIndex, ok := (*toBeAdded)[tokenID]; ok {
			index.postingList = MergePostings(index.postingList, toBeAddedIndex.postingList)
			index.docsCount += toBeAddedIndex.docsCount
			//TODO: 不确定要不要加 先todo
			// index.positionCount += toBeAddedIndex.positionCount
			delete(*toBeAdded, tokenID)
		}
	}
	for tokenID, index := range *toBeAdded {
		(*base)[tokenID] = index
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
		fmt.Println("updatePostings p is nil")
		return nil
	}
	// 拉取数据库数据
	oldPostings, size, err := fetchPostings(p.TokenID)
	if err != nil {
		return fmt.Errorf("updatePostings fetchPostings err: %v", err)
	}
	// merge
	if size > 0 {
		p.postingList = MergePostings(oldPostings, p.postingList)
		p.docsCount += size
	}
	// 开始写入数据库
	buf := encodePostings(p.postingList, p.docsCount)

	storage.DBUpdatePostings(p.TokenID, p.docsCount, buf, int64(buf.Len()))

	return nil
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
	return nil

}

func (e *Engine) token2PostingsLists(docID int64, token []byte,position int64,) error {
