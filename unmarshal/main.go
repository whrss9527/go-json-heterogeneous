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
