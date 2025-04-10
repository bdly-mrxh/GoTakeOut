package utils

import (
	"context"
	"crypto"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"time"
)

// 创建微信支付Http客户端
func getClient() (*core.Client, error) {
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(global.Config.Wechat.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}
	wechatPayCertificate, err := utils.LoadCertificateWithPath(global.Config.Wechat.WeChatPayCertFilePath)
	if err != nil {
		return nil, err
	}
	return core.NewClient(
		context.Background(),
		// 一次性设置 签名/验签/敏感字段加解密，并注册 平台证书下载器，自动定时获取最新的平台证书
		option.WithMerchantCredential(global.Config.Wechat.MchID, global.Config.Wechat.MchSerialNumber, mchPrivateKey),
		option.WithWechatPayCertificate([]*x509.Certificate{wechatPayCertificate}),
		//option.WithWechatPayCipher(
		//	encryptors.NewWechatPayEncryptor(downloader.MgrInstance().GetCertificateVisitor(global.Config.Wechat.Mchid)),
		//	decryptors.NewWechatPayDecryptor(mchPrivateKey),
		//),
	)
}

// 发送jsapi请求
func jsapiPost(orderNum, description, openid string, total decimal.Decimal) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}
	svc := jsapi.JsapiApiService{Client: client}
	resp, _, err := svc.Prepay(context.Background(),
		jsapi.PrepayRequest{
			Appid:       core.String(global.Config.Wechat.AppID),
			Mchid:       core.String(global.Config.Wechat.MchID),
			Description: core.String(description),
			OutTradeNo:  core.String(orderNum),
			NotifyUrl:   core.String(global.Config.Wechat.NotifyUrl),
			Amount: &jsapi.Amount{
				Total:    core.Int64(total.Mul(decimal.NewFromInt(100)).Round(2).IntPart()),
				Currency: core.String("CNY"),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openid),
			},
		})
	if err != nil {
		return "", err
	}
	if resp.PrepayId == nil {
		return "", errs.New(constant.CodeInternalError, constant.MsgNotFound)
	}
	return *resp.PrepayId, nil
}

// Pay 获取返回给客户端的支付信息
func Pay(orderNum string, total decimal.Decimal, description, openid string) ([]byte, error) {
	prePayID, err := jsapiPost(orderNum, description, openid, total)
	if err != nil {
		return nil, err
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := GenerateRandomNumericString(32)
	list := []any{global.Config.Wechat.AppID, timeStamp, nonceStr, "prepay_id=" + prePayID}
	var stringBuilder strings.Builder
	for _, item := range list {
		stringBuilder.WriteString(fmt.Sprintf("%v\n", item))
	}
	msg := stringBuilder.String()

	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(global.Config.Wechat.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}
	packageSign, err := utils.SignSHA256WithRSA(msg, mchPrivateKey)
	if err != nil {
		return nil, err
	}
	// 构造数据给微信小程序
	jo := map[string]any{
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + prePayID,
		"signType":  "RSA",
		"paySign":   packageSign,
	}
	// 转换为 JSON
	joJSON, err := json.Marshal(jo)
	if err != nil {
		return nil, err
	}
	return joJSON, nil
}

func GenerateRandomNumericString(length int) string {
	var randomString string

	for i := 0; i < length; i++ {
		randomDigit := rand.Intn(10) // 生成一个 0 到 9 的随机数字
		randomString += strconv.Itoa(randomDigit)
	}

	return randomString
}

// 签名
func signMessage(message string, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(message))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(crand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	// Base64 编码签名
	return base64.StdEncoding.EncodeToString(signature), nil
}

var wechatClientV3 *wechat.ClientV3

var wcOnce sync.Once

// GetWechatClientV3 wechatClientV3
func GetWechatClientV3() *wechat.ClientV3 {
	var err error
	wcOnce.Do(func() {
		wechatClientV3, err = wechat.NewClientV3(global.Config.Wechat.MchID, global.Config.Wechat.MchSerialNumber, global.Config.Wechat.ApiV3Key, ReadPem())
		if err != nil {
			log.Panic(err.Error())
		}
		// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
		err = wechatClientV3.AutoVerifySign()
		if err != nil {
			log.Panic(err.Error())
		}
	})
	return wechatClientV3
}

func ReadPem() string {
	privateKey, err := os.ReadFile(global.Config.Wechat.PrivateKeyFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(privateKey)
}
