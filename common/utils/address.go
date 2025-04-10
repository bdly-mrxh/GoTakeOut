package utils

import (
	"encoding/json"
	"fmt"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
)

const baiduUrl = "https://api.map.baidu.com/geocoding/v3"

// CheckOutOfRange 检查是否超出配送范围
func CheckOutOfRange(address string) error {
	m := map[string]string{
		"address": global.Config.Shop.Address,
		"output":  "json",
		"ak":      global.Config.Baidu.AK,
	}

	// 获取店铺的经纬度
	shopCoordinate, err := DoGET(baiduUrl, m)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgServerError)
	}

	// 解析JSON响应
	var shopResult map[string]any
	if err = json.Unmarshal([]byte(shopCoordinate), &shopResult); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}

	// 检查状态码
	if shopResult["status"].(string) != "0" {
		return errs.New(constant.CodeInternalError, "店铺地址解析错误")
	}

	// 解析店铺位置数据
	result := shopResult["result"].(map[string]any)
	location := result["location"].(map[string]any)
	lat := fmt.Sprintf("%f", location["lat"].(float64))
	lng := fmt.Sprintf("%f", location["lng"].(float64))
	shopLngLat := lat + "," + lng

	// 更新参数获取用户地址经纬度
	m["address"] = address
	userCoordinate, err := DoGET(baiduUrl, m)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgServerError)
	}

	// 解析用户地址JSON响应
	var userResult map[string]any
	if err = json.Unmarshal([]byte(userCoordinate), &userResult); err != nil {
		return err
	}

	// 检查状态码
	if userResult["status"].(string) != "0" {
		return errs.New(constant.CodeInternalError, "收货地址解析失败")
	}

	// 解析用户位置数据
	userResultData := userResult["result"].(map[string]any)
	userLocation := userResultData["location"].(map[string]any)
	userLat := fmt.Sprintf("%f", userLocation["lat"].(float64))
	userLng := fmt.Sprintf("%f", userLocation["lng"].(float64))
	userLngLat := userLat + "," + userLng

	// 路线规划参数
	m["origin"] = shopLngLat
	m["destination"] = userLngLat
	m["steps_info"] = "0"

	// 获取路线规划
	routeJSON, err := DoGET("https://api.map.baidu.com/directionlite/v1/driving", m)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgServerError)
	}

	// 解析路线规划JSON
	var routeResult map[string]any
	if err = json.Unmarshal([]byte(routeJSON), &routeResult); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
	}

	// 检查状态码
	if routeResult["status"].(string) != "0" {
		return errs.New(constant.CodeInternalError, "配送路线规划失败")
	}

	// 解析距离数据
	routeData := routeResult["result"].(map[string]any)
	routes := routeData["routes"].([]any)
	route := routes[0].(map[string]any)
	distance := int(route["distance"].(float64))

	// 检查是否超出配送范围
	if distance > 5000 {
		return errs.New(constant.CodeBusinessError, "超出配送范围")
	}

	return nil
}
