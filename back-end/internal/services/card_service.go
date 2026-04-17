package services

import (
	"database/sql"
	"errors"

	"go-do-list/internal/models"
)

type CardService struct {
	db *sql.DB
}

func NewCardService(db *sql.DB) *CardService {
	return &CardService{db: db}
}

func (s *CardService) CreateCard(columnID, title, description string) (*models.Card, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	var maxPos sql.NullInt64
	s.db.QueryRow(`SELECT MAX(position) FROM cards WHERE column_id = $1`, columnID).Scan(&maxPos)
	pos := int(maxPos.Int64) + 1

	var card models.Card
	card.Labels = []models.Label{}
	err := s.db.QueryRow(
		`INSERT INTO cards (column_id, title, description, position)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, column_id, title, description, position, created_at`,
		columnID, title, description, pos,
	).Scan(&card.ID, &card.ColumnID, &card.Title, &card.Description, &card.Position, &card.CreatedAt)
	return &card, err
}

func (s *CardService) UpdateCard(id, title, description string) (*models.Card, error) {
	var card models.Card
	err := s.db.QueryRow(
		`UPDATE cards SET
		   title = COALESCE(NULLIF($2,''), title),
		   description = CASE WHEN $3 = '' THEN description ELSE $3 END
		 WHERE id = $1
		 RETURNING id, column_id, title, description, position, created_at`,
		id, title, description,
	).Scan(&card.ID, &card.ColumnID, &card.Title, &card.Description, &card.Position, &card.CreatedAt)
	if err != nil {
		return nil, errors.New("card not found")
	}
	return &card, nil
}

func (s *CardService) DeleteCard(id string) error {
	_, err := s.db.Exec(`DELETE FROM cards WHERE id = $1`, id)
	return err
}

// MoveCard moves a card to a different column at the given position.
func (s *CardService) MoveCard(cardID, targetColumnID string, position int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current column
	var currentColumnID string
	err = tx.QueryRow(`SELECT column_id FROM cards WHERE id = $1`, cardID).Scan(&currentColumnID)
	if err != nil {
		return errors.New("card not found")
	}

	if currentColumnID != targetColumnID {
		// Remove from old column: shift down positions
		_, err = tx.Exec(
			`UPDATE cards SET position = position - 1
			 WHERE column_id = $1 AND position > (SELECT position FROM cards WHERE id = $2)`,
			currentColumnID, cardID,
		)
		if err != nil {
			return err
		}

		// Make room in target column
		_, err = tx.Exec(
			`UPDATE cards SET position = position + 1
			 WHERE column_id = $1 AND position >= $2`,
			targetColumnID, position,
		)
		if err != nil {
			return err
		}

		// Update card
		_, err = tx.Exec(
			`UPDATE cards SET column_id = $1, position = $2 WHERE id = $3`,
			targetColumnID, position, cardID,
		)
		if err != nil {
			return err
		}
	} else {
		// Same column reorder
		var currentPos int
		tx.QueryRow(`SELECT position FROM cards WHERE id = $1`, cardID).Scan(&currentPos)

		if currentPos < position {
			_, err = tx.Exec(
				`UPDATE cards SET position = position - 1
				 WHERE column_id = $1 AND position > $2 AND position <= $3 AND id != $4`,
				currentColumnID, currentPos, position, cardID,
			)
		} else {
			_, err = tx.Exec(
				`UPDATE cards SET position = position + 1
				 WHERE column_id = $1 AND position >= $2 AND position < $3 AND id != $4`,
				currentColumnID, position, currentPos, cardID,
			)
		}
		if err != nil {
			return err
		}

		_, err = tx.Exec(`UPDATE cards SET position = $1 WHERE id = $2`, position, cardID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ReorderCards sets positions 0..n for the given card IDs in a column.
func (s *CardService) ReorderCards(columnID string, cardIDs []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, id := range cardIDs {
		_, err := tx.Exec(
			`UPDATE cards SET position = $1 WHERE id = $2 AND column_id = $3`,
			i, id, columnID,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *CardService) AddLabel(cardID, labelID string) error {
	_, err := s.db.Exec(
		`INSERT INTO card_labels (card_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		cardID, labelID,
	)
	return err
}

func (s *CardService) RemoveLabel(cardID, labelID string) error {
	_, err := s.db.Exec(
		`DELETE FROM card_labels WHERE card_id = $1 AND label_id = $2`, cardID, labelID,
	)
	return err
}
