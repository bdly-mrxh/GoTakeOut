package utils

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"takeout/common/global"
)

// Refund 微信支付退款
func Refund(outTradeNo, outRefundNo string, refund decimal.Decimal, total decimal.Decimal) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	svc := refunddomestic.RefundsApiService{Client: client}
	_, _, err = svc.Create(context.Background(),
		refunddomestic.CreateRequest{
			OutTradeNo:  core.String(outTradeNo),
			OutRefundNo: core.String(outRefundNo),
			NotifyUrl:   core.String(global.Config.Wechat.RefundNotifyUrl),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(refund.Mul(decimal.NewFromInt(100)).Round(2).IntPart()),
				Total:    core.Int64(total.Mul(decimal.NewFromInt(100)).Round(2).IntPart()),
			},
		},
	)
	return err
}
