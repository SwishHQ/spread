package service

import (
	"context"
	"errors"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppService interface {
	CreateApp(ctx context.Context, appName string, os string) (*model.App, error)
	GetAppByName(ctx context.Context, appName string) (*model.App, error)
	GetApps(ctx context.Context) ([]*model.App, error)
	GetAppById(ctx context.Context, id string) (*model.App, error)
}

type appServiceImpl struct {
	appRepository repository.AppRepository
}

func NewAppService(appRepository repository.AppRepository) AppService {
	return &appServiceImpl{appRepository: appRepository}
}

func (appService *appServiceImpl) CreateApp(ctx context.Context, appName string, os string) (*model.App, error) {
	if os != "ios" && os != "android" {
		return nil, errors.New("invalid os")
	}
	existingApp, err := appService.appRepository.GetByName(ctx, appName)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingApp != nil {
		return nil, errors.New("app with name " + appName + " already exists")
	}
	app := model.App{
		Name:      appName,
		OS:        os,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdApp, err := appService.appRepository.Insert(ctx, &app)
	if err != nil {
		return nil, err
	}
	return createdApp, nil
}

func (appService *appServiceImpl) GetAppByName(ctx context.Context, appName string) (*model.App, error) {
	existingApp, err := appService.appRepository.GetByName(ctx, appName)
	if err != nil {
		return nil, err
	}
	return existingApp, nil
}

func (appService *appServiceImpl) GetApps(ctx context.Context) ([]*model.App, error) {
	apps, err := appService.appRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (appService *appServiceImpl) GetAppById(ctx context.Context, id string) (*model.App, error) {
	appId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("error converting id")
	}
	existingApp, err := appService.appRepository.GetById(ctx, appId)
	if err != nil {
		return nil, err
	}
	if existingApp == nil {
		return nil, errors.New("app not found")
	}
	return existingApp, nil
}
