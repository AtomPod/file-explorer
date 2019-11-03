package services

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/phantom-atom/file-explorer/internal/locker"

	"github.com/phantom-atom/file-explorer/internal/log"
	"github.com/phantom-atom/file-explorer/models"
	"github.com/phantom-atom/file-explorer/repository"

	"github.com/phantom-atom/file-explorer/config"
)

var (
	//ErrParentNotADirectory 父路径不是一个文件夹
	ErrParentNotADirectory = errors.New("父路径不是一个文件夹")
	//ErrFileNotFound 文件未找到
	ErrFileNotFound = errors.New("文件未找到")
	//ErrFileAlreadyExists 文件已经存在
	ErrFileAlreadyExists = errors.New("文件已经存在")
	//ErrDirectoryNotFound 目录未找到
	ErrDirectoryNotFound = errors.New("文件夹不存在")
	//ErrFileShareInvalid 文件分享码无效
	ErrFileShareInvalid = errors.New("文件分享码无效")
	//ErrFileIsMissing 文件丢失
	ErrFileIsMissing = errors.New("文件丢失")
	//ErrCannotDownloadDirectory 不能下载文件夹
	ErrCannotDownloadDirectory = errors.New("文件夹不支持下载")
)

//PathError os.PathError
type PathError struct {
	Op   string
	Path string
	Err  error
}

//Error error.Error实现
func (p *PathError) Error() string {
	return "[" + p.Op + "] (" + p.Path + ") " + p.Err.Error()
}

