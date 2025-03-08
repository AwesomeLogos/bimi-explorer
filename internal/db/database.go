package db

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/AwesomeLogos/bimi-explorer/internal/common"
	"github.com/AwesomeLogos/bimi-explorer/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func Connect() (*pgx.Conn, error) {
	dbUrl := os.Getenv("DB_URL")
	u, urlErr := url.Parse(dbUrl)
	if urlErr != nil {
		common.Logger.Error("Unable to parse database URL", "err", urlErr, "url", dbUrl)
		return nil, urlErr
	}
	u.User = url.UserPassword(u.User.Username(), os.Getenv("DB_PASSWORD"))
	conn, connErr := pgx.Connect(context.Background(), u.String())
	if connErr != nil {
		//LATER: maybe log masked password?
		common.Logger.Error("Unable to connect to database", "err", connErr, "url", dbUrl)
		return nil, connErr
	}

	return conn, nil
}

func CountDomains() (int64, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountDomains(context.Background())
	if queryErr != nil {
		common.Logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func CountInvalidDomains() (int64, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountInvalidDomains(context.Background())
	if queryErr != nil {
		common.Logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func CountUnvalidatedDomains() (int64, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return 0, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.CountUnvalidatedDomains(context.Background())
	if queryErr != nil {
		common.Logger.Error("Unable to count", "err", queryErr)
		return 0, queryErr
	}

	return row, nil
}

func GetDomain(domain string) (generated.Domain, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return generated.Domain{}, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	row, queryErr := queries.GetDomain(context.Background(), domain)
	if queryErr != nil {
		common.Logger.Error("Unable to get domain", "err", queryErr)
		return generated.Domain{}, queryErr
	}

	return row, nil
}

func ListDomains(limit int32, offset int32) ([]generated.Domain, error) {

	conn, connErr := Connect()
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
		common.Logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func ListInvalidDomains(limit int32, offset int32) ([]generated.Domain, error) {

	conn, connErr := Connect()
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
		common.Logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func ListRandomDomains(limit int32) ([]generated.Domain, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	queries := generated.New(conn)
	rows, queryErr := queries.ListRandom(context.Background(), limit)
	if queryErr != nil {
		common.Logger.Error("Unable to select random", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func ListUnvalidatedDomains(limit int32, offset int32) ([]generated.Domain, error) {

	conn, connErr := Connect()
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
		common.Logger.Error("Unable to insert into database", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func SearchDomains(searchTerm string) ([]generated.Domain, error) {

	conn, connErr := Connect()
	if connErr != nil {
		// already logged
		return nil, connErr
	}
	defer conn.Close(context.Background())

	if !strings.Contains(searchTerm, "%") {
		searchTerm = "%" + searchTerm + "%"
	}

	queries := generated.New(conn)
	rows, queryErr := queries.SearchDomains(context.Background(), searchTerm)
	if queryErr != nil {
		common.Logger.Error("Unable to select random", "err", queryErr)
		return nil, queryErr
	}

	return rows, nil
}

func UpdateValidation(domain string, valid bool, reason string) error {

	conn, connErr := Connect()
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
		common.Logger.Error("Unable to update validity in database", "err", queryErr)
		return queryErr
	}

	return nil
}

func UpsertDomain(domain, imgurl string) error {

	conn, connErr := Connect()
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
		common.Logger.Error("Unable to insert into database", "err", queryErr)
		return queryErr
	}

	return nil
}
