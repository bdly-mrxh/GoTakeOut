package dto

// CategoryDTO 分类创建请求DTO
type CategoryDTO struct {
	ID   int    `json:"id"`                      // 主键
	Type string `json:"type"`                    // 类型 1 菜品分类 2 套餐分类
	Name string `json:"name" binding:"required"` // 分类名称
	Sort string `json:"sort" binding:"required"` // 排序
}

// CategoryPageDTO 分类分页查询DTO
type CategoryPageDTO struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`     // 页码
	PageSize int    `form:"pageSize"` // 每页记录数
	Type     int    `form:"type"`
}
