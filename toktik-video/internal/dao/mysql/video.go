package mysql

type VideoDao struct {
	conn *GormConn
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: NewGormConn(),
	}
}
