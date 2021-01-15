package iface

import "net"

type IServer interface {
	Start()
	Stop()
	Listener() net.Listener
	//Accept() (ISession,error)
	//GetSession(sessionid uint64) ISession
	//GetSessionId(session ISession) uint64
}
