package service

import (
	"github.com/lib/pq"
)

func isUniqueViolation(err error) bool {
	pgerr, ok := err.(*pq.Error)
	return ok && pgerr.Code == "23505"
}
