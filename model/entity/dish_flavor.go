package entity

// DishFlavor 菜品口味关系数据模型
type DishFlavor struct {
	ID     int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	DishID int    `json:"dishId" gorm:"column:dish_id;not null"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// TableName 设置表名
func (DishFlavor) TableName() string {
	return "dish_flavor"
}
