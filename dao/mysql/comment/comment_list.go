package comment

import (
	"sync"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"gorm.io/gorm"
)

var userMap map[uint]model.UserDetail

func GetArticleComments(pcl *model.ParamCommentList) ([]model.ResponseCommentList, error) {
	// 查找与这篇文章所有相关用户的 IDList
	idList := make([]uint, 0)
	err := global.DB.Table("comment_models").
		Where("article_id = ?", pcl.ArticleID).
		Pluck("user_id", &idList).Error // 使用 Pluck 只提取 user_id 字段
	if err != nil {
		global.Log.Errorf("select idList err:%s\n", err.Error())
		return []model.ResponseCommentList{}, err
	}

	// 去重 idList
	uniqueIDList := make([]uint, 0)
	idSet := make(map[uint]struct{}) // 使用 map 去重
	for _, id := range idList {
		if _, exists := idSet[id]; !exists {
			idSet[id] = struct{}{}
			uniqueIDList = append(uniqueIDList, id)
		}
	}
	// 查找去重后的 idList 用户信息列表
	udl, err := user.GetUserDetailListByIDList(uniqueIDList)
	if err != nil {
		return []model.ResponseCommentList{}, err
	}
	// 存储起来，方便后面拼接数据
	userMap = make(map[uint]model.UserDetail, len(*udl))
	for _, ud := range *udl {
		userMap[ud.ID] = ud
	}
	// 获取根评论
	rootComments, err := GetRootComments(pcl.ArticleID)
	if err != nil {
		return nil, err
	}
	return rootComments, nil
}

func GetRootComments(articleID uint) ([]model.ResponseCommentList, error) {
	// 查找根评论
	var rootComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_id = ? AND parent_comment_id = -1", articleID).
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
		userDetail := userMap[comment.UserID]
		// 递归获取子评论
		subComments, err := getSubComments(comment.ID, articleID)
		if err != nil {
			return nil, err
		}

		// 构建响应评论
		responseComment := model.ResponseCommentList{
			MODEL:           comment.MODEL,
			Content:         comment.Content,
			ParentCommentID: comment.ParentCommentID,
			ArticleID:       comment.ArticleID,
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

func getSubComments(parentCommentID uint, articleID uint) ([]model.ResponseCommentList, error) {
	// 查找子评论
	var subComments []model.CommentModel
	err := global.DB.Table("comment_models").
		Where("article_id = ? AND parent_comment_id = ?", articleID, parentCommentID).
		Order("create_time ASC"). // 按创建时间排序
		Find(&subComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, err
	}

	// 如果没有子评论，直接返回空
	if len(subComments) == 0 {
		return nil, nil
	}

	// 并发处理子评论
	// wg 用来等待所有并发任务完成
	var wg sync.WaitGroup
	responseSubComments := make([]model.ResponseCommentList, 0, len(subComments))
	// ch 用于发送数据给 responseSubComments
	ch := make(chan model.ResponseCommentList)

	// 使用 Goroutine 并发获取每个子评论的详情
	for _, comment := range subComments {
		wg.Add(1)
		go func(comment model.CommentModel) {
			defer wg.Done()

			// 获取用户详情
			userDetail := userMap[comment.UserID]

			// 递归获取子评论
			subSubComments, err := getSubComments(comment.ID, articleID)
			if err != nil {
				global.Log.Errorf("获取子评论失败: %s", err)
				return
			}

			// 构建响应评论数据
			responseComment := model.ResponseCommentList{
				MODEL:           comment.MODEL,
				Content:         comment.Content,
				ParentCommentID: comment.ParentCommentID,
				ArticleID:       comment.ArticleID,
				DiggCount:       comment.DiggCount,
				CommentCount:    comment.CommentCount,
				SubComments:     subSubComments,
				UserDetail:      &userDetail,
			}

			// 通过 channel 发送结果
			ch <- responseComment
		}(comment)
	}

	// 等待所有 goroutine 完成并关闭 channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	// 从 channel 中读取并返回结果
	for comment := range ch {
		responseSubComments = append(responseSubComments, comment)
	}

	return responseSubComments, nil
}
func DeleteArticleComments(uID uint, pi *model.ParamID, articleID uint) (string, error) {
	// 创建一个局部变量来存储需要删除的评论 ID 列表
	var deleteCommentIDList []uint

	// 从给定的评论 ID 开始递归收集所有需要删除的评论
	if err := RecursiveCollectComments(pi.ID, &deleteCommentIDList); err != nil {
		global.Log.Errorf("Delete err:%s\n", err.Error())
		return "", err
	}

	// 将根评论 ID 添加到删除列表中
	deleteCommentIDList = append(deleteCommentIDList, pi.ID)

	// 计算需要递减的评论数量
	deleteCount := len(deleteCommentIDList)

	// 更新根评论的父评论计数，根据删除数量递减
	if err := decrementParentCommentCount(pi.ID, deleteCount); err != nil {
		return "", err
	}

	// 更新文章的评论计数
	if err := decrementArticleCommentCount(articleID, deleteCount); err != nil {
		return "", err
	}

	// 删除所有收集到的评论
	if err := global.DB.Where("id IN (?)", deleteCommentIDList).Delete(&model.CommentModel{}).Error; err != nil {
		return "", err
	}

	// 删除 Redis 中的评论缓存
	if err := redis.DeleteArticleComments(uID, deleteCommentIDList); err != nil {
		return "", err
	}

	return "删除成功", nil
}

// RecursiveCollectComments 递归收集评论及其所有子评论的 ID
func RecursiveCollectComments(commentID uint, deleteCommentIDList *[]uint) error {
	// 查询当前评论的子评论
	var subComments []model.CommentModel
	if err := global.DB.Where("parent_comment_id = ?", commentID).Find(&subComments).Error; err != nil {
		return err
	}

	// 递归收集每个子评论的 ID
	for _, subComment := range subComments {
		*deleteCommentIDList = append(*deleteCommentIDList, subComment.ID)
		if err := RecursiveCollectComments(subComment.ID, deleteCommentIDList); err != nil {
			return err
		}
	}

	return nil
}

// decrementParentCommentCount 递减父评论的 CommentCount
func decrementParentCommentCount(commentID uint, deleteCount int) error {
	// 获取当前评论的父评论 ID
	var parentID int
	if err := global.DB.Model(&model.CommentModel{}).
		Where("id = ?", commentID).
		Select("parent_comment_id").
		Scan(&parentID).Error; err != nil {
		return err
	}

	// 递减父评论的 CommentCount
	if parentID != -1 {
		return global.DB.Model(&model.CommentModel{}).
			Where("id = ?", parentID).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", deleteCount)).Error
	}
	return nil
}

// decrementArticleCommentCount 递减 article_models 表中的 comment_count
func decrementArticleCommentCount(articleID uint, count int) error {
	return global.DB.Model(&model.ArticleModel{}).
		Where("id = ?", articleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", count)).Error
}
