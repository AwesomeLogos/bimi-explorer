-- name: UpsertDomain :one
INSERT INTO domain 
        (domain, imgurl)
    VALUES
        (sqlc.arg(domain), sqlc.arg(imgurl))
    ON CONFLICT (domain) DO UPDATE SET 
        imgurl = sqlc.arg(imgurl),
        updated = NOW()
    RETURNING *;

-- name: ListSampleDomains :many
SELECT * FROM domain LIMIT sqlc.arg(maxlimit);