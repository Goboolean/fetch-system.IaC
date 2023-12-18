-- name: GetAllProducts :many
SELECT id, symbol, locale, market, name, description FROM product_meta;

-- name: GetProductById :one
SELECT id, symbol, locale, market, name, description FROM product_meta
WHERE id = $1;

-- name: GetProductsByCondition :many
SELECT id, symbol, locale, market, name, description FROM product_meta
WHERE locale = $1 AND market = $2;

-- name: InsertProducts :copyfrom
INSERT INTO product_meta (id, symbol, locale, market, name, description)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteAllProducts :exec
DELETE FROM product_meta;