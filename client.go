package archiverclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rxdn/gdl/objects/channel/message"
	"net/http"
	"strings"
	"time"
)

type ArchiverClient struct {
	endpoint string
}

func NewArchiverClient(endpoint string) ArchiverClient {
	endpoint = strings.TrimSuffix(endpoint, "/")

	return ArchiverClient{
		endpoint: endpoint,
	}
}

func (c *ArchiverClient) Get(guildId uint64, ticketId int) ([]message.Message, error) {
	endpoint := fmt.Sprintf("%s/?guild=%d&id=%d", c.endpoint, guildId, ticketId)
	httpClient := newHttpClient()
	res, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var decoded map[string]string
		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			return nil, err
		}

		return nil, errors.New(decoded["message"])
	} else {
		var messages []message.Message
		if err := json.NewDecoder(res.Body).Decode(&messages); err != nil {
			return nil, err
		}

		return messages, nil
	}
}

func (c *ArchiverClient) Store(messages []message.Message, guildId uint64, ticketId int, premium bool) error {
	encoded, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/?guild=%d&id=%d", c.endpoint, guildId, ticketId)
	if premium {
		endpoint += "&premium"
	}

	httpClient := newHttpClient()
	res, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var decoded map[string]string
		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			return err
		}

		return errors.New(decoded["message"])
	}

	return nil
}

func (c *ArchiverClient) Encode(messages []message.Message, guildId uint64, ticketId int) ([]byte, error) {
	encoded, err := json.Marshal(messages); if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/encode", c.endpoint)
	httpClient := newHttpClient()
	res, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var decoded map[string]string
		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			return nil, err
		}

		return nil, errors.New(decoded["message"])
	} else {
		var buff bytes.Buffer
		if _, err := buff.ReadFrom(res.Body); err != nil {
			return nil, err
		}

		return buff.Bytes(), nil
	}
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Second * 3,
		},
	}
}
