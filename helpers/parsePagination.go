package helpers

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ParsePagination menguraikan parameter halaman dan per halaman dari query
func ParsePagination(c echo.Context) (perPage int, page int, offset int, sort string) {
	var err error

	// Menguraikan per_page
	perPage, err = strconv.Atoi(c.QueryParam("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 10
		log.Println("Failed to parse per_page query parameter or per_page <= 0. Defaulting to 10")
	}

	// Menguraikan page
	page, err = strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
		log.Println("Failed to parse page query parameter or page <= 0. Defaulting to 1")
	}

	// Menghitung offset berdasarkan halaman dan jumlah item per halaman
	offset = (page - 1) * perPage

	// Menguraikan sort
	sort = c.QueryParam("sort")
	if sort == "" {
		sort = "created_at DESC"
	}

	return perPage, page, offset, sort
}
