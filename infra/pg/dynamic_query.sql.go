package pg

import (
	"context"
	"fmt"
	"strings"
)

type SortOption struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // "asc" or "desc"
}

type ListMataKuliahRequest struct {
	Limit      int32        `json:"limit"`
	Offset     int32        `json:"offset"`
	MataKuliah *string      `json:"mata_kuliah"`
	IsActive   *bool        `json:"is_active"`
	Sort       []SortOption `json:"sort"`
}

// Safe whitelist for sorting

var allowedSortColumns = map[string]bool{
	"mata_kuliah": true,
	"created_at":  true,
	"is_active":   true,
}

func buildOrderBy(sort []SortOption) string {
	if len(sort) == 0 {
		return "ORDER BY created_at DESC" // default
	}
	clauses := []string{}
	for _, s := range sort {
		col := strings.ToLower(s.Column)
		dir := strings.ToUpper(s.Direction)
		if !allowedSortColumns[col] {
			continue
		}
		if dir != "ASC" && dir != "DESC" {
			dir = "ASC"
		}
		clauses = append(clauses, fmt.Sprintf("%s %s", col, dir))
	}
	if len(clauses) == 0 {
		return "ORDER BY created_at DESC"
	}
	return "ORDER BY " + strings.Join(clauses, ", ")
}

// Dynamic list query

func (q *Queries) ListMataKuliahDynamic(ctx context.Context, arg ListMataKuliahRequest) ([]MataKuliah, error) {
	orderBy := buildOrderBy(arg.Sort)

	sql := `
	SELECT id, mata_kuliah, is_active, deleted_by, deleted_at,
	       updated_note, updated_by, updated_at, created_by, created_at
	FROM mata_kuliah
	WHERE deleted_at IS NULL
	`

	params := []interface{}{}
	if arg.MataKuliah != nil {
		sql += fmt.Sprintf(" AND mata_kuliah ILIKE '%%' || $%d || '%%' ", len(params)+1)
		params = append(params, *arg.MataKuliah)
	}
	if arg.IsActive != nil {
		sql += fmt.Sprintf(" AND is_active = $%d ", len(params)+1)
		params = append(params, *arg.IsActive)
	}

	sql += " " + orderBy
	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, arg.Limit, arg.Offset)

	rows, err := q.db.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MataKuliah
	for rows.Next() {
		var i MataKuliah
		if err := rows.Scan(
			&i.ID,
			&i.MataKuliah,
			&i.IsActive,
			&i.DeletedBy,
			&i.DeletedAt,
			&i.UpdatedNote,
			&i.UpdatedBy,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, rows.Err()
}

//count

func (q *Queries) CountMataKuliahDynamic(ctx context.Context, arg ListMataKuliahRequest) (int64, error) {
	sql := `
	SELECT COUNT(*)
	FROM mata_kuliah
	WHERE deleted_at IS NULL
	`

	params := []interface{}{}
	if arg.MataKuliah != nil {
		sql += fmt.Sprintf(" AND mata_kuliah ILIKE '%%' || $%d || '%%' ", len(params)+1)
		params = append(params, *arg.MataKuliah)
	}
	if arg.IsActive != nil {
		sql += fmt.Sprintf(" AND is_active = $%d ", len(params)+1)
		params = append(params, *arg.IsActive)
	}

	var count int64
	err := q.db.QueryRow(ctx, sql, params...).Scan(&count)
	return count, err
}

// example multi sorted column

// -- name: ListMataKuliah :many
// SELECT
//   id,
//   mata_kuliah,
//   is_active,
//   deleted_by,
//   deleted_at,
//   updated_note,
//   updated_by,
//   updated_at,
//   created_by,
//   created_at
// FROM mata_kuliah
// WHERE deleted_at IS NULL
//   AND (sqlc.narg('mata_kuliah')::text IS NULL OR mata_kuliah ILIKE '%' || sqlc.narg('mata_kuliah') || '%')
//   AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
// ORDER BY
//   -- dynamic sort using CASE on multiple columns
//   CASE WHEN sqlc.narg('sort_mata_kuliah') = 'asc'  THEN mata_kuliah END ASC,
//   CASE WHEN sqlc.narg('sort_mata_kuliah') = 'desc' THEN mata_kuliah END DESC,
//   CASE WHEN sqlc.narg('sort_created_at') = 'asc'  THEN created_at END ASC,
//   CASE WHEN sqlc.narg('sort_created_at') = 'desc' THEN created_at END DESC
// LIMIT sqlc.arg('limit')
// OFFSET sqlc.arg('offset');

// -- name: CountMataKuliah :one
// SELECT COUNT(*)::bigint
// FROM mata_kuliah
// WHERE deleted_at IS NULL
//   AND (sqlc.narg('mata_kuliah')::text IS NULL OR mata_kuliah ILIKE '%' || sqlc.narg('mata_kuliah') || '%')
//   AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean);
