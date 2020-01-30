package database

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// MongoClient manage all mongodb action
type MongoClient struct {
	Client *mongo.Client
	Config *MongoDB
}

// NewMongoDB function return a new mongo client based on singleton pattern
func NewMongoDB(config *MongoDB) IDatabase {
	currentSession := &MongoClient{nil, nil}

	// Setup client options
	clientOptions := options.Client().ApplyURI(getConnectionURI(config))

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
	}

	currentSession.Client = client
	currentSession.Config = config
	log.Println("Connected to MongoDB Server")

	return currentSession
}

// getConnectionURL return mongo connection URI
func getConnectionURI(config *MongoDB) (URI string) {
	host := strings.Join(config.Hosts, ",")
	opt := strings.Join(config.Options, "?")
	if config.User == "" && config.Password == "" {
		return fmt.Sprintf("%v?%v", host, opt)
	}
	URI = fmt.Sprintf("mongodb+srv://%v:%v@%v/%v", config.User, config.Password, host, opt)

	return URI
}

// createSession return a new mongo session & transaction
func (m *MongoClient) createSession() (session mongo.Session) {
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

// GetALL ...
func (m *MongoClient) GetALL(databaseName, collectionName, lastID, pageSize string, dataModel interface{}) (results []interface{}, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	if collectionName == "" && lastID == "" && pageSize == "" {
		return nil, errors.New("collectionName, lastID and pageSize must not empty")
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		id, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			fmt.Printf("%d can not convert to ObjectID", id)
		}

		filter := bson.M{
			"_id": bson.M{"$gt": id},
		}

		// Convert pageSize from string to int64
		limit, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			fmt.Printf("%d can not convert to int64", limit)
		}

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})

		collection := m.Client.Database(databaseName).Collection(collectionName)
		cur, err := collection.Find(ctx, filter, findOptions)
		defer cur.Close(ctx)
		if err != nil {
			return err
		}

		// Decode cursor
		for cur.Next(ctx) {
			if err := cur.Decode(&dataModel); err != nil {
				return err
			}
			results = append(results, dataModel)
		}
		if err = cur.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at GetALL method: ", err)
		return nil, err
	}

	return nil, nil
}

// GetByField ...
func (m *MongoClient) GetByField(databaseName, collectionName, field, value string, dataModel reflect.Type) (result interface{}, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		filter := bson.M{
			field: value,
		}

		collection := m.Client.Database(databaseName).Collection(collectionName)
		SR := collection.FindOne(ctx, filter)
		if SR.Err() != nil {
			return SR.Err()
		}

		result = reflect.New(dataModel).Interface()
		err = SR.Decode(result)
		if err == nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at GetUser method: ", err)
		return nil, err
	}

	return result, nil
}

// Create ...
func (m *MongoClient) Create(databaseName, collectionName string, dataModel interface{}) (result interface{}, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.InsertOne(ctx, dataModel)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at SaveUser function: ", err)
		return nil, err
	}

	return result, nil
}

// Update ...
func (m *MongoClient) Update(databaseName, collectionName string, ID, dataModel interface{}) (result interface{}, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	collection := m.Client.Database(databaseName).Collection(collectionName)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		id, ok := ID.(primitive.ObjectID)
		if !ok {
			return errors.New("can't convert userID type interface to primitive.ObjectID at DeleteUser function")
		}
		filter := bson.M{
			"_id": id,
		}
		update := bson.M{"$set": dataModel}

		result, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at UpdateUser function: ", err)
		return nil, err
	}

	return result, nil
}

// Delete ...
func (m *MongoClient) Delete(databaseName, collectionName string, ID interface{}) (result interface{}, err error) {
	session := m.createSession()
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {
		id, ok := ID.(primitive.ObjectID)
		if !ok {
			return errors.New("can't convert userID type interface to primitive.ObjectID at DeleteUser function")
		}
		filter := bson.M{
			"_id": id,
		}

		collection := m.Client.Database(databaseName).Collection(collectionName)
		result, err = collection.DeleteOne(ctx, filter)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("Error when try to use with session at DeleteUser function: ", err)
		return nil, err
	}

	return result, nil
}
