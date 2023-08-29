package mysql

type ChatDao struct {
	conn *GormConn
}

func NewChatDao() *ChatDao {
	return &ChatDao{
		conn: NewGormConn(),
	}
}
