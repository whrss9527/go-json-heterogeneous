# 强类型处理处理http请求时的异构json
两种方法⬇️
### 1. 重写 json 解析器，将 json 解析为强类型对象
```go
package main

import (
	`encoding/json`
	`fmt`
)

type IdP string

const (
	IdP_DEVICE   IdP = "device"
	IdP_WEIXIN   IdP = "weixin"
	IdP_APPLE    IdP = "apple"
	IdP_GOOGLE   IdP = "google"
	IdP_FACEBOOK IdP = "facebook"
)

func main() {
	jsonString := `{"idp":"device","credential":{"device_id":"123456"}}`
	var verifyCredentialRequest VerifyCredentialRequest
	if err := json.Unmarshal([]byte(jsonString), &verifyCredentialRequest); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", verifyCredentialRequest)
	println(verifyCredentialRequest.Credential.(*DeviceCredential).DeviceId)
}

type VerifyCredentialRequest struct {
	IdP        string      `json:"idp" required:"true"`
	Credential interface{} `json:"credential" required:"true"`
}

type DeviceCredential struct {
	DeviceId string `json:"device_id" required:"true"`
}

type WeixinCredential struct {
	Code string `json:"code" required:"true"`
}

type AppleCredential struct {
	IdentityToken string `json:"identity_token" required:"true"`
	FullName      string `json:"full_name" required:"true"`
}

type GoogleCredential struct {
	IdToken string `json:"id_token" required:"true"`
}

type FacebookCredential struct {
	UserAccessToken string `json:"user_access_token" required:"true"`
}

// 辅助结构体
type helperVerifyCredentialRequest struct {
	IdP        string          `json:"idp"`
	Credential json.RawMessage `json:"credential"`
}

// 实现 UnmarshalJSON 方法
func (v *VerifyCredentialRequest) UnmarshalJSON(data []byte) error {
	var helper helperVerifyCredentialRequest
	if err := json.Unmarshal(data, &helper); err != nil {
		return err
	}
	v.IdP = helper.IdP
	var cred interface{}
	switch helper.IdP {
	case string(IdP_DEVICE):
		cred = &DeviceCredential{}

	case string(IdP_WEIXIN):
		cred = &WeixinCredential{}

	case string(IdP_APPLE):
		cred = &AppleCredential{}

	case string(IdP_GOOGLE):
		cred = &GoogleCredential{}

	case string(IdP_FACEBOOK):
		cred = &FacebookCredential{}

	default:
		return fmt.Errorf("unsupported IdP type: %s", helper.IdP)
	}
	if err := json.Unmarshal(helper.Credential, &cred); err != nil {
		return err
	}
	v.Credential = cred
	return nil
}

```


### 2. 使用 mapstructure 进行转化
```go
package main

import (
	`encoding/json`
	`fmt`

	"github.com/mitchellh/mapstructure"
)

func main() {
	// 示例 JSON
	jsonStr := `{"idp": "device", "credential": {"device_id": "12345"}}`

	// 解析 JSON 到 map
	var request map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &request); err != nil {
		fmt.Println("Error:", err)
		return
	}

	idpType, ok := request["idp"].(string)
	if !ok {
		fmt.Println("Error: idp type is not a string")
		return
	}

	credentialData, ok := request["credential"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: credential data type is not correct")
		return
	}

	credential, err := parseCredential(credentialData, idpType)
	if err != nil {
		fmt.Println("Error parsing credential:", err)
		return
	}

	fmt.Printf("Parsed Credential: %+v\n", credential)
}

// 使用 mapstructure 解析凭证
func parseCredential(data map[string]interface{}, idpType string) (interface{}, error) {
	var result interface{}
	switch idpType {
	case "device":
		var cred DeviceCredential
		if err := mapstructure.Decode(data, &cred); err != nil {
			return nil, err
		}
		result = cred
	case "apple":
		var cred AppleCredential
		if err := mapstructure.Decode(data, &cred); err != nil {
			return nil, err
		}
		result = cred
	// 添加其他 case 以处理其他 IdP 类型
	default:
		return nil, fmt.Errorf("unsupported IdP type: %s", idpType)
	}
	return result, nil
}

type VerifyCredentialRequest struct {
	IdP        string      `json:"idp" required:"true"`
	Credential interface{} `json:"credential" required:"true"`
}

type DeviceCredential struct {
	DeviceId string `json:"device_id" required:"true" mapstructure:"device_id"`
}

type WeixinCredential struct {
	Code string `json:"code" required:"true"`
}

type AppleCredential struct {
	IdentityToken string `json:"identity_token" required:"true"`
	FullName      string `json:"full_name" required:"true"`
}

type GoogleCredential struct {
	IdToken string `json:"id_token" required:"true"`
}

type FacebookCredential struct {
	UserAccessToken string `json:"user_access_token" required:"true"`
}

```