package routes

import ( 
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"
  "time"
"github.com/gorilla/mux"
  "github.com/golang-jwt/jwt/v5"
  "golang.org/x/crypto/bcrypt"
)


type user struct {
  AccountID string
  Email string
  Username string
  Password string
  Banned bool
  }


type TokenResponse struct {
  AccessToken  string  `json:"access_token"`
  
