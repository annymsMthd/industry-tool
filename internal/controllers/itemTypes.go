package controllers

import (
	"context"

	"github.com/annymsMthd/industry-tool/internal/models"
	"github.com/annymsMthd/industry-tool/internal/web"
)

type ItemTypeRepository interface {
	SearchItemTypes(ctx context.Context, query string, limit int) ([]models.EveInventoryType, error)
	GetItemTypeByName(ctx context.Context, typeName string) (*models.EveInventoryType, error)
}

type ItemTypes struct {
	repository ItemTypeRepository
}

func NewItemTypes(router Routerer, repository ItemTypeRepository) *ItemTypes {
	c := &ItemTypes{
		repository: repository,
	}

	router.RegisterRestAPIRoute("/v1/item-types/search", web.AuthAccessUser, c.SearchItemTypes, "GET")

	return c
}

// SearchItemTypes searches for item types by name
func (c *ItemTypes) SearchItemTypes(args *web.HandlerArgs) (any, *web.HttpError) {
	query := args.Request.URL.Query().Get("q")
	if query == "" {
		return []models.EveInventoryType{}, nil
	}

	items, err := c.repository.SearchItemTypes(args.Request.Context(), query, 20)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: 500,
			Error:      err,
		}
	}

	if items == nil {
		items = []models.EveInventoryType{}
	}

	return items, nil
}
