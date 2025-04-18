package external

import (
	"context"
	"encoding/json"
	"errors"
	"BE-DroneMEQ-Violation-Management/domain-layer/client"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"net/http"
)

type UserManagementService interface {
	LoginAndGetToken(ctx context.Context, username string, password string) (string, string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	GetMetadataByToken(ctx context.Context, token string) (external.MetadataDTO, error)
}

type userManagementService struct {
	client *client.UserManagementClient
}

func NewUserManagementService(client *client.UserManagementClient) UserManagementService {
	return &userManagementService{client: client}
}

// LoginAndGetToken logs in and retrieves the token and username
func (s *userManagementService) LoginAndGetToken(ctx context.Context, username, password string) (string, string, error) {
	return s.client.GetToken(username, password)
}

func (s *userManagementService) VerifyToken(ctx context.Context, token string) (string, error) {
	return s.client.VerifyToken(token)
}

// GetMetadataByToken retrieves user metadata using the token
func (s *userManagementService) GetMetadataByToken(ctx context.Context, token string) (external.MetadataDTO, error) {
	req, err := http.NewRequest("GET", s.client.BaseURL+"/api/user-management/user-by-token", nil)
	if err != nil {
		return external.MetadataDTO{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		return external.MetadataDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return external.MetadataDTO{}, errors.New("failed to retrieve user metadata")
	}

	var metadataResponse struct {
		Data external.MetadataDTO `json:"data"` // Make sure this matches the API response structure.
	}

	err = json.NewDecoder(resp.Body).Decode(&metadataResponse)
	if err != nil {
		return external.MetadataDTO{}, err
	}

	return metadataResponse.Data, nil
}