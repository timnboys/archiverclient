package archiverclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TicketsBot/logarchiver"
	"github.com/rxdn/gdl/objects/channel/message"
	"net/http"
	"strings"
	"time"
)

type ArchiverClient struct {
	endpoint   string
	httpClient *http.Client
}

var ErrExpired = errors.New("log has expired")

func NewArchiverClient(endpoint string) ArchiverClient {
	endpoint = strings.TrimSuffix(endpoint, "/")

	return ArchiverClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				TLSHandshakeTimeout: time.Second * 3,
			},
		},
	}
}

func (c *ArchiverClient) Get(guildId uint64, ticketId int) ([]message.Message, error) {
	endpoint := fmt.Sprintf("%s/?guild=%d&id=%d", c.endpoint, guildId, ticketId)
	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, ErrExpired
		}

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

	res, err := c.httpClient.Post(endpoint, "application/json", bytes.NewReader(encoded))
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

func (c *ArchiverClient) GetModmail(guildId uint64, uuid string) ([]message.Message, error) {
	endpoint := fmt.Sprintf("%s/modmail?guild=%d&uuid=%s", c.endpoint, guildId, uuid)
	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return nil, ErrExpired
		}

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

func (c *ArchiverClient) StoreModmail(messages []message.Message, guildId uint64, uuid string, premium bool) error {
	encoded, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/modmail?guild=%d&id=%s", c.endpoint, guildId, uuid)
	if premium {
		endpoint += "&premium"
	}

	res, err := c.httpClient.Post(endpoint, "application/json", bytes.NewReader(encoded))
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

func (c *ArchiverClient) GetModmailKeys(guildId uint64) ([]logarchiver.StoredObject, error) {
	endpoint := fmt.Sprintf("%s/modmail/all?guild=%d", c.endpoint, guildId)

	res, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var decoded map[string]string
		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			return nil, err
		}

		return nil, errors.New(decoded["error"])
	}

	var decoded []logarchiver.StoredObject
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func (c *ArchiverClient) Encode(messages []message.Message, title string) ([]byte, error) {
	encoded, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/encode?title=%s", c.endpoint, title)
	res, err := c.httpClient.Post(endpoint, "application/json", bytes.NewReader(encoded))
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
