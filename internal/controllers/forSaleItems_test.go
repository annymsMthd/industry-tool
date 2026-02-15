package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/annymsMthd/industry-tool/internal/controllers"
	"github.com/annymsMthd/industry-tool/internal/models"
	"github.com/annymsMthd/industry-tool/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ForSaleItemsRepository
type MockForSaleItemsRepository struct {
	mock.Mock
}

func (m *MockForSaleItemsRepository) GetByUser(ctx context.Context, userID int64) ([]*models.ForSaleItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ForSaleItem), args.Error(1)
}

func (m *MockForSaleItemsRepository) GetBrowsableItems(ctx context.Context, buyerUserID int64, sellerUserIDs []int64) ([]*models.ForSaleItem, error) {
	args := m.Called(ctx, buyerUserID, sellerUserIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ForSaleItem), args.Error(1)
}

func (m *MockForSaleItemsRepository) Upsert(ctx context.Context, item *models.ForSaleItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockForSaleItemsRepository) Delete(ctx context.Context, itemID int64, userID int64) error {
	args := m.Called(ctx, itemID, userID)
	return args.Error(0)
}

func (m *MockForSaleItemsRepository) GetByID(ctx context.Context, itemID int64) (*models.ForSaleItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ForSaleItem), args.Error(1)
}

func (m *MockForSaleItemsRepository) GetUserIDByCharacterID(ctx context.Context, characterID int64) (int64, error) {
	args := m.Called(ctx, characterID)
	return args.Get(0).(int64), args.Error(1)
}

func Test_ForSaleItemsController_GetMyListings_Success(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	expectedItems := []*models.ForSaleItem{
		{
			ID:                1,
			UserID:            userID,
			TypeID:            34,
			TypeName:          "Tritanium",
			OwnerType:         "character",
			OwnerID:           456,
			OwnerName:         "Test Character",
			LocationID:        30000142,
			LocationName:      "Jita",
			QuantityAvailable: 1000,
			PricePerUnit:      50,
			IsActive:          true,
		},
	}

	mockRepo.On("GetByUser", mock.Anything, userID).Return(expectedItems, nil)

	req := httptest.NewRequest("GET", "/v1/for-sale", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.GetMyListings(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	items := result.([]*models.ForSaleItem)
	assert.Len(t, items, 1)
	assert.Equal(t, "Tritanium", items[0].TypeName)
	assert.Equal(t, int64(1000), items[0].QuantityAvailable)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_GetMyListings_RepositoryError(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("GetByUser", mock.Anything, userID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/v1/for-sale", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.GetMyListings(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 500, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_CreateListing_Success(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(item *models.ForSaleItem) bool {
		return item.UserID == userID &&
			item.TypeID == 34 &&
			item.QuantityAvailable == 1000 &&
			item.PricePerUnit == 50
	})).Return(nil)

	body := map[string]interface{}{
		"typeId":            34,
		"ownerType":         "character",
		"ownerId":           456,
		"locationId":        30000142,
		"quantityAvailable": 1000,
		"pricePerUnit":      50,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/for-sale", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.CreateListing(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_CreateListing_InvalidJSON(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	req := httptest.NewRequest("POST", "/v1/for-sale", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.CreateListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func Test_ForSaleItemsController_CreateListing_MissingRequiredFields(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)


	// Missing typeId
	body := map[string]interface{}{
		"ownerType":         "character",
		"ownerId":           456,
		"locationId":        30000142,
		"quantityAvailable": 1000,
		"pricePerUnit":      50,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/for-sale", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.CreateListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
	assert.Contains(t, httpErr.Error.Error(), "typeId is required")
}

func Test_ForSaleItemsController_CreateListing_InvalidQuantity(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)


	// Quantity <= 0
	body := map[string]interface{}{
		"typeId":            34,
		"ownerType":         "character",
		"ownerId":           456,
		"locationId":        30000142,
		"quantityAvailable": 0,
		"pricePerUnit":      50,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/for-sale", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.CreateListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
	assert.Contains(t, httpErr.Error.Error(), "quantityAvailable must be greater than 0")
}

func Test_ForSaleItemsController_UpdateListing_Success(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	itemID := int64(1)

	existingItem := &models.ForSaleItem{
		ID:                itemID,
		UserID:            userID,
		TypeID:            34,
		TypeName:          "Tritanium",
		OwnerType:         "character",
		OwnerID:           456,
		LocationID:        30000142,
		QuantityAvailable: 1000,
		PricePerUnit:      50,
		IsActive:          true,
	}

	mockRepo.On("GetByID", mock.Anything, itemID).Return(existingItem, nil)
	mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(item *models.ForSaleItem) bool {
		return item.ID == itemID &&
			item.QuantityAvailable == 2000 &&
			item.PricePerUnit == 75
	})).Return(nil)

	body := map[string]interface{}{
		"quantityAvailable": 2000,
		"pricePerUnit":      75,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/v1/for-sale/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.UpdateListing(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_UpdateListing_NotFound(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	itemID := int64(999)

	mockRepo.On("GetByID", mock.Anything, itemID).Return(nil, errors.New("for-sale item not found"))

	body := map[string]interface{}{
		"quantityAvailable": 2000,
		"pricePerUnit":      75,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/v1/for-sale/999", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "999"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.UpdateListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 404, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_UpdateListing_NotOwner(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	otherUserID := int64(999)
	itemID := int64(1)

	existingItem := &models.ForSaleItem{
		ID:                itemID,
		UserID:            otherUserID, // Different owner
		TypeID:            34,
		QuantityAvailable: 1000,
		PricePerUnit:      50,
	}

	mockRepo.On("GetByID", mock.Anything, itemID).Return(existingItem, nil)

	body := map[string]interface{}{
		"quantityAvailable": 2000,
		"pricePerUnit":      75,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/v1/for-sale/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.UpdateListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 403, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_DeleteListing_Success(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	itemID := int64(1)

	mockRepo.On("Delete", mock.Anything, itemID, userID).Return(nil)

	req := httptest.NewRequest("DELETE", "/v1/for-sale/1", nil)

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.DeleteListing(args)

	assert.Nil(t, httpErr)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_DeleteListing_NotFound(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	itemID := int64(999)

	mockRepo.On("Delete", mock.Anything, itemID, userID).Return(errors.New("for-sale item not found or user is not the owner"))

	req := httptest.NewRequest("DELETE", "/v1/for-sale/999", nil)

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "999"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.DeleteListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 404, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ForSaleItemsController_DeleteListing_InvalidID(t *testing.T) {
	mockRepo := new(MockForSaleItemsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)


	req := httptest.NewRequest("DELETE", "/v1/for-sale/invalid", nil)

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "invalid"},
	}

	controller := controllers.NewForSaleItems(mockRouter, mockRepo, &MockContactPermissionsRepository{})
	result, httpErr := controller.DeleteListing(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}
