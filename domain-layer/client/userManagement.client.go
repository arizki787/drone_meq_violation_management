package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type UserManagementClient struct {
	BaseURL    string
	HttpClient *http.Client
}

func NewUserManagementClient(baseURL string) *UserManagementClient {
	return &UserManagementClient{
		BaseURL:    baseURL,
		HttpClient: NewHTTPClient(),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func (c *UserManagementClient) GetToken(username, password string) (string, string, error) {
	reqBody, err := json.Marshal(LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}

	resp, err := c.HttpClient.Post(c.BaseURL+"/api/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("failed to login")
	}

	var loginResponse LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return "", "", err
	}

	return loginResponse.Token, loginResponse.Username, nil
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// Given token by calling User Management
func (c *UserManagementClient) VerifyToken(token string) (string, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/user-management/user-by-token", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid token")
	}

	var response struct {
		Username string `json:"username"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Username, nil
}