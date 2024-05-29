package route

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/pprof"

	"gitlab.com/merakilab9/meracore/ginext"
	"gitlab.com/merakilab9/meracore/logger"
	"gitlab.com/merakilab9/meracore/service"

	"gitlab.com/merakilab9/meradia/conf"
	"gitlab.com/merakilab9/meradia/pkg/repo/pg"

	handlerMeradia "gitlab.com/merakilab9/meradia/pkg/handler"
	serviceMeradia "gitlab.com/merakilab9/meradia/pkg/service"
)

type Service struct {
	*service.BaseApp
}

func NewService() *Service {

	s := &Service{
		service.NewApp("meradia Service", "v1.0"),
	}
	db := s.GetDB()

	if !conf.LoadEnv().DbDebugEnable {
		db = db.Debug()
	}
	endpoint := aws.Endpoint{
		PartitionID:   "aws",
		URL:           conf.LoadEnv().AWSMediaDomain,
		SigningRegion: conf.LoadEnv().AWSUpCloudRegion,
	}

	endpointResolver := aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			return endpoint, nil
		},
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.LoadEnv().AWSUpCloudAccessKeyID, conf.LoadEnv().AWSUpCloudSecretKey, "")),
		config.WithRegion(conf.LoadEnv().AWSUpCloudRegion),
		config.WithEndpointResolver(endpointResolver),
	)

	if err != nil {
		logger.WithCtx(context.Background(), "init server").Error(err.Error())
	}

	client := s3.NewFromConfig(cfg)

	repoPG := pg.NewPGRepo(db)
	storageService := serviceMeradia.NewS3Service(client, conf.LoadEnv().AWSBucket)
	mediaService := serviceMeradia.NewMediaService(storageService, repoPG)
	mediaHandle := handlerMeradia.NewMediaHandlers(mediaService)
	migrateHandle := handlerMeradia.NewMigrationHandler(db)

	pprof.Register(s.Router)

	v1Api := s.Router.Group("/api/v1")
	v1Api.POST("/media/pre-upload", ginext.WrapHandler(mediaHandle.PreUpload))
	v1Api.POST("/media/pos-upload", ginext.WrapHandler(mediaHandle.PosUpload))
	v1Api.POST("/media/upload", ginext.WrapHandler(mediaHandle.Upload))

	v1Api.POST("/internal/migrate", migrateHandle.Migrate)

	return s
}
