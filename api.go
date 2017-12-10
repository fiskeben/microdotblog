package microdotblog

import "time"

// Feed represents an entire feed with posts and metadata.
type Feed struct {
	Version             string `json:"version"`
	Title               string `json:"title"`
	HomepageURL         string `json:"home_page_url"`
	FeedURL             string `json:"feed_url"`
	Items               []Post `json:"items"`
	Author              Author `json:"author"`
	MicroblogProperties struct {
		About          string `json:"about"`
		ID             int64  `json:"id,string"`
		Username       string `json:"username"`
		Bio            string `json:"bio"`
		IsFollowing    bool   `json:"is_following"`
		IsYou          bool   `json:"is_you"`
		FollowingCount int    `json:"following_count"`
	} `json:"_microblog"`
}

// Post represents a single post.
type Post struct {
	ID                  int64     `json:"id,string"`
	URL                 string    `json:"url"`
	ContentHTML         string    `json:"content_html"`
	DatePublished       time.Time `json:"date_published"`
	Author              Author    `json:"author"`
	MicroblogProperties struct {
		IsDeletable  bool   `json:"is_deletable"`
		IsFavorite   bool   `json:"is_favourite"`
		DateRelative string `json:"date_relative"`
	} `json:"_microblog"`
}

// Photo represents a photo.
// TODO: implement as a io.Reader (or ReadCloser)
type Photo struct {
	path string
}

// Author is a represetation of the author of a post.
type Author struct {
	Name                string `json:"name"`
	URL                 string `json:"url"`
	Avatar              string `json:"avatar"`
	MicroblogProperties struct {
		Username    string `json:"username"`
		IsFollowing bool   `json:"is_following"`
	} `json:"_microblog"`
}

// Check is returned when checking for new posts.
type Check struct {
	Count        int `json:"count"`
	CheckSeconds int `json:"check_seconds"`
}

// APIClient gives access to the API.
type APIClient interface {
	// GetPosts gets all posts from a feed.
	GetPosts() (*Feed, error)

	// GetMentions gets a feed with mentions of the current user.
	GetMentions() (*Feed, error)

	// GetFavourites gets a feed of the current user's favourites.
	GetFavourites() (*Feed, error)

	// Discover returns a feed of curated posts.
	Discover() (*Feed, error)

	// GetUserPosts gets the timeline of the specified user.
	GetUserPosts(username string) (*Feed, error)

	// GetConversation gets all replies to a post.
	GetConversation(ID int64) (*Feed, error)

	// Check looks for posts newer than the post with the specified ID.
	Check(sinceID int64) (*Check, error)

	// Favourite marks the post with the given ID a favourite.
	Favourite(ID int64) error

	// Unfavourite removes favourite flag from the post with the specified ID.
	Unfavourite(ID int64) error

	// Reply sends a new Post with the specified message as a reply to the Post
	// with the given ID.
	Reply(ID int64, message string) (*Post, error)

	// DeletePost removes the Post with the given ID.
	DeletePost(ID int64) error

	// Follow will let the current user start following the user with the
	// given username.
	Follow(username string) error

	// Unfollow will remove the user with the specified username from the list
	// of users the current user follows.
	Unfollow(username string) error

	// Post posts a new update to the blog.
	Post(message string) (*Post, error)

	// PostPhoto posts a new update including a photo.
	PostPhoto(message string, photo Photo) (*Post, error)
}
