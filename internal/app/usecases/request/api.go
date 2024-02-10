package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const _asyncServerURL = "http://localhost:5000"

var ErrResponseAsyncServerNoOk = errors.New("async server response status not 200")

type bodyForAsyncServer struct {
	RequestID int `json:"request_id"`
}

func sendRequestToAsyncServer(reqID int) error {
	data, err := json.Marshal(bodyForAsyncServer{reqID})
	if err != nil {
		return fmt.Errorf("send request to async server: %w", err)
	}

	res, err := http.Post(_asyncServerURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("post request to async server: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrResponseAsyncServerNoOk
	}
	return nil
}
