package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/buguzei/medods-task/internal/models"
	"github.com/buguzei/medods-task/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo(ctx context.Context, cfg *config.Config) *Mongo {
	opts := options.Client().ApplyURI(cfg.MongoConn)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("cannot connect to mongoDB:", err)
	}

	client.Database("app")

	log.Println("mongoDB connected successfully!")

	return &Mongo{Client: client}
}

type UsersDoc struct {
	UserID       string `bson:"userID"`
	RefreshToken string `bson:"refreshToken"`
}

func (m Mongo) IsRefreshTokenExists(ctx context.Context, refreshToken string) (bool, string, error) {
	collection := m.Client.Database("app").Collection("users")

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return false, "", fmt.Errorf("error of finding: %w", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var res bson.M

		if err = cur.Decode(&res); err != nil {
			return false, "", fmt.Errorf("error of decoding res: %w", err)
		}

		hashedToken := res["refreshToken"].(string)

		err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshToken))
		if err != nil {
			continue
		}

		return true, res["userID"].(string), nil
	}

	return false, "", nil
}

func (m Mongo) UpdateRefreshToken(ctx context.Context, user models.User, refreshToken string) error {
	collection := m.Client.Database("app").Collection("users")

	update := bson.D{
		{"$set", bson.D{
			{"refreshToken", refreshToken},
		}},
	}

	res, err := collection.UpdateOne(ctx, bson.M{"userID": user.GUID}, update)
	if err != nil {
		return fmt.Errorf("error of updating: %w", err)
	}

	if res.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (m Mongo) IsUserExist(ctx context.Context, user models.User) (bool, error) {
	collection := m.Client.Database("app").Collection("users")

	res := collection.FindOne(ctx, bson.M{"userID": user.GUID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, res.Err()
	}

	return true, nil
}

func (m Mongo) NewRefreshToken(ctx context.Context, user models.User, refreshToken string) error {
	collection := m.Client.Database("app").Collection("users")

	_, err := collection.InsertOne(ctx, bson.M{
		"userID":       user.GUID,
		"refreshToken": refreshToken,
	})
	if err != nil {
		return fmt.Errorf("error of inserting: %w", err)
	}

	var doc UsersDoc

	res := collection.FindOne(ctx, bson.M{
		"refreshToken": refreshToken,
	})

	err = res.Decode(&doc)
	if err != nil {
		return fmt.Errorf("error of decoding res: %w", err)
	}

	return nil
}
