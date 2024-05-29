package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/merakilab9/meracore/ginext"
	"gitlab.com/merakilab9/meracore/logger"
	"gitlab.com/merakilab9/meradia/conf"
	"gitlab.com/merakilab9/meradia/pkg/model"
	"gitlab.com/merakilab9/meradia/pkg/repo/pg"
	"gitlab.com/merakilab9/meradia/pkg/utils"
	"mime/multipart"
	"net/http"
	"strings"
)

type Sizer interface {
	Size() int64
}

type MediaService struct {
	repo  pg.PGInterface
	awsS3 S3ServiceInterface
}

func NewMediaService(s3 S3ServiceInterface, repo pg.PGInterface) MediaInterface {
	return &MediaService{awsS3: s3, repo: repo}
}

type MediaInterface interface {
	PreUpload(ctx context.Context, currentUser uuid.UUID, req model.PreUploadMediaDataRequest) (res interface{}, err error)
	Upload(ctx context.Context, file multipart.File, url string, filesize int64) (rs interface{}, err error)
	PosUpload(ctx context.Context, currentUser uuid.UUID, req model.PreUploadMediaDataRequest) (rs *model.PostUploadResponse, err error)
}

func (s *MediaService) Upload(ctx context.Context, file multipart.File, url string, filesize int64) (rs interface{}, err error) {
	log := logger.WithCtx(ctx, "MediaService.Upload")

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		log.WithError(err).Error("Error when GetMediaDirectory")
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}
	if _, err := file.Seek(0, 0); err != nil {
		log.WithError(err).Error("Error when GetMediaDirectory")
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}

	if filesize == 0 {
		filesize = file.(Sizer).Size()
	}
	req := &model.UploadDataRequest{
		UploadURL:     url,
		ContentType:   http.DetectContentType(fileHeader),
		ContentLength: filesize,
	}

	res, err := s.awsS3.Upload(file, req.UploadURL, req.ContentType, req.ContentLength)
	if err != nil {
		log.WithError(err).Error("Error when GetMediaDirectory")
		return nil, ginext.NewError(http.StatusInternalServerError, utils.MessageError()[http.StatusInternalServerError])
	}
	return res, nil
}

