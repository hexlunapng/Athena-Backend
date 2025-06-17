package managers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	db "Athena-Backend/database"
	"Athena-Backend/database/models"
	"context"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAccountRoutes(r *mux.Router) {
	r.HandleFunc("/account/api/oauth/token", oauthTokenHandler).Methods("POST")
	r.HandleFunc("/account/api/oauth/verify", oauthVerifyHandler).Methods("GET")
	r.HandleFunc("/account/api/oauth/sessions/kill", killSessionHandler).Methods("DELETE")
	r.HandleFunc("/account/api/oauth/sessions/kill/{token}", killSessionHandler).Methods("DELETE")

	r.HandleFunc("/account/api/public/account/{accountId}", accountByIDHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/displayName/{displayName}", accountByDisplayNameHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/", accountBatchHandler).Methods("GET")

	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth", deviceAuthListHandler).Methods("GET")
	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth", deviceAuthCreateHandler).Methods("POST")
	r.HandleFunc("/account/api/public/account/{accountId}/deviceAuth/{deviceId}", deviceAuthDeleteHandler).Methods("DELETE")
}

func findUserByUsername(username string) (*models.User, error) {
	collection := db.GetMongoCollection("users")
	filter := bson.M{"username": username}

	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func findUserByAccountID(accountID string) (*models.User, error) {
	collection := db.GetMongoCollection("users")
	filter := bson.M{"accountId": accountID}

	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func oauthTokenHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		GrantType    string `json:"grant_type"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Code         string `json:"code"`
		AccountID    string `json:"account_id"`
		ExchangeCode string `json:"exchange_code"`
	}

	var req Req

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}
	} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		req.GrantType = r.FormValue("grant_type")
		req.Username = r.FormValue("username")
		req.Password = r.FormValue("password")
		req.Code = r.FormValue("code")
		req.AccountID = r.FormValue("account_id")
		req.ExchangeCode = r.FormValue("exchange_code")
	} else {
		http.Error(w, "Unsupported Content-Type", http.StatusBadRequest)
		return
	}

	var user *models.User
	var displayName, accountId string
	var err error

	switch req.GrantType {
	case "client_credentials", "refresh_token":
		http.Error(w, "Grant type not supported yet", http.StatusBadRequest)
		return
	case "password":
		if req.Username == "" {
			http.Error(w, "username is required", http.StatusBadRequest)
			return
		}
		if req.Password == "" {
			http.Error(w, "password is required", http.StatusBadRequest)
			return
		}

		user, err = findUserByUsername(req.Username)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if user.Password != req.Password {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		displayName = user.Username
		accountId = user.AccountID
	case "authorization_code":
		displayName = req.Code
		accountId = req.Code
	case "device_auth":
		if req.AccountID == "" {
			http.Error(w, "account_id is required", http.StatusBadRequest)
			return
		}
		user, err = findUserByAccountID(req.AccountID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		displayName = user.Username
		accountId = user.AccountID
	case "exchange_code":
		displayName = req.ExchangeCode
		accountId = req.ExchangeCode
	default:
		http.Error(w, "Unsupported grant_type", http.StatusBadRequest)
		return
	}

	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	accessToken := hex.EncodeToString(randomBytes)

	response := map[string]interface{}{
		"access_token":       accessToken,
		"expires_in":         28800,
		"expires_at":         "9999-12-31T23:59:59.999Z",
		"token_type":         "bearer",
		"account_id":         accountId,
		"client_id":          "ec684b8c687f479fadea3cb2ad83f5c6",
		"internal_client":    true,
		"client_service":     "fortnite",
		"refresh_token":      "STATIC_REFRESH_TOKEN",
		"refresh_expires":    115200,
		"refresh_expires_at": "9999-12-31T23:59:59.999Z",
		"displayName":        displayName,
		"app":                "fortnite",
		"in_app_id":          accountId,
		"device_id":          "static-device-id",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func oauthVerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "bearer ")

	user, err := findUserByAccountID(token)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"access_token":       token,
		"expires_in":         28800,
		"expires_at":         "9999-12-31T23:59:59.999Z",
		"token_type":         "bearer",
		"refresh_token":      "STATIC_REFRESH_TOKEN",
		"refresh_expires":    115200,
		"refresh_expires_at": "9999-12-31T23:59:59.999Z",
		"account_id":         user.AccountID,
		"client_id":          "3446cd72694c4a4485d81b77adbb2141",
		"internal_client":    true,
		"client_service":     "fortnite",
		"displayName":        user.Username,
		"app":                "fortnite",
		"in_app_id":          user.AccountID,
		"device_id":          "static-device-id",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func killSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func accountByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["accountId"]

	user, err := findUserByAccountID(id)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            user.AccountID,
		"displayName":   user.Username,
		"externalAuths": map[string]interface{}{},
	})
}

func accountByDisplayNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["displayName"]

	user, err := findUserByUsername(name)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            user.AccountID,
		"displayName":   user.Username,
		"externalAuths": map[string]interface{}{},
	})
}

func accountBatchHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["accountId"]
	if !ok || len(ids) == 0 {
		http.Error(w, "Missing accountId", http.StatusBadRequest)
		return
	}

	var response []map[string]interface{}
	for _, id := range ids {
		user, err := findUserByAccountID(id)
		if err != nil || user == nil {
			continue
		}
		displayName := user.Username
		response = append(response, map[string]interface{}{
			"id":            user.AccountID,
			"displayName":   displayName,
			"externalAuths": map[string]interface{}{},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deviceAuthListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func deviceAuthCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"accountId": vars["accountId"],
		"deviceId":  "null",
		"secret":    "null",
	})
}

func deviceAuthDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
