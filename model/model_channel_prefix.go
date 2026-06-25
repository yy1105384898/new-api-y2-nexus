package model

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"

	"gorm.io/gorm"
)

type ModelChannelPrefix struct {
	Id          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix      string         `json:"prefix" gorm:"size:64;not null;uniqueIndex:uk_model_channel_prefix"`
	Note        string         `json:"note" gorm:"size:255"`
	Enabled     bool           `json:"enabled" gorm:"default:true;not null"`
	SortOrder   int            `json:"sort_order" gorm:"default:0;not null"`
	CreatedTime int64          `json:"created_time" gorm:"bigint"`
	UpdatedTime int64          `json:"updated_time" gorm:"bigint"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func NormalizeModelChannelPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	if !strings.HasSuffix(prefix, "-") {
		prefix += "-"
	}
	return prefix
}

func (item *ModelChannelPrefix) Insert() error {
	item.Prefix = NormalizeModelChannelPrefix(item.Prefix)
	if item.Prefix == "" {
		return errors.New("prefix is required")
	}
	now := common.GetTimestamp()
	item.CreatedTime = now
	item.UpdatedTime = now
	return DB.Create(item).Error
}

func (item *ModelChannelPrefix) Update() error {
	item.Prefix = NormalizeModelChannelPrefix(item.Prefix)
	if item.Prefix == "" {
		return errors.New("prefix is required")
	}
	item.UpdatedTime = common.GetTimestamp()
	return DB.Model(item).Where("id = ?", item.Id).Updates(map[string]interface{}{
		"prefix":       item.Prefix,
		"note":         item.Note,
		"enabled":      item.Enabled,
		"sort_order":   item.SortOrder,
		"updated_time": item.UpdatedTime,
	}).Error
}

func GetAllModelChannelPrefixes() ([]ModelChannelPrefix, error) {
	var items []ModelChannelPrefix
	err := DB.Order("sort_order asc, prefix asc").Find(&items).Error
	return items, err
}

func GetEnabledModelChannelPrefixes() ([]ModelChannelPrefix, error) {
	var items []ModelChannelPrefix
	err := DB.Where("enabled = ?", true).Order("sort_order asc, prefix asc").Find(&items).Error
	return items, err
}

func GetModelChannelPrefixByID(id int) (*ModelChannelPrefix, error) {
	var item ModelChannelPrefix
	err := DB.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func IsModelChannelPrefixDuplicated(id int, prefix string) (bool, error) {
	prefix = NormalizeModelChannelPrefix(prefix)
	if prefix == "" {
		return false, nil
	}
	var count int64
	err := DB.Model(&ModelChannelPrefix{}).
		Where("prefix = ? AND id <> ?", prefix, id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteModelChannelPrefix(id int) error {
	result := DB.Delete(&ModelChannelPrefix{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
