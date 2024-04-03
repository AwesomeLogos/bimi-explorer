package main

import (
	"context"
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
	_, queryErr := queries.UpsertDomain(context.Background(), generated.UpsertDomainParams{
		Domain: domain,
		Imgurl: pgtype.Text{String: imgurl, Valid: true},
	})
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return queryErr
	}

	return nil
}

func countDomains() (int64, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountDomains(context.Background())
	if queryErr != nil {
		logger.Error("Unable count rows", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func listDomains(limit int32, offset int32) ([]generated.ListDomainsRow, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	rows, queryErr := queries.ListDomains(context.Background(), generated.ListDomainsParams{
		Thelimit:  limit,
		Theoffset: offset,
	})
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}
