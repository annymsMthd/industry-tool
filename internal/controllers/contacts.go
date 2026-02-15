package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/annymsMthd/industry-tool/internal/models"
	"github.com/annymsMthd/industry-tool/internal/web"
	"github.com/pkg/errors"
)

type ContactsRepository interface {
	GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error)
	Create(ctx context.Context, requesterID, recipientID int64) (*models.Contact, error)
	UpdateStatus(ctx context.Context, contactID int64, recipientUserID int64, status string) (*models.Contact, error)
	Delete(ctx context.Context, contactID int64, userID int64) error
	GetUserIDByCharacterName(ctx context.Context, characterName string) (int64, error)
	GetUserIDByCharacterID(ctx context.Context, characterID int64) (int64, error)
}

type ContactPermissionsInitializer interface {
	InitializePermissionsForContact(ctx context.Context, tx *sql.Tx, contactID, userID1, userID2 int64) error
}

type Contacts struct {
	repository            ContactsRepository
	permissionsRepository ContactPermissionsInitializer
	db                    *sql.DB
}

func NewContacts(router Routerer, repository ContactsRepository, permissionsRepository ContactPermissionsInitializer, db *sql.DB) *Contacts {
	controller := &Contacts{
		repository:            repository,
		permissionsRepository: permissionsRepository,
		db:                    db,
	}

	router.RegisterRestAPIRoute("/v1/contacts", web.AuthAccessUser, controller.GetContacts, "GET")
	router.RegisterRestAPIRoute("/v1/contacts", web.AuthAccessUser, controller.CreateContact, "POST")
	router.RegisterRestAPIRoute("/v1/contacts/{id}/accept", web.AuthAccessUser, controller.AcceptContact, "POST")
	router.RegisterRestAPIRoute("/v1/contacts/{id}/reject", web.AuthAccessUser, controller.RejectContact, "POST")
	router.RegisterRestAPIRoute("/v1/contacts/{id}", web.AuthAccessUser, controller.DeleteContact, "DELETE")

	return controller
}

// getUserID converts the session user (which is the main character's EVE ID) to a user ID
// Since users.id IS the main character's EVE character ID, we can use it directly
func (c *Contacts) getUserID(ctx context.Context, sessionUser int64) (int64, error) {
	// The session user is the main character's EVE ID, which is also the user ID
	return sessionUser, nil
}

func (c *Contacts) GetContacts(args *web.HandlerArgs) (any, *web.HttpError) {
	userID, err := c.getUserID(args.Request.Context(), *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 401,
			Error:      errors.Wrap(err, "failed to identify user"),
		}
	}

	contacts, err := c.repository.GetByUser(args.Request.Context(), userID)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to get contacts"),
		}
	}

	return contacts, nil
}

func (c *Contacts) CreateContact(args *web.HandlerArgs) (any, *web.HttpError) {
	d := json.NewDecoder(args.Request.Body)
	var req struct {
		CharacterName string `json:"characterName"`
	}
	err := d.Decode(&req)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "failed to decode json"),
		}
	}

	// Get requester user ID
	requesterUserID, err := c.getUserID(args.Request.Context(), *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 401,
			Error:      errors.Wrap(err, "failed to identify user"),
		}
	}

	// Look up recipient user ID by character name
	recipientUserID, err := c.repository.GetUserIDByCharacterName(args.Request.Context(), req.CharacterName)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 404,
			Error:      errors.Wrap(err, "character not found"),
		}
	}

	// Prevent self-contact
	if recipientUserID == requesterUserID {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.New("cannot add yourself as a contact"),
		}
	}

	contact, err := c.repository.Create(args.Request.Context(), requesterUserID, recipientUserID)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to create contact request"),
		}
	}

	return contact, nil
}

func (c *Contacts) AcceptContact(args *web.HandlerArgs) (any, *web.HttpError) {
	contactIDStr := args.Params["id"]
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "invalid contact ID"),
		}
	}

	// Get user ID
	userID, err := c.getUserID(args.Request.Context(), *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 401,
			Error:      errors.Wrap(err, "failed to identify user"),
		}
	}

	// Update contact status
	contact, err := c.repository.UpdateStatus(args.Request.Context(), contactID, userID, "accepted")
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to accept contact"),
		}
	}

	// Initialize permissions for both users (in separate transaction)
	tx, err := c.db.BeginTx(args.Request.Context(), nil)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to begin transaction for permissions"),
		}
	}
	defer tx.Rollback()

	err = c.permissionsRepository.InitializePermissionsForContact(
		args.Request.Context(),
		tx,
		contact.ID,
		contact.RequesterUserID,
		contact.RecipientUserID,
	)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to initialize permissions"),
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to commit permissions transaction"),
		}
	}

	return contact, nil
}

func (c *Contacts) RejectContact(args *web.HandlerArgs) (any, *web.HttpError) {
	contactIDStr := args.Params["id"]
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "invalid contact ID"),
		}
	}

	// Get user ID
	userID, err := c.getUserID(args.Request.Context(), *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 401,
			Error:      errors.Wrap(err, "failed to identify user"),
		}
	}

	contact, err := c.repository.UpdateStatus(args.Request.Context(), contactID, userID, "rejected")
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to reject contact"),
		}
	}

	return contact, nil
}

func (c *Contacts) DeleteContact(args *web.HandlerArgs) (any, *web.HttpError) {
	contactIDStr := args.Params["id"]
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "invalid contact ID"),
		}
	}

	// Get user ID
	userID, err := c.getUserID(args.Request.Context(), *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 401,
			Error:      errors.Wrap(err, "failed to identify user"),
		}
	}

	err = c.repository.Delete(args.Request.Context(), contactID, userID)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to delete contact"),
		}
	}

	return nil, nil
}
