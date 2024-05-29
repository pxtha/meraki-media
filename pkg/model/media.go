package model

import (
	"github.com/google/uuid"
	"gitlab.com/merakilab9/meracore/ginext"
	"net/http"
)

type Media struct {
	BaseModel
	UploadUrl   string `json:"upload_url"`
	FileType    string `json:"file_type"`
	FileName    string `json:"file_name"`
	Description string `json:"description"`
	URL         string `json:"url" gorm:"-"`
}

func (r *Media) TableName() string {
	return "media"
}

type PreUploadDataRequest struct {
	CreatorID uuid.UUID `json:"creator_id" valid:"Required"`
	MediaType string    `json:"media_type" valid:"Required"`
	Name      string    `json:"name" valid:"Required"`
}

type PreUploadMediaDataRequest struct {
	CreatorID   uuid.UUID `json:"creator_id"`
	MediaType   string    `json:"media_type" valid:"Required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TypeToCrop  string    `json:"type_to_crop" `
	Root        string    `json:"root"`
}

type PreUploadDataResponse struct {
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
}

type UploadDataRequest struct {
	UploadURL     string `json:"upload_url" valid:"Required"`
	ContentType   string `json:"content_type" valid:"Required"`
	ContentLength int64  `json:"content_length"`
	FileName      string `json:"file_name"`
}

type UploadInfo struct {
	BaseModel
	UploadURL      string `json:"upload_url" valid:"Required"`
	ContentType    string `json:"content_type" valid:"Required"`
	ContentLength  int64  `json:"content_length"`
	FileNameOrigin string `json:"file_name_origin"`
	FileNameS3     string `json:"file_name_s3"`
	Platform       string `json:"platform"`
}

type UploadDataResponse struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type MediaRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	UpdaterID   *uuid.UUID `json:"updater_id,omitempty"`
	CreatorID   *uuid.UUID `json:"creator_id,omitempty"`
	UploadUrl   *string    `json:"upload_url,omitempty"`
	FileType    *string    `json:"file_type,omitempty"`
	FileName    *string    `json:"file_name,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type PostUploadResponse struct {
	BaseModel
	Url         ImageUrl `json:"url,omitempty"`
	FileType    string   `json:"file_type,omitempty"`
	FileName    string   `json:"file_name,omitempty"`
	Description string   `json:"description,omitempty"`
}

type ImageUrl struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Origin    string `json:"origin,omitempty"`
	Small     string `json:"small,omitempty"`
	File      string `json:"file,omitempty"`
}

type MediaListRequest struct {
	CreatorID *string `form:"creator_id"`
}

type MediaFilter struct {
	MediaListRequest
	Pager *ginext.Pager
}

type MediaFilterResult struct {
	Filter  *MediaFilter
	Records []*Media
}

type UrlRes struct {
	PushURL string `json:"push_url"`
	Name    string `json:"name"`
}
