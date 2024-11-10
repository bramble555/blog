package redis

import (
	"strconv"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func PostArticleDig(id uint) (string, error) {
	idStr := strconv.Itoa(int(id)) // 将 uint 转换为 string
	num, _ := global.Redis.HGet(getKeyName(KeyZSetArticleDig), idStr).Int()
	num++ // 点赞数自增

	// 更新 Redis
	err := global.Redis.HSet(getKeyName(KeyZSetArticleDig), idStr, num).Err()
	if err != nil {
		global.Log.Errorf("redis HSet error: %s\n", err.Error())
		return "", err
	}

	// 同步到 MySQL
	err = updateArticleDigCountInDB(id, num)
	if err != nil {
		return "", err
	}

	return "点赞成功", nil
}

// 更新文章点赞数到 MySQL
func updateArticleDigCountInDB(id uint, digCount int) error {
	err := global.DB.Model(&model.ArticleModel{}).Where("id = ?", id).
		Update("digg_count", digCount).Error
	if err != nil {
		global.Log.Errorf("updateArticleDigCountInDB err:%s\n", err.Error())
		return err
	}
	return nil
}
