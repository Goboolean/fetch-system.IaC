-- name: GetAllProducts :many
SELECT id, platform, symbol, locale, market, name, description FROM product_meta;

-- name: GetProductById :one
SELECT id, platform, symbol, locale, market, name, description FROM product_meta
WHERE id = $1;

-- name: GetProductsByCondition :many
SELECT id, platform, symbol, locale, market, name, description FROM product_meta
WHERE platform = $1 AND market = $2;

-- name: InsertProducts :copyfrom
INSERT INTO product_meta (id, platform, symbol, locale, market, name, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ProductExists :one
SELECT EXISTS(SELECT 1 FROM product_meta WHERE id = $1);

-- name: CountProducts :one
SELECT COUNT(*) FROM product_meta
WHERE platform = $1 AND market = $2;

-- name: DeleteAllProducts :exec
DELETE FROM product_meta;