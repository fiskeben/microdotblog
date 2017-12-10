# A Go client for micro.blog

This is a Go implementation of the
[micro.blog](https://micro.blog)
[JSON API](http://help.micro.blog/2017/api-json/).

Currently only the `GET` methods are implemented.

## Usage

You need an API key.
You can [get one here](https://micro.blog/account/apps).

```go
package main

import (
    micro "github.com/fiskeben/microdotblog"
    "fmt"
)

func main() {
    client := micro.NewAPIClient("your-api-key")
    feed, err := client.GetPosts()
    if err != nil {
        panic(err)
    }
    fmt.Printf("It got a feed called %s with %d posts in it", feed.Title, len(feed.Posts))
}
```

## TODO

* [x] Implement remaining methods.
* [x] Testing. Currently only have tests that go directly to micro.blog.
* [x] Better errors. Right now raw http and unmarshalling errors are returned.
* [ ] Even better errors. Make it easier to test the type of error.
* [ ] Fix some issues with `DELETE` method.

## Follow me

You can [follow me on micro.blog here](https://micro.blog/ricco).
