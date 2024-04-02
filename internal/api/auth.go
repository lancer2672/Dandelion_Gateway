package api

import (
	"encoding/json"
	"io"

	"github.com/lancer2672/Dandelion_Gateway/internal/constants"
	"github.com/lancer2672/Dandelion_Gateway/internal/helper"
)

func GetUserCredential(userId string) (data map[string]any, err error) {
	err = helper.RetryHandler(func() error {
		resp, err := helper.HttpClient.Get(constants.AUTH_PATH + "credential/" + userId)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return err
		}
		data = responseData["data"].(map[string]any)
		return nil
	})
	return
}
