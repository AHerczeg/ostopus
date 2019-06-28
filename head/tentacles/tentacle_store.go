package tentacles

import (
	"github.com/AHerczeg/ostopus/head/api/model"
)

var (
	tentacles *store
)

type tentacleStore interface {
	SaveTentacle(tentacle model.Tentacle)
	HasTentacle(int64) bool
	RemoveTentacle(int64) bool
}

type store struct {
	tentacles map[int64]model.Tentacle
}

func Tentacles() *store {
	if tentacles == nil {
		tentacles = &store{
			tentacles: make(map[int64]model.Tentacle),
		}
	}
	return tentacles
}

func (s *store) SaveTentacle(tentacle model.Tentacle) {
	s.tentacles[tentacle.ID] = tentacle
}

func (s *store) GetAllTentacles() []model.Tentacle {
	var tentacles []model.Tentacle
	for _, tentacle := range s.tentacles {
		tentacles = append(tentacles, tentacle)
	}
	return tentacles
}

func (s *store) HasTentacle(id int64) bool {
	_, ok := s.tentacles[id]
	return ok
}

func (s *store) RemoveTentacle(id int64) bool {
	_, ok := s.tentacles[id]
	delete(s.tentacles, id)
	return ok
}

func (s *store) GetTentacle(id int64) (model.Tentacle, bool) {
	t, ok := s.tentacles[id]
	return t, ok
}
