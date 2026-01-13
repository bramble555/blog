package comment

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

var userMap map[int64]model.UserDetail // Changed to int64 for Snowflake IDs

func GetArticleComments(pcl *model.ParamCommentList) ([]model.ResponseCommentList, error) {
	// 查找与这篇文章所有相关用户的 SNList
	snList := make([]int64, 0)
	err := global.DB.Table("comment_models").
		Where("article_sn = ?", pcl.ArticleSN).
		Pluck("user_sn", &snList).Error // 使用 Pluck 只提取 user_sn 字段
	if err != nil {
		global.Log.Errorf("select snList err:%s\n", err.Error())
		return []model.ResponseCommentList{}, err
	}

	// 去重 snList
	uniqueSNList := make([]int64, 0)
	snSet := make(map[int64]struct{}) // 使用 map 去重
	for _, sn := range snList {
		if _, exists := snSet[sn]; !exists {
			snSet[sn] = struct{}{}
			uniqueSNList = append(uniqueSNList, sn)
		}
	}
	// 查找去重后的 snList 用户信息列表
	udl, err := user.GetUserDetailListBySNList(uniqueSNList)
	if err != nil {
		return []model.ResponseCommentList{}, err
	}
	// 存储起来，方便后面拼接数据
	userMap = make(map[int64]model.UserDetail, len(*udl)) // Changed to int64 for Snowflake IDs
	for _, ud := range *udl {
		userMap[ud.SN] = ud
	}
	// 获取根评论
	rootComments, err := GetRootComments(pcl.ArticleSN)
	if err != nil {
		return nil, err
	}
	return rootComments, nil
}

func GetRootComments(articleSN int64) ([]model.ResponseCommentList, error) {
	// 查找根评论
	var rootComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_sn = ? AND parent_comment_sn = -1", articleSN).
		Order("create_time ASC"). // 按创建时间排序
		Find(&rootComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, err
	}

	// 构建响应数据
	responseCommentsList := make([]model.ResponseCommentList, 0, len(rootComments))
	for _, comment := range rootComments {
		// 查找用户详情，直接通过评论中的 UserID 获取
		userDetail := userMap[comment.UserSN]
		// 递归获取子评论
		subComments, err := getSubComments(comment.SN, articleSN)
		if err != nil {
			return nil, err
		}

		// 构建响应评论
		responseComment := model.ResponseCommentList{
			MODEL:           comment.MODEL,
			Content:         comment.Content,
			ParentCommentSN: comment.ParentCommentSN,
			ArticleSN:       comment.ArticleSN,
			DiggCount:       comment.DiggCount,
			CommentCount:    comment.CommentCount,
			SubComments:     subComments, // 添加子评论
			UserDetail:      &userDetail,
		}

		// 将根评论添加到响应列表
		responseCommentsList = append(responseCommentsList, responseComment)
	}

	return responseCommentsList, nil
}

func getSubComments(parentCommentSN int64, articleSN int64) ([]model.ResponseCommentList, error) {
	// 查找子评论
	var subComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_sn = ? AND parent_comment_sn = ?", articleSN, parentCommentSN).
		Order("create_time ASC").
		Find(&subComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, err
	}

	// 构建响应数据
	responseCommentsList := make([]model.ResponseCommentList, 0, len(subComments))
	for _, comment := range subComments {
		// 查找用户详情
		userDetail := userMap[comment.UserSN]

		// 递归获取子评论
		nestedSubComments, err := getSubComments(comment.SN, articleSN)
		if err != nil {
			return nil, err
		}

		// 构建响应评论
		responseComment := model.ResponseCommentList{
			MODEL:           comment.MODEL,
			Content:         comment.Content,
			ParentCommentSN: comment.ParentCommentSN,
			ArticleSN:       comment.ArticleSN,
			DiggCount:       comment.DiggCount,
			CommentCount:    comment.CommentCount,
			SubComments:     nestedSubComments, // 递归添加子评论
			UserDetail:      &userDetail,
		}

		// 将子评论添加到响应列表
		responseCommentsList = append(responseCommentsList, responseComment)
	}

	return responseCommentsList, nil
}
func DeleteArticleComments(uSN int64, role int64, psn *model.ParamSN, articleSN int64) (string, error) {
	// 验证权限：如果是管理员(1)可以直接删除；如果是用户，必须是自己的评论
	if role != int64(ctype.PermissionAdmin) {
		var commentUserSN int64
		// 查询评论的所有者
		err := global.DB.Table("comment_models").Where("sn = ?", psn.SN).Pluck("user_sn", &commentUserSN).Error
		if err != nil {
			return "", err
		}
		if commentUserSN != uSN {
			return "", errors.New("无权删除该评论")
		}
	}
	// 创建一个局部变量来存储需要删除的评论 SN 列表
	var deleteCommentSNList []int64

	// 从给定的评论 SN 开始递归收集所有需要删除的评论
	if err := RecursiveCollectComments(psn.SN, &deleteCommentSNList); err != nil {
		global.Log.Errorf("Delete err:%s\n", err.Error())
		return "", err
	}

	// 将根评论 SN 添加到删除列表中
	deleteCommentSNList = append(deleteCommentSNList, psn.SN)

	// 计算需要递减的评论数量
	deleteCount := len(deleteCommentSNList)

	// 更新根评论的父评论计数，根据删除数量递减
	if err := decrementParentCommentCount(psn.SN, deleteCount); err != nil {
		return "", err
	}

	// 更新文章的评论计数
	if err := decrementArticleCommentCount(articleSN, deleteCount); err != nil {
		return "", err
	}

	// 删除所有收集到的评论
	if err := global.DB.Table("comment_models").Where("sn IN (?)", deleteCommentSNList).Delete(&model.CommentModel{}).Error; err != nil {
		return "", err
	}

	// 删除 Redis 中的评论缓存
	if err := redis.DeleteArticleComments(uSN, deleteCommentSNList); err != nil {
		return "", err
	}

	return "删除成功", nil
}

// RecursiveCollectComments 递归收集评论及其所有子评论的 SN
func RecursiveCollectComments(commentSN int64, deleteCommentSNList *[]int64) error {
	// 查询当前评论的子评论
	var subComments []model.CommentModel
	if err := global.DB.Table("comment_models").Where("parent_comment_sn = ?", commentSN).Find(&subComments).Error; err != nil {
		return err
	}

	// 递归收集每个子评论的 SN
	for _, subComment := range subComments {
		*deleteCommentSNList = append(*deleteCommentSNList, subComment.SN)
		if err := RecursiveCollectComments(subComment.SN, deleteCommentSNList); err != nil {
			return err
		}
	}

	return nil
}

// decrementParentCommentCount 递减父评论的 CommentCount
func decrementParentCommentCount(commentSN int64, deleteCount int) error {
	// 获取当前评论的父评论 SN
	var parentSN int64
	if err := global.DB.Table("comment_models").
		Where("sn = ?", commentSN).
		Select("parent_comment_sn").
		Scan(&parentSN).Error; err != nil {
		return err
	}

	// 递减父评论的 CommentCount
	if parentSN != -1 {
		return global.DB.Table("comment_models").
			Where("sn = ?", parentSN).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", deleteCount)).Error
	}
	return nil
}

// decrementArticleCommentCount 递减 article_models 表中的 comment_count
func decrementArticleCommentCount(articleSN int64, count int) error {
	return global.DB.Table("article_models").
		Where("sn = ?", articleSN).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", count)).Error
}
