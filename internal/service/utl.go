package service

import (
	"github.com/jackc/pgconn"
)

func isUniqueViolation(err error) bool {
	pgerr, ok := err.(*pgconn.PgError)
	return ok && pgerr.Code == "23505"
}