//NewPathError 新建PathError
func NewPathError(op, path string, err error) *PathError {
	return &PathError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

//File multipart.File
type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

//FileService 文件服务
type FileService struct {
	config      func() *config.Config
	uuid        func() string
	dataContext repository.DataContext
	namedLocker locker.NamedLocker
}

//NewFileService 创建FileService
func NewFileService(
	configFunc func() *config.Config,
	uuid func() string,
	dataContext repository.DataContext,
	namedLocker locker.NamedLocker,
) *FileService {
	return &FileService{
		config:      configFunc,
		uuid:        uuid,
		dataContext: dataContext,
		namedLocker: namedLocker,
	}
}

func (f *FileService) createFileModel(file *models.File) (*models.File, error) {

	if file == nil {
		return nil, errors.New("FileService: invalid argument 'file' in createFileModel")
	}

	if file.FID == "" {
		fid := f.uuid()
		if fid == "" {
			return nil, errors.New("FileService: cannot create uuid in createFileModel")
		}
		file.FID = fid
	}

	repository, err := f.dataContext.File()
	if err != nil {
		return nil, err
	}

	directory := "/"

	if file.Owner != file.PFID {
		directoryMod, err := repository.GetFileByID(file.Owner, file.FID)
		if err == nil {
			return nil, err
		}

		if directoryMod == nil {
			return nil, NewPathError("create", file.PFID, ErrParentNotADirectory)
		}

		if !directoryMod.IsDir {
			return nil, NewPathError("create", file.PFID, ErrParentNotADirectory)
		}

		directory = path.Join(directoryMod.Directory, directoryMod.Filename)
	}

	sameMod, err := repository.GetFileByPFIDAndName(file.Owner, file.PFID, file.Filename)
	if err != nil {
		return nil, err
	}
	if sameMod != nil {
		absPath := path.Join(sameMod.Directory, sameMod.Filename)
		return nil, NewPathError("create", absPath, ErrFileAlreadyExists)
	}

	file.Directory = directory

	if err := repository.CreateFile(file); err != nil {
		return nil, err
	}

	return file, nil
}

//CreateDirectory 创建一个文件夹，文件属于owner，存在于directoryID文件夹中
func (f *FileService) CreateDirectory(owner string, directoryID string, name string) (*models.File, error) {
	if owner == "" {
		return nil, invalidArgument("FileService", "owner", "CreateDirectory")
	}

	if name == "" {
		return nil, invalidArgument("FileService", "name", "CreateDirectory")
	}

	if directoryID == "" {
		directoryID = owner
	}

	file := &models.File{
		Owner:    owner,
		PFID:     directoryID,
		Filename: name,
		IsDir:    true,
	}

	f.namedLocker.Lock(owner)
	defer f.namedLocker.UnLock(owner)
	return f.createFileModel(file)
}

func (f *FileService) saveFile(file File, name string) (absolutePath string, err error) {
	basePath := f.config().FileService.FileAbsolutePath()
	absolutePath = filepath.Join(basePath, name)

	if err := os.MkdirAll(basePath, 0666); err != nil {
		return "", err
	}

	diskFile, err := os.Create(absolutePath)
	if err != nil {
		return "", err
	}
	defer diskFile.Close()

	if _, err := io.Copy(diskFile, file); err != nil {
		if e := os.Remove(absolutePath); e != nil {
			log.Error("msg", "occur a error when delete file", "error", e.Error())
		}
		return "", err
	}
	return
}

//CreateFile 创建一个文件，文件属于owner，存在于directoryID文件夹中
func (f *FileService) CreateFile(owner string,
	directoryID string,
	name string,
	size int64,
	file File) (*models.File, error) {

	if owner == "" {
		return nil, invalidArgument("FileService", "owner", "CreateFile")
	}

	if name == "" {
		return nil, invalidArgument("FileService", "name", "CreateFile")
	}

	if file == nil {
		return nil, invalidArgument("FileService", "file", "CreateFile")
	}

	if directoryID == "" {
		directoryID = owner
	}

	fid := f.uuid()
	if fid == "" {
		return nil, errors.New("FileService: cannot create uuid in CreateFile")
	}

	fileMod := &models.File{
		Owner:    owner,
		PFID:     directoryID,
		Filename: name,
		IsDir:    false,
		FID:      fid,
		Size:     size,
	}

	absolutePath, err := f.saveFile(file, fid)
	if err != nil {
		return nil, err
	}

	f.namedLocker.Lock(owner)
	defer f.namedLocker.UnLock(owner)
	newFileMod, err := f.createFileModel(fileMod)
	if err != nil {
		if e := os.Remove(absolutePath); e != nil {
			log.Error("msg", "occur a error when delete file", "error", e.Error())
		}
		return nil, err
	}
	return newFileMod, nil
}

//DeleteFile 删除文件，文件属于owner
func (f *FileService) DeleteFile(owner string, fid string) error {
	if owner == "" {
		return invalidArgument("FileService", "owner", "DeleteFile")
	}
	if fid == "" {
		return invalidArgument("FileService", "fid", "DeleteFile")
	}

	f.namedLocker.Lock(owner)
	defer f.namedLocker.UnLock(owner)

	var commited = false
	UOW, err := f.dataContext.Unit()
	if err != nil {
		return err
	}
	defer func() {
		if !commited {
			if err := UOW.Rollback(); err != nil {
				log.Warn("msg", "rollback failed in FileService.DeleteFile", "error", err.Error())
			}
		}
	}()

	fileRepository, err := UOW.File()
	if err != nil {
		return err
	}

	deleteFile, err := fileRepository.GetFileByID(owner, fid)
	if err != nil {
		return err
	}

	if deleteFile == nil {
		return NewPathError("delete", fid, ErrFileNotFound)
	}

	if err := fileRepository.DeleteFile(deleteFile); err != nil {
		return err
	}

	if err := UOW.Commit(); err != nil {
		return err
	}
	commited = true
	return nil
}

//Download 下载文件，文件属于owner，编号为fid
func (f *FileService) Download(owner string, fid string) (File, *models.File, error) {
	if owner == "" {
		return nil, nil, invalidArgument("FileService", "owner", "Download")
	}
	if fid == "" {
		return nil, nil, invalidArgument("FileService", "fid", "Download")
	}

	fileRepository, err := f.dataContext.File()
	if err != nil {
		return nil, nil, err
	}

	downloadFile, err := fileRepository.GetFileByID(owner, fid)

	if err != nil {
		return nil, nil, err
	}

	if downloadFile == nil {
		return nil, nil, NewPathError("download", fid, ErrFileNotFound)
	}

	if downloadFile.IsDir {
		return nil, nil, NewPathError("download", fid, ErrCannotDownloadDirectory)
	}

	absolutePath := filepath.Join(f.config().FileService.FileAbsolutePath(), downloadFile.FID)
	file, err := os.OpenFile(absolutePath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			if err := fileRepository.DeleteFile(downloadFile); err != nil {
				return nil, nil, err
			}
			return nil, nil, NewPathError("download", fid, ErrFileIsMissing)
		}
		return nil, nil, err
	}
	return file, downloadFile, nil
}

