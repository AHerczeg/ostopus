package tentacles

import (
	"ostopus/shared/tentacle"
)

var (
	tentacles *store
)

type tentacleStore interface {
	SaveTentacle(tentacle.Tentacle)
	HasTentacle(string) bool
	RemoveTentacle(string)bool
}

type store struct {
	tentacles map[string]tentacle.Tentacle
}

func Tentacles() *store {
	if tentacles == nil {
		tentacles = &store{
			tentacles: make(map[string]tentacle.Tentacle),
		}
	}
	return tentacles
}

func (s *store) SaveTentacle(tentacle tentacle.Tentacle) {
	s.tentacles[tentacle.Name] = tentacle
}

func (s *store) GetAllTentacles() []tentacle.Tentacle {
	var tentacles []tentacle.Tentacle
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
