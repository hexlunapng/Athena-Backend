package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"Athena-Backend/database/models"

	"go.uber.org/zap"
)

const (
	profileTag = "PROFILE"
)

func CreateProfile(accountId, username string, logger *zap.Logger) (*models.Profiles, error) {
	logger = logger.With(zap.String("tag", profileTag))

	dir := os.Getenv("PROFILE_DIR")
	if dir == "" {
		logger.Error("PROFILE_DIR environment variable not set")
		return nil, ErrProfileDirNotSet
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		logger.Error("failed to resolve absolute path for PROFILE_DIR", zap.String("dir", dir), zap.Error(err))
		return nil, err
	}

	profiles := make(map[string]interface{})

	files, err := os.ReadDir(absDir)
	if err != nil {
		logger.Error("read profiles dir failed", zap.Error(err))
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(absDir, file.Name())
		logger.Info("loading profile file", zap.String("file", filePath))

		data, err := loadProfileModule(filePath)
		if err != nil {
			logger.Warn("failed to load profile file, skipping", zap.String("file", filePath), zap.Error(err))
			continue
		}

		now := time.Now().UTC().Format(time.RFC3339)
		data["accountId"] = accountId
		data["created"] = now
		data["updated"] = now

		profileId, ok := data["profileId"].(string)
		if !ok || profileId == "" {
			logger.Warn("profile missing profileId, skipping", zap.String("file", filePath))
			continue
		}

		profiles[profileId] = data
		logger.Info("loaded profile", zap.String("profileId", profileId))
	}

	p := models.UserProfiles(accountId, profiles)
	now := time.Now().UTC()
	p.Updated = &now

	err = p.Save()
	if err != nil {
		logger.Error("failed to save profiles", zap.Error(err))
		return nil, err
	}

	logger.Info("Profiles fully loaded.")
	logger.Info("Someone with the accountid: " + accountId + " and username: " + username + " has created an account.")

	return p, nil
}

var ErrProfileDirNotSet = fmt.Errorf("PROFILE_DIR env var not set")

func loadProfileModule(path string) (map[string]interface{}, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	return data, err
}

func ValidateProfile(profileId string, profiles map[string]interface{}) bool {
	if profileId == "" || profiles == nil {
		return false
	}
	_, ok := profiles[profileId]
	return ok
}
