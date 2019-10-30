package repository

import (
	"github.com/phantom-atom/file-explorer/models"
)

//FileRepository 文件仓库接口
type FileRepository interface {
	CreateFile(*models.File) error
	DeleteFile(*models.File) error
	UpdateFile(*models.File) error
	GetFileList(limit int, offset int) ([]*models.File, error)
	GetFileByPFIDAndName(owner string, pfid string, name string) (*models.File, error)
	GetFilesByPFID(owner string, pfid string, limit int, offset int) ([]*models.File, error)
	GetFileByID(owner string, fid string) (*models.File, error)
	GetFileByOwner(owner string, isdir bool, limit int, offset int) ([]*models.File, error)
}
