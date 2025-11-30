package graphql

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
)

// Upload mirrors the gqlgen Upload scalar.
type Upload struct {
	File     http.File
	Filename string
	Size     int64
}

// QueryResolver exposes read operations.
type QueryResolver interface {
	Publications(ctx context.Context) ([]*lcp.Publication, error)
	Licenses(ctx context.Context, publicationID *string) ([]*lcp.License, error)
}

// MutationResolver exposes write operations.
type MutationResolver interface {
	UploadPublication(ctx context.Context, title string, file Upload) (*lcp.Publication, error)
	CreateLicense(ctx context.Context, input struct {
		PublicationID string
		UserID        string
		Passphrase    string
		Hint          string
		RightPrint    *int
		RightCopy     *int
		StartDate     *string
		EndDate       *string
	}) (*lcp.License, error)
	RevokeLicense(ctx context.Context, id string) (bool, error)
}

// ResolverRoot collects the resolvers used by the handler.
type ResolverRoot interface {
	Query() QueryResolver
	Mutation() MutationResolver
}

// NewHandler returns a minimal HTTP handler that defers to the provided
// resolvers. The current implementation returns a 501 status because the
// GraphQL transport layer is intentionally lightweight for this milestone.
func NewHandler(resolver ResolverRoot) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "GraphQL handler not implemented yet"})
	})
}
