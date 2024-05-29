package pg

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/merakilab9/meracore/logger"
	"gitlab.com/merakilab9/meradia/pkg/model"
	"gorm.io/gorm"
	"strconv"
)

func (r *RepoPG) CreateMedia(ctx context.Context, ob *model.Media, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	return tx.Debug().Create(ob).Error
}

func (r *RepoPG) FilterMedia(ctx context.Context, f *model.MediaFilter) (*model.MediaFilterResult, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.WithContext(ctx).Model(&model.Media{})
	log := logger.WithCtx(ctx, "Media.Filter")

	if f.CreatorID != nil {
		tx = tx.Where("creator_id = ?", f.CreatorID)
	}

	//f.Pager.SortableFields = []string{"id","created_at","priority"}
	//f.Pager.Sort = "id,-created_at"

	result := &model.MediaFilterResult{
		Filter:  f,
		Records: []*model.Media{},
	}

	pager := result.Filter.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		log.WithError(tx.Error).Error("error while filter Media")
	}

	return result, tx.Error
}

func (r *RepoPG) GetOneMedia(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*model.Media, error) {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	o := &model.Media{}

	err := tx.First(&o, id).Error
	return o, err
}

func (r *RepoPG) GetRandomMedia(ctx context.Context, number int, des string, tx *gorm.DB) ([]*model.Media, error) {
	var o []*model.Media
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	query := `
		SELECT *
		FROM "media"
		WHERE description = ?
		ORDER BY random() LIMIT ` + strconv.Itoa(number)

	err := tx.Raw(query, des).Scan(&o).Error

	return o, err
}

func (r *RepoPG) UpdateMedia(ctx context.Context, update *model.Media, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	return tx.WithContext(ctx).Where("id = ?", update.ID).Save(&update).Error
}

func (r *RepoPG) DeleteMedia(ctx context.Context, d *model.Media, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	return tx.WithContext(ctx).Delete(&d).Error
}
