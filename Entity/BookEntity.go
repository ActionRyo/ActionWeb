// 书本信息类模型
package Entity

type BookInfo struct {
	TableID  int
	BookCode string
	BookName string
	UserID   int
}

type ListBooks struct {
	ArrBooks []*BookInfo
}
