package sqlparser

import querypb "github.com/youtube/vitess/go/vt/proto/query"

// return a sql string with the params stripped out for display
func RedactSQLQuery(sql string) (string, error) {
	bv := map[string]*querypb.BindVariable{}
	sqlStripped, comments := SplitTrailingComments(sql)

	stmt, err := Parse(sqlStripped)
	if err != nil {
		return "", err
	}

	prefix := "redacted"
	Normalize(stmt, bv, prefix)

	return String(stmt) + comments, nil
}
