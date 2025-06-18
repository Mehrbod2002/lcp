package lcp

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
	"github.com/jackc/pgx/v5"
)

type PublicationRepository interface {
	Save(ctx context.Context, pub *lcp.Publication) error
	FindAll(ctx context.Context) ([]*lcp.Publication, error)
}

type publicationRepository struct {
	db *pgx.Conn
}

func NewPublicationRepository(db *pgx.Conn) PublicationRepository {
	return &publicationRepository{db}
}

func (r *publicationRepository) Save(ctx context.Context, pub *lcp.Publication) error {
	query := squirrel.Insert("publications").
		Columns("id", "title", "file_path", "encrypted_path", "created_at").
		Values(pub.ID, pub.Title, pub.FilePath, pub.EncryptedPath, pub.CreatedAt).
		Suffix("RETURNING id").
		RunWith(r.db)

	_, err := query.ExecContext(ctx)
	return err
}

func (r *publicationRepository) FindAll(ctx context.Context) ([]*lcp.Publication, error) {
	query := squirrel.Select("id", "title", "file_path", "encrypted_path", "created_at").
		From("publications").
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publications []*lcp.Publication
	for rows.Next() {
		var pub lcp.Publication
		err := rows.Scan(&pub.ID, &pub.Title, &pub.FilePath, &pub.EncryptedPath, &pub.CreatedAt)
		if err != nil {
			return nil, err
		}
		publications = append(publications, &pub)
	}
	return publications, nil
}
