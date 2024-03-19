package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/AwesomeLogos/bimi-explorer/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func connect() (*pgx.Conn, error) {
	dbUrl := os.Getenv("DB_URL")
	u, urlErr := url.Parse(dbUrl)
	if urlErr != nil {
		logger.Error("Unable to parse database URL", "err", urlErr, "url", dbUrl)
		return nil, urlErr
	}
	u.User = url.UserPassword(u.User.Username(), os.Getenv("DB_PASSWORD"))
	conn, connErr := pgx.Connect(context.Background(), u.String())
	if connErr != nil {
		//LATER: maybe log masked password?
		logger.Error("Unable to connect to database", "err", connErr, "url", dbUrl)
		return nil, connErr
	}

	return conn, nil
}

func upsertDomain(domain, imgurl string) error {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.UpsertDomain(context.Background(), generated.UpsertDomainParams{
		Domain: domain,
		Imgurl: pgtype.Text{String: imgurl, Valid: true},
	})
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return queryErr
	}

	fmt.Printf("upsert=%v\n", row)
	return nil
}

func listSampleDomains(limit int32) ([]generated.Domain, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	rows, queryErr := queries.ListSampleDomains(context.Background(), limit)
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	fmt.Printf("upsert=%v\n", rows)
	return rows, nil
}
