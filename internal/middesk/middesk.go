package middesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultURL = "https://api.middesk.com/v1"

type MiddeskBusinessId struct {
	Id string `json:"id"`
}

type MiddeskBusinessFormation struct {
	EntityType string `json:"entity_type"`
}

type MiddeskBusinessRegistration struct {
	EntityType string `json:"entity_type"`
}

type MiddeskBusinessRaw struct {
	Id            string                        `json:"id"`
	ExternalId    *string                       `json:"external_id"`
	CreatedAt     time.Time                     `json:"created_at"`
	Status        string                        `json:"status"`
	Formation     MiddeskBusinessFormation      `json:"formation"`
	Registrations []MiddeskBusinessRegistration `json:"registrations"`
}

type MiddeskBusiness struct {
	Id                       string
	ExternalId               *string
	CreatedAt                time.Time
	Status                   string
	FormationEntityType      string
	RegistrationsEntityTypes []string
}

type MiddeskBusinesses struct {
	Object     string              `json:"object"`
	Data       []MiddeskBusinessId `json:"data"`
	Url        string              `json:"url"`
	HasMore    bool                `json:"has_more"`
	TotalCount int                 `json:"total_count"`
}

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewClient(API_KEY string) (*Client, error) {
	return &Client{
		baseURL: defaultURL,
		apiKey:  API_KEY,
		client:  http.DefaultClient,
	}, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return c.client.Do(req)
}

func (c *Client) GetMiddeskBusinessIds(page_number int, per_page int) (*MiddeskBusinesses, error) {
	url := fmt.Sprintf("%s/businesses?page=%d&per_page=%d", c.baseURL, page_number, per_page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	businesses := &MiddeskBusinesses{}
	if err := json.NewDecoder(resp.Body).Decode(businesses); err != nil {
		return nil, err
	}
	return businesses, nil
}

func (c *Client) GetMiddeskBusiness(id string) (*MiddeskBusinessRaw, error) {
	url := fmt.Sprintf("%s/businesses/%s", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	business := &MiddeskBusinessRaw{}
	if err := json.NewDecoder(resp.Body).Decode(business); err != nil {
		return nil, err
	}
	return business, nil
}
