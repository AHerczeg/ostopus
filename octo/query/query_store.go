package query

type QueryStore struct {
	queries map[string]string
}

func NewQueryStore() QueryStore {
	var qs QueryStore

}

func GetQuery(option string) (string, bool) {
	query, ok := queryStore[option]
	return query, ok
}

func HasQuery(option string) bool {
	_, ok := queryStore[option]
	return ok
}
