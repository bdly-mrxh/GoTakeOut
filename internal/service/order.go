package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/internal/websocket"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"
	"takeout/model/wrap"
	"time"
)

type OrderService struct {
	orderDAO        dao.OrderDAO
	addressBookDAO  dao.AddressBookDAO
	shoppingCartDAO dao.ShoppingCartDAO
	orderDetailDAO  dao.OrderDetailDAO
	userDAO         dao.UserDAO
}

// Submit 提交订单
func (s *OrderService) Submit(ctx *gin.Context, submitDTO *dto.OrderSubmitDTO) (*vo.OrderSubmitVO, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	var submitVO vo.OrderSubmitVO
	err = global.DB.Transaction(func(db *gorm.DB) error {
		address, e := s.addressBookDAO.GetByID(db, submitDTO.AddressBookID)
		if e != nil {
			if errors.Is(e, gorm.ErrRecordNotFound) {
				return errs.Wrap(e, constant.CodeBusinessError, constant.MsgAddressBookIsNull)
			}
			return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}

		// 判断是否可以配送
		//e = utils.CheckOutOfRange(address.CityName + address.DistrictName + address.Detail)
		//if e != nil {
		//	return errs.Wrap(e, constant.CodeBusinessError, "无法配送")
		//}

		cartList, e := s.shoppingCartDAO.List(db, &entity.ShoppingCart{UserID: userID})
		if e != nil {
			return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		if len(cartList) == 0 {
			return errs.New(constant.CodeBusinessError, constant.MsgShoppingCartIsNull)
		}

		// 插入一条订单数据
		order := &entity.Order{}
		e = utils.CopyProperties(submitDTO, order)
		if e != nil {
			return errs.Wrap(e, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
		}
		order.Number = strconv.FormatInt(time.Now().UnixMilli(), 10)
		order.Status = constant.PendingPayment
		order.PayStatus = constant.UnPaid
		order.Phone = address.Phone
		order.Consignee = address.Consignee
		order.UserID = userID
		order.Address = address.Detail
		e = s.orderDAO.Insert(db, order)
		if e != nil {
			return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}

		// 向订单明细表插入数据
		var orderDetailList []*entity.OrderDetail
		// index, value := range ☆
		for _, cart := range cartList {
			orderDetail := &entity.OrderDetail{}
			e = utils.CopyProperties(cart, orderDetail)
			if e != nil {
				return errs.Wrap(e, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
			}
			orderDetail.OrderID = order.ID
			orderDetailList = append(orderDetailList, orderDetail)
		}
		e = s.orderDetailDAO.BatchInsert(db, orderDetailList)
		if e != nil {
			return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}

		// 清空购物车数据
		e = s.shoppingCartDAO.CleanByUserID(db, userID)
		if e != nil {
			return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}

		submitVO.ID = order.ID
		submitVO.OrderTime = order.OrderTime
		submitVO.OrderAmount = order.Amount
		submitVO.OrderNumber = order.Number
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &submitVO, nil
}

// RealPayment 实际微信支付
func (s *OrderService) RealPayment(ctx *gin.Context, payDTO *dto.OrderPaymentDTO) (*vo.OrderPaymentVO, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.userDAO.GetByID(global.DB, userID)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	jsonObject, err := utils.Pay(payDTO.OrderNumber, decimal.NewFromFloat(0.01), "外卖订单", user.OpenID)
	if err != nil { // 检查错误
		return nil, errs.Wrap(err, constant.CodeInternalError, "支付失败")
	}
	// 使用map解析
	var jsonMap map[string]any
	err = json.Unmarshal(jsonObject, &jsonMap)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}
	if code, ok := jsonMap["code"]; ok {
		if code.(string) == "ORDERPAID" {
			return nil, errs.New(constant.CodeBusinessError, constant.MsgOrderPaid)
		}
	}
	// 反序列化为VO
	var payVO vo.OrderPaymentVO
	if err = json.Unmarshal(jsonObject, &payVO); err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}
	payVO.PackageStr = jsonMap["package"].(string)
	return &payVO, nil
}

// PaySuccess 支付成功后修改订单状态
func (s *OrderService) PaySuccess(no string) error { // no是订单号
	order, err := s.orderDAO.GetByNumber(global.DB, no)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	order.Status = constant.ToBeConfirmed
	order.PayStatus = constant.Paid
	order.CheckoutTime = wrap.LocalTime(time.Now())
	err = s.orderDAO.Update(global.DB, order)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	// 通知商户
	m := map[string]any{"type": constant.NotifyOrder, "orderId": order.ID, "content": "订单号：" + order.Number}
	websocket.SendToAllClients(m)
	return nil
}

// Payment 绕过微信支付
func (s *OrderService) Payment(dto *dto.OrderPaymentDTO) (*vo.OrderPaymentVO, error) {
	//userID, err := utils.GetId(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//user, err := s.userDAO.GetByID(global.DB, userID)
	//if err != nil {
	//	return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	//}

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := utils.GenerateRandomNumericString(32)
	jo := map[string]any{
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + strconv.Itoa(123),
		"signType":  "RSA",
		"paySign":   "success",
	}
	// 转换为 JSON
	jsonObject, err := json.Marshal(jo)
	if err != nil {
		return nil, err
	}
	var jsonMap map[string]any
	err = json.Unmarshal(jsonObject, &jsonMap)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}
	if code, ok := jsonMap["code"]; ok {
		if code.(string) == "ORDERPAID" {
			return nil, errs.New(constant.CodeBusinessError, constant.MsgOrderPaid)
		}
	}
	// 反序列化为VO
	var payVO vo.OrderPaymentVO
	if err = json.Unmarshal(jsonObject, &payVO); err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}
	payVO.PackageStr = jsonMap["package"].(string)

	// 直接更改支付状态
	err = s.PaySuccess(dto.OrderNumber)
	if err != nil {
		return nil, err
	}

	return &payVO, nil
}

// Page 分页查询
func (s *OrderService) Page(ctx *gin.Context, queryDTO *dto.OrderPageQueryDTO) (*vo.PageResult, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	queryDTO.UserID = userID
	total, list, err := s.orderDAO.Page(global.DB, queryDTO)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	var orderVOs []*vo.OrderVO
	if len(list) > 0 {
		for _, order := range list {
			// 查询订单明细
			orderDetails, e := s.orderDetailDAO.GetByOrderID(global.DB, order.ID)
			if e != nil {
				return nil, errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
			}
			orderVO := &vo.OrderVO{}
			if e = utils.CopyProperties(order, orderVO); e != nil {
				return nil, errs.Wrap(e, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
			}
			orderVO.OrderDetailList = orderDetails
			orderVOs = append(orderVOs, orderVO)
		}
	}
	return &vo.PageResult{Total: total, Records: orderVOs}, nil
}

// Detail 查询订单详细信息
func (s *OrderService) Detail(id int) (*vo.OrderVO, error) {
	// 查询订单
	order, err := s.orderDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgNotFound)
		}
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	// 查询订单详细
	orderDetail, err := s.orderDetailDAO.GetByOrderID(global.DB, id)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	orderVO := &vo.OrderVO{}
	if err = utils.CopyProperties(order, orderVO); err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	orderVO.OrderDetailList = orderDetail
	return orderVO, nil
}

