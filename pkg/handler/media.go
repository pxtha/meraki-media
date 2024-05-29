package handler

import (
	"github.com/google/uuid"
	"gitlab.com/merakilab9/meracore/ginext"
	"gitlab.com/merakilab9/meracore/logger"

	"gitlab.com/merakilab9/meradia/pkg/model"
	"gitlab.com/merakilab9/meradia/pkg/service"
	"gitlab.com/merakilab9/meradia/pkg/utils"
	"net/http"
)

type MediaHandlers struct {
	service service.MediaInterface
}

func NewMediaHandlers(service service.MediaInterface) *MediaHandlers {
	return &MediaHandlers{service: service}
}

// ---------------------------Media---------------------------------------------
//===================================================================================

// PreUpload
// @Tags MediaHandlers
// @Summary PreUpload
// @Description PreUpload
// @Accept  json
// @Produce  json
// @Param x-user-id header string false "user id"
// @Param data body model.PreUploadMediaDataRequest true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/media/pre-upload [post]
func (h *MediaHandlers) PreUpload(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "PreUpload")

	req := model.PreUploadMediaDataRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).WithField("request", req).Error("Error when check valid request")
		return nil, ginext.NewError(http.StatusNotFound, err.Error())
	}

	rs, err := h.service.PreUpload(r.GinCtx, uuid.New(), req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, rs), nil
}

func (h *MediaHandlers) Upload(c *ginext.Request) (*ginext.Response, error) {
	file, _, err := c.GinCtx.Request.FormFile("file")
	if err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	rs, err := h.service.Upload(c.GinCtx, file, c.GinCtx.Request.Form.Get("upload_url"), 0)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, rs), nil
}

// PosUpload
// @Tags MediaHandlers
// @Summary PosUpload
// @Description PosUpload
// @Accept  json
// @Produce  json
// @Param x-user-id header string false "user id"
// @Param data body model.PreUploadMediaDataRequest true "body data"
// @Success 200 {object} interface{}
// @Router /api/v1/media/pos-upload [post]
func (h *MediaHandlers) PosUpload(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, "PosUpload")
	req := model.PreUploadMediaDataRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, utils.MessageError()[http.StatusBadRequest])
	}
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).WithField("request", req).Error("Error when check valid request")
		return nil, ginext.NewError(http.StatusNotFound, utils.MessageError()[http.StatusNotFound])
	}
	rs, err := h.service.PosUpload(r.GinCtx, uuid.New(), req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusOK, rs), nil
}
