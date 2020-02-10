package util

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/tespo/satya/v2/types"
)

//
// GetPaginationParams returns the valuse for query params
// related to pagination
//
func GetPaginationParams(r *http.Request) (int, int) {
	limitParam := r.URL.Query()["limit"]
	pageParam := r.URL.Query()["page"]
	limit, page := 25, 0
	if len(limitParam) == 1 {
		limitAmount, err := strconv.Atoi(r.URL.Query()["limit"][0])
		if err == nil {
			limit = limitAmount
		}
	}
	if len(pageParam) == 1 {
		pageAmount, err := strconv.Atoi(r.URL.Query()["page"][0])
		if err == nil {
			page = pageAmount
		}
	}
	return limit, page
}

//
// PaginationResponder builds a pagination response with a
// given data slice
//
func PaginationResponder(w http.ResponseWriter, r *http.Request, data interface{}) {
	val := reflect.ValueOf(data)
	var length int
	if val.Kind() == reflect.Slice {
		length = val.Len()
	}

	limit, page := GetPaginationParams(r)
	var next, previous string
	if page > 0 {
		previous = r.URL.Host + r.URL.Path + "?limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page-1)
	}
	if length == limit {
		next = r.URL.Host + r.URL.Path + "?limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page+1)
	}
	paginatedReponse := types.PaginatedResponse{
		Data: data,
		Pagination: types.Pagination{
			Next:     next,
			Previous: previous,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paginatedReponse); err != nil {
		panic(err)
	}
}

//
// SetDBPagination returns a gorm.DB with offset and limit set
//
func SetDBPagination(db *gorm.DB, r *http.Request) *gorm.DB {
	limit, page := GetPaginationParams(r)
	return db.Offset(page * limit).Limit(limit)
}
