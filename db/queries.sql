-- name: UpsertDomain :one
INSERT INTO domain 
        (domain, imgurl)
    VALUES
        (sqlc.arg(domain), sqlc.arg(imgurl))
    ON CONFLICT (domain) DO UPDATE SET 
        imgurl = sqlc.arg(imgurl),
        updated = NOW()
    RETURNING *;

-- name: ListDomains :many
SELECT domain, imgurl FROM domain ORDER BY domain LIMIT sqlc.arg(theLimit) OFFSET sqlc.arg(theOffset);

-- name: CountDomains :one
SELECT COUNT(*) AS "Count" FROM domain;