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
		logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func countInvalidDomains() (int64, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountInvalidDomains(context.Background())
	if queryErr != nil {
		logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func countUnvalidatedDomains() (int64, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountUnvalidatedDomains(context.Background())
	if queryErr != nil {
		logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func getDomain(domain string) (generated.Domain, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return generated.Domain{}, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.GetDomain(context.Background(), domain)
	if queryErr != nil {
		logger.Error("Unable to get domain", "err", queryErr)
		return generated.Domain{}, queryErr
	}

	return row, nil
}

func listDomains(limit int32, offset int32) ([]generated.Domain, error) {

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

func listInvalidDomains(limit int32, offset int32) ([]generated.Domain, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	rows, queryErr := queries.ListInvalidDomains(context.Background(), generated.ListInvalidDomainsParams{
		Thelimit:  limit,
		Theoffset: offset,
	})
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func listUnvalidatedDomains(limit int32, offset int32) ([]generated.Domain, error) {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	rows, queryErr := queries.ListUnvalidatedDomains(context.Background(), generated.ListUnvalidatedDomainsParams{
		Thelimit:  limit,
		Theoffset: offset,
	})
	if queryErr != nil {
		logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func updateValidation(domain string, valid bool, reason string) error {

	conn, connErr := connect()
	if connErr != nil {
		// already logged
		return connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	queryErr := queries.UpdateValidation(context.Background(), generated.UpdateValidationParams{
		Domain: domain,
		Valid:  pgtype.Bool{Bool: valid, Valid: true},
		Reason: pgtype.Text{String: reason, Valid: true},
	})
	if queryErr != nil {
		logger.Error("Unable to update validity in database", "err", queryErr)
		return queryErr
	}

	return nil
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
