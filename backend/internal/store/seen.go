package store

import (
	"sync"
)

const recentBufferSize = 80 // не показывать повторно пока не покажем столько других

// SeenStore хранит множества просмотренных video ID по ключу зрителя (в памяти).
// recent — кольцевой буфер последних показанных ID, чтобы не давать повторы подряд.
type SeenStore struct {
	mu     sync.RWMutex
	m      map[string]map[uint]struct{}
	recent map[string][]uint // viewerKey -> последние N показанных ID
}

func NewSeenStore() *SeenStore {
	return &SeenStore{
		m:      make(map[string]map[uint]struct{}),
		recent: make(map[string][]uint),
	}
}

// MarkSeen добавляет videoIDs к множеству просмотренных для viewerKey.
func (s *SeenStore) MarkSeen(viewerKey string, videoIDs []uint) {
	if len(videoIDs) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.m[viewerKey] == nil {
		s.m[viewerKey] = make(map[uint]struct{})
	}
	for _, id := range videoIDs {
		s.m[viewerKey][id] = struct{}{}
	}
	// добавляем в буфер «недавно показанных»
	buf := append(s.recent[viewerKey], videoIDs...)
	if len(buf) > recentBufferSize {
		buf = buf[len(buf)-recentBufferSize:]
	}
	s.recent[viewerKey] = buf
}

// GetRecent возвращает множество недавно показанных ID для viewerKey (не показывать их снова в случайной выборке).
func (s *SeenStore) GetRecent(viewerKey string) map[uint]struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	buf, ok := s.recent[viewerKey]
	if !ok || len(buf) == 0 {
		return nil
	}
	out := make(map[uint]struct{}, len(buf))
	for _, id := range buf {
		out[id] = struct{}{}
	}
	return out
}

// GetSeen возвращает копию множества просмотренных ID для viewerKey.
func (s *SeenStore) GetSeen(viewerKey string) map[uint]struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, ok := s.m[viewerKey]
	if !ok || len(set) == 0 {
		return nil
	}
	out := make(map[uint]struct{}, len(set))
	for id := range set {
		out[id] = struct{}{}
	}
	return out
}

// ClearSeen очищает список просмотренных для viewerKey (чтобы начать показывать заново, например в случайном порядке).
func (s *SeenStore) ClearSeen(viewerKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, viewerKey)
	delete(s.recent, viewerKey)
}
