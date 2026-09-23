package helpers

import (
	"net/http"
	"strconv"
)

type paginationParams struct {
	Page  int32
	Limit int32
}

func ParsePaginationParams(r *http.Request) paginationParams {
	queryParams := r.URL.Query()

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1 // default
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // default
	}

	return paginationParams{
		Page:  int32(page),
		Limit: int32(limit),
	}
}

func GetPathValue(r *http.Request, key string) string {
	return r.PathValue(key)
}
