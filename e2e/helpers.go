package e2e

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	db *mongo.Database
}

func NewMongo() (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	uri := fmt.Sprintf("%v://%v:%v@%v/%v", "mongodb", "transactor", "transactor", "mongodb:27017", "transactor")
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongodb: %w", err)
	}

	return &Mongo{
		db: mongoClient.Database("transactor"),
	}, nil
}

type registrationBounty struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Identity string             `bson:"identity" json:"identity"`
}

func (m *Mongo) InsertRegistrationBounty(identity common.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	toInsert := registrationBounty{
		ID:       primitive.NewObjectID(),
		Identity: identity.Hex(),
	}
	_, err := m.db.Collection("registration_bounties").InsertOne(ctx, toInsert)
	return err
}
