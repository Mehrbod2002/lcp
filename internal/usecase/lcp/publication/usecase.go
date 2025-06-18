package publication

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/lcp-project/internal/domain/lcp"
	"github.com/yourusername/lcp-project/internal/lcp/encrypt"
)

type PublicationUsecase interface {
	UploadAndEncrypt(ctx context.Context, title string, file io.Reader) (*lcp.Publication, error)
	GetAll(ctx context.Context) ([]*lcp.Publication, error)
}

type publicationUsecase struct {
	repo lcp.PublicationRepository
	enc  encrypt.Encrypter
}

func NewPublicationUsecase(repo lcp.PublicationRepository, enc encrypt.Encrypter) PublicationUsecase {
	return &publicationUsecase{repo, enc}
}

func (u *publicationUsecase) UploadAndEncrypt(ctx context.Context, title string, file io.Reader) (*lcp.Publication, error) {
	// Save file temporarily
	tempPath := "/tmp/" + title + ".tmp"
	out, err := os.Create(tempPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}

	// Encrypt using lcpencrypt
	encryptedPath, err := u.enc.Encrypt(tempPath, title)
	if err != nil {
		return nil, err
	}

	// Store publication metadata
	pub := &lcp.Publication{
		ID:            uuid.New().String(),
		Title:         title,
		FilePath:      tempPath,
		EncryptedPath: encryptedPath,
		CreatedAt:     time.Now(),
	}
	err = u.repo.Save(ctx, pub)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (u *publicationUsecase) GetAll(ctx context.Context) ([]*lcp.Publication, error) {
	return u.repo.FindAll(ctx)
}
