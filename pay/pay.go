package pay

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/dcsunny/qqchat/context"
	"github.com/dcsunny/qqchat/util"
)

const (
	payGateway  = "https://qpay.qq.com/cgi-bin/pay/qpay_unified_order.cgi"
	mchTransUri = "https://api.qpay.qq.com/cgi-bin/epay/qpay_epay_b2c.cgi"
	sendRedUri  = "https://api.qpay.qq.com/cgi-bin/hongbao/qpay_hb_mch_send.cgi"
)

// Pay struct extends context
type Pay struct {
	*context.Context
}

// 传入的参数，用于生成 prepay_id 的必需参数
// PayParams was NEEDED when request unifiedorder
type Params struct {
	TotalFee     int
	CreateIP     string
	Body         string
	OutTradeNo   string
	TradeType    string
	MiniAppParam string
	OpenID       string
	ContractCode string
	PromotionTag string
	Attach       string
	//以下红包使用
	Wishing   string
	SendName  string
	ActName   string
	IconID    int
	BannerID  int
	NotifyUrl string
}

// PayConfig 是传出用于 jsdk 用的参数
type PayConfig struct {
	AppID     string `xml:"appId" json:"appId"`
	TimeStamp string `xml:"timeStamp" json:"timeStamp"`
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`
	Package   string `xml:"package" json:"package"`
	SignType  string `xml:"signType" json:"signType"`
	PaySign   string `xml:"paySign" json:"paySign"`
}

type PreOrder struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	TradeType  string `xml:"trade_type,omitempty"`
	PrePayID   string `xml:"prepay_id,omitempty"`
	CodeURL    string `xml:"code_url,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

type AppPayConfig struct {
	AppID     string `xml:"appid" json:"appid"`
	PartnerID string `xml:"partnerid" json:"partnerid"`
	PrePayID  string `xml:"prepayid" json:"prepayid"`
	Package   string `xml:"package" json:"package"`
	NonceStr  string `xml:"noncestr" json:"noncestr"`
	Timestamp string `xml:"timestamp" json:"timestamp"`
	Sign      string `xml:"sign" json:"sign"`
}

// payResult 是 unifie order 接口的返回
type payResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	TradeType  string `xml:"trade_type,omitempty"`
	PrePayID   string `xml:"prepay_id,omitempty"`
	CodeURL    string `xml:"code_url,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

//payRequest 接口请求参数
type payRequest struct {
	AppID          string `xml:"appid,omitempty"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Body           string `xml:"body"`
	Attach         string `xml:"attach,omitempty"`        //附加数据
	OutTradeNo     string `xml:"out_trade_no"`            //商户订单号
	FeeType        string `xml:"fee_type"`                //标价币种
	TotalFee       int    `xml:"total_fee"`               //标价金额
	SpbillCreateIp string `xml:"spbill_create_ip"`        //终端IP
	TimeStart      string `xml:"time_start,omitempty"`    //交易起始时间
	TimeExpire     string `xml:"time_expire,omitempty"`   //交易结束时间
	LimitPay       string `xml:"limit_pay,omitempty"`     //
	ContractCode   string `xml:"contract_code,omitempty"` //商户侧记录的用户代扣协议序列号，支付中开通代扣必传
	PromotionTag   string `xml:"promotion_tag,omitempty"` //指定本单参与某个QQ钱包活动或活动档位的标识，包含两个标识：sale_tag --- 不同活动的匹配标志 level_tag --- 同一活动不同优惠档位的标志，可不填  格式如下（本字段参与签名）：promotion_tag=level_tag=xxx&sale_tag=xxx
	TradeType      string `xml:"trade_type"`              //交易类型
	NotifyUrl      string `xml:"notify_url"`              //通知地址
	DeviceInfo     string `xml:"device_info,omitempty"`
	MiniAppParam   string `xml:"mini_app_param"`
}

type NotifyResult struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	Appid         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	DeviceInfo    string `xml:"device_info"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	FeeType       string `xml:"fee_type"`
	TotalFee      string `xml:"total_fee"`
	CashFee       string `xml:"cash_fee"`
	CouponFee     string `xml:"coupon_fee"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEnd       string `xml:"time_end"`
	Openid        string `xml:"openid"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"`
}

type NotifyReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// NewPay return an instance of Pay package
func NewPay(ctx *context.Context) *Pay {
	pay := Pay{Context: ctx}
	return &pay
}

func (pcf *Pay) PrePayIdByJs(p *Params) (prePayID string, err error) {
	p.TradeType = "JSAPI"
	return pcf.PrePayId(p)
}

