package checkrun

import "github.com/mbell-nzl/fid/internal/journal"

func (db store) newCheckID(generate func(string, journal.IDExistsFunc) (string, error)) (string, error) {
	return generate("CHK", func(candidate string) (bool, error) {
		var n int
		err := db.conn.QueryRow(`SELECT count(*) FROM (SELECT id FROM check_requests UNION SELECT id FROM check_starts) WHERE id=?`, candidate).Scan(&n)
		return n > 0, err
	})
}
