package entity

// AddressBook 地址簿数据模型
type AddressBook struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int    `json:"userId" gorm:"column:user_id;not null"`
	Consignee    string `json:"consignee"`
	Sex          string `json:"sex"`
	Phone        string `json:"phone" gorm:"not null"`
	ProvinceCode string `json:"provinceCode" gorm:"column:province_code"`
	ProvinceName string `json:"provinceName" gorm:"column:province_name"`
	CityCode     string `json:"cityCode" gorm:"column:city_code"`
	CityName     string `json:"cityName" gorm:"column:city_name"`
	DistrictCode string `json:"districtCode" gorm:"column:district_code"`
	DistrictName string `json:"districtName" gorm:"column:district_name"`
	Detail       string `json:"detail"`
	Label        string `json:"label"`
	IsDefault    int    `json:"isDefault" gorm:"column:is_default;default:0"`
}

// TableName 设置表名
func (AddressBook) TableName() string {
	return "address_book"
}
