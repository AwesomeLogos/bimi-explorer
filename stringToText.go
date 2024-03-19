package main

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func stringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}
