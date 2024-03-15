package repository

import (
	"App/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type GroupPostgres struct {
	conn *sqlx.DB
}

func NewGroupPostgres(conn *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{conn}
}

func (r *GroupPostgres) Create(name string, userID int) (int, error) {
	stmt := "INSERT INTO groups (name, user_id) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.conn.QueryRow(stmt, name, userID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *GroupPostgres) GetAll(userID int) ([]models.GroupOut, error) {
	query := `
	SELECT g.id, g.name, COUNT(s.id)
	FROM groups g
	LEFT JOIN students s ON g.id = s.group_id
	WHERE user_id = $1
	GROUP BY g.id, g.name
	ORDER BY g.id DESC;
	`

	rows, err := r.conn.Query(query, userID)

	defer rows.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var groups []models.GroupOut

	for rows.Next() {
		var id int
		var name string
		var count int

		if err := rows.Scan(&id, &name, &count); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.GroupOut{
			ID:    id,
			Name:  name,
			Count: count,
		}
		groups = append(groups, result)
	}

	return groups, nil
}

func (r *GroupPostgres) Get(groupID, userID int) (models.GroupOut, error) {
	query := "SELECT id, name FROM groups WHERE id = $1 AND user_id = $2"

	var group models.GroupOut

	err := r.conn.QueryRow(query, groupID, userID).Scan(&group.ID, &group.Name)

	if err != nil {
		log.Err(err).Send()
		return models.GroupOut{}, err
	}

	return group, nil
}

func (r *GroupPostgres) UpdateTitle(groupID, userID int, name string) error {
	query := "UPDATE groups SET name = $1 WHERE id = $2 AND user_id = $3"

	result, err := r.conn.Exec(query, name, groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		log.Error().Msg("Rows affected: " + string(count))
		log.Err(err).Send()
		return err
	}

	return nil
}
