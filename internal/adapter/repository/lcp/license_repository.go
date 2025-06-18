package lcp

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Mehrbod2002//internal/domain/lcp"
	"github.com/jackc/pgx/v5"
)

type LicenseRepository interface {
	Save(ctx context.Context, license *lcp.License) error
	FindByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error)
}

type licenseRepository struct {
	db *pgx.Conn
}

func NewLicenseRepository(db *pgx.Conn) LicenseRepository {
	return &licenseRepository{db}
}

func (r *licenseRepository) Save(ctx context.Context, license *lcp.License) error {
	query := squirrel.Insert("licenses").
		Columns("id", "publication_id", "user_id", "passphrase", "hint", "publication_url", "right_print", "right_copy", "start_date", "end_date", "created_at").
		Values(
			license.ID,
			license.PublicationID,
			license.UserID,
			license.Passphrase,
			license.Hint,
			license.PublicationURL,
			license.RightPrint,
			license.RightCopy,
			license.StartDate,
			license.EndDate,
			license.CreatedAt,
		).
		Suffix("RETURNING id").
		RunWith(r.db)

	_, err := query.ExecContext(ctx)
	return err
}

func (r *licenseRepository) FindByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error) {
	query := squirrel.Select("id", "publication_id", "user_id", "passphrase", "hint", "publication_url", "right_print", "right_copy", "start_date", "end_date", "created_at").
		From("licenses")
	if publicationID != nil {
		query = query.Where(squirrel.Eq{"publication_id": *publicationID})
	}
	query = query.RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var licenses []*lcp.License
	for rows.Next() {
		var lic lcp.License
		err := rows.Scan(
			&lic.ID,
			&lic.PublicationID,
			&lic.UserID,
			&lic.Passphrase,
			&lic.Hint,
			&lic.PublicationURL,
			&lic.RightPrint,
			&lic.RightCopy,
			&lic.StartDate,
			&lic.EndDate,
			&lic.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, &lic)
	}
	return licenses, nil
}
