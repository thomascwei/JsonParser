-- name: ListTemplates :many
SELECT id, template_name, description
FROM template
ORDER BY id;

-- name: GetTemplate :one
SELECT *
FROM template
WHERE id = ?;

-- name: ListTemplatesForParse :many
SELECT id, parse_type, value_extract, go_struct
FROM template
ORDER BY id;

-- name: GetObjectIDs :many
SELECT *
FROM objectids
WHERE template_id = ?;

-- name: CreateTemplate :execresult
INSERT INTO template (template_name, description, parse_type, auto_gen_object_id, origin_object_id,
                      value_extract, go_struct)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: CreateObjectIDs :execresult
INSERT INTO objectids (template_id, serial, object_id)
VALUES (?, ?, ?);

-- name: UpdateTemplate :exec
UPDATE template
SET template_name=?,
    description=?,
    parse_type=?,
    auto_gen_object_id=?,
    origin_object_id=?,
    value_extract=?,
    go_struct=?,
    updated_at=?
WHERE id = ?;

-- name: UpdateObjectIDs :exec
UPDATE objectids
set template_id=?,
    serial=?,
    object_id=?,
    updated_at=?
WHERE id = ?;

-- name: DeleteTemplate :exec
DELETE
FROM template
WHERE id = ?;

-- name: DeleteObjectIDs :exec
DELETE
FROM objectids
WHERE template_id = ?;