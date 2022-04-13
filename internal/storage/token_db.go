package storage

// TokenDB token存储器
type TokenDB struct {
}

// GetToken 获取token的绑定doc count
// 从tokens表中获取指定词元的编号
func (t *TokenDB) GetToken(token string, insertFlag uint64) (uint64, error) {
	// 写入token库标识
	if insertFlag > 0 {

	}
	return 0, nil
}