//MoveFile 移动文件位置，文件编号为fid，新文件夹newPFID
func (f *FileService) MoveFile(owner string, fid string, newPFID string) error {
	if owner == "" {
		return invalidArgument("FileService", "owner", "Download")
	}

	if fid == "" {
		return invalidArgument("FileService", "fid", "Download")
	}

	if newPFID == "" {
		return invalidArgument("FileService", "newPFID", "Download")
	}

	f.namedLocker.Lock(owner)
	defer f.namedLocker.UnLock(owner)

	var commited = false
	UOW, err := f.dataContext.Unit()
	if err != nil {
		return err
	}
	defer func() {
		if !commited {
			if err := UOW.Rollback(); err != nil {
				log.Warn("msg", "rollback failed in FileService.MoveFile", "error", err.Error())
			}
		}
	}()

	fileRepository, err := UOW.File()
	if err != nil {
		return err
	}

	moveFile, err := fileRepository.GetFileByID(owner, fid)
	if err != nil {
		return err
	}

	if moveFile == nil {
		return NewPathError("move", fid, ErrFileNotFound)
	}

	if moveFile.PFID == newPFID {
		return nil
	}

	if err := f.move(moveFile, newPFID, fileRepository); err != nil {
		return err
	}

	if err := UOW.Commit(); err != nil {
		return err
	}
	commited = true
	return nil
}

//RenameFile 修改文件名称，文件编号为fid，新文件名称newName
func (f *FileService) RenameFile(owner string, fid string, newName string) error {
	if owner == "" {
		return invalidArgument("FileService", "owner", "Rename")
	}

	if fid == "" {
		return invalidArgument("FileService", "fid", "Rename")
	}

	if newName == "" {
		return invalidArgument("FileService", "newName", "Rename")
	}

	f.namedLocker.Lock(owner)
	defer f.namedLocker.UnLock(owner)

	var commited = false
	UOW, err := f.dataContext.Unit()
	if err != nil {
		return err
	}
	defer func() {
		if !commited {
			if err := UOW.Rollback(); err != nil {
				log.Warn("msg", "rollback failed in FileService.MoveFile", "error", err.Error())
			}
		}
	}()

	fileRepository, err := UOW.File()
	if err != nil {
		return err
	}

	renameFile, err := fileRepository.GetFileByID(owner, fid)
	if err != nil {
		return err
	}

	if renameFile == nil {
		return NewPathError("rename", fid, ErrFileNotFound)
	}

	if renameFile.Filename == newName {
		return nil
	}

	if err := f.rename(renameFile, newName, fileRepository); err != nil {
		return err
	}

	if err := UOW.Commit(); err != nil {
		return err
	}
	commited = true
	return nil
}

