package simple

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"

	"github.com/phantom-atom/file-explorer/internal/log"
	"github.com/phantom-atom/file-explorer/models"
)

func (r *dbRepository) CreateFile(f *models.File) error {
	return r.db.Create(f).Error
}

func (r *dbRepository) DeleteFile(f *models.File) error {
	if f.IsDir {
		files, err := r.GetFilesByPFID(f.Owner, f.FID, 0, 0)
		if err != nil {
			return err
		}

		for _, file := range files {
			r.DeleteFile(file)
		}

		err = r.db.Delete(f).Error
		if err != nil {
			return err
		}
	} else {
		err := r.db.Delete(f).Error
		if err != nil {
			return err
		}

		absolutePath := filepath.Join(r.config().FileService.FileAbsolutePath(), f.FID)
		if err := os.Remove(absolutePath); err != nil {
			log.Error("msg", "occur a error when delete file", "error", err.Error())
		}
	}
	return nil
}

func (r *dbRepository) UpdateFile(f *models.File) error {
	return r.db.Save(f).Error
}

func (r *dbRepository) GetFileList(limit int, offset int) ([]*models.File, error) {
	files := make([]*models.File, 0)

	db := r.db
	if limit != 0 {
		db = db.Limit(limit)
	}

	if offset != 0 {
		db = db.Offset(offset)
	}

	err := db.Find(&files).Error

	if err == gorm.ErrRecordNotFound {
		files = nil
		err = nil
	}
	return files, err
}

func (r *dbRepository) GetFileByPFIDAndName(owner string, pfid string, name string) (*models.File, error) {
	file := &models.File{}
	db := r.db
	db = db.Where("pfid = ? AND owner = ? AND filename = ?", pfid, owner, name)
	err := db.First(file).Error
	if err == gorm.ErrRecordNotFound {
		file = nil
		err = nil
	}
	return file, err
}

func (r *dbRepository) GetFilesByPFID(owner string, pfid string, limit int, offset int) ([]*models.File, error) {
	files := make([]*models.File, 0)
	db := r.db
	db = db.Where("owner = ? AND pfid = ?", owner, pfid)

	if limit != 0 {
		db = db.Limit(limit)
	}
	if offset != 0 {
		db = db.Offset(offset)
	}

	err := db.Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		files = nil
		err = nil
	}
	return files, err
}

func (r *dbRepository) GetFileByID(owner string, fid string) (*models.File, error) {
	file := &models.File{}
	db := r.db
	db = db.Where("fid = ? AND owner = ?", fid, owner)
	err := db.First(file).Error
	if err == gorm.ErrRecordNotFound {
		file = nil
		err = nil
	}
	return file, err
}

func (r *dbRepository) GetFileByOwner(owner string, isdir bool, limit int, offset int) ([]*models.File, error) {

	files := make([]*models.File, 0)
	db := r.db
	db = db.Where("owner = ?", owner)

	if isdir {
		db = db.Where("isdir = ?", isdir)
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	if offset != 0 {
		db = db.Offset(offset)
	}

	err := db.Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		files = nil
		err = nil
	}
	return files, err
}
