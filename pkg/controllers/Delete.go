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

type ItemID struct {
	DeleteOne string
}

type MapItemID struct {
	DeleteArr []string
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func DeleteOneRequestHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		bodySting := string(info)
		if bodySting != "" {

			var newId ItemID
			json.Unmarshal(info, &newId)

			theCollection := config.GetCollection()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			primitiveID, err := primitive.ObjectIDFromHex(newId.DeleteOne)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				res, err := theCollection.DeleteOne(ctx, bson.M{"_id": primitiveID})
				if err != nil {
					var newErrorMessage ErrorMessage

					newErrorMessage.Message = fmt.Sprintf("DeleteOne() ERROR: %v", err)
					json.NewEncoder(w).Encode(newErrorMessage)
				} else {
					json.NewEncoder(w).Encode(res)
				}

			}
		}
	}
}

func DeleteListRequestHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	info, err := io.ReadAll(r.Body)
	if err != nil {
		var newErrorMessage ErrorMessage

		newErrorMessage.Message = fmt.Sprintf("Error reading request body: %v", err)
		json.NewEncoder(w).Encode(newErrorMessage)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var newIdMap MapItemID
		json.Unmarshal(info, &newIdMap)
		idMap := newIdMap.DeleteArr

		theCollection := config.GetCollection()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userMassage := []interface{}{}

		for i := range idMap {
			idPermitive, err := primitive.ObjectIDFromHex(idMap[i])
			if err != nil {
				userMassage = append(userMassage, err)
			} else {
				res, err := theCollection.DeleteOne(ctx, bson.M{"_id": idPermitive})
				if err != nil {
					userMassage = append(userMassage, err)
				} else {
					userMassage = append(userMassage, res)
				}
			}
		}

		json.NewEncoder(w).Encode(userMassage)

	}
}

func DeleteAllRequestHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	theCollection := config.GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deteleResult, err := theCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(deteleResult)
	}
	// 	// ** Note that: Returned ID can be In String or primitive.ObjectID type.
	// 	var currentStringId string
	// 	var currentObjectId primitive.ObjectID
	// 	var resultArr []*mongo.DeleteResult
	// 	var errorArr []error
	// 	type FinalResult struct {
	// 		Result []*mongo.DeleteResult
	// 		Error  []error
	// 	}
	// 	for _, v := range content {
	// 		fmt.Println(v["_id"])
	// 		// ** If id is of type string
	// 		if reflect.TypeOf(v["_id"]) == reflect.TypeOf(currentStringId) {
	// 			currentStringId = v["_id"].(string)
	// 			idPrimitive, err := primitive.ObjectIDFromHex(currentStringId)
	// 			fmt.Println(idPrimitive)
	// 			if err != nil {
	// 				fmt.Println("Error with primitive string", err)
	// 			} else {
	// 				res, err := theCollection.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	// 				if err != nil {
	// 					fmt.Println("error deleting one", err)
	// 					errorArr = append(errorArr, err)
	// 				} else {
	// 					resultArr = append(resultArr, res)
	// 				}
	// 			}
	// 			// ** If id is of type string
	// 		} else if reflect.TypeOf(v["_id"]) == reflect.TypeOf(currentObjectId) {
	// 			currentObjectId := v["_id"]
	// 			stringObjectId := currentObjectId.(primitive.ObjectID).Hex()
	// 			idPrimitive, err := primitive.ObjectIDFromHex(stringObjectId)
	// 			fmt.Println(idPrimitive)
	// 			if err != nil {
	// 				fmt.Println("Error deleting one", err)
	// 			} else {
	// 				res, err := theCollection.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	// 				if err != nil {
	// 					fmt.Println("error deleting one", err)
	// 					errorArr = append(errorArr, err)
	// 				} else {
	// 					resultArr = append(resultArr, res)
	// 				}
	// 			}
	// 		}
	// 	}
	// 	var newFinalResult FinalResult
	// 	newFinalResult.Result = resultArr
	// 	newFinalResult.Error = errorArr
	// 	json.NewEncoder(w).Encode(newFinalResult)
}
