package forms

import "mime/multipart"

//FileID 文件ID表单
type FileID struct {
	ID string `uri:"id" binding:"required,uuid"`
}

//FileMove 文件移动表单
type FileMove struct {
	FileID
	DirectoryID string `json:"directory_id" form:"directory_id" binding:"uuid"`
}

//FileRename 文件修改名称表单
type FileRename struct {
	FileID
	NewName string `json:"new_name" form:"new_name" binding:"required"`
}

//FileMkdir 文件夹创建表单
type FileMkdir struct {
	DirectoryID string `json:"directory_id" form:"directory_id" binding:"omitempty,uuid"`
	Name        string `json:"name" form:"name" binding:"required"`
}

//FileQuery 文件查询表单
type FileQuery struct {
	FileID
	Limit  int `json:"limit" form:"limit" binding:"omitempty"`
	Offset int `json:"offset" form:"offset" binding:"omitempty"`
}

//FileUpload 文件上传表单
type FileUpload struct {
	DirectoryID string                `form:"directory_id" binding:"omitempty,uuid"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

//FileDownload 文件下载
type FileDownload struct {
	ID string `uri:"id" binding:"required,uuid"`
}
