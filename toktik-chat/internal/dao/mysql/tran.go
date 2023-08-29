package mysql

// Transaction 事务的操作 一定跟数据库有关 注入数据库的连接 gorm.db
type Transaction interface {
	Action(func(conn DbConn) error) error
}

type DbConn interface {
	Begin()
	Rollback()
	Commit()
}

type TransactionImpl struct {
	conn DbConn
}

func (t TransactionImpl) Action(f func(conn DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: NewTran(),
	}
}
