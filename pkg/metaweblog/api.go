package metaweblog

import "github.com/kolo/xmlrpc"

type Client struct {
	*xmlrpc.Client
}

// NewClient url is like https://rpc.cnblogs.com/metaweblog/yahuian
func NewClient(url string) (*Client, error) {
	c, err := xmlrpc.NewClient(url, nil)
	if err != nil {
		return nil, err
	}

	return &Client{Client: c}, nil
}

type FileData struct {
	Bits xmlrpc.Base64 `xmlrpc:"bits"`
	Name string        `xmlrpc:"name"`
	Type string        `xmlrpc:"type"`
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) NewMediaObject(blogid, username, password string, file FileData) (string, error) {
	param := []any{
		blogid,
		username,
		password,
		file,
	}

	reply := struct {
		URL string `xmlrpc:"url"`
	}{}

	if err := c.Call("metaWeblog.newMediaObject", param, &reply); err != nil {
		return "", err
	}

	return reply.URL, nil
}
