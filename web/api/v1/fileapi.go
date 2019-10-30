package v1

import (
	"io"
	"net/http"

	"github.com/phantom-atom/file-explorer/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/services"
	"github.com/phantom-atom/file-explorer/web/forms"
)

func fileErrorToAPIResult(err error) *APIResult {
	pathErr, ok := err.(*services.PathError)
	if !ok {
		return Internal(err, nil)
	}

	switch pathErr.Err {
	case services.ErrParentNotADirectory, services.ErrFileIsMissing, services.ErrCannotDownloadDirectory:
		return FailedPrecondition(err, nil)
	case services.ErrFileNotFound, services.ErrDirectoryNotFound:
		return NotFound(err, nil)
	case services.ErrFileAlreadyExists:
		return AlreadyExists(err, nil)
	default:
		return Internal(err, nil)
	}
}

//FileUpload 上传文件API
//POST /api/v1/file/upload
func (api *API) FileUpload(c *gin.Context, form *forms.FileUpload) *APIResult {
	owner := c.GetString("userID")

	multipartFile, err := form.File.Open()
	if err != nil {
		return InvalidArgument(err, nil)
	}

	if form.DirectoryID == "" {
		form.DirectoryID = owner
	}

	createdFile, err := api.fileServ.CreateFile(owner,
		form.DirectoryID,
		form.File.Filename,
		form.File.Size,
		multipartFile)

	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(createdFile, nil)
}

//FileDownload 下载文件API
//GET /api/v1/file/{id}/download
func (api *API) FileDownload(c *gin.Context, form *forms.FileDownload) *APIResult {
	owner := c.GetString("userID")

	file, fileInfo, err := api.fileServ.Download(owner, form.ID)
	if err != nil {
		return fileErrorToAPIResult(err)
	}

	return OK(nil, func(c *gin.Context, data interface{}) error {
		defer func() {
			if err := file.Close(); err != nil {
				log.Error("msg", "occur a error when delete file", "error", err.Error())
			}
		}()

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileInfo.Filename)
		if _, err := io.Copy(c.Writer, file); err != nil {
			return err
		}
		c.Status(http.StatusOK)
		return nil
	})
}

//FileGetRootList 获取文件信息
//GET /api/v1/file/root
func (api *API) FileGetRootList(c *gin.Context) *APIResult {
	owner := c.GetString("userID")

	file, err := api.fileServ.GetFileByID(owner, owner)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(file, nil)
}

//FileGetInfo 获取文件信息
//GET /api/v1/file/{id}
func (api *API) FileGetInfo(c *gin.Context, form *forms.FileID) *APIResult {
	owner := c.GetString("userID")

	file, err := api.fileServ.GetFileByID(owner, form.ID)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(file, nil)
}

//FileGetList 获取文件夹文件列表
//GET /api/v1/file/{id}/list
func (api *API) FileGetList(c *gin.Context, form *forms.FileQuery) *APIResult {
	owner := c.GetString("userID")

	files, err := api.fileServ.GetFileByPID(owner, form.ID, form.Limit, form.Offset)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(files, nil)
}

//FileDelete 删除文件API
//DELETE /api/v1/file/{id}
func (api *API) FileDelete(c *gin.Context, form *forms.FileID) *APIResult {
	owner := c.GetString("userID")

	err := api.fileServ.DeleteFile(owner, form.ID)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(nil, nil)
}

//FileMkdir 创建文件夹API
//POST /api/v1/file/mkdir
func (api *API) FileMkdir(c *gin.Context, form *forms.FileMkdir) *APIResult {
	owner := c.GetString("userID")

	if form.DirectoryID == "" {
		form.DirectoryID = owner
	}

	createdDir, err := api.fileServ.CreateDirectory(owner,
		form.DirectoryID, form.Name)

	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(createdDir, nil)
}

//FileMove 移动文件API
//POST /api/v1/file/{id}/move
func (api *API) FileMove(c *gin.Context, form *forms.FileMove) *APIResult {
	owner := c.GetString("userID")

	if form.DirectoryID == "@" {
		form.DirectoryID = owner
	}

	err := api.fileServ.MoveFile(owner, form.ID, form.DirectoryID)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(nil, nil)
}

//FileRename 修改文件名称API
//POST /api/v1/file/{id}/rename
func (api *API) FileRename(c *gin.Context, form *forms.FileRename) *APIResult {
	owner := c.GetString("userID")

	err := api.fileServ.RenameFile(owner, form.ID, form.NewName)
	if err != nil {
		return fileErrorToAPIResult(err)
	}
	return OK(nil, nil)
}
