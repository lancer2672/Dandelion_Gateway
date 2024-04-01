package api

import (
	"fmt"

	"github.com/lancer2672/Dandelion_Gateway/internal/constants"
	"github.com/lancer2672/Dandelion_Gateway/internal/helper"
)

func GetUserCredential(apikey string) (user map[string]string, credential map[string]string, err error) {
	err = helper.RetryHandler(func() error {

		resp, err := helper.HttpClient.Get(constants.AUTH_PATH + "checkapikey?apikey=" + apikey)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		fmt.Println("resp", resp)
		// var result struct {
		// 	Role       Role
		// 	Permission Permission
		// }
		// err = json.NewDecoder(resp.Body).Decode(&result)
		// if err != nil {
		// 	return nil
		// }
		// role = &result.Role
		// permission = &result.Permission
		return nil
	})
	return
}
