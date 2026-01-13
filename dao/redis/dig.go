package redis

import (
	"strconv"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
)

func PostArticleDig(sn int64) (string, error) {
	snStr := strconv.FormatInt(sn, 10)
	num, _ := global.Redis.HGet(getKeyName(KeyZSetArticleDig), snStr).Int()
	num++ // 点赞数自增

	// 更新 Redis
	err := global.Redis.HSet(getKeyName(KeyZSetArticleDig), snStr, num).Err()
	if err != nil {
		global.Log.Errorf("redis HSet error: %s\n", err.Error())
		return "", err
	}

	// 同步到 MySQL
	err = updateArticleDigCountInDB(sn, num)
	if err != nil {
		return "", err
	}

	return "点赞成功", nil
}
func PostArticleCommentDig(sn int64) (string, error) {
	snStr := strconv.FormatInt(sn, 10)
	num, _ := global.Redis.HGet(getKeyName(KeyZSetCommentDigg), snStr).Int()
	num++ // 点赞数自增

	// 更新 Redis
	err := global.Redis.HSet(getKeyName(KeyZSetCommentDigg), snStr, num).Err()
	if err != nil {
		global.Log.Errorf("redis HSet error: %s\n", err.Error())
		return "", err
	}

	// 同步到 MySQL
	err = updateArticleCommentDiggCountInDB(sn, num)
	if err != nil {
		return "", err
	}

	return "点赞成功", nil
}
func DeleteArticleComments(uSN int64, deleteCommentSNList []int64) error {
	// 将所有 int64 SN 转换为字符串切片，以便用于 Redis 删除
	var snStrList []string
	for _, sn := range deleteCommentSNList {
		snStrList = append(snStrList, strconv.FormatInt(sn, 10))
	}

	// 批量删除 Redis 中的键
	err := global.Redis.HDel(getKeyName(KeyZSetCommentDigg), snStrList...).Err()
	if err != nil {
		global.Log.Errorf("HDel redis err:%s\n", err.Error())
		return err
	}
	return nil
}

// 更新评论点赞数到 MySQL
func updateArticleCommentDiggCountInDB(sn int64, digCount int) error {
	err := global.DB.Model(&model.CommentModel{}).Where("sn = ?", sn).
		Update("digg_count", digCount).Error
	if err != nil {
		global.Log.Errorf("updateArticleDiggCountInDB err:%s\n", err.Error())
		return err
	}
	return nil
}

// 更新文章点赞数到 MySQL
func updateArticleDigCountInDB(sn int64, digCount int) error {
	err := global.DB.Model(&model.ArticleModel{}).Where("sn = ?", sn).
		Update("digg_count", digCount).Error
	if err != nil {
		global.Log.Errorf("updateArticleDigCountInDB err:%s\n", err.Error())
		return err
	}
	return nil
}
