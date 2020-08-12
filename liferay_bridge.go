package liferay_bridge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/yukiz97/go_config"
	"github.com/yukiz97/go_utils"
)

var urlLiferayAPI string
var tokenAuthen string
var companyID string

func GetUserIDByScreenNane(screenName string) int {
	strURL := urlLiferayAPI + "user/get-user-id-by-screen-name/company-id/" + companyID + "/screen-name/" + screenName
	reqURL, _ := url.Parse(strURL)

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	userID, _ := strconv.Atoi(string(byteData))

	return userID
}

func AuthenLiferay(userName string, password string) (bool, map[string]interface{}) {
	tokenAuthen := go_utils.EncodeStringToBase64(userName + ":" + password)

	reqURL, _ := url.Parse(urlLiferayAPI + "user/get-current-user")

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	mapData := make(map[string]interface{})

	err = json.Unmarshal(byteData, &mapData)

	_, authenSuccess := mapData["userId"]

	return authenSuccess, mapData
}

func GetUserByID(userID int) (bool, map[string]interface{}) {
	reqURL, _ := url.Parse(urlLiferayAPI + "user/get-user-by-id/user-id/" + strconv.Itoa(userID))

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	mapData := make(map[string]interface{})

	err = json.Unmarshal(byteData, &mapData)

	_, authenSuccess := mapData["userId"]

	return authenSuccess, mapData
}

func GetUserRoles(userIdentity interface{}) []interface{} {
	userID := getUserIDByUserIdentityType(userIdentity)

	strURL := urlLiferayAPI + "role/get-user-roles/user-id/" + fmt.Sprint(userID)
	reqURL, _ := url.Parse(strURL)

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	var data []interface{}

	err = json.Unmarshal(byteData, &data)

	if err != nil {
		panic(err)
	}

	return data
}

func ChangePassword(screenName string, password string) bool {
	userID := GetUserIDByScreenNane(screenName)

	strURL := urlLiferayAPI + "user/update-password/user-id/" + fmt.Sprint(userID) + "/password1/" + password + "/password2/" + password + "/password-reset/false"
	reqURL, _ := url.Parse(strURL)

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	mapData := make(map[string]interface{})

	err = json.Unmarshal(byteData, &mapData)

	_, updateSuccess := mapData["userId"]

	return updateSuccess
}

func GetUserOrganizations(userIdentity interface{}) []interface{} {
	userID := getUserIDByUserIdentityType(userIdentity)

	strURL := urlLiferayAPI + "organization/get-user-organizations/user-id/" + fmt.Sprint(userID)
	reqURL, _ := url.Parse(strURL)

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	var data []interface{}

	err = json.Unmarshal(byteData, &data)

	if err != nil {
		panic(err)
	}

	return data
}

func GetSubOrganizations(parentIdOrg int) []interface{} {
	strURL := urlLiferayAPI + "organization/get-organizations/company-id/" + companyID + "/parent-organization-id/" + fmt.Sprint(parentIdOrg)
	reqURL, _ := url.Parse(strURL)

	request := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"Basic " + tokenAuthen},
		},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("Error:", err)
	}

	byteData, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	var data []interface{}

	err = json.Unmarshal(byteData, &data)

	if err != nil {
		panic(err)
	}

	return data
}

func getUserIDByUserIdentityType(userIdentity interface{}) int {
	var userID int
	switch go_utils.TypeOf(userIdentity) {
	case "string":
		userID = GetUserIDByScreenNane(userIdentity.(string))
		break
	case "int":
		userID = userIdentity.(int)
		break
	}

	return userID
}

func InitLiferayBridge() {
	var isErr bool = false
	args := os.Args
	if len(args) == 1 {
		if args[1] != "" {
			go_config.InitConfig(args[1], Configuration{})
			urlLiferayAPI = go_config.GetConfigPropery("Liferay_requesturl")
			tokenAuthen = go_config.GetConfigPropery("Liferay_authen_token")
			companyID = go_config.GetConfigPropery("Company_ID")
		} else {
			isErr = true
		}
	} else {
		isErr = true
	}

	if isErr {
		panic("Tham số không chính xác vui lòng kiểm tra lại!! /1 = configpath/")
	}
}
