package service

type DebugService struct {
}

func NewDebugService() *DebugService {
	return &DebugService{}
}

func (s *DebugService) List() (map[string]interface{}, error) {
	// 获取 meta 信息

	// 获取数据库信息 用于跟 meta 信息对比

	// 获取 term list

	// 对应 term 的 doc list

	// 对应 doc 的内容

	// 正排信息存储展示 编解码

	// return

	return nil, nil
}
