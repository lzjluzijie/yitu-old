package onedrive

import (
	"fmt"
	"io"
	"net/http"
)

func (n *Node) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", n.AccessToken))
	req.Header.Add("User-Agent", "6tu")
	return
}
