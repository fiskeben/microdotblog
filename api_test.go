package microdotblog

import (
	"bytes"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const posts string = `
{
	"version": "https://jsonfeed.org/version/1",
	"title": "Micro.blog - Ricco Førgaard",
	"home_page_url": "https://micro.blog/",
	"feed_url": "https://micro.blog/posts/ricco",
	"_microblog": {
	  "about": "https://micro.blog/about/api",
	  "id": "735",
	  "username": "ricco",
	  "bio": "I like to code stuff, grow stuff, brew stuff, cook stuff, solder stuff, and stuff.",
	  "is_following": true,
	  "is_you": true,
	  "following_count": 33
	},
	"author": {
	  "name": "Ricco Førgaard",
	  "url": "https://github.com/fiskeben",
	  "avatar": "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96"
	},
	"items": [
	  {
		"id": "218679",
		"content_html": "<p>I’m testing my micro.blog API client right now, so you may see some wierd posts :)</p>\n",
		"url": "http://micro.fiskeben.dk/2017/12/09/im-testing-my.html",
		"date_published": "2017-12-09T18:46:00+00:00",
		"author": {
		  "name": "Ricco Førgaard",
		  "url": "https://github.com/fiskeben",
		  "avatar": "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96",
		  "_microblog": {
			"username": "ricco"
		  }
		},
		"_microblog": {
		  "date_relative": "7:46 pm",
		  "is_favorite": false,
		  "is_deletable": true
		}
	  },
	  {
		"id": "218680",
		"content_html": "<p>Another test of my Go client library.</p>\n",
		"url": "http://micro.fiskeben.dk/2017/12/09/another-test-of.html",
		"date_published": "2017-12-09T18:44:00+00:00",
		"author": {
		  "name": "Ricco Førgaard",
		  "url": "https://github.com/fiskeben",
		  "avatar": "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96",
		  "_microblog": {
			"username": "ricco"
		  }
		},
		"_microblog": {
		  "date_relative": "7:44 pm",
		  "is_favorite": false,
		  "is_deletable": true
		}
	  },
	  {
		"id": "218675",
		"content_html": "<p>Just testing how to post from my Go micro.blog API.</p>\n",
		"url": "http://micro.fiskeben.dk/2017/12/09/just-testing-how.html",
		"date_published": "2017-12-09T18:31:00+00:00",
		"author": {
		  "name": "Ricco Førgaard",
		  "url": "https://github.com/fiskeben",
		  "avatar": "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96",
		  "_microblog": {
			"username": "ricco"
		  }
		},
		"_microblog": {
		  "date_relative": "7:31 pm",
		  "is_favorite": false,
		  "is_deletable": true
		}
	  }
	]
}`

type mockClient struct {
	responseData http.Response
}

type body struct {
	buf *bytes.Buffer
}

func (b body) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

func (b body) Close() error {
	return nil
}

func (m mockClient) Do(req *http.Request) (*http.Response, error) {
	return &m.responseData, nil
}

func makeMockClient(token, responseData string) APIClient {
	body := body{bytes.NewBufferString(responseData)}
	response := http.Response{Body: body, StatusCode: 200, Status: "OK"}
	c := apiClient{
		httpClient: aClient{
			httpClient: mockClient{responseData: response},
			token:      token,
		},
	}

	return c
}

func makeFailingMockClient(statusCode int, status string) APIClient {
	body := body{bytes.NewBufferString(status)}
	response := http.Response{Body: body, StatusCode: statusCode, Status: status}
	c := apiClient{
		httpClient: aClient{
			httpClient: mockClient{responseData: response},
			token:      "",
		},
	}

	return c
}

func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	val := f.Interface()

	switch someV := val.(type) {
	case int:
		return strconv.FormatInt(int64(someV), 10)
	case int64:
		return strconv.FormatInt(someV, 10)
	case string:
		return f.String()
	case bool:
		if f.Bool() {
			return "true"
		}
		return "false"
	case time.Time:
		return someV.String()
	}
	return ""
}