func (f *FileService) rename(file *models.File, newName string,
	repos repository.FileRepository) error {
	matchedFile, err := repos.GetFileByPFIDAndName(file.Owner, file.PFID, newName)
	if err != nil {
		return err
	}

	if matchedFile != nil {
		return NewPathError("rename", newName, ErrFileAlreadyExists)
	}

	file.Filename = newName
	if file.IsDir {
		subfiles, err := repos.GetFilesByPFID(file.Owner, file.FID, 0, 0)
		if err != nil {
			return err
		}

		newDirectory := path.Join(file.Directory, newName)
		for _, subfile := range subfiles {
			if err := f.adjustFilePath(subfile, newDirectory, repos); err != nil {
				return err
			}
		}
	}
	return repos.UpdateFile(file)
}

func (f *FileService) move(file *models.File, newFPID string,
	repos repository.FileRepository) error {
	matchedFile, err := repos.GetFileByPFIDAndName(file.Owner, newFPID, file.Filename)
	if err != nil {
		return err
	}

	if matchedFile != nil {
		return NewPathError("move", file.Filename, ErrFileAlreadyExists)
	}

	var parentFile *models.File
	if newFPID != file.Owner {
		var err error
		parentFile, err = repos.GetFileByID(file.Owner, newFPID)
		if err != nil {
			return err
		}
		if parentFile == nil {
			return NewPathError("move", newFPID, ErrParentNotADirectory)
		}
	} else {
		parentFile = &models.File{
			Directory: "/",
			Filename:  "",
		}
	}

	file.PFID = newFPID
	directory := path.Join(parentFile.Directory, parentFile.Filename)
	return f.adjustFilePath(file, directory, repos)
}

//GetFileByID 根据FID获取文件信息
func (f *FileService) GetFileByID(owner string, fid string) (*models.File, error) {

	if owner == "" {
		return nil, invalidArgument("FileService", "owner", "GetFileByID")
	}

	if fid == "" {
		return nil, invalidArgument("FileService", "fid", "GetFileByID")
	}

	fileRepository, err := f.dataContext.File()
	if err != nil {
		return nil, err
	}

	file, err := fileRepository.GetFileByID(owner, fid)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, NewPathError("get", fid, ErrFileNotFound)
	}

	return file, nil
}

//GetFileByPID 根据目录文件ID获取文件信息
func (f *FileService) GetFileByPID(owner string, pid string, limit int, offset int) ([]*models.File, error) {

	if owner == "" {
		return nil, invalidArgument("FileService", "owner", "GetFileByID")
	}

	if pid == "" {
		return nil, invalidArgument("FileService", "pid", "GetFileByID")
	}

	fileRepository, err := f.dataContext.File()
	if err != nil {
		return nil, err
	}

	files, err := fileRepository.GetFilesByPFID(owner, pid, limit, offset)
	if err != nil {
		return nil, err
	}

	if files == nil {
		return nil, NewPathError("get", pid, ErrFileNotFound)
	}

	return files, nil
}

//GetFileListsByOnwer 根据拥有者获取所有文件信息
func (f *FileService) GetFileListsByOnwer(owner string) ([]*models.File, error) {

	if owner == "" {
		return nil, invalidArgument("FileService", "owner", "GetFileByID")
	}

	fileRepository, err := f.dataContext.File()
	if err != nil {
		return nil, err
	}

	files, err := fileRepository.GetFileByOwner(owner, false, 0, 0)
	if err != nil {
		return nil, err
	}

	if files == nil {
		return nil, NewPathError("get", owner, ErrFileNotFound)
	}

	return files, nil
}

func (f *FileService) adjustFilePath(file *models.File, dir string, frepository repository.FileRepository) error {
	if file.IsDir {
		subfiles, err := frepository.GetFilesByPFID(file.Owner, file.FID, 0, 0)
		if err != nil {
			return err
		}

		newDirectory := path.Join(dir, file.Filename)
		for _, subfile := range subfiles {
			if err := f.adjustFilePath(subfile, newDirectory, frepository); err != nil {
				return err
			}
		}
	}
	file.Directory = dir
	return frepository.UpdateFile(file)
}

//Close 关闭服务
func (f *FileService) Close() error {
	return nil
}
