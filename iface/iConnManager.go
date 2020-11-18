package iface

type IconnManager interface {
	Add(conn Iconnection)

	Get(connID uint32) (Iconnection, error)
	Remove(conn Iconnection)
	Len() int

	ClearConn()
}
