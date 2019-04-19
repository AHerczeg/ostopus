package shared

type Query struct {
	Query		string
	Frequency	string
}

func (q Query) Validate() bool  {
	if !validateQuery(q.Query) || !validateFrequency(q.Frequency) {
		return false
	}

	return true
}

func validateQuery(query string) bool {
	// TODO

	return true
}

func validateFrequency(frequency string) bool {
	// TODO

	return true
}