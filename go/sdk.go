package sdk

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeiXinSDK struct {
	AppId     string
	AppSecret string
}

type SdkQuery struct {
	AppId     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
	TimeStamp string `json:"timestamp"`
}

var (
	access_token_api = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	jsapi_ticket_api = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

func NewSDK(appid string, appsecret string) (*WeiXinSDK, error) {

	wxSdk := new(WeiXinSDK)

	wxSdk.AppId = appid
	wxSdk.AppSecret = appsecret

	return wxSdk, nil
}

/**
* @todo get jssdk singature params example
 */
func (this *WeiXinSDK) GetJssdkSignatureParams(noncestr string, timestamps string, url string) (querys SdkQuery, err error) {

	var (
		accessTokenBody  []byte
		accessTokenRes   map[string]interface{}
		accessToken      string
		ticketBody       []byte
		ticketRes        map[string]interface{}
		jsapiTicket      string
		signatureJoinStr string
		signatureJoinRes string
		errMsg           string
	)

	if accessTokenBody, err = this.GetAccessToken(); err != nil {
		return
	}

	if err = json.Unmarshal(accessTokenBody, &accessTokenRes); err != nil {
		return
	}

	if _, ok := accessTokenRes["access_token"]; !ok {
		errMsg = accessTokenRes["errmsg"].(string)
		err = errors.New(errMsg)
		return
	}

	accessToken = accessTokenRes["access_token"].(string)

	if ticketBody, err = this.GetJsapiTicket(accessToken); err != nil {
		return
	}

	if err = json.Unmarshal(ticketBody, &ticketRes); err != nil {
		return
	}

	if _, ok := ticketRes["ticket"]; !ok {
		errMsg = ticketRes["errmsg"].(string)
		err = errors.New(errMsg)
		return
	}

	jsapiTicket = ticketRes["ticket"].(string)

	signatureJoinStr = fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", jsapiTicket, noncestr, timestamps, url)

	signatureJoinRes = signature(signatureJoinStr)

	querys.AppId = this.AppId
	querys.NonceStr = noncestr
	querys.Signature = signatureJoinRes
	querys.TimeStamp = timestamps

	return
}

/**
* @todo signature
**/
func signature(joinstr string) string {

	h := sha1.New()

	h.Write([]byte(joinstr))

	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
* @todo get access_token
* @tips please custom the cache
**/
func (this *WeiXinSDK) GetAccessToken() ([]byte, error) {

	var (
		url  string
		body []byte
		req  *http.Response
		err  error
	)

	url = fmt.Sprintf(access_token_api, this.AppId, this.AppSecret)

	if req, err = http.Get(url); err != nil {
		return nil, err
	}

	defer req.Body.Close()

	if body, err = ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	}

	return body, nil
}

/**
* @todo get jsapi_ticket
* @tips please custom the cache
**/
func (this *WeiXinSDK) GetJsapiTicket(access_token string) ([]byte, error) {

	var (
		url  string
		body []byte
		req  *http.Response
		err  error
	)

	url = fmt.Sprintf(jsapi_ticket_api, access_token)

	if req, err = http.Get(url); err != nil {
		return nil, err
	}

	defer req.Body.Close()

	if body, err = ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	}

	return body, nil
}
