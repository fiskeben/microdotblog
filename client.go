package microdotblog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewAPIClient creates a new client with a default HTTP client.
// Pass an access token here.
func NewAPIClient(token string) APIClient {
	c := apiClient{
		httpClient: aClient{
			httpClient: http.DefaultClient,
			token:      token,
		},
	}

	return c
}

type internalClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type aClient struct {
	httpClient internalClient
	token      string
}

type apiClient struct {
	httpClient aClient
}

func (a apiClient) GetPosts() (*Feed, error) {
	data, err := a.httpClient.getAndRead("https://micro.blog/posts/all")
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) GetMentions() (*Feed, error) {
	data, err := a.httpClient.getAndRead("https://micro.blog/posts/mentions")
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) GetFavourites() (*Feed, error) {
	data, err := a.httpClient.getAndRead("https://micro.blog/posts/favorites")
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) Discover() (*Feed, error) {
	data, err := a.httpClient.getAndRead("https://micro.blog/posts/discover")
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) GetUserPosts(username string) (*Feed, error) {
	url := fmt.Sprintf("https://micro.blog/posts/%s", username)
	data, err := a.httpClient.getAndRead(url)
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) GetConversation(ID int64) (*Feed, error) {
	url := fmt.Sprintf("https://micro.blog/posts/conversation?id=%d", ID)
	data, err := a.httpClient.getAndRead(url)
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) Check(sinceID int64) (*Check, error) {
	url := fmt.Sprintf("https://micro.blog/posts/check?since_id=%d", sinceID)
	data, err := a.httpClient.getAndRead(url)
	if err != nil {
		return nil, err
	}
	c := &Check{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a apiClient) Favourite(ID int64) error {
	return nil
}

func (a apiClient) Unfavourite(ID int64) error {
	return nil
}

func (a apiClient) Reply(ID int64, message string) (*Post, error) {
	return &Post{}, nil
}

func (a apiClient) DeletePost(ID int64) error {
	return nil
}

func (a apiClient) Follow(username string) error {
	return nil
}

func (a apiClient) Unfollow(username string) error {
	return nil
}

func (a aClient) getAndRead(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", a.token)

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func feedFromResponse(data []byte) (*Feed, error) {
	f := &Feed{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
