package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}
