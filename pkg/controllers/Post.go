package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"doe-base/idopost-backend/pkg/config"
	"doe-base/idopost-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostRequestHandle(w http.ResponseWriter, r *http.Request) {
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
			formatedString := utils.StringJsonFormatter(info)
			theCollection := config.GetCollection()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			InsertToMongoDB(formatedString, theCollection, ctx, w)
		}
	}

}

func InsertToMongoDB(formatedString string, theCollection *mongo.Collection, ctx context.Context, w http.ResponseWriter) {
	chars := formatedString
	jsonMap := make(map[string]interface{})

	// ** Handle One Post
	if chars[0] == '{' && chars[len(chars)-1] == '}' {
		// utils.CreateObjectID(chars)
		err2 := json.Unmarshal([]byte(formatedString), &jsonMap)
		if err2 != nil {
			json.NewEncoder(w).Encode(err2)
		} else {
			// ** convert _id of type string to objID
			for k := range jsonMap {
				if k == "_id" {
					str := jsonMap[k].(string)
					objID, err := primitive.ObjectIDFromHex(str)
					if err != nil {
						json.NewEncoder(w).Encode(err)
					} else {
						jsonMap[k] = objID
					}
				}
			}
			result, err := theCollection.InsertOne(ctx, jsonMap)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(result)
			}
		}

	}

	// ** Handle List Post
	if chars[0] == '[' && chars[len(chars)-1] == ']' {
		finishedSlice := utils.CreateJsonMaps(chars)
		userMassage := []interface{}{}

		for i := range finishedSlice {
			err2 := json.Unmarshal([]byte(finishedSlice[i]), &jsonMap)
			if err2 != nil {
				log.Fatal("Error decoding formated string", err2)
			} else {
				for k := range jsonMap {
					if k == "_id" {
						str := jsonMap[k].(string)
						objID, err := primitive.ObjectIDFromHex(str)
						if err != nil {
							json.NewEncoder(w).Encode(err)
						} else {
							jsonMap[k] = objID
						}
					}
				}
				result, err := theCollection.InsertOne(ctx, jsonMap)
				if err != nil {
					fmt.Println("Error Inserting to DataBase", err.Error())
					userMassage = append(userMassage, err)
				} else {
					fmt.Println("added successfully to collection")
					userMassage = append(userMassage, result)
				}
			}
		}
		json.NewEncoder(w).Encode(userMassage)
	}
}
