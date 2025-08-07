package util

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/peconote/peconote/internal/domain/model"
)

func BuildLinkHeader(base string, p model.Pagination, tag *string) string {
	var links []string
	tagParam := ""
	if tag != nil {
		tagParam = "&tag=" + url.QueryEscape(*tag)
	}
	if p.Page < p.TotalPages {
		next := fmt.Sprintf("%s?page=%d&page_size=%d%s", base, p.Page+1, p.PageSize, tagParam)
		links = append(links, fmt.Sprintf("<%s>; rel=\"next\"", next))
	}
	if p.Page > 1 {
		prev := fmt.Sprintf("%s?page=%d&page_size=%d%s", base, p.Page-1, p.PageSize, tagParam)
		links = append(links, fmt.Sprintf("<%s>; rel=\"prev\"", prev))
	}
	return strings.Join(links, ", ")
}
