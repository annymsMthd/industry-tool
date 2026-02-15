package controllers

import (
	"net/http"
	"strconv"

	"github.com/annymsMthd/industry-tool/internal/repositories"
	"github.com/annymsMthd/industry-tool/internal/web"
	"github.com/pkg/errors"
)

type Analytics struct {
	repository *repositories.SalesAnalytics
}

func NewAnalytics(router Routerer, repository *repositories.SalesAnalytics) *Analytics {
	c := &Analytics{
		repository: repository,
	}

	router.RegisterRestAPIRoute("/v1/analytics/sales", web.AuthAccessUser, c.GetSalesMetrics, "GET")
	router.RegisterRestAPIRoute("/v1/analytics/top-items", web.AuthAccessUser, c.GetTopItems, "GET")
	router.RegisterRestAPIRoute("/v1/analytics/buyers", web.AuthAccessUser, c.GetBuyerAnalytics, "GET")
	router.RegisterRestAPIRoute("/v1/analytics/item-history", web.AuthAccessUser, c.GetItemSalesHistory, "GET")

	return c
}

// GetSalesMetrics returns sales metrics with time series data
func (c *Analytics) GetSalesMetrics(args *web.HandlerArgs) (any, *web.HttpError) {
	periodStr := args.Request.URL.Query().Get("period")
	period := parsePeriod(periodStr)

	metrics, err := c.repository.GetSalesMetrics(args.Request.Context(), *args.User, period)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: http.StatusInternalServerError,
			Error:      errors.Wrap(err, "failed to get sales metrics"),
		}
	}

	return metrics, nil
}

// GetTopItems returns the top selling items
func (c *Analytics) GetTopItems(args *web.HandlerArgs) (any, *web.HttpError) {
	periodStr := args.Request.URL.Query().Get("period")
	period := parsePeriod(periodStr)

	limitStr := args.Request.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	items, err := c.repository.GetTopItems(args.Request.Context(), *args.User, period, limit)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: http.StatusInternalServerError,
			Error:      errors.Wrap(err, "failed to get top items"),
		}
	}

	return items, nil
}

// GetBuyerAnalytics returns analytics about buyers
func (c *Analytics) GetBuyerAnalytics(args *web.HandlerArgs) (any, *web.HttpError) {
	periodStr := args.Request.URL.Query().Get("period")
	period := parsePeriod(periodStr)

	limitStr := args.Request.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	buyers, err := c.repository.GetBuyerAnalytics(args.Request.Context(), *args.User, period, limit)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: http.StatusInternalServerError,
			Error:      errors.Wrap(err, "failed to get buyer analytics"),
		}
	}

	return buyers, nil
}

// GetItemSalesHistory returns sales history for a specific item
func (c *Analytics) GetItemSalesHistory(args *web.HandlerArgs) (any, *web.HttpError) {
	typeIDStr := args.Request.URL.Query().Get("typeId")
	if typeIDStr == "" {
		return nil, &web.HttpError{
			StatusCode: http.StatusBadRequest,
			Error:      errors.New("typeId parameter is required"),
		}
	}

	typeID, err := strconv.ParseInt(typeIDStr, 10, 64)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: http.StatusBadRequest,
			Error:      errors.New("invalid typeId parameter"),
		}
	}

	periodStr := args.Request.URL.Query().Get("period")
	period := parsePeriod(periodStr)

	item, err := c.repository.GetItemSalesHistory(args.Request.Context(), *args.User, typeID, period)
	if err != nil {
		return nil, &web.HttpError{
			StatusCode: http.StatusInternalServerError,
			Error:      errors.Wrap(err, "failed to get item sales history"),
		}
	}

	return item, nil
}

// parsePeriod converts period string to number of days
// Supported formats: "7d", "30d", "90d", "1y", "all" (0 = all time)
func parsePeriod(periodStr string) int {
	switch periodStr {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	case "1y", "365d":
		return 365
	case "all", "":
		return 0
	default:
		// Try to parse as number followed by 'd'
		if len(periodStr) > 1 && periodStr[len(periodStr)-1] == 'd' {
			if days, err := strconv.Atoi(periodStr[:len(periodStr)-1]); err == nil && days > 0 {
				return days
			}
		}
		return 30 // Default to 30 days
	}
}
