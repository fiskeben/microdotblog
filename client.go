package microdotblog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
	endpoint := fmt.Sprintf("https://micro.blog/posts/%s", username)
	data, err := a.httpClient.getAndRead(endpoint)
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) GetConversation(ID int64) (*Feed, error) {
	endpoint := fmt.Sprintf("https://micro.blog/posts/conversation?id=%d", ID)
	data, err := a.httpClient.getAndRead(endpoint)
	if err != nil {
		return nil, err
	}
	return feedFromResponse(data)
}

func (a apiClient) Check(sinceID int64) (*Check, error) {
	endpoint := fmt.Sprintf("https://micro.blog/posts/check?since_id=%d", sinceID)
	data, err := a.httpClient.getAndRead(endpoint)
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
	endpoint := fmt.Sprintf("https://micro.blog/posts/favorites?id=%d", ID)

	_, err := a.httpClient.postAndRead(endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a apiClient) Unfavourite(ID int64) error {
	endpoint := fmt.Sprintf("https://micro.blog/posts/favorites/%d", ID)
	if err := a.httpClient.delete(endpoint); err != nil {
		return err
	}
	return nil
}

func (a apiClient) Reply(ID int64, message string) (*Post, error) {
	endpoint := fmt.Sprintf("https://micro.blog/posts/reply")
	data := url.Values{}
	data.Add("id", strconv.FormatInt(ID, 10))
	data.Add("text", message)

	return a.sendPost(endpoint, data.Encode())
}

func (a apiClient) DeletePost(ID int64) error {
	endpoint := fmt.Sprintf("https://micro.blog/posts/%d", ID)
	if err := a.httpClient.delete(endpoint); err != nil {
		return err
	}
	return nil
}

func (a apiClient) Follow(username string) error {
	endpoint := fmt.Sprintf("https://micro.blog/users/follow?username=%s", username)
	if _, err := a.httpClient.postAndRead(endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (a apiClient) Unfollow(username string) error {
	endpoint := fmt.Sprintf("https://micro.blog/users/unfollow?username=%s", username)
	if _, err := a.httpClient.postAndRead(endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (a apiClient) Followers(username string) ([]User, error) {
	endpoint := fmt.Sprintf("https://micro.blog/users/following/%s", username)
	bytes, err := a.httpClient.getAndRead(endpoint)
	if err != nil {
		return nil, err
	}

	var users = []User{}
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a apiClient) Post(message string) (*Post, error) {
	endpoint := "https://micro.blog/micropub"

	data := url.Values{}
	data.Set("h", "entry")
	data.Set("content", message)

	return a.sendPost(endpoint, data.Encode())
}

func (a apiClient) sendPost(endpoint, payload string) (*Post, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(payload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.httpClient.token)

	res, err := a.httpClient.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err = newAPIError(res.StatusCode, res.Body); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Until Micro.blog starts returning the created post
	// there's no need to read the response :(
	// bytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, err
	// }

	return &Post{}, nil
}

func (a apiClient) PostPhoto(message string, photo Photo) (*Post, error) {
	return &Post{}, nil
}

func (a aClient) getAndRead(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", a.token)

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = newAPIError(res.StatusCode, res.Body); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (a aClient) postAndRead(endpoint string, payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", a.token)

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = newAPIError(res.StatusCode, res.Body); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (a aClient) delete(endpoint string) error {
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	res, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	if err = newAPIError(res.StatusCode, res.Body); err != nil {
		return err
	}
	return nil
}

func feedFromResponse(data []byte) (*Feed, error) {
	f := &Feed{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