func (pcf *Pay) PrePayIdByApp(p *Params) (prePayID string, err error) {
	p.TradeType = "APP"
	return pcf.PrePayId(p)
}

func (pcf *Pay) PrePayOrderByJs(p *Params) (payOrder PreOrder, err error) {
	p.TradeType = "JSAPI"
	return pcf.PrePayOrder(p)
}

func (pcf *Pay) PrePayOrderByApp(p *Params) (payOrder PreOrder, err error) {
	p.TradeType = "APP"
	return pcf.PrePayOrder(p)
}

func (pcf *Pay) PrePayOrderByMiniApp(p *Params) (payOrder PreOrder, err error) {
	p.TradeType = "MINIAPP"
	return pcf.PrePayOrder(p)
}

func (pcf *Pay) PrePayOrder(p *Params) (payOrder PreOrder, err error) {
	nonceStr := util.RandomStr(32)

	request := payRequest{
		AppID:          pcf.AppID,
		MchID:          pcf.PayMchID,
		NonceStr:       nonceStr,
		Body:           p.Body,
		OutTradeNo:     p.OutTradeNo,
		TotalFee:       p.TotalFee,
		SpbillCreateIp: p.CreateIP,
		NotifyUrl:      pcf.PayNotifyURL,
		TradeType:      p.TradeType,
		MiniAppParam:   p.MiniAppParam,
		Attach:         p.Attach,
	}
	sign, err := pcf.Sign(&request, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return payOrder, err
	}
	request.Sign = sign
	rawRet, err := util.PostXML(payGateway, request, "payRequest", nil)
	if err != nil {
		return PreOrder{}, errors.New(err.Error())
	}
	err = xml.Unmarshal(rawRet, &payOrder)
	if err != nil {
		return payOrder, errors.New(err.Error())
	}
	if payOrder.ReturnCode == "SUCCESS" {
		//pay success
		if payOrder.ResultCode == "SUCCESS" {
			return payOrder, nil
		}
		return payOrder, errors.New(payOrder.ErrCode + payOrder.ErrCodeDes)
	} else {
		return payOrder, errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "] [sign : " + sign + "]")
	}
}

// PrePayId will request wechat merchant api and request for a pre payment order id
func (pcf *Pay) PrePayId(p *Params) (prePayID string, err error) {
	order, err := pcf.PrePayOrder(p)
	if err != nil {
		return
	}
	if order.PrePayID == "" {
		err = errors.New("empty prepayid")
	}
	prePayID = order.PrePayID
	return
}

func (pcf *Pay) JSPayParams(prePayID string) PayConfig {
	payConf := PayConfig{
		AppID:     pcf.AppID,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  util.RandomStr(32),
		Package:   fmt.Sprintf("prepay_id=%s", prePayID),
		SignType:  "MD5",
	}
	str := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s", payConf.AppID, payConf.NonceStr, payConf.Package, payConf.SignType, payConf.TimeStamp, pcf.PayKey)
	payConf.PaySign = util.MD5Sum(str)
	return payConf
}

