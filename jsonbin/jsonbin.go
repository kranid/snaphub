package jsonbin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Resource)
}

type MetaData struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Private   bool      `json:"private"`
}

type BinCreateResp struct {
	MetaData MetaData        `json:"metadata"`
	Record   json.RawMessage `json:"record"`
}

type SnapStore struct {
	ApiKey    string
	AccessKey string
	rootUrl   string
}

func JsonBinNew(apiKey, accessKey string) SnapStore {
	return SnapStore{
		AccessKey: accessKey,
		ApiKey:    apiKey,
		rootUrl:   "https://api.jsonbin.io/v3/b",
	}
}

func (jb SnapStore) Create(body io.Reader, name string) (string, error) {
	req, _ := http.NewRequest(http.MethodPost, jb.rootUrl, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", jb.ApiKey)
	req.Header.Set("X-Access-Key", jb.AccessKey)
	req.Header.Set("X-Bin-Name", name)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error %s", resp.Status)
	}
	var result BinCreateResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.MetaData.Id, nil
}

func (jb SnapStore) Get(id string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/%s", jb.rootUrl, id)
	fmt.Printf("url: %s \n", url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-Master-Key", jb.ApiKey)
	req.Header.Set("X-Access-Key", jb.AccessKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, &NotFoundError{Resource: resp.Status}
		}
		return nil, fmt.Errorf("error! JsonBin has returned %s", resp.Status)
	}
	bin := BinCreateResp{}
	err = json.NewDecoder(resp.Body).Decode(&bin)
	return bin.Record, err
}
