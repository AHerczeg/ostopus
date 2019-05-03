package query

var standardQueries = map[string]string{
	"kernel_info":			"SELECT * FROM kernel_info;",
	"users":				"SELECT * FROM users;",
	"malicious_actors":		"SELECT name, path, pid FROM processes WHERE on_disk = 0;",
	"users_with_groups":	"SELECT u.uid, u.gid, u.username, g.name, u.description FROM users u LEFT JOIN groups g ON (u.gid = g. gid);",
	"empty_groups":			"SELECT groups.gid, groups.name FROM groups LEFT JOIN users ON (groups.gid = users.gid) WHERE users.ui d IS NULL;",
	"largest_process_10":	"SELECT pid, name, uid, resident_size FROM processes ORDER BY resident_size DESC LIMIT 10;",
	"most_active_10":		"SELECT count(pid) as total, name FROM processes GROUP BY name ORDER BY total DESC LIMIT 10;",
	"current_logins":		"SELECT * FROM logged_in_users;",
	"last_logins":			"SELECT * FROM last;",
	"iptables":				"SELECT * FROM iptables;",
	"crons":				"SELECT command, path FROM crontab;",
}
