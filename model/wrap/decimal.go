package wrap

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/shopspring/decimal"
)

// 可以设置 MarshalJSONWithoutQuotes = true 实现，传给前端一定要是 number 类型

// Decimal 用于自定义序列化 （未使用）
type Decimal decimal.Decimal

// MarshalJSON 自定义 decimal.Decimal 的 JSON 序列化为 float64
// Jackson BigDecimal 序列化是
func (d Decimal) MarshalJSON() ([]byte, error) {
	// 这里可以根据需求返回浮动精度的数字
	return json.Marshal(decimal.Decimal(d).InexactFloat64())
}

// UnmarshalJSON 自定义反序列化
func (d *Decimal) UnmarshalJSON(b []byte) error {
	return (*decimal.Decimal)(d).UnmarshalJSON(b)
}

// Scan 直接调用 decimal.Decimal 的 Scan 方法
func (d *Decimal) Scan(value any) error {
	return (*decimal.Decimal)(d).Scan(value)
}

// Value 直接调用 decimal.Decimal 的 Value 方法
func (d Decimal) Value() (driver.Value, error) {
	return decimal.Decimal(d).Value()
}

// String 直接调用 decimal.Decimal 的 String 方法
func (d Decimal) String() string {
	return decimal.Decimal(d).String()
}
