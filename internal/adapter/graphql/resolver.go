package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/yourusername/lcp-project/internal/domain/lcp"
	"github.com/yourusername/lcp-project/internal/usecase/lcp/license"
	"github.com/yourusername/lcp-project/internal/usecase/lcp/publication"
)

// Resolver is the main resolver struct
type Resolver struct {
	PublicationUsecase publication.PublicationUsecase
	LicenseUsecase     license.LicenseUsecase
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Publications(ctx context.Context) ([]*lcp.Publication, error) {
	return r.PublicationUsecase.GetAll(ctx)
}

func (r *queryResolver) Licenses(ctx context.Context, publicationID *string) ([]*lcp.License, error) {
	return r.LicenseUsecase.GetByPublication(ctx, publicationID)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) UploadPublication(ctx context.Context, title string, file graphql.Upload) (*lcp.Publication, error) {
	return r.PublicationUsecase.UploadAndEncrypt(ctx, title, file.File)
}

func (r *mutationResolver) CreateLicense(ctx context.Context, input struct {
	PublicationID string
	UserID        string
	Passphrase    string
	Hint          string
	RightPrint    *int
	RightCopy     *int
	StartDate     *string
	EndDate       *string
}) (*lcp.License, error) {
	return r.LicenseUsecase.Create(ctx, &lcp.LicenseInput{
		PublicationID: input.PublicationID,
		UserID:        input.UserID,
		Passphrase:    input.Passphrase,
		Hint:          input.Hint,
		RightPrint:    input.RightPrint,
		RightCopy:     input.RightCopy,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
	})
}

func (r *mutationResolver) RevokeLicense(ctx context.Context, id string) (bool, error) {
	err := r.LicenseUsecase.Revoke(ctx, id)
	return err == nil, err
}
