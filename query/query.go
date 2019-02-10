package query

var queryStore = map[string]string {
	"kernel_info": "SELECT * FROM kernel_info;",
}

func GetQuery(option string) (string, bool) {
	query, ok := queryStore[option]
	return query, ok
}

func HasQuery(option string) bool {
	_, ok := queryStore[option]
	return ok
}