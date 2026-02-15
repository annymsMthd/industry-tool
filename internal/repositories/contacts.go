package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/annymsMthd/industry-tool/internal/models"
	"github.com/pkg/errors"
)

type Contacts struct {
	db *sql.DB
}

func NewContacts(db *sql.DB) *Contacts {
	return &Contacts{db: db}
}

// GetByUser returns all contacts for a user (both sent and received, all statuses)
func (r *Contacts) GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error) {
	query := `
		SELECT
			c.id,
			c.requester_user_id,
			c.recipient_user_id,
			u_req.name AS requester_name,
			u_rec.name AS recipient_name,
			c.status,
			c.requested_at,
			c.responded_at
		FROM contacts c
		JOIN users u_req ON c.requester_user_id = u_req.id
		JOIN users u_rec ON c.recipient_user_id = u_rec.id
		WHERE c.requester_user_id = $1 OR c.recipient_user_id = $1
		ORDER BY c.requested_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query contacts")
	}
	defer rows.Close()

	contacts := []*models.Contact{}
	for rows.Next() {
		var contact models.Contact
		err = rows.Scan(
			&contact.ID,
			&contact.RequesterUserID,
			&contact.RecipientUserID,
			&contact.RequesterName,
			&contact.RecipientName,
			&contact.Status,
			&contact.RequestedAt,
			&contact.RespondedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan contact")
		}
		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

// Create sends a contact request
func (r *Contacts) Create(ctx context.Context, requesterID, recipientID int64) (*models.Contact, error) {
	query := `
		INSERT INTO contacts (requester_user_id, recipient_user_id, status, requested_at)
		VALUES ($1, $2, 'pending', NOW())
		RETURNING id, requester_user_id, recipient_user_id, status, requested_at, responded_at
	`

	var contact models.Contact
	err := r.db.QueryRowContext(ctx, query, requesterID, recipientID).Scan(
		&contact.ID,
		&contact.RequesterUserID,
		&contact.RecipientUserID,
		&contact.Status,
		&contact.RequestedAt,
		&contact.RespondedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create contact request")
	}

	return &contact, nil
}

// UpdateStatus accepts or rejects a contact request and returns the updated contact
func (r *Contacts) UpdateStatus(ctx context.Context, contactID int64, recipientUserID int64, status string) (*models.Contact, error) {
	query := `
		UPDATE contacts
		SET status = $1, responded_at = $2, updated_at = NOW()
		WHERE id = $3 AND recipient_user_id = $4
		RETURNING id, requester_user_id, recipient_user_id, status, requested_at, responded_at
	`

	var contact models.Contact
	err := r.db.QueryRowContext(ctx, query, status, time.Now(), contactID, recipientUserID).Scan(
		&contact.ID,
		&contact.RequesterUserID,
		&contact.RecipientUserID,
		&contact.Status,
		&contact.RequestedAt,
		&contact.RespondedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("contact not found or user is not the recipient")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to update contact status")
	}

	return &contact, nil
}

// Delete removes a contact relationship
func (r *Contacts) Delete(ctx context.Context, contactID int64, userID int64) error {
	query := `
		DELETE FROM contacts
		WHERE id = $1 AND (requester_user_id = $2 OR recipient_user_id = $2)
	`

	result, err := r.db.ExecContext(ctx, query, contactID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete contact")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("contact not found or user is not part of this contact")
	}

	return nil
}

// GetUserIDByCharacterName finds a user ID by their character name
func (r *Contacts) GetUserIDByCharacterName(ctx context.Context, characterName string) (int64, error) {
	query := `
		SELECT DISTINCT c.user_id
		FROM characters c
		JOIN users u ON c.user_id = u.id
		WHERE LOWER(c.name) = LOWER($1)
		LIMIT 1
	`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, characterName).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, errors.New("character not found - they may not have added their character to this tool yet")
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to find character")
	}

	return userID, nil
}

// GetUserIDByCharacterID finds a user ID by their character ID (from session)
func (r *Contacts) GetUserIDByCharacterID(ctx context.Context, characterID int64) (int64, error) {
	query := `
		SELECT user_id
		FROM characters
		WHERE id = $1
		LIMIT 1
	`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, characterID).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, errors.New("character not found in database")
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to find character by ID")
	}

	return userID, nil
}
