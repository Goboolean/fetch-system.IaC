-- name: GetAllMetadata :many
SELECT id, symbol, locale, market FROM product_meta;

-- name: InsertManyMetadata :copyfrom
INSERT INTO product_meta (id, symbol, locale, market) VALUES ($1, $2, $3, $4);
