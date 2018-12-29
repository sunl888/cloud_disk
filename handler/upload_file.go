package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/wq1019/cloud_disk/model"
	"github.com/wq1019/cloud_disk/service"
	"github.com/zm-dev/go-file-uploader"
	"io"
	"io/ioutil"
	"mime"
	"os"
)

type uploadFile struct {
	u go_file_uploader.Uploader
}

func copy2TmpFile(file io.Reader) (tmpFileName string, err error) {
	tmpFile, err := ioutil.TempFile("", "cloud-")
	_, err = io.Copy(tmpFile, file)
	if cerr := tmpFile.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return
	}
	return tmpFile.Name(), nil
}

func (uf *uploadFile) UploadFile(c *gin.Context) {
	l := &struct {
		FolderId int64 `json:"folder_id" form:"folder_id"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BadRequest("请指定上传的文件夹", err))
		return
	}
	uploadFile, fh, err := c.Request.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest("请上传文件", err))
		return
	}
	defer uploadFile.Close()

	tmpFileName, err := copy2TmpFile(uploadFile)
	if err != nil {
		defer os.Remove(tmpFileName)
	}
	uFile, err := uf.u.Upload(go_file_uploader.FileHeader{Filename: fh.Filename, Size: fh.Size / 1024, File: uploadFile}, "")
	if err != nil {
		_ = c.Error(errors.InternalServerError("上传失败", err))
		return
	}
	fileModel := convert2FileModel(uFile)
	fileModel.Extra = mime.TypeByExtension(uFile.Format)
	fileModel.FolderId = l.FolderId
	err = service.UpdateFile(c, fileModel.Id, fileModel)
	if err != nil {
		_ = c.Error(err)
		return
	}
	//type Cover struct {
	//	Hash string `json:"hash"`
	//	Url  string `json:"url"`
	//}
	//
	//var coverHashUrl []*Cover
	//if len(data.CoverHash) > 0 {
	//	coverHashUrl = make([]*Cover, 0, len(data.CoverHash))
	//	for _, h := range data.CoverHash {
	//		coverHashUrl = append(coverHashUrl, &Cover{
	//			Hash: h,
	//			Url:  uf.imageUrl.Generate(h),
	//		})
	//	}
	//} else {
	//	coverHashUrl = make([]*Cover, 0)
	//}
	//c.JSON(200, gin.H{"hash": fileModel.Hash, "cover_hash": coverHashUrl})
}
func convert2FileModel(upload *go_file_uploader.FileModel) *model.File {
	return &model.File{
		Hash:      upload.Hash,
		Format:    upload.Format,
		Filename:  upload.Filename,
		Size:      upload.Size,
		CreatedAt: upload.CreatedAt,
		UpdatedAt: upload.UpdatedAt,
	}
}

/*func (uf *uploadFile) ShowFile(c *gin.Context) {
	hash := strings.TrimSpace(c.Param("hash"))
	s := uf.u.Store()
	fileModel, err := s.FileLoad(hash)
	if s.FileIsNotExistError(err) {
		c.Error(errors.NotFound("file not found"))
		return
	}
	c.JSON(200, fileModel)
}*/

func NewUploadFile(u go_file_uploader.Uploader) *uploadFile {
	return &uploadFile{u: u}
}
