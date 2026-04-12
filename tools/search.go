package tools

import "context"

type ParsedSearchResult struct {
	Title    string
	Content  string
	Comments []string
}

type SearchProvider interface {
	Search(ctx context.Context, query string) ([]ParsedSearchResult, error)
}
