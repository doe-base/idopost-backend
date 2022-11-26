package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"doe-base/idopost-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAccount struct {
	AccountName    string
	PasswordName   string
	AtlasName      string
	DatabaseName   string
	CollectionName string
	ModelName      string
}

type MongoAccoutURL struct {
	AccountURL  string
	AccountName string
}
type MongoAccoutURLDetails struct {
	DatabaseName   string
	CollectionName string
}

type SuccessResponse struct {
	Name       string `json:"name"`
	Collection string `json:"collection"`
}
type ErrorMessage struct {
	Message string `json:"message"`
}

type ConnectionResponse struct {
	ConnectionStatus bool
}

type Response struct {
	UserDetails SuccessResponse `json:"collectioncontent"`
	DBContent   []bson.M        `json:"dbcontent"`
}

var Client *mongo.Client

func HandleMongodbConnectFormSubmit(w http.ResponseWriter, r *http.Request) {
	//** Allow CORS By * or specific origin
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var newMongoAccount MongoAccount
		json.Unmarshal(info, &newMongoAccount)

		// Convert string to lowercase
		AccountName := strings.ToLower(newMongoAccount.AccountName)
		AtlasName := strings.ToLower(newMongoAccount.AtlasName)
		DatabaseName := newMongoAccount.DatabaseName
		CollectionName := strings.ToLower(newMongoAccount.CollectionName)
		// Set Environment Variables for User Details
		os.Setenv("DB_USERNAME", AccountName)
		os.Setenv("DB_ATLAS", AtlasName)
		os.Setenv("DB_NAME", DatabaseName)
		os.Setenv("COLLECTION_NAME", CollectionName)
		var userDbUrl = "mongodb+srv://" + os.Getenv("DB_USERNAME") + ":" + newMongoAccount.PasswordName + "@" + os.Getenv("DB_ATLAS") + ".wvunv.mongodb.net/" + os.Getenv("DB_NAME") + "?retryWrites=true&w=majority"

		SuccessRes := SuccessResponse{
			Name:       os.Getenv("DB_USERNAME"),
			Collection: os.Getenv("COLLECTION_NAME"),
		}
		var newErrorMessage ErrorMessage

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()
		// ** Establish connection to client mongo database
		mongoDBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(userDbUrl))
		if err != nil {
			newErrorMessage.Message = err.Error()
			json.NewEncoder(w).Encode(newErrorMessage)
		} else {
			Client = mongoDBClient
			res, err := json.Marshal(SuccessRes)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				w.Write(res)
			}
		}
	}
}

func HandleMongodbURLConnectFormSubmit(w http.ResponseWriter, r *http.Request) {
	//** Allow CORS By * or specific origin
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		bodyString := string(info)
		if bodyString != "" {

			var newMongoAccount MongoAccoutURL
			json.Unmarshal(info, &newMongoAccount)

			userDbUrl := newMongoAccount.AccountURL
			userAccountName := newMongoAccount.AccountName
			os.Setenv("DB_USERNAME", userAccountName)
			ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancle()
			var newErrorMessage ErrorMessage

			// ** Establish connection to client mongo database
			mongoDBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(userDbUrl))
			if err != nil {
				newErrorMessage.Message = err.Error()
				json.NewEncoder(w).Encode(newErrorMessage)
			} else {
				Client = mongoDBClient
				newconnectionResponse := ConnectionResponse{
					ConnectionStatus: true,
				}
				res, err := json.Marshal(newconnectionResponse)
				if err != nil {
					json.NewEncoder(w).Encode(err)
				} else {
					w.Write(res)
				}
			}
		}
	}
}

func HandleMongodbURLConnectFormDetailsSubmit(w http.ResponseWriter, r *http.Request) {
	//** Allow CORS By * or specific origin
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		bodyString := string(info)
		if bodyString != "" {
			var newMongoAccoutURLDetails MongoAccoutURLDetails
			json.Unmarshal(info, &newMongoAccoutURLDetails)

			databaseName := newMongoAccoutURLDetails.DatabaseName
			collectionName := newMongoAccoutURLDetails.CollectionName

			os.Setenv("DB_NAME", databaseName)
			os.Setenv("COLLECTION_NAME", collectionName)

			accountName := os.Getenv("DB_USERNAME")

			ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancle()
			fmt.Println(databaseName)
			fmt.Println(collectionName)
			theCollection := Client.Database(databaseName).Collection(collectionName)

			cursor, err := theCollection.Find(ctx, bson.M{})
			if err != nil {
				newconnectionResponse := ConnectionResponse{
					ConnectionStatus: false,
				}
				json.NewEncoder(w).Encode(newconnectionResponse)
			} else {
				var content []bson.M
				if err = cursor.All(ctx, &content); err != nil {
					newconnectionResponse := ConnectionResponse{
						ConnectionStatus: false,
					}
					json.NewEncoder(w).Encode(newconnectionResponse)
				} else {
					var newResponse Response
					newResponse.UserDetails.Name = accountName
					newResponse.UserDetails.Collection = collectionName
					newResponse.DBContent = content
					json.NewEncoder(w).Encode(newResponse)
				}
			}

		}
	}

}

func HandleMongodbLogOut(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	err := Client.Disconnect(ctx)

	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

}

// * Returns accountName, databaseName, collection
func GetClientDetails() (string, string, string) {
	accountName := os.Getenv("DB_USERNAME")
	databaseName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	return accountName, databaseName, collectionName
}
