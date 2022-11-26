package config

import (
	"encoding/json"
	"net/http"

	"doe-base/idopost-backend/pkg/utils"
)

type StatusResponse struct {
	ConnectionStatus bool `json:"connection_status"`
	UserInfo         SuccessResponse
}

// ** Get user connection status
// ** Get user name and collectoin name in "SuccessResponse"
func ConnectionStatus(w http.ResponseWriter, r *http.Request) {
	//** EnableCors & SetResponse Header
	utils.EnableCors(w, r)

	var connectionStatus StatusResponse
	if Client != nil {
		accountName, _, collectionName := GetClientDetails()
		connectionStatus.ConnectionStatus = true
		connectionStatus.UserInfo.Name = accountName
		connectionStatus.UserInfo.Collection = collectionName

		json.NewEncoder(w).Encode(connectionStatus)
	} else {
		connectionStatus.ConnectionStatus = false
		connectionStatus.UserInfo.Name = ""
		connectionStatus.UserInfo.Collection = ""

		json.NewEncoder(w).Encode(connectionStatus)
	}
}
