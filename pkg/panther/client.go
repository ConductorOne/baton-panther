package panther

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	token      string
	url        string
}

func NewClient(httpClient *http.Client, token, url string) *Client {
	return &Client{
		httpClient: httpClient,
		token:      token,
		url:        url,
	}
}

type UsersResponse struct {
	GraphQLError
	Data struct {
		Users []User `json:"users"`
	} `json:"data"`
}

type RolesResponse struct {
	GraphQLError
	Data struct {
		Roles []Role `json:"roles"`
	} `json:"data"`
}

type GraphQLError struct {
	// sometimes message is not in errors array (e.g when the token is invalid).
	Message string `json:"message,omitempty"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// GetUsers returns all users from Panther.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	query := `query users {
		users {
			id
			givenName
			familyName
			email
			status
			role {
				name
				id
			}
		}
  }`

	var res UsersResponse
	if err := c.doRequest(ctx, query, &res); err != nil {
		return nil, err
	}

	if res.Message != "" {
		return nil, fmt.Errorf(res.Message)
	}

	if res.Errors != nil {
		return nil, fmt.Errorf(res.Errors[0].Message)
	}

	return res.Data.Users, nil
}

// GetRoles returns all Panther roles.
func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	query := `query roles {
		roles {
			id
			name
		}
	}`

	var res RolesResponse
	if err := c.doRequest(ctx, query, &res); err != nil {
		return nil, err
	}

	if res.Message != "" {
		return nil, fmt.Errorf(res.Message)
	}

	if res.Errors != nil {
		return nil, fmt.Errorf(res.Errors[0].Message)
	}

	return res.Data.Roles, nil
}

func (c *Client) doRequest(ctx context.Context, query string, response interface{}) error {
	body, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("X-API-Key", c.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	return nil
}
