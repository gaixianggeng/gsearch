package storage

import (
	"bytes"
)

// DB 存储器
type DB struct{}

// AddDoc 通过写入正排数据，获取docid
func (d *DB) AddDoc(title, body []byte) (int64, error) {

	return 0, nil
}

// /**
//  * 将倒排列表存储到数据库中
//  * @param[in] env 存储着应用程序运行环境的结构体
//  * @param[in] token_id 词元编号
//  * @param[in] docs_count 倒排列表中的文档数
//  * @param[in] postings 待存储的倒排列表
//  * @param[in] postings_size 倒排列表的字节数
//  */
// int db_update_postings(const wiser_env *env, int token_id, int docs_count, void *postings, int postings_size) {
//     int rc;
//     sqlite3_reset(env->update_postings_st);
//     sqlite3_bind_int(env->update_postings_st, 1, docs_count);
//     sqlite3_bind_blob(env->update_postings_st, 2, postings, (unsigned int)postings_size, SQLITE_STATIC);
//     sqlite3_bind_int(env->update_postings_st, 3, token_id);
// query:
//     rc = sqlite3_step(env->update_postings_st);

//     switch (rc) {
//         case SQLITE_BUSY:
//             goto query;
//         case SQLITE_ERROR:
//             print_error("ERROR: %s", sqlite3_errmsg(env->db));
//             break;
//         case SQLITE_MISUSE:
//             print_error("MISUSE: %s", sqlite3_errmsg(env->db));
//             break;
//     }
//     return rc;
// }

// DBUpdatePostings 倒排列表存储到数据库中
func DBUpdatePostings(tokenID int64, docsCount int64, buf *bytes.Buffer, size int64) error {
	return nil

}