func (pcf *Pay) AppPayParams(prePayID string) AppPayConfig {
	payConf := AppPayConfig{
		AppID:     pcf.AppID,
		PartnerID: pcf.PayMchID,
		PrePayID:  prePayID,
		Package:   "Sign=WXPay",
		NonceStr:  util.RandomStr(32),
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	sign, err := pcf.Sign(&payConf, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return payConf
	}
	payConf.Sign = sign
	return payConf
}

func (pcf *Pay) Sign(variable interface{}, key string) (sign string, err error) {
	ss := &SignStruct{
		ToLower: false,
		Tag:     "xml",
	}
	sign, err = ss.Sign(variable, nil, key)
	return
}

func (pcf *Pay) SignByJson(variable interface{}, key string) (sign string, err error) {
	ss := &SignStruct{
		ToLower: false,
		Tag:     "json",
	}
	sign, err = ss.Sign(variable, nil, key)
	return
}

type MchTransfersParams struct {
	InputCharset   string `xml:"input_charset"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	Memo           string `xml:"memo"`
	AppID          string `xml:"appid"`
	OpenID         string `xml:"openid"`
	OpUserID       string `xml:"op_user_id"`
	OpUserPasswd   string `xml:"op_user_passwd"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
}

func (pcf *Pay) MchPay(p *Params) error {
	nonceStr := util.RandomStr(32)
	params := &MchTransfersParams{
		InputCharset:   "UTF-8",
		AppID:          pcf.AppID,
		OpenID:         p.OpenID,
		MchID:          pcf.PayMchID,
		OutTradeNo:     p.OutTradeNo,
		NonceStr:       nonceStr,
		TotalFee:       p.TotalFee,
		Memo:           p.Body,
		SpbillCreateIp: p.CreateIP,
		OpUserID:       pcf.PayOpUserID,
		OpUserPasswd:   pcf.PayOpUserPwd,
	}
	sign, err := pcf.Sign(params, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	params.Sign = sign
	client, err := util.NewTLSHttpClient([]byte(pcf.PayCertPEMBlock), []byte(pcf.PayKeyPEMBlock))
	if err != nil {
		return err
	}
	rawRet, err := util.PostXML(mchTransUri, params, "MchTransfersParams", client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	payRet := payResult{}
	err = xml.Unmarshal(rawRet, &payRet)
	if err != nil {
		fmt.Println("xmlUnmarshalError,res:" + string(rawRet))
		return err
	}
	if payRet.ReturnCode == "SUCCESS" {
		if payRet.ResultCode == "SUCCESS" {
			return nil
		}
		return errors.New(payRet.ErrCodeDes)
	}
	return errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "]")
}

type RedParams struct {
	Charset     int    `json:"charset" xml:"charset"`           //必填 1 utf8 , 2 gbk
	NonceStr    string `json:"nonce_str" xml:"nonce_str"`       //必填
	Sign        string `json:"sign" xml:"sign"`                 //必填
	MchBillno   string `json:"mch_billno" xml:"mch_billno"`     //必填 订单号
	MchID       string `json:"mch_id" xml:"mch_id"`             //必填
	MchName     string `json:"mch_name" xml:"mch_name"`         //必填//商户名称，会展示在红包领取页面上
	QqAppID     string `json:"qqappid" xml:"qqappid"`           //必填
	ReOpenID    string `json:"re_openid" xml:"re_openid"`       //必填
	TotalAmount int    `json:"total_amount" xml:"total_amount"` //必填
	TotalNum    int    `json:"total_num" xml:"total_num"`       //必填
	Wishing     string `json:"wishing" xml:"wishing"`           //必填
	ActName     string `json:"act_name" xml:"act_name"`         //必填
	IconID      int    `json:"icon_id" xml:"icon_id"`           //必填
	BannerID    int    `json:"banner_id,omitempty" xml:"banner_id,omitempty"`
	NotifyUrl   string `json:"notify_url,omitempty" xml:"notify_url,omitempty"`
	NotSendMsg  int    `json:"not_send_msg,omitempty" xml:"not_send_msg,omitempty"`
	MinValue    int    `json:"min_value" xml:"min_value"` //必填 1
	MaxValue    int    `json:"max_value" xml:"max_value"` //必填 100
}

type RedResult struct {
	Retcode string `json:"retcode"`
	Retmsg  string `json:"retmsg"`
	Listid  string `json:"listid"`
}

func (pcf *Pay) SendRed(p *Params) (string, error) {
	nonceStr := util.RandomStr(32)
	params := &RedParams{
		Charset:     1,
		NonceStr:    nonceStr,
		MchBillno:   p.OutTradeNo,
		MchID:       pcf.PayMchID,
		QqAppID:     pcf.AppID,
		MchName:     p.SendName,
		ReOpenID:    p.OpenID,
		TotalAmount: p.TotalFee,
		TotalNum:    1,
		Wishing:     p.Wishing,
		ActName:     p.ActName,
		IconID:      p.IconID,
		BannerID:    p.BannerID,
		MinValue:    1,
		MaxValue:    p.TotalFee + 1,
		NotifyUrl:   p.NotifyUrl,
	}
	sign, err := pcf.SignByJson(params, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	params.Sign = sign
	client, err := util.NewTLSHttpClient([]byte(pcf.PayCertPEMBlock), []byte(pcf.PayKeyPEMBlock))
	if err != nil {
		return "", err
	}

	j, _ := json.Marshal(params)
	var paramsMap map[string]interface{}
	err = json.Unmarshal(j, &paramsMap)
	if err != nil {
		return "", err
	}
	urlParams := url.Values{}
	for k, v := range paramsMap {
		urlParams.Add(k, fmt.Sprint(v))
	}
	link := sendRedUri + "?"
	link = link + urlParams.Encode()
	rawRet, err := util.HTTPGetV2(link, client)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	payRet := RedResult{}
	err = json.Unmarshal(rawRet, &payRet)
	if err != nil {
		fmt.Println("jsonUnmarshalError,res:" + string(rawRet))
		return "", err
	}
	if payRet.Retcode == "0" {
		return payRet.Listid, nil
	}
	return "", errors.New(payRet.Retmsg)
}
