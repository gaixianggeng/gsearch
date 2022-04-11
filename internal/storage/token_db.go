package storage

// TokenDB token存储器
type TokenDB struct {
}

// GetTokenID 获取tokenid和出现次数
// 从tokens表中获取指定词元的编号
func (t *TokenDB) GetTokenID(token []rune, insertFlag uint64) (uint64, uint64, error) {
	// 写入token库标识
	if insertFlag > 0 {

	}
	return 0, 0, nil
}
