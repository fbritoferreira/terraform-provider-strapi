package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type StrapiClient struct {
	Endpoint   string
	APIToken   string
	HTTPClient *http.Client
}

func New(endpoint, apiToken string) *StrapiClient {
	return &StrapiClient{
		Endpoint:   endpoint,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

// GetContentTypes retrieves all content types from Strapi
func (c *StrapiClient) GetContentTypes() (interface{}, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// CreateContentType creates a new content type
func (c *StrapiClient) CreateContentType(data map[string]interface{}) (interface{}, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// UpdateContentType updates an existing content type
func (c *StrapiClient) UpdateContentType(uid string, data map[string]interface{}) (interface{}, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// DeleteContentType deletes a content type
func (c *StrapiClient) DeleteContentType(uid string) error {
	// TODO: Implement API call
	return fmt.Errorf("not implemented")
}

// AdminUser represents a Strapi admin user
type AdminUser struct {
	ID                int    `json:"id,omitempty"`
	Email             string `json:"email"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname,omitempty"`
	Password          string `json:"password,omitempty"`
	IsActive          *bool  `json:"isActive,omitempty"`
	PreferedLanguage  string `json:"preferedLanguage,omitempty"`
	Roles             []int  `json:"roles,omitempty"`
	RegistrationToken string `json:"registrationToken,omitempty"`
}

// GetAdminUsers retrieves all admin users from Strapi
func (c *StrapiClient) GetAdminUsers() ([]AdminUser, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/admin/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get admin users: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data struct {
			Results []AdminUser `json:"results"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Results, nil
}

// GetAdminUser retrieves an admin user by ID
func (c *StrapiClient) GetAdminUser(id int) (*AdminUser, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/admin/users/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get admin user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data AdminUser `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// CreateAdminUser creates a new admin user
func (c *StrapiClient) CreateAdminUser(user AdminUser) (*AdminUser, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/admin/users", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create admin user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data AdminUser `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// UpdateAdminUser updates an existing admin user
func (c *StrapiClient) UpdateAdminUser(id int, user AdminUser) (*AdminUser, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.Endpoint+"/admin/users/"+strconv.Itoa(id), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update admin user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data AdminUser `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// DeleteAdminUser deletes an admin user
func (c *StrapiClient) DeleteAdminUser(id int) error {
	req, err := http.NewRequest("DELETE", c.Endpoint+"/admin/users/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete admin user: %s - %s", resp.Status, string(body))
	}

	return nil
}

// User represents a Strapi content API user
type User struct {
	ID         int                    `json:"id"`
	DocumentID string                 `json:"documentId"`
	Username   string                 `json:"username"`
	Email      string                 `json:"email"`
	Confirmed  bool                   `json:"confirmed"`
	Blocked    bool                   `json:"blocked"`
	Role       map[string]interface{} `json:"role,omitempty"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
}

// GetUsers retrieves all users from Strapi
func (c *StrapiClient) GetUsers() ([]User, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get users: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data []User `json:"data"`
		Meta struct {
			Pagination struct {
				Total int `json:"total"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetUser retrieves a user by ID
func (c *StrapiClient) GetUser(id int) (*User, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/users/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// CreateUser creates a new user
func (c *StrapiClient) CreateUser(user User) (*User, error) {
	payload := map[string]interface{}{
		"username":  user.Username,
		"email":     user.Email,
		"confirmed": user.Confirmed,
		"blocked":   user.Blocked,
	}

	if user.Role != nil {
		if roleConnect, ok := user.Role["connect"].([]interface{}); ok && len(roleConnect) > 0 {
			if roleData, ok := roleConnect[0].(map[string]interface{}); ok {
				if id, ok := roleData["id"].(int); ok {
					payload["role"] = id
				} else if idFloat, ok := roleData["id"].(float64); ok {
					payload["role"] = int(idFloat)
				}
			}
		}
	}

	requestBody := map[string]interface{}{
		"data": payload,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/api/users", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// UpdateUser updates an existing user
func (c *StrapiClient) UpdateUser(id int, user User) (*User, error) {
	payload := map[string]interface{}{
		"username":  user.Username,
		"email":     user.Email,
		"confirmed": user.Confirmed,
		"blocked":   user.Blocked,
	}

	if user.Role != nil {
		if roleConnect, ok := user.Role["connect"].([]interface{}); ok && len(roleConnect) > 0 {
			if roleData, ok := roleConnect[0].(map[string]interface{}); ok {
				if idRole, ok := roleData["id"].(int); ok {
					payload["role"] = idRole
				} else if idFloat, ok := roleData["id"].(float64); ok {
					payload["role"] = int(idFloat)
				}
			}
		}
	}

	requestBody := map[string]interface{}{
		"data": payload,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.Endpoint+"/api/users/"+strconv.Itoa(id), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update user: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// DeleteUser deletes a user
func (c *StrapiClient) DeleteUser(id int) error {
	req, err := http.NewRequest("DELETE", c.Endpoint+"/api/users/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user: %s - %s", resp.Status, string(body))
	}

	return nil
}

// Role represents a Strapi role
type Role struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Permissions map[string]interface{} `json:"permissions,omitempty"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

// GetRoles retrieves all roles from Strapi
func (c *StrapiClient) GetRoles() ([]Role, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/users-permissions/roles", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get roles: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Roles []Role `json:"roles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Roles, nil
}

// GetRole retrieves a role by ID
func (c *StrapiClient) GetRole(id int) (*Role, error) {
	req, err := http.NewRequest("GET", c.Endpoint+"/api/users-permissions/roles/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get role: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Role Role `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Role, nil
}

// FindRoleByName retrieves a role by name
func (c *StrapiClient) FindRoleByName(name string) (*Role, error) {
	roles, err := c.GetRoles()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if strings.EqualFold(role.Name, name) {
			return &role, nil
		}
	}

	return nil, fmt.Errorf("role not found: %s", name)
}

// CreateRole creates a new role
func (c *StrapiClient) CreateRole(role Role) (*Role, error) {
	payload := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"type":        role.Type,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/api/users-permissions/roles", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create role: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Role Role `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Role, nil
}

// UpdateRole updates an existing role
func (c *StrapiClient) UpdateRole(id int, role Role) (*Role, error) {
	payload := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
	}

	if role.Type != "" {
		payload["type"] = role.Type
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.Endpoint+"/api/users-permissions/roles/"+strconv.Itoa(id), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update role: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Role Role `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Role, nil
}

// DeleteRole deletes a role
func (c *StrapiClient) DeleteRole(id int) error {
	req, err := http.NewRequest("DELETE", c.Endpoint+"/api/users-permissions/roles/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete role: %s - %s", resp.Status, string(body))
	}

	return nil
}
