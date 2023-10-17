package equipment

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Usecase struct {
	repo        equipment.Repository
	minioClient *minio.Client
	bucketName  string
	minioURL    string
}

func New(repo equipment.Repository, cfg *minioConfig) (*Usecase, error) {
	u, err := url.Parse(cfg.apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse minio api url: %w", err)
	}

	cl, err := minio.New(u.Host, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.accessKeyID, cfg.secretAccessKey, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("build minio client with endpoint %s: %w", u.Host, err)
	}

	return &Usecase{
		repo:        repo,
		minioClient: cl,
		bucketName:  cfg.bucketName,
		minioURL:    cfg.apiURL,
	}, nil
}

func (u *Usecase) AddNewEquipment(ctx context.Context, title, description string, body io.Reader,
	mimeType string, size int64, pictureName string) error {

	if before, _, ok := strings.Cut(mimeType, "/"); !ok || before != "image" {
		return fmt.Errorf("bad content tipe %s", mimeType)
	}
	ind := strings.LastIndex(pictureName, ".")
	if ind == -1 || ind+1 == len(pictureName) {
		return fmt.Errorf("bad expansion file %s", pictureName)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("new uuid: %w", err)
	}
	filename := id.String() + pictureName[ind:]

	_, err = u.minioClient.PutObject(ctx, u.bucketName, filename, body, size, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return fmt.Errorf("put file in minio bucket %s: %w", u.bucketName, err)
	}

	err = u.repo.AddEquipment(&ds.Equipment{
		Title:       title,
		Description: description,
		Picture:     fmt.Sprintf("%s/%s/%s", u.minioURL, u.bucketName, filename),
		Status:      "active",
		Count:       1,
	})
	return err
}

func (u *Usecase) GetListEquipments() ([]ds.Equipment, error) {
	return u.repo.GetAllEquipments()
}

func (u *Usecase) SearchEquipmentsByTitle(title string) ([]ds.Equipment, error) {
	return u.repo.SearchEquipmentsByTitle(title)
}

func (u *Usecase) DeleteEquipmentByID(id int) error {
	return u.DeleteEquipmentByID(id)
}

func (u *Usecase) EditEquipmentByID(id int) error {
	return nil
}
