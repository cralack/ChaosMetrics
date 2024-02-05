package request

// Pager 分页机制
type Pager struct {
	Total int64 // 共有多少数据，只有返回值使用
	Start int   // 起始位置
	Size  int   // 每个页面个数
}
