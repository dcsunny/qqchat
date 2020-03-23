package safe

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/qqchat/define"
	"github.com/dcsunny/qqchat/util"
)

const (
	imgSecCheckUrl     = "https://api.q.qq.com/api/json/security/ImgSecCheck"
	msgSecCheckUrl     = "https://api.q.qq.com/api/json/security/MsgSecCheck"
	mediaCheckAsyncUrl = "https://api.q.qq.com/api/json/security/MediaCheckAsync"
)

func (s *WxSafe) ImgSecCheck(filename string, fileBytes []byte) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", imgSecCheckUrl, accessToken)

	fields := []util.MultipartFormField{
		{
			IsFile:    true,
			Fieldname: "media",
			Filename:  filename,
			Value:     fileBytes,
		},
		{
			IsFile:    false,
			Fieldname: "appid",
			Value:     []byte(s.AppID),
		},
	}

	var response []byte
	response, err = util.PostMultipartForm(fields, uri)
	if err != nil {
		return
	}
	fmt.Println(string(response))
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("ImgSecCheck error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (s *WxSafe) MsgSecCheck(content string) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", msgSecCheckUrl, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{
		"appid":   s.AppID,
		"content": content,
	})
	if err != nil {
		return
	}
	fmt.Println(string(response))
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("MsgSecCheck error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return err
}

func (s *WxSafe) MediaCheckAsync(mediaUrl string, mediaType string) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", mediaCheckAsyncUrl, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{
		"appid":      s.AppID,
		"media_url":  mediaUrl,
		"media_type": mediaType,
	})
	if err != nil {
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	fmt.Println(string(response))
	if result.ErrCode != 0 {
		err = fmt.Errorf("MediaCheckAsync error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return err
}
