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
func PostArticleCommentDig(id uint) (string, error) {
	idStr := strconv.Itoa(int(id)) // 将 uint 转换为 string
	num, _ := global.Redis.HGet(getKeyName(KeyZSetCommentDigg), idStr).Int()
	num++ // 点赞数自增

	// 更新 Redis
	err := global.Redis.HSet(getKeyName(KeyZSetCommentDigg), idStr, num).Err()
	if err != nil {
		global.Log.Errorf("redis HSet error: %s\n", err.Error())
		return "", err
	}

	// 同步到 MySQL
	err = updateArticleCommentDiggCountInDB(id, num)
	if err != nil {
		return "", err
	}

	return "点赞成功", nil
}
func DeleteArticleComments(uID uint, deleteCommentIDList []uint) error {
	// 将所有 uint ID 转换为字符串切片，以便用于 Redis 删除
	var idStrList []string
	for _, id := range deleteCommentIDList {
		idStrList = append(idStrList, strconv.Itoa(int(id)))
	}

	// 批量删除 Redis 中的键
	err := global.Redis.HDel(getKeyName(KeyZSetCommentDigg), idStrList...).Err()
	if err != nil {
		global.Log.Errorf("HDel redis err:%s\n", err.Error())
		return err
	}
	return nil
}

// 更新评论点赞数到 MySQL
func updateArticleCommentDiggCountInDB(id uint, digCount int) error {
	err := global.DB.Model(&model.CommentModel{}).Where("id = ?", id).
		Update("digg_count", digCount).Error
	if err != nil {
		global.Log.Errorf("updateArticleDiggCountInDB err:%s\n", err.Error())
		return err
	}
	return nil
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
