package repository

import (
	"context"
	"errors"
	"testcontainers-mongo-demo/config"
	"testcontainers-mongo-demo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoUserRepository(cfg *config.DatabaseConfig) (*MongoUserRepository, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, errors.New("error connecting to mongo: " + err.Error())
	}

	return &MongoUserRepository{
		client:     client,
		database:   cfg.Database,
		collection: cfg.Collection,
	}, nil
}

func (r *MongoUserRepository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

func (r *MongoUserRepository) Close(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}

func (r *MongoUserRepository) Create(ctx context.Context, user *models.User) error {
	coll := r.client.Database(r.database).Collection(r.collection)

	// Generate ID if not set
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Ensure the ID is set from the result if needed
	if result.InsertedID != nil {
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			user.ID = oid
		}
	}

	return nil
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	var user models.User
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *models.User) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

func (r *MongoUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
