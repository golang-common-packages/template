package datastore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/golang-common-packages/template/model"
)

// MongoDBDatastore manage all mongodb action
type MongoDBDatastore struct {
	Client *mongo.Client
	Config *model.Database
}

/*
	@sessionMapping: Mapping between model.Database and MongoDBDatastore for singleton pattern
	@ctx: returns a non-nil, empty Context when it's unclear which Context to use or it is not yet available
*/
var (
	sessionMapping = make(map[string]*MongoDBDatastore)
	ctx            = context.Background()
)

// NewMongoDBDatastore function return a new mongo client based on singleton pattern
func NewMongoDBDatastore(config *model.Service) Datastore {
	hash := config.Hash()
	currentSession := sessionMapping[hash]
	if currentSession == nil {
		currentSession = &MongoDBDatastore{nil, nil}

		// Setup client options
		clientOptions := options.Client().ApplyURI(getConnectionURI(config.Database.MongoDB))

		// Establish MongoDB connection
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Println("Error when try to connect to Mongodb server: ", err)
			panic(err)
		}

		// Check the connection status
		if err := client.Ping(ctx, nil); err != nil {
			log.Println("Can not ping to Mongodb server: ", err)
			panic(err)
		} else {
			currentSession.Client = client
			currentSession.Config = &config.Database
			sessionMapping[hash] = currentSession
			log.Println("Connected to Mongodb Server")
		}
	}

	return currentSession
}

// getConnectionURL return mongo connection URI
func getConnectionURI(config model.MongoDB) (URI string) {
	host := strings.Join(config.Hosts, ",")
	opt := strings.Join(config.Options, "?")
	if config.User == "" && config.Password == "" {
		return fmt.Sprintf("%v?%v", host, opt)
	}
	URI = fmt.Sprintf("mongodb+srv://%v:%v@%v/%v", config.User, config.Password, host, opt)

	return URI
}

// createSession return a new mongo session & transaction
func (m *MongoDBDatastore) createSession() (session mongo.Session) {
	session, err := m.Client.StartSession()
	if err != nil {
		log.Println("Error when try to start session: ", err)
		panic(err)
	}

	if err := session.StartTransaction(); err != nil {
		log.Println("Error when try to start transaction: ", err)
		panic(err)
	}

	return session
}

// GetUsers function return all users
func (m *MongoDBDatastore) GetUsers(lastID, pageSize string) (users []model.User, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		// Setup filter
		filter := bson.M{}
		if lastID != "" {
			id, err := primitive.ObjectIDFromHex(lastID)
			if err != nil {
				fmt.Printf("%d can not convert to ObjectID", id)
			}

			filter = bson.M{
				"_id": bson.M{"$gt": id},
			}
		}

		// Convert pageSize from string to int64
		limit, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			fmt.Printf("%d can not convert to int64", limit)
		}

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

		cur, err := collection.Find(ctx, filter, findOptions)
		defer cur.Close(ctx)
		if err != nil {
			return err
		}

		// Decode cursor
		for cur.Next(ctx) {
			result := model.User{}
			if err := cur.Decode(&result); err != nil {
				return err
			}
			result.Password = nil // Remove password
			users = append(users, result)
		}
		if err = cur.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session: ", err)
		return users, err
	}

	return users, nil
}

// GetUser return user info bases on username but without sensitive information
func (m *MongoDBDatastore) GetUser(username string) (user model.User, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		filter := bson.M{
			"username": username,
		}
		if err = collection.FindOne(ctx, filter).Decode(&user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at GetUser function: ", err)
		return user, err
	}

	return user, nil
}

// SaveUser function store user to collection
func (m *MongoDBDatastore) SaveUser(user model.User) (err error) {
	currentTime := time.Now()
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		// Set unique to collectin index for prevent duplicates value
		_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bsonx.Doc{
					{"email", bsonx.String("")},
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bsonx.Doc{
					{"username", bsonx.String("")},
				},
				Options: options.Index().SetUnique(true),
			},
		})
		if err != nil {
			return err
		}

		user.ID = primitive.NewObjectID()
		user.Created = currentTime
		user.Updated = currentTime
		user.Expiration = currentTime.Add(time.Hour * 24 * 90)
		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at SaveUser function: ", err)
		return err
	}

	return nil
}

// UpdateUser function update user info base on user object
func (m *MongoDBDatastore) UpdateUser(user *model.User) (err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		filter := bson.M{
			"_id": user.ID,
		}
		update := bson.M{"$set": bson.M{
			"updated": time.Now(),
		}}

		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at UpdateUser function: ", err)
		return err
	}

	return nil
}

// DeleteUser function delete user based on userID
func (m *MongoDBDatastore) DeleteUser(userID interface{}) (err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	userID, ok := userID.(primitive.ObjectID)
	if !ok {
		return errors.New("Can't convert userID type interface to primitive.ObjectID at DeleteUser function")
	}

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		filter := bson.M{
			"_id": userID,
		}
		_, err = collection.DeleteOne(ctx, filter)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at DeleteUser function: ", err)
		return err
	}

	return nil
}

// ActiveUser function active user based on username
func (m *MongoDBDatastore) ActiveUser(username string) (err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.User)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		filter := bson.M{
			"username": username,
		}
		update := bson.M{"$set": bson.M{
			"isactive": true,
			"updated":  time.Now(),
		}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at ActiveUser function: ", err)
		return err
	}

	return nil
}

// GetDocuments function for document pagination (Find and SetLimit)
func (m *MongoDBDatastore) GetDocuments(lastID, pageSize string) (documents []model.Document, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.Document)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		// Setup filter
		filter := bson.M{}
		if lastID != "" {
			id, err := primitive.ObjectIDFromHex(lastID)
			if err != nil {
				fmt.Printf("%d can not convert to ObjectID", id)
			}

			filter = bson.M{
				"_id": bson.M{"$gt": id},
			}
		}

		// Convert pageSize from string to int64
		limit, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			fmt.Printf("%d can not convert to int64", limit)
		}

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

		cur, err := collection.Find(ctx, filter, findOptions)
		defer cur.Close(ctx)
		if err != nil {
			return err
		}

		// Decode cursor
		for cur.Next(ctx) {
			result := model.Document{}
			if err := cur.Decode(&result); err != nil {
				return err
			}
			documents = append(documents, result)
		}
		if err = cur.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at PagingDocumentByFind function: ", err)
		return documents, err
	}

	return documents, nil
}

// SaveDocuments fucntion store document to collection
func (m *MongoDBDatastore) SaveDocuments(document model.Document) (err error) {
	currentTime := time.Now()
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(m.Config.MongoDB.DB).Collection(m.Config.Collection.Document)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		document.Created = currentTime
		document.Updated = currentTime

		_, err = collection.InsertOne(ctx, document)
		if err != nil {
			log.Println(err)
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at SaveDocuments function: ", err)
		return err
	}

	return nil
}
