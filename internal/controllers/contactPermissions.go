package controllers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/annymsMthd/industry-tool/internal/models"
	"github.com/annymsMthd/industry-tool/internal/web"
	"github.com/pkg/errors"
)

type ContactPermissionsRepository interface {
	GetByContact(ctx context.Context, contactID int64, userID int64) ([]*models.ContactPermission, error)
	Upsert(ctx context.Context, perm *models.ContactPermission) error
	GetUserPermissionsForService(ctx context.Context, viewerUserID int64, serviceType string) ([]int64, error)
	CheckPermission(ctx context.Context, grantingUserID, receivingUserID int64, serviceType string) (bool, error)
}

type ContactPermissions struct {
	repository ContactPermissionsRepository
}

func NewContactPermissions(router Routerer, repository ContactPermissionsRepository) *ContactPermissions {
	controller := &ContactPermissions{
		repository: repository,
	}

	router.RegisterRestAPIRoute("/v1/contacts/{id}/permissions", web.AuthAccessUser, controller.GetPermissions, "GET")
	router.RegisterRestAPIRoute("/v1/contacts/{id}/permissions", web.AuthAccessUser, controller.UpdatePermission, "POST")

	return controller
}

func (c *ContactPermissions) GetPermissions(args *web.HandlerArgs) (any, *web.HttpError) {
	contactIDStr := args.Params["id"]
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "invalid contact ID"),
		}
	}

	permissions, err := c.repository.GetByContact(args.Request.Context(), contactID, *args.User)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to get permissions"),
		}
	}

	return permissions, nil
}

func (c *ContactPermissions) UpdatePermission(args *web.HandlerArgs) (any, *web.HttpError) {
	contactIDStr := args.Params["id"]
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "invalid contact ID"),
		}
	}

	d := json.NewDecoder(args.Request.Body)
	var req struct {
		ServiceType     string `json:"serviceType"`
		CanAccess       bool   `json:"canAccess"`
		ReceivingUserID int64  `json:"receivingUserId"`
	}
	err = d.Decode(&req)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 400,
			Error:      errors.Wrap(err, "failed to decode json"),
		}
	}

	// User is granting permission to another user
	perm := &models.ContactPermission{
		ContactID:       contactID,
		GrantingUserID:  *args.User,
		ReceivingUserID: req.ReceivingUserID,
		ServiceType:     req.ServiceType,
		CanAccess:       req.CanAccess,
	}

	err = c.repository.Upsert(args.Request.Context(), perm)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      errors.Wrap(err, "failed to update permission"),
		}
	}

	return nil, nil
}
