package dao

import (
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/model/dto"
	"takeout/model/entity"
	"time"
)

type OrderDAO struct{}

// Insert 新增订单
func (d *OrderDAO) Insert(db *gorm.DB, order *entity.Order) error {
	return db.Model(&entity.Order{}).Create(order).Error
}

// GetByNumber 根据订单号获取订单信息
func (d *OrderDAO) GetByNumber(db *gorm.DB, no string) (*entity.Order, error) {
	var order entity.Order
	result := db.Model(&entity.Order{}).Where("number = ?", no).First(&order)
	return &order, result.Error
}

// Update 更新订单状态
func (d *OrderDAO) Update(db *gorm.DB, order *entity.Order) error {
	return db.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(order).Error
}

// Page 分页查询
func (d *OrderDAO) Page(db *gorm.DB, queryDTO *dto.OrderPageQueryDTO) (int64, []*entity.Order, error) {
	var (
		err   error
		total int64
		list  []*entity.Order
	)
	query := db.Model(&entity.Order{})
	if queryDTO.Number != "" {
		query = query.Where("number like ?", "%"+queryDTO.Number+"%")
	}
	if queryDTO.Phone != "" {
		query = query.Where("phone like ?", "%"+queryDTO.Phone+"%")
	}
	if queryDTO.UserID != 0 {
		query = query.Where("user_id = ?", queryDTO.UserID)
	}
	if queryDTO.Status != 0 {
		query = query.Where("status = ?", queryDTO.Status)
	}
	if queryDTO.BeginTime != "" {
		query = query.Where("order_time >= ?", queryDTO.BeginTime)
	}
	if queryDTO.EndTime != "" {
		query = query.Where("order_time <= ?", queryDTO.EndTime)
	}

	if err = query.Order("order_time desc").Error; err != nil {
		return 0, nil, err
	}
	if err = query.Count(&total).Limit(queryDTO.PageSize).
		Offset((queryDTO.Page - 1) * queryDTO.PageSize).
		Find(&list).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil, nil
		}
		return 0, nil, err
	}
	return total, list, nil
}

// GetByID 根据ID查询订单
func (d *OrderDAO) GetByID(db *gorm.DB, id int) (*entity.Order, error) {
	var order entity.Order
	result := db.Model(&entity.Order{}).Where("id = ?", id).First(&order)
	return &order, result.Error
}

// CountStatus 统计某种状态的订单数量
func (d *OrderDAO) CountStatus(db *gorm.DB, status int) (int64, error) {
	var cnt int64
	res := db.Model(&entity.Order{}).Where("status = ?", status).Count(&cnt)
	return cnt, res.Error
}

// GetByStatusLT 查询在某个时间之前某种状态的订单列表
func (d *OrderDAO) GetByStatusLT(db *gorm.DB, status int, t time.Time) ([]*entity.Order, error) {
	var list []*entity.Order
	result := db.Model(&entity.Order{}).Where("status = ? and order_time < ?", status, t).Find(&list)
	return list, result.Error
}

// GetAmount 统计营业额
func (d *OrderDAO) GetAmount(db *gorm.DB, begin *time.Time, end *time.Time) (decimal.Decimal, error) {
	var amount decimal.Decimal
	res := db.Table("orders").Select("ifnull(sum(amount), 0)").Where("order_time >= ? and order_time <= ? and status = ?", begin, end, constant.Completed).Scan(&amount)
	return amount, res.Error
}

// GetCount 统计订单数量
func (d *OrderDAO) GetCount(db *gorm.DB, begin *time.Time, end *time.Time, status int) (int64, error) {
	var cnt int64
	query := db.Model(&entity.Order{})
	if begin != nil {
		query = query.Where("order_time >= ?", begin)
	}
	if end != nil {
		query = query.Where("order_time <= ?", end)
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&cnt).Error
	return cnt, err
}

// GetSalesTop10 查询销量前十
func (d *OrderDAO) GetSalesTop10(db *gorm.DB, begin *time.Time, end *time.Time) ([]dto.GoodsSalesDTO, error) {
	var list []dto.GoodsSalesDTO
	err := db.Table("order_detail od").
		Joins("left join orders o on o.id = od.order_id").
		Select("od.name, sum(od.number) as number").
		Where("o.status = ? and o.order_time >= ? and o.order_time <= ?", constant.Completed, begin, end). // 如果参数是指针类型，Gorm 会自动解引用以获取其值
		Group("od.name").Order("number desc").Offset(0).Limit(10).Find(&list).Error
	return list, err
}
