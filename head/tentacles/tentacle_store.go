package tentacles

import (
	"ostopus/shared"
)

var (
	tentacles *store
)

type tentacleStore interface {
	SaveTentacle(shared.Tentacle)
	HasTentacle(string) bool
	RemoveTentacle(string) bool
}

type store struct {
	tentacles map[string]shared.Tentacle
}

func Tentacles() *store {
	if tentacles == nil {
		tentacles = &store{
			tentacles: make(map[string]shared.Tentacle),
		}
	}
	return tentacles
}

func (s *store) SaveTentacle(tentacle shared.Tentacle) {
	s.tentacles[tentacle.Name] = tentacle
}

func (s *store) GetAllTentacles() []shared.Tentacle {
	var tentacles []shared.Tentacle
	for _, tentacle := range s.tentacles {
		tentacles = append(tentacles, tentacle)
	}
	return tentacles
}

func (s *store) HasTentacle(name string) bool {
	_, ok := s.tentacles[name]
	return ok
}

func (s *store) RemoveTentacle(name string) bool {
	_, ok := s.tentacles[name]
	delete(s.tentacles, name)
	return ok
}

func (s *store) GetTentacle(name string) (shared.Tentacle, bool) {
	t, ok :=  s.tentacles[name]
	return t, ok
}