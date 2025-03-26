package chat

import "sync"

type shard struct {
	m     map[uint]iWsConn
	mutex *sync.RWMutex
}

// getAll retrieves all connections from the shard
// Returns a slice of all connections in the shard
func (s *shard) getAll() []iWsConn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	cons := make([]iWsConn, 0, len(s.m))
	for _, conn := range s.m {
		cons = append(cons, conn)
	}
	return cons
}

// getTotalCount returns the total number of connections in the shard
// Returns the total number of connections in the shard
func (s *shard) getTotalCount() uint {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return uint(len(s.m))
}

// get retrieves a connection from the shard by user ID
// Returns the connection and a boolean indicating if it exists
func (s *shard) get(uid uint) (conn iWsConn, exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conn, exist = s.m[uid]
	return
}

// set adds a connection to the shard
// Adds a connection to the shard
func (s *shard) set(conn iWsConn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[conn.getUserId()] = conn
}

// remove removes a connection from the shard by user ID
// Removes a connection from the shard
func (s *shard) remove(uid uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.m, uid)
}