func (s *MediaService) PosUpload(ctx context.Context, currentUser uuid.UUID, req model.PreUploadMediaDataRequest) (rs *model.PostUploadResponse, err error) {
	log := logger.WithCtx(ctx, "MediaService.PosUpload")

	req.CreatorID = currentUser

	preUploadKey, err := utils.GetMediaDirectory(req.MediaType, req.CreatorID.String(), req.Name)
	if err != nil {
		log.WithError(err).Error("Error when GetMediaDirectory")
		return nil, ginext.NewError(http.StatusForbidden, utils.MessageError()[http.StatusForbidden])
	}

	uploadURL := preUploadKey

	uploadInfo := &model.Media{
		BaseModel: model.BaseModel{
			CreatorID: &req.CreatorID,
		},
		UploadUrl:   uploadURL,
		FileType:    req.MediaType,
		FileName:    req.Name,
		Description: req.Description,
	}

	if err := s.repo.CreateMedia(ctx, uploadInfo, nil); err != nil {
		log.WithError(err).Error("Error when CreateMedia")
		return nil, ginext.NewError(http.StatusForbidden, utils.MessageError()[http.StatusForbidden])
	}
	thumbnail := ""
	smallImg := ""
	file := ""
	origin := ""
	switch req.TypeToCrop {
	case "file":
		file = conf.LoadEnv().AWSMediaFullDomain + uploadInfo.UploadUrl
	case "avatar":
		cropConfig := "250x250,sc"
		origin = conf.LoadEnv().MediaBaseProxyURL + "800x800,sc,s" + "/" + uploadInfo.UploadUrl
		thumbnail = conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + "/" + uploadInfo.UploadUrl
		break
	case "moment":
		cropConfig := "250x250,sc"
		origin = conf.LoadEnv().MediaBaseProxyURL + "800x800,sc,s" + "/" + uploadInfo.UploadUrl
		thumbnail = conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + "/" + uploadInfo.UploadUrl
		break
	case "post":
		cropConfig := "250x250,sc"
		origin = conf.LoadEnv().MediaBaseProxyURL + "800x800,sc,s" + "/" + uploadInfo.UploadUrl
		thumbnail = conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + "/" + uploadInfo.UploadUrl
		break
	case "cover":
		cropConfig := "450x250,sc"
		origin = conf.LoadEnv().MediaBaseProxyURL + "800x800,sc,s" + "/" + uploadInfo.UploadUrl
		thumbnail = conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + "/" + uploadInfo.UploadUrl
		break
	case "chat":
		cropConfig := "300x300,sc"
		smallConfig := "100x100,sc"
		thumbnail = conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + "/" + uploadInfo.UploadUrl
		smallImg = conf.LoadEnv().MediaBaseProxyURL + smallConfig + ",s" + "/" + uploadInfo.UploadUrl
		origin = conf.LoadEnv().MediaBaseProxyURL + "800x800,sc,s" + "/" + uploadInfo.UploadUrl
		break
	default:
		return rs, ginext.NewError(http.StatusForbidden, fmt.Errorf("not valid type to crop").Error())
	}

	//macOrigin := hmac.New(sha256.New, []byte(conf.LoadEnv().MediaSecretKey))
	//macOrigin.Write([]byte(conf.LoadEnv().MediaBaseURL + uploadInfo.UploadUrl + "#0x0"))
	//resultSignOrigin := macOrigin.Sum(nil)
	//origin := conf.LoadEnv().MediaBaseProxyURL + "0x0,s" + base64.URLEncoding.EncodeToString(resultSignOrigin) + "/" + uploadInfo.UploadUrl

	//

	//macThumb := hmac.New(sha256.New, []byte(conf.LoadEnv().MediaSecretKey))
	//macThumb.Write([]byte(conf.LoadEnv().MediaBaseURL + uploadInfo.UploadUrl + "#" + cropConfig))
	//resultSign := macThumb.Sum(nil)
	//thumbnail := conf.LoadEnv().MediaBaseProxyURL + cropConfig + ",s" + base64.URLEncoding.EncodeToString(resultSign) + "/" + uploadInfo.UploadUrl

	return &model.PostUploadResponse{
		BaseModel: model.BaseModel{
			ID:        uploadInfo.ID,
			CreatedAt: uploadInfo.CreatedAt,
			UpdatedAt: uploadInfo.UpdatedAt,
		},
		Url: model.ImageUrl{
			Thumbnail: thumbnail,
			Origin:    origin,
			Small:     smallImg,
			File:      file,
		},
		FileType:    uploadInfo.FileType,
		FileName:    uploadInfo.FileName,
		Description: uploadInfo.Description,
	}, nil
}

func (s *MediaService) PreUpload(ctx context.Context, currentUser uuid.UUID, req model.PreUploadMediaDataRequest) (rs interface{}, err error) {
	log := logger.WithCtx(ctx, "MediaService.PreUpload")

	req.CreatorID = currentUser
	req.Name = uuid.New().String() + "." + strings.ToLower(req.MediaType)

	preUploadKey, err := utils.GetMediaDirectory(req.MediaType, req.CreatorID.String(), req.Name)
	if err != nil {
		log.WithError(err).WithField("name", req.Name).Error("Error when GetMediaDirectory")
		return nil, ginext.NewError(http.StatusForbidden, utils.MessageError()[http.StatusForbidden])
	}

	res, err := s.awsS3.PreUploadMedia(preUploadKey)
	if err != nil {
		log.WithError(err).WithField("preUploadKey", preUploadKey).Error("Error when PreUploadMedia")
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}

	return model.UrlRes{
		PushURL: res,
		Name:    req.Name,
	}, nil
}
