package chat

import "sync"

type shard struct {
	m     map[uint]iWsConn
	mutex *sync.RWMutex
}

func (s *shard) getAll() []iWsConn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conns := make([]iWsConn, 0, len(s.m))
	for _, conn := range s.m {
		conns = append(conns, conn)
	}
	return conns
}

func (s *shard) getTotalCount() uint {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return uint(len(s.m))
}

func (s *shard) get(uid uint) (conn iWsConn, exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conn, exist = s.m[uid]
	return
}
func (s *shard) set(conn iWsConn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[conn.GetUserId()] = conn
}
func (s *shard) remove(uid uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.m, uid)
}