func TestGetPosts(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	feed, err := c.GetPosts()
	if err != nil {
		t.Error(err)
	}

	feedPropertiesTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"Version", "https://jsonfeed.org/version/1"},
		{"Title", "Micro.blog - Ricco Førgaard"},
		{"HomepageURL", "https://micro.blog/"},
		{"FeedURL", "https://micro.blog/posts/ricco"},
	}

	for _, tc := range feedPropertiesTestCases {
		if value := getField(feed, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

	microblogpropertiesTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"ID", "735"},
		{"Username", "ricco"},
		{"About", "https://micro.blog/about/api"},
		{"Bio", "I like to code stuff, grow stuff, brew stuff, cook stuff, solder stuff, and stuff."},
		{"IsFollowing", "true"},
		{"IsYou", "true"},
		{"FollowingCount", "33"},
	}

	props := feed.MicroblogProperties
	for _, tc := range microblogpropertiesTestCases {
		if value := getField(props, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

	authorPropertiesTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"Name", "Ricco Førgaard"},
		{"URL", "https://github.com/fiskeben"},
		{"Avatar", "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96"},
	}

	author := feed.Author
	for _, tc := range authorPropertiesTestCases {
		if value := getField(author, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

	if len(feed.Items) != 3 {
		t.Errorf("Not 3 items in feed (was %d)", len(feed.Items))
	}

	postTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"ID", "218679"},
		{"ContentHTML", "<p>I’m testing my micro.blog API client right now, so you may see some wierd posts :)</p>\n"},
		{"URL", "http://micro.fiskeben.dk/2017/12/09/im-testing-my.html"},
		{"DatePublished", "2017-12-09 18:46:00 +0000 +0000"},
	}

	post := feed.Items[0]
	for _, tc := range postTestCases {
		if value := getField(post, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

	postAuthorTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"Name", "Ricco Førgaard"},
		{"URL", "https://github.com/fiskeben"},
		{"Avatar", "https://www.gravatar.com/avatar/5a1b964e34cc2b1d3e233fa387b791b0?s=96"},
	}

	postAuthor := post.Author
	for _, tc := range postAuthorTestCases {
		if value := getField(postAuthor, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

	if postAuthor.MicroblogProperties.Username != "ricco" {
		t.Errorf("Post author metadata username is not equal to '%s (was '%s')", "ricco", postAuthor.MicroblogProperties.Username)
	}

	postMetadataTestCases := []struct {
		PropertyName  string
		ExpectedValue interface{}
	}{
		{"DateRelative", "7:46 pm"},
		{"IsFavorite", "false"},
		{"IsDeletable", "true"},
	}

	postMetadata := post.MicroblogProperties
	for _, tc := range postMetadataTestCases {
		if value := getField(postMetadata, tc.PropertyName); value != tc.ExpectedValue {
			t.Errorf("%s is not equal to %v (was '%s')", tc.PropertyName, tc.ExpectedValue, value)
		}
	}

}

func TestGetMentions(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	feed, err := c.GetMentions()
	if err != nil {
		t.Error(err)
	}
	if len(feed.Items) != 3 {
		t.Errorf("Returned feed doesn't look right")
	}
}

func TestGetFavourites(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	faves, err := c.GetFavourites()
	if err != nil {
		t.Error(err)
	}
	if len(faves.Items) != 3 {
		t.Errorf("Returned feed doesn't look right")
	}
}

func TestDiscover(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	interesting, err := c.Discover()
	if err != nil {
		t.Error(err)
	}
	if len(interesting.Items) != 3 {
		t.Errorf("Returned feed doesn't look right")
	}
}

func TestGetUserPosts(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	userPosts, err := c.GetUserPosts("manton")
	if err != nil {
		t.Error(err)
	}
	if len(userPosts.Items) != 3 {
		t.Errorf("Returned feed doesn't look right")
	}
}

func TestGetConversation(t *testing.T) {
	c := makeMockClient("ABCD12345", posts)
	convo, err := c.GetConversation(215436)
	if err != nil {
		t.Error(err)
	}
	if len(convo.Items) != 3 {
		t.Errorf("Returned feed doesn't look right")
	}
}

func TestCheck(t *testing.T) {
	responseBody := `{"count":5,"check_seconds":120}`
	c := makeMockClient("ABCD12345", responseBody)
	check, err := c.Check(215436)
	if err != nil {
		t.Error(err)
	}

	if check.Count != 5 {
		t.Errorf("Check call failed, got '%d' items, expected 5", check.Count)
	}

	if check.CheckSeconds != 120 {
		t.Errorf("Check call failed, got '%d' seconds, expected 120", check.CheckSeconds)
	}
}

func TestFavourite(t *testing.T) {
	c := makeMockClient("ABCD12345", "")
	err := c.Favourite(1234)
	if err != nil {
		t.Error(err)
	}

}

func TestUnfavourite(t *testing.T) {
	c := makeMockClient("ABCD12345", "")
	err := c.Unfavourite(1234)
	if err != nil {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	c := makeMockClient("ABCD12345", "")
	_, err := c.Post("Just testing how to post from my Go micro.blog API.")
	if err != nil {
		t.Error(err)
	}
}

func TestDeletePost(t *testing.T) {
	c := makeMockClient("ABCD12345", "")
	err := c.DeletePost(1234)
	if err != nil {
		t.Error(err)
	}
}

func TestFollow(t *testing.T) {
	c := makeMockClient("ABCD12345", "")
	if err := c.Follow("manton"); err != nil {
		t.Error(err)
	}
}

func TestNotFound(t *testing.T) {
	c := makeFailingMockClient(404, "Not found")
	_, err := c.GetPosts()

	if _, ok := err.(NotFound); !ok {
		t.Errorf("Expected HTTP 404 not found, got %v", err)
	}
}

func TestForbidden(t *testing.T) {
	c := makeFailingMockClient(403, "Forbidden")
	_, err := c.Post("I'm not allowd to do this")
	if _, ok := err.(Forbidden); !ok {
		t.Errorf("Expected HTTP 403 forbidden, got %v", err)
	}
}

func TestNotAuthorized(t *testing.T) {
	c := makeFailingMockClient(401, "Not authorized")
	err := c.DeletePost(123)
	if _, ok := err.(NotAuthorized); !ok {
		t.Errorf("Expected HTTP 401 not authorized, got %v", err)
	}
}

func TestInternalServerError(t *testing.T) {
	c := makeFailingMockClient(504, "Gateway timed out")
	_, err := c.GetConversation(1234)
	if _, ok := err.(ServerError); !ok {
		t.Errorf("Expected internal server error, got %v", err)
	}
}

func TestClientError(t *testing.T) {
	c := makeFailingMockClient(418, "I'm a teapot")
	err := c.Follow("manton")
	if _, ok := err.(ClientError); !ok {
		t.Errorf("Expected client error, got %v", err)
	}
}
