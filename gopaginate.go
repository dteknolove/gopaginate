package gopaginate

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"regexp"
	"strconv"
)

func New(db *pgxpool.Pool, sql string, r *http.Request) (int16, int16, int16, int16, int16, int16, int16, error) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	cleanPage := regexp.MustCompile("[^0-9]+")
	cleanLimit := regexp.MustCompile("[^0-9]+")

	finalPageParam := cleanPage.ReplaceAllString(pageParam, "")
	finalLimitParam := cleanLimit.ReplaceAllString(limitParam, "")

	var errPagination error
	if finalPageParam == "" || finalPageParam == "0" {
		finalPageParam = "1"
	}
	if finalLimitParam == "" || finalLimitParam == "0" {
		finalLimitParam = "10"
	}

	currentPage, errParsePage := strconv.ParseInt(finalPageParam, 10, 64)
	if errParsePage != nil {
		errPagination = errors.New("error parse page")
	}
	limit, errParseLimit := strconv.ParseInt(finalLimitParam, 10, 64)
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
	finalCurrentPage := int16(currentPage)
	finalLimit := int16(limit)
	finalOffset := int16(offSet)

	return finalCurrentPage, finalLimit, finalOffset, totalPages, isNextPage, isPrevPage, totalData, errPagination
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