// CancelByUser 用户取消订单
func (s *OrderService) CancelByUser(id int) error {
	order, err := s.orderDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 业务错误：未查询到订单
			return errs.Wrap(err, constant.CodeBusinessError, constant.MsgOrderNotFound)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if order.Status > constant.ToBeConfirmed {
		return errs.New(constant.CodeBusinessError, constant.MsgOrderStatusError)
	}
	o := &entity.Order{ID: id}
	if order.Status == constant.ToBeConfirmed {
		// 已付款的订单需要退款
		// _ = utils.Refund(order.Number, order.Number, decimal.NewFromFloat(0.01), decimal.NewFromFloat(0.01))
		o.Status = constant.ReFund
	}
	o.Status = constant.Cancelled
	o.CancelReason = "用户取消"
	o.CancelTime = wrap.LocalTime(time.Now())
	if err = s.orderDAO.Update(global.DB, o); err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Repetition 再来一单
func (s *OrderService) Repetition(ctx *gin.Context, id int) error {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return err
	}
	orderDetail, err := s.orderDetailDAO.GetByOrderID(global.DB, id)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	var cartList []*entity.ShoppingCart
	if len(orderDetail) > 0 {
		for _, od := range orderDetail {
			cart := &entity.ShoppingCart{}
			if err = utils.CopyProperties(od, cart); err != nil {
				return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
			}
			cart.ID = 0
			cart.UserID = userID
			cartList = append(cartList, cart)
		}
		err = s.shoppingCartDAO.BatchInsert(global.DB, cartList)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	}
	return nil
}

