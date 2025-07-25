package repository

import (
	"context"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BundleRepository interface {
	CreateBundle(ctx context.Context, bundle *model.Bundle) (*model.Bundle, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*model.Bundle, error)
	GetByHashAndVersionId(ctx context.Context, hash string, versionId primitive.ObjectID) (*model.Bundle, error)
	UpdateVersionIdById(ctx context.Context, id primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error)
	GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Bundle, error)
	GetNextSeqByEnvironmentIdAndVersionId(ctx context.Context, environmentId primitive.ObjectID, versionId primitive.ObjectID) (int64, error)
	GetBySequenceIdEnvironmentIdAndVersionId(ctx context.Context, sequenceId int64, environmentId primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error)
	GetByLabelAndEnvironmentId(ctx context.Context, label string, environmentId primitive.ObjectID) (*model.Bundle, error)
	GetAllByVersionId(ctx context.Context, versionId primitive.ObjectID) ([]*model.Bundle, error)
	UpdateIsMandatoryById(ctx context.Context, id primitive.ObjectID, isMandatory bool) (*model.Bundle, error)
	UpdateIsValid(ctx context.Context, id primitive.ObjectID, isValid bool) (*model.Bundle, error)
	AddActive(ctx context.Context, id primitive.ObjectID) error
	AddFailed(ctx context.Context, id primitive.ObjectID) error
	AddInstalled(ctx context.Context, id primitive.ObjectID) error
	DecrementActive(ctx context.Context, id primitive.ObjectID) error
}

type bundleRepository struct {
	Connection *mongo.Database
}

func NewBundleRepository(db *mongo.Database) BundleRepository {
	return &bundleRepository{Connection: db}
}

func (bundleRepository *bundleRepository) CreateBundle(ctx context.Context, bundle *model.Bundle) (*model.Bundle, error) {
	bundle.CreatedAt = time.Now()
	bundle.UpdatedAt = time.Now()
	collection := bundleRepository.Connection.Collection("bundles")
	insrtedBundle, err := collection.InsertOne(ctx, &bundle, options.InsertOne())
	if err != nil {
		return nil, err
	}
	bundle.Id = insrtedBundle.InsertedID.(primitive.ObjectID)
	return bundle, nil
}

func (bundleRepository *bundleRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundle model.Bundle
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (bundleRepository *bundleRepository) GetByHashAndVersionId(ctx context.Context, hash string, versionId primitive.ObjectID) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundle model.Bundle
	err := collection.FindOne(ctx, bson.M{"hash": hash, "versionId": versionId}).Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (bundleRepository *bundleRepository) GetNextSeqByEnvironmentIdAndVersionId(ctx context.Context, environmentId primitive.ObjectID, versionId primitive.ObjectID) (int64, error) {
	var result struct {
		SequenceId int64 `bson:"sequenceId"`
	}
	collection := bundleRepository.Connection.Collection("bundles")
	filter := bson.M{"environmentId": environmentId, "versionId": versionId}
	err := collection.FindOne(ctx, filter, options.FindOne().SetSort(bson.M{"sequenceId": -1})).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, err
	}
	return result.SequenceId + 1, nil
}

func (bundleRepository *bundleRepository) UpdateVersionIdById(ctx context.Context, id primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"versionId": versionId}})
	if err != nil {
		return nil, err
	}
	return &model.Bundle{Id: id, VersionId: versionId}, nil
}

func (bundleRepository *bundleRepository) GetByEnvironmentAndVersion(ctx context.Context, environment string, version string) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundle model.Bundle
	err := collection.FindOne(ctx, bson.M{"environment": environment, "version": version}).Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (bundleRepository *bundleRepository) GetByLabelAndEnvironmentId(ctx context.Context, label string, environmentId primitive.ObjectID) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundle model.Bundle
	err := collection.FindOne(ctx, bson.M{"label": label, "environmentId": environmentId}).Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (bundleRepository *bundleRepository) GetBySequenceIdEnvironmentIdAndVersionId(ctx context.Context, sequenceId int64, environmentId primitive.ObjectID, versionId primitive.ObjectID) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundle model.Bundle
	err := collection.FindOne(ctx, bson.M{"sequenceId": sequenceId, "environmentId": environmentId, "versionId": versionId}).Decode(&bundle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bundle, nil
}

func (bundleRepository *bundleRepository) AddActive(ctx context.Context, id primitive.ObjectID) error {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"active": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (bundleRepository *bundleRepository) AddFailed(ctx context.Context, id primitive.ObjectID) error {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"failed": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (bundleRepository *bundleRepository) AddInstalled(ctx context.Context, id primitive.ObjectID) error {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"installed": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (bundleRepository *bundleRepository) DecrementActive(ctx context.Context, id primitive.ObjectID) error {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"active": -1}})
	if err != nil {
		return err
	}
	return nil
}

func (bundleRepository *bundleRepository) GetAllByVersionId(ctx context.Context, versionId primitive.ObjectID) ([]*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	var bundles []*model.Bundle
	cursor, err := collection.Find(ctx, bson.M{"versionId": versionId})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &bundles)
	if err != nil {
		return nil, err
	}
	return bundles, nil
}

func (bundleRepository *bundleRepository) UpdateIsMandatoryById(ctx context.Context, id primitive.ObjectID, isMandatory bool) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isMandatory": isMandatory}})
	if err != nil {
		return nil, err
	}
	return &model.Bundle{Id: id, IsMandatory: isMandatory}, nil
}

func (bundleRepository *bundleRepository) UpdateIsValid(ctx context.Context, id primitive.ObjectID, isValid bool) (*model.Bundle, error) {
	collection := bundleRepository.Connection.Collection("bundles")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isValid": isValid}})
	if err != nil {
		return nil, err
	}
	return &model.Bundle{Id: id, IsValid: isValid}, nil
}
