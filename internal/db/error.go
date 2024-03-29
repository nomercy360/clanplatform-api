package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

const UniqueViolationConstraint = "unique_violation"
const ForeignKeyViolationConstraint = "foreign_key_violation"
const CheckViolationConstraint = "check_violation"
const ExclusionViolationConstraint = "exclusion_violation"

func ErrorIs(err error, c string) bool {
	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)

	return ok && pgErr.Code.Name() == c
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
