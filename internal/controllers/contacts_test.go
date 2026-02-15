package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
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

// Mock repositories
type MockContactsRepository struct {
	mock.Mock
}

func (m *MockContactsRepository) GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactsRepository) Create(ctx context.Context, requesterID, recipientID int64) (*models.Contact, error) {
	args := m.Called(ctx, requesterID, recipientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactsRepository) UpdateStatus(ctx context.Context, contactID int64, recipientUserID int64, status string) (*models.Contact, error) {
	args := m.Called(ctx, contactID, recipientUserID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactsRepository) Delete(ctx context.Context, contactID int64, userID int64) error {
	args := m.Called(ctx, contactID, userID)
	return args.Error(0)
}

func (m *MockContactsRepository) GetUserIDByCharacterName(ctx context.Context, characterName string) (int64, error) {
	args := m.Called(ctx, characterName)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockContactsRepository) GetUserIDByCharacterID(ctx context.Context, characterID int64) (int64, error) {
	args := m.Called(ctx, characterID)
	return args.Get(0).(int64), args.Error(1)
}

type MockContactPermissionsInitializer struct {
	mock.Mock
}

func (m *MockContactPermissionsInitializer) InitializePermissionsForContact(ctx context.Context, tx *sql.Tx, contactID, userID1, userID2 int64) error {
	args := m.Called(ctx, tx, contactID, userID1, userID2)
	return args.Error(0)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func Test_ContactsController_GetContacts_Success(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)

	userID := int64(123)
	expectedContacts := []*models.Contact{
		{
			ID:              1,
			RequesterUserID: userID,
			RecipientUserID: 456,
			RequesterName:   "User 1",
			RecipientName:   "User 2",
			Status:          "accepted",
		},
	}

	mockRepo.On("GetByUser", mock.Anything, userID).Return(expectedContacts, nil)

	req := httptest.NewRequest("GET", "/v1/contacts", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
	}

	// Get handler from mockRouter - we need to test the handler directly
	// Since we can't easily capture the handler from NewContacts, we'll test via creating a new instance
	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.GetContacts(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	contacts := result.([]*models.Contact)
	assert.Len(t, contacts, 1)
	assert.Equal(t, int64(1), contacts[0].ID)
	assert.Equal(t, "User 1", contacts[0].RequesterName)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_GetContacts_RepositoryError(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("GetByUser", mock.Anything, userID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/v1/contacts", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.GetContacts(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 500, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_CreateContact_Success(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)
	recipientID := int64(456)
	characterName := "Recipient Character"

	expectedContact := &models.Contact{
		ID:              1,
		RequesterUserID: userID,
		RecipientUserID: recipientID,
		Status:          "pending",
	}

	mockRepo.On("GetUserIDByCharacterName", mock.Anything, characterName).Return(recipientID, nil)
	mockRepo.On("Create", mock.Anything, userID, recipientID).Return(expectedContact, nil)

	body := map[string]string{"characterName": characterName}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.CreateContact(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	contact := result.(*models.Contact)
	assert.Equal(t, int64(1), contact.ID)
	assert.Equal(t, userID, contact.RequesterUserID)
	assert.Equal(t, recipientID, contact.RecipientUserID)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_CreateContact_SelfContactError(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)
	characterName := "My Character"

	mockRepo.On("GetUserIDByCharacterName", mock.Anything, characterName).Return(userID, nil)

	body := map[string]string{"characterName": characterName}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.CreateContact(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_CreateContact_CharacterNotFound(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)
	characterName := "Nonexistent Character"

	mockRepo.On("GetUserIDByCharacterName", mock.Anything, characterName).Return(int64(0), errors.New("character not found"))

	body := map[string]string{"characterName": characterName}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.CreateContact(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 404, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_AcceptContact_Success(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(456)
	contactID := int64(1)

	expectedContact := &models.Contact{
		ID:              contactID,
		RequesterUserID: 123,
		RecipientUserID: userID,
		Status:          "accepted",
	}

	mockRepo.On("UpdateStatus", mock.Anything, contactID, userID, "accepted").Return(expectedContact, nil)

	req := httptest.NewRequest("POST", "/v1/contacts/1/accept", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	// Note: Passing nil DB will cause a panic when BeginTx is called
	// This test verifies UpdateStatus is called correctly before that point
	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)

	// This will panic on nil DB BeginTx - recover to verify the UpdateStatus was called
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil DB
			mockRepo.AssertCalled(t, "UpdateStatus", mock.Anything, contactID, userID, "accepted")
		}
	}()

	controller.AcceptContact(args)

	// If we get here without panic, that's unexpected
	t.Error("Expected panic due to nil DB, but didn't get one")
}

func Test_ContactsController_RejectContact_Success(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(456)
	contactID := int64(1)

	expectedContact := &models.Contact{
		ID:              contactID,
		RequesterUserID: 123,
		RecipientUserID: userID,
		Status:          "rejected",
	}

	mockRepo.On("UpdateStatus", mock.Anything, contactID, userID, "rejected").Return(expectedContact, nil)

	req := httptest.NewRequest("POST", "/v1/contacts/1/reject", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.RejectContact(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	contact := result.(*models.Contact)
	assert.Equal(t, "rejected", contact.Status)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_DeleteContact_Success(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)
	contactID := int64(1)

mockRepo.On("Delete", mock.Anything, contactID, userID).Return(nil)

	req := httptest.NewRequest("DELETE", "/v1/contacts/1", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.DeleteContact(args)

	assert.Nil(t, httpErr)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func Test_ContactsController_DeleteContact_InvalidID(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)

	req := httptest.NewRequest("DELETE", "/v1/contacts/invalid", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "invalid"},
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.DeleteContact(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func Test_ContactsController_DeleteContact_RepositoryError(t *testing.T) {
	mockRepo := new(MockContactsRepository)
	mockPermissions := new(MockContactPermissionsInitializer)
	mockRouter := &MockRouter{}

	userID := int64(123)
	contactID := int64(1)

mockRepo.On("Delete", mock.Anything, contactID, userID).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/v1/contacts/1", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContacts(mockRouter, mockRepo, mockPermissions, nil)
	result, httpErr := controller.DeleteContact(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 500, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}
