package notify

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"go.uber.org/zap"
	"net/http"
	"takeout/common/constant"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/common/utils"
	"takeout/internal/service"
)

type NotifyController struct {
	orderService service.OrderService
}

func NewNotifyController() *NotifyController {
	return &NotifyController{}
}

// PaySuccess 解析微信回调请求的参数到 V3NotifyReq 结构体
func (c *NotifyController) PaySuccess(ctx *gin.Context) {
	notifyReq, err := wechat.V3ParseNotify(ctx.Request)
	if err != nil {
		logger.Error("V3ParseNotify ERR", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "回调内容异常"})
		return
	}
	// 获取微信平台证书
	wxClient := utils.GetWechatClientV3()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPK(wxClient.WxPublicKey())
	if err != nil {
		logger.Error("VerifySignByPKMap ERR", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容验证失败"})
		return
	}
	// 普通支付通知解密
	result, rErr := notifyReq.DecryptPayCipherText(global.Config.Wechat.ApiV3Key)
	if rErr != nil {
		logger.Error("DecryptPayCipherText Error", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
		return
	}
	if result != nil && result.TradeState == "SUCCESS" {
		logger.Info("商户平台订单号", zap.String("out_trade_no", result.OutTradeNo))
		logger.Info("微信支付交易号", zap.String("transaction_id", result.TransactionId))

		err = c.orderService.PaySuccess(result.OutTradeNo)
		logger.Error(constant.Update, zap.Error(err))
	}

	// 此写法是 gin 框架返回微信的写法
	ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: gopay.SUCCESS})
}

// RefundSuccess 微信支付退款回调
func (c *NotifyController) RefundSuccess(ctx *gin.Context) {
	notifyReq, err := wechat.V3ParseNotify(ctx.Request)
	if err != nil {
		logger.Error("V3ParseNotify ERR", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "回调内容异常"})
		return
	}
	// 获取微信平台证书
	wxClient := utils.GetWechatClientV3()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPK(wxClient.WxPublicKey())
	if err != nil {
		logger.Error("VerifySignByPKMap ERR", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容验证失败"})
		return
	}
	// 普通支付通知解密
	result, rErr := notifyReq.DecryptRefundCipherText(global.Config.Wechat.ApiV3Key)
	if rErr != nil {
		logger.Error("DecryptPayCipherText Error", zap.Error(err))
		ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
		return
	}
	if result != nil && result.RefundStatus == "SUCCESS" {
		logger.Info("退款成功", zap.String("退单号", result.OutRefundNo))
	}

	ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: gopay.SUCCESS})
}
