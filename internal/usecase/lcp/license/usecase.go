package license

import (
	"context"
	"time"

	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
	"github.com/google/uuid"
	"github.com/readium/readium-lcp-server/lcpserver"
)

type LicenseUsecase interface {
	Create(ctx context.Context, input *lcp.LicenseInput) (*lcp.License, error)
	GetByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error)
	Revoke(ctx context.Context, id string) error
}

type licenseUsecase struct {
	repo lcp.LicenseRepository
	lcp  *lcpserver.LCPServer
}

func NewLicenseUsecase(repo lcp.LicenseRepository, lcp *lcpserver.LCPServer) LicenseUsecase {
	return &licenseUsecase{repo, lcp}
}

func (u *licenseUsecase) Create(ctx context.Context, input *lcp.LicenseInput) (*lcp.License, error) {
	license := &lcp.License{
		ID:             uuid.New().String(),
		PublicationID:  input.PublicationID,
		UserID:         input.UserID,
		Passphrase:     input.Passphrase,
		Hint:           input.Hint,
		PublicationURL: "http://yourdomain.com/storage/" + input.PublicationID,
		RightPrint:     input.RightPrint,
		RightCopy:      input.RightCopy,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		CreatedAt:      time.Now(),
	}

	// Generate LCP license using lcpserver
	err := u.lcp.GenerateLicense(license)
	if err != nil {
		return nil, err
	}

	// Save license to database
	err = u.repo.Save(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

func (u *licenseUsecase) GetByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error) {
	return u.repo.FindByPublication(ctx, publicationID)
}

func (u *licenseUsecase) Revoke(ctx context.Context, id string) error {
	return u.lcp.RevokeLicense(id)
}
