package wrap

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const TimeFormat = "2006-01-02 15:04"

// LocalTime 自定义时间类型，用于处理时间格式化
type LocalTime time.Time

// MarshalJSON 实现json.Marshaler接口
func (t LocalTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

func ParseFlexibleTime(str string) (time.Time, error) {
	// 设置时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	var t time.Time
	for _, f := range formats {
		t, err = time.ParseInLocation(f, str, loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间: %s", str)
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	pt, err := ParseFlexibleTime(str)
	if err != nil {
		return err
	}

	*t = LocalTime(pt)
	return nil
}

// Value 实现driver.Valuer接口
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tm := time.Time(t)
	if tm.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tm, nil
}

// Scan 实现sql.Scanner接口
func (t *LocalTime) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// String 实现fmt.Stringer接口
func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

// Time 转换为标准time.Time
func (t LocalTime) Time() time.Time {
	return time.Time(t)
}

// UnmarshalText 解析文本
func (t *LocalTime) UnmarshalText(text []byte) error {
	// 解码URL编码的字符串
	str, err := url.QueryUnescape(string(text))
	if err != nil {
		return fmt.Errorf("URL解码失败: %v", err)
	}

	if str == "" {
		*t = LocalTime(time.Time{})
		return nil
	}

	// 支持的时间格式
	formats := []string{
		"2006-01-02 15:04:05", // YYYY-MM-DD HH:MM:SS
		"2006-01-02T15:04:05", // YYYY-MM-DDThh:mm:ss
		time.RFC3339,          // YYYY-MM-DDThh:mm:ssZ 或 YYYY-MM-DDThh:mm:ss+hh:mm
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}

	var parseErr error
	for _, format := range formats {
		pt, e := time.ParseInLocation(format, str, loc)
		if e == nil {
			*t = LocalTime(pt)
			return nil
		}
		parseErr = e
	}

	return fmt.Errorf("无法解析时间 '%s': %v", str, parseErr)
}

// MarshalText 格式化为文本
func (t LocalTime) MarshalText() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte(""), nil
	}
	return []byte(tt.Format(TimeFormat)), nil
}

// IsZero 判断是否为零值
func (t LocalTime) IsZero() bool {
	return time.Time(t).IsZero()
}
