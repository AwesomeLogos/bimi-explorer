-- name: CountDomains :one
SELECT 
        COUNT(*) AS "Count" 
    FROM 
        domain
    WHERE
        valid IS NULL OR valid = TRUE;

-- name: CountInvalidDomains :one
SELECT 
        COUNT(*) AS "Count" 
    FROM 
        domain
    WHERE
        valid = FALSE;

-- name: CountUnvalidatedDomains :one
SELECT 
        COUNT(*) AS "Count" 
    FROM 
        domain
    WHERE
        valid IS NULL;

-- name: GetDomain :one
SELECT 
        * 
    FROM 
        domain 
    WHERE 
        domain = sqlc.arg(domain);

-- name: ListDomains :many
SELECT 
        *
    FROM 
        domain 
    WHERE
        valid IS NULL OR valid = TRUE
    ORDER BY 
        domain 
    LIMIT 
        sqlc.arg(theLimit) 
    OFFSET 
        sqlc.arg(theOffset);

-- name: ListRandom :many
SELECT
        *
    FROM 
        domain TABLESAMPLE SYSTEM_ROWS(2 * sqlc.arg(theLimit)) 
    WHERE
        valid = TRUE
    LIMIT 
        sqlc.arg(theLimit);

-- name: ListInvalidDomains :many
SELECT 
        *
    FROM 
        domain 
    WHERE
        valid = FALSE
    ORDER BY 
        domain 
    LIMIT 
        sqlc.arg(theLimit) 
    OFFSET 
        sqlc.arg(theOffset);

-- name: ListUnvalidatedDomains :many
SELECT 
        *
    FROM 
        domain 
    WHERE
        valid IS NULL
    ORDER BY 
        domain 
    LIMIT 
        sqlc.arg(theLimit) 
    OFFSET 
        sqlc.arg(theOffset);

-- name: SearchDomains :many
SELECT 
        *
    FROM 
        domain 
    WHERE TRUE
        AND domain ILIKE sqlc.arg(search)
        AND valid = TRUE
    ORDER BY 
        domain 
    LIMIT 
        100; 

-- name: UpdateValidation :exec
UPDATE 
        domain 
    SET 
        valid = sqlc.arg(valid),
        reason = sqlc.arg(reason),
        updated = NOW()
    WHERE 
        domain = sqlc.arg(domain);

-- name: UpsertDomain :one
INSERT INTO domain 
        (domain, imgurl)
    VALUES
        (sqlc.arg(domain), sqlc.arg(imgurl))
    ON CONFLICT (domain) DO UPDATE SET 
        imgurl = sqlc.arg(imgurl),
        updated = NOW()
    RETURNING *;

