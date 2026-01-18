package model

import (
	"time"

	"github.com/bramble555/blog/global"
	"gorm.io/gorm"
)

type MODEL struct {
	ID         int64     `gorm:"primaryKey" json:"-"`                      // Hide internal auto-increment ID
	SN         int64     `gorm:"type:bigint;uniqueIndex" json:"sn,string"` // Snowflake ID exposed as "sn"
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"-"`
}

// BeforeCreate hook to auto-generate Snowflake ID
// 防御性赋值, 创建 article_tag_model 的时候,必须立马需要 SN
func (m *MODEL) BeforeCreate(tx *gorm.DB) error {
	// 防御性赋值
	if m.SN == 0 {
		m.SN = global.Snowflake.GetID()
	}
	return nil
}

type ParamDeleteList struct {
	SNList []string `json:"sn_list" binding:"required"`
}
type ParamList struct {
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}
type ParamSN struct {
	SN int64 `json:"sn,string" binding:"required"`
}

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Count int64 `json:"count"`
}

// 默认按照最新时间排序
const OrderByTime = "create_time DESC, sn DESC"
