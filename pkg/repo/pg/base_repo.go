package pg

import (
	"context"
	"gitlab.com/merakilab9/meradia/pkg/model"
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source base_repo.go -destination mocks/base_repo.go

const (
	generalQueryTimeout = 60 * time.Second
	defaultPageSize     = 30
	maxPageSize         = 1000
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatorID uuid.UUID  `json:"creator_id"`
	UpdaterID uuid.UUID  `json:"updater_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type RepoPG struct {
	DB    *gorm.DB
	debug bool
}

func (r *RepoPG) GetRepo() *gorm.DB {
	return r.DB
}

func NewPGRepo(db *gorm.DB) PGInterface {
	return &RepoPG{DB: db}
}

type PGInterface interface {
	// DB
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)

	CreateMedia(ctx context.Context, ob *model.Media, tx *gorm.DB) error
	FilterMedia(ctx context.Context, f *model.MediaFilter) (*model.MediaFilterResult, error)
	GetOneMedia(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*model.Media, error)
	UpdateMedia(ctx context.Context, update *model.Media, tx *gorm.DB) error
	DeleteMedia(ctx context.Context, d *model.Media, tx *gorm.DB) error
	GetRandomMedia(ctx context.Context, number int, des string, tx *gorm.DB) ([]*model.Media, error)
}

func (r *RepoPG) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}
