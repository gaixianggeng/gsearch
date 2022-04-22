package storage

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

const bucketName = "forward"

const forwardCountKey = "forwardCount"

// ForwardDB 存储器
type ForwardDB struct {
	db *bolt.DB
}

// Add add forward data
func (f *ForwardDB) Add(doc *Document) error {
	key := strconv.Itoa(int(doc.DocID))
	body, _ := json.Marshal(doc)
	return Put(f.db, bucketName, []byte(key), body)
}

// Count 获取文档总数
func (f *ForwardDB) Count() (uint64, error) {
	body, err := Get(f.db, bucketName, []byte(forwardCountKey))
	if err != nil {
		return 0, err
	}
	c, err := strconv.Atoi(string(body))
	return uint64(c), err
}

// UpdateCount 获取文档总数
func (f *ForwardDB) UpdateCount(count uint64) error {
	return Put(f.db, bucketName, []byte(forwardCountKey), []byte(strconv.Itoa(int(count))))
}

// Get get forward data
func (f *ForwardDB) Get(docID uint64) ([]byte, error) {
	key := strconv.Itoa(int(docID))
	return Get(f.db, bucketName, []byte(key))
}

// Close --
func (f *ForwardDB) Close() {
	f.db.Close()
}

// NewForwardDB --
func NewForwardDB(dbName string) *ForwardDB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &ForwardDB{db}
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
