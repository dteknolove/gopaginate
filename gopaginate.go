package gopaginate

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
)

func New(db *pgxpool.Pool, sql string, r *http.Request) (int16, int16, int16, int16, int16, int16, int16, error) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	var errPagination error
	if pageParam == "" || pageParam == "0" {
		pageParam = "1"
	}
	if limitParam == "" || limitParam == "0" {
		limitParam = "10"
	}

	currentPage, errParsePage := strconv.ParseInt(pageParam, 10, 64)
	if errParsePage != nil {
		errPagination = errors.New("error parse page")
	}
	limit, errParseLimit := strconv.ParseInt(limitParam, 10, 64)
	if errParseLimit != nil {
		errPagination = errors.New("error parse limit")
	}
	offSet := (currentPage - 1) * limit

	totalCount, errCount := getTotalUsersCount(r.Context(), db, sql)
	if errCount != nil {
		errPagination = errors.New("error parse counting page")
	}

	totalPages := calculateTotalPages(totalCount, int16(limit))
	nextPage := currentPage + 1
	prevPage := currentPage - 1
	isNextPage := getNextPage(totalPages, int16(nextPage))
	isPrevPage := getPrevPage(int16(prevPage))
	totalData := totalCount

	return int16(currentPage), int16(limit), int16(offSet), totalPages, isNextPage, isPrevPage, totalData, errPagination
}
func getTotalUsersCount(ctx context.Context, db *pgxpool.Pool, sql string) (int16, error) {
	var count int16
	errCount := db.QueryRow(ctx, sql).Scan(&count)
	return count, errCount
}

func calculateTotalPages(totalCount, limit int16) int16 {
	if totalCount/limit == 0 {
		return totalCount / limit
	}
	return totalCount/limit + 1
}

func getNextPage(totalPages, nextPage int16) int16 {
	if nextPage <= totalPages {
		return nextPage
	}
	return -1
}

func getPrevPage(prevPage int16) int16 {
	if prevPage > 0 {
		return prevPage
	}
	return -1
}
