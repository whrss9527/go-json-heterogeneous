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
