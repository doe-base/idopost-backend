package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"doe-base/idopost-backend/pkg/config"
	"doe-base/idopost-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetOne struct {
	GetOne string
}
type GetList struct {
	GetListArr []string
}

// ** Get All Items Request Handler Function **
func GetAllRequest(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	theCollection := config.GetCollection()
	cursor, err := theCollection.Find(ctx, bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		var content []bson.M
		if err = cursor.All(ctx, &content); err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(content)
		}
	}
}

// ** Get One Request Handler Function  **\\
func GetOneRequest(w http.ResponseWriter, r *http.Request) {
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
			var newGetOneId GetOne
			json.Unmarshal(info, &newGetOneId)

			theCollection := config.GetCollection()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// result := theCollection.FindOne(ctx, bson.M{"_id": objID})
			objID, err := primitive.ObjectIDFromHex(newGetOneId.GetOne)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				var content bson.M
				filter := bson.D{{Key: "_id", Value: objID}}
				result := theCollection.FindOne(ctx, filter)
				err := result.Decode(&content)
				if err != nil {
					json.NewEncoder(w).Encode(err)
				} else {
					json.NewEncoder(w).Encode(content)
				}
			}
		}
	}

}

// ** Get List Request Handler function
func GetListRequest(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var newGetListId GetList
		json.Unmarshal(info, &newGetListId)
		deleteList := newGetListId.GetListArr
		theCollection := config.GetCollection()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		userMassage := []interface{}{}

		for i := range deleteList {
			currentId := deleteList[i]
			objID, err := primitive.ObjectIDFromHex(currentId)
			if err != nil {
				userMassage = append(userMassage, err)
			} else {
				filter := bson.D{{Key: "_id", Value: objID}}
				var content bson.M
				err := theCollection.FindOne(ctx, filter).Decode(&content)
				if err != nil {
					userMassage = append(userMassage, err)
				} else {
					userMassage = append(userMassage, content)
				}

			}
		}
		json.NewEncoder(w).Encode(userMassage)
	}
}