// Search 条件搜索订单
func (s *OrderService) Search(queryDTO *dto.OrderPageQueryDTO) (*vo.PageResult, error) {
	total, orders, err := s.orderDAO.Page(global.DB, queryDTO)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	// make([]*vo.OrderVO, 0)切片不为 nil，是一个空切片; var orderVOs []*vo.OrderVO 切片值为 nil; 这里前端不能是null，必须是[]
	orderVOs := make([]*vo.OrderVO, 0)
	for _, order := range orders {
		orderVO := &vo.OrderVO{}
		if err = utils.CopyProperties(order, orderVO); err != nil {
			return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
		}
		// 获取orderDishes字符串
		orderDetail, e := s.orderDetailDAO.GetByOrderID(global.DB, order.ID)
		if e != nil {
			return nil, errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		var strs []string
		for _, od := range orderDetail {
			str := od.Name + "*" + strconv.Itoa(od.Number) + ";"
			strs = append(strs, str)
		}
		orderDishes := strings.Join(strs, "")
		orderVO.OrderDishes = orderDishes
		orderVOs = append(orderVOs, orderVO)
	}
	return &vo.PageResult{Total: total, Records: orderVOs}, nil
}

// Statistics 统计各个订单状态
func (s *OrderService) Statistics() (*vo.OrderStatisticsVO, error) {
	toBeConfirmed, err := s.orderDAO.CountStatus(global.DB, constant.ToBeConfirmed)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	confirmed, err := s.orderDAO.CountStatus(global.DB, constant.Confirmed)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	deliveryInProgress, err := s.orderDAO.CountStatus(global.DB, constant.DeliveryInProgress)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return &vo.OrderStatisticsVO{ToBeConfirmed: toBeConfirmed, Confirmed: confirmed, DeliveryInProgress: deliveryInProgress}, nil
}

// Confirm 商家接单
func (s *OrderService) Confirm(dto *dto.OrderConfirmDTO) error {
	order := &entity.Order{
		ID:     dto.OrderID,
		Status: constant.Confirmed,
	}
	err := s.orderDAO.Update(global.DB, order)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Reject 商家拒单
func (s *OrderService) Reject(dto *dto.OrderRejectionDTO) error {
	order, err := s.orderDAO.GetByID(global.DB, dto.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, constant.CodeBusinessError, constant.MsgOrderStatusError)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if order.PayStatus == constant.Paid {
		// 退款
		// _ = utils.Refund(order.Number, order.Number, decimal.NewFromFloat(0.01), decimal.NewFromFloat(0.01))
	}
	o := &entity.Order{
		ID:              order.ID,
		Status:          constant.Cancelled,
		RejectionReason: dto.RejectionReason,
		CancelTime:      wrap.LocalTime(time.Now()),
	}
	err = s.orderDAO.Update(global.DB, o)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Cancel 商家取消订单
func (s *OrderService) Cancel(dto *dto.OrderCancelDTO) error {
	order, err := s.orderDAO.GetByID(global.DB, dto.OrderID)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if order.PayStatus == constant.Paid {
		// 退款
		// _ = utils.Refund(order.Number, order.Number, decimal.NewFromFloat(0.01), decimal.NewFromFloat(0.01))
	}

	o := &entity.Order{
		ID:           order.ID,
		Status:       constant.Cancelled,
		CancelReason: dto.CancelReason,
		CancelTime:   wrap.LocalTime(time.Now()),
	}
	err = s.orderDAO.Update(global.DB, o)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Delivery 派送订单
func (s *OrderService) Delivery(id int) error {
	order, err := s.orderDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, constant.CodeBusinessError, constant.MsgOrderStatusError)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	// logger.Infof("%v", *order)
	// 必须是接单的订单才能派送
	if order.Status != constant.Confirmed {
		return errs.New(constant.CodeBusinessError, constant.MsgOrderStatusError)
	}
	o := &entity.Order{
		ID:     order.ID,
		Status: constant.DeliveryInProgress,
	}
	err = s.orderDAO.Update(global.DB, o)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Complete 完成订单
func (s *OrderService) Complete(id int) error {
	order, err := s.orderDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, constant.CodeBusinessError, constant.MsgOrderStatusError)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if order.Status != constant.DeliveryInProgress {
		return errs.New(constant.CodeBusinessError, constant.MsgOrderStatusError)
	}
	o := &entity.Order{
		ID:           order.ID,
		Status:       constant.Completed,
		DeliveryTime: wrap.LocalTime(time.Now()),
	}
	err = s.orderDAO.Update(global.DB, o)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Reminder 用户催单
func (s *OrderService) Reminder(id int) error {
	order, err := s.orderDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(constant.CodeBusinessError, constant.MsgOrderNotFound)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	m := map[string]any{"type": constant.UserRemind, "orderId": id, "content": "订单号：" + order.Number}
	//js, err := json.Marshal(m)
	//if err != nil {
	//	return errs.Wrap(err, constant.CodeInternalError, constant.MsgMarshalFail)
	//}
	websocket.SendToAllClients(m)
	return nil
}
