package onedrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RenameResponse struct {
	URL string `json:"@microsoft.graph.downloadUrl"`
}

func (node *Node) Rename(id, name string) (err error) {
	req, err := node.NewRequest("PATCH", "https://graph.microsoft.com/v1.0/me/drive/items/"+id, bytes.NewBufferString(fmt.Sprintf(`{"name": "%s"}`, name)))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	renameResponse := &RenameResponse{}
	err = json.Unmarshal(data, renameResponse)
	if err != nil {
		return
	}
	return
}
