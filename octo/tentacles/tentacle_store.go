package tentacles

var (
	tentacles *store
)

type store struct {
	tentacles map[string]Tentacle
}

func Tentacles() *store {
	if tentacles == nil {
		tentacles = &store{
			tentacles: make(map[string]Tentacle),
		}
	}
	return tentacles
}

func (s *store) SaveTentacle(tentacle Tentacle) {
	s.tentacles[tentacle.Name] = tentacle
}

func (s *store) GetAllTentacles() []Tentacle {
	var tentacles []Tentacle
	for _, tentacle := range s.tentacles {
		tentacles = append(tentacles, tentacle)
	}
	return tentacles
}
