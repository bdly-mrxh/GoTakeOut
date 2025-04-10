package dto

// AddressBookDTO 地址簿传输数据模型
type AddressBookDTO struct {
	ID           int    `json:"id"`
	UserID       int    `json:"userId"`
	Consignee    string `json:"consignee"`
	Phone        string `json:"phone"`
	Sex          string `json:"sex"`
	ProvinceCode string `json:"provinceCode"`
	ProvinceName string `json:"provinceName"`
	CityCode     string `json:"cityCode"`
	CityName     string `json:"cityName"`
	DistrictCode string `json:"districtCode"`
	DistrictName string `json:"districtName"`
	Detail       string `json:"detail"`
	Label        int    `json:"label"`
	IsDefault    int    `json:"isDefault"`
}
