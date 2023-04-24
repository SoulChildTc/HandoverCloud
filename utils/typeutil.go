package utils

import "database/sql"

func ScanNullString(src string) sql.NullString {
	nullStr := sql.NullString{}
	_ = nullStr.Scan(src)
	return nullStr
}
