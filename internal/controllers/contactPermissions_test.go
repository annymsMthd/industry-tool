package controllers_test

import (
	"bytes"
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

func Test_ContactPermissionsController_GetPermissions_Success(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	contactID := int64(1)

	expectedPermissions := []*models.ContactPermission{
		{
			ID:              1,
			ContactID:       contactID,
			GrantingUserID:  123,
			ReceivingUserID: 456,
			ServiceType:     "for_sale_browse",
			CanAccess:       true,
		},
		{
			ID:              2,
			ContactID:       contactID,
			GrantingUserID:  456,
			ReceivingUserID: 123,
			ServiceType:     "for_sale_browse",
			CanAccess:       false,
		},
	}

	mockRepo.On("GetByContact", mock.Anything, contactID, userID).Return(expectedPermissions, nil)

	req := httptest.NewRequest("GET", "/v1/contacts/1/permissions", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.GetPermissions(args)

	assert.Nil(t, httpErr)
	assert.NotNil(t, result)

	permissions := result.([]*models.ContactPermission)
	assert.Len(t, permissions, 2)
	assert.Equal(t, int64(1), permissions[0].ID)
	assert.True(t, permissions[0].CanAccess)
	assert.False(t, permissions[1].CanAccess)

	mockRepo.AssertExpectations(t)
}

func Test_ContactPermissionsController_GetPermissions_InvalidContactID(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	req := httptest.NewRequest("GET", "/v1/contacts/invalid/permissions", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "invalid"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.GetPermissions(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func Test_ContactPermissionsController_GetPermissions_RepositoryError(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	contactID := int64(1)

	mockRepo.On("GetByContact", mock.Anything, contactID, userID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/v1/contacts/1/permissions", nil)
	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.GetPermissions(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 500, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ContactPermissionsController_UpdatePermission_Success(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)
	contactID := int64(1)
	receivingUserID := int64(456)
	serviceType := "for_sale_browse"
	canAccess := true

	expectedPerm := &models.ContactPermission{
		ContactID:       contactID,
		GrantingUserID:  userID,
		ReceivingUserID: receivingUserID,
		ServiceType:     serviceType,
		CanAccess:       canAccess,
	}

	mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(perm *models.ContactPermission) bool {
		return perm.ContactID == expectedPerm.ContactID &&
			perm.GrantingUserID == expectedPerm.GrantingUserID &&
			perm.ReceivingUserID == expectedPerm.ReceivingUserID &&
			perm.ServiceType == expectedPerm.ServiceType &&
			perm.CanAccess == expectedPerm.CanAccess
	})).Return(nil)

	body := map[string]interface{}{
		"serviceType":     serviceType,
		"receivingUserId": receivingUserID,
		"canAccess":       canAccess,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts/1/permissions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, httpErr)
	assert.Nil(t, result) // Controller returns nil on success

	mockRepo.AssertExpectations(t)
}

func Test_ContactPermissionsController_UpdatePermission_InvalidContactID(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	body := map[string]interface{}{
		"serviceType":     "for_sale_browse",
		"receivingUserId": 456,
		"canAccess":       true,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts/invalid/permissions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "invalid"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func Test_ContactPermissionsController_UpdatePermission_InvalidJSON(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	req := httptest.NewRequest("POST", "/v1/contacts/1/permissions", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 400, httpErr.StatusCode)
}

func Test_ContactPermissionsController_UpdatePermission_RepositoryError(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("Upsert", mock.Anything, mock.Anything).Return(errors.New("database error"))

	body := map[string]interface{}{
		"serviceType":     "for_sale_browse",
		"receivingUserId": 456,
		"canAccess":       true,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts/1/permissions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, result)
	assert.NotNil(t, httpErr)
	assert.Equal(t, 500, httpErr.StatusCode)

	mockRepo.AssertExpectations(t)
}

func Test_ContactPermissionsController_UpdatePermission_GrantAccess(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(perm *models.ContactPermission) bool {
		return perm.CanAccess == true
	})).Return(nil)

	body := map[string]interface{}{
		"serviceType":     "for_sale_browse",
		"receivingUserId": 456,
		"canAccess":       true,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts/1/permissions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, httpErr)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func Test_ContactPermissionsController_UpdatePermission_RevokeAccess(t *testing.T) {
	mockRepo := new(MockContactPermissionsRepository)
	mockRouter := &MockRouter{}

	userID := int64(123)

	mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(perm *models.ContactPermission) bool {
		return perm.CanAccess == false
	})).Return(nil)

	body := map[string]interface{}{
		"serviceType":     "for_sale_browse",
		"receivingUserId": 456,
		"canAccess":       false,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/contacts/1/permissions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	args := &web.HandlerArgs{
		Request: req,
		User:    &userID,
		Params:  map[string]string{"id": "1"},
	}

	controller := controllers.NewContactPermissions(mockRouter, mockRepo)
	result, httpErr := controller.UpdatePermission(args)

	assert.Nil(t, httpErr)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}
