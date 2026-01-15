package comment

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/digg"
	"github.com/bramble555/blog/dao/mysql/user"

	// "github.com/bramble555/blog/dao/redis"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

func GetArticleComments(pcl *model.ParamCommentList, uSN int64) (*model.ResponseCommentListWrapper, error) {
	// 查找与这篇文章所有相关用户的 SNList
	snList := make([]int64, 0)
	err := global.DB.Table("comment_models").
		Where("article_sn = ?", pcl.ArticleSN).
		Pluck("user_sn", &snList).Error // 使用 Pluck 只提取 user_sn 字段
	if err != nil {
		global.Log.Errorf("select snList err:%s\n", err.Error())
		return nil, err
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
		return nil, err
	}
	// 存储起来，方便后面拼接数据
	userMap := make(map[int64]model.UserDetail, len(*udl)) // Changed to local variable
	for _, ud := range *udl {
		userMap[ud.SN] = ud
	}

    if pcl.Page <= 0 {
        pcl.Page = 1
    }
    if pcl.Size <= 0 {
        pcl.Size = 10
    }

	// 获取根评论
	rootComments, count, err := GetRootComments(pcl, uSN, userMap)
	if err != nil {
		return nil, err
	}
	return &model.ResponseCommentListWrapper{
		List:  rootComments,
		Count: count,
	}, nil
}

// GetAllComments 获取所有文章的根评论（管理员后台用途）
func GetAllComments(pcl *model.ParamCommentList, uSN int64) (*model.ResponseCommentListWrapper, error) {
    if pcl.Page <= 0 {
        pcl.Page = 1
    }
    if pcl.Size <= 0 {
        pcl.Size = 20 // 后台默认每页20条
    }
    if pcl.Order == "" {
        pcl.Order = model.OrderByTime
    }

    // 查询所有根评论
    var rootComments []model.CommentModel
    var count int64
    db := global.DB.Table("comment_models").Where("parent_comment_sn = -1")
    if err := db.Count(&count).Error; err != nil {
        return nil, err
    }

    offset := (pcl.Page - 1) * pcl.Size
    if err := db.Order(pcl.Order).Limit(pcl.Size).Offset(offset).Find(&rootComments).Error; err != nil {
        global.Log.Errorf("comment_models Find all root err: %s\n", err.Error())
        return nil, err
    }

    // 预取用户详情：减少重复查询
    snSet := make(map[int64]struct{}, len(rootComments))
    for _, c := range rootComments {
        snSet[c.UserSN] = struct{}{}
    }
    userSNs := make([]int64, 0, len(snSet))
    for sn := range snSet {
        userSNs = append(userSNs, sn)
    }
    udl, err := user.GetUserDetailListBySNList(userSNs)
    if err != nil {
        return nil, err
    }
    userMap := make(map[int64]model.UserDetail, len(*udl))
    for _, ud := range *udl {
        userMap[ud.SN] = ud
    }

    // 构建响应列表（不递归子评论以提升性能）
    resp := make([]model.ResponseCommentList, 0, len(rootComments))
    for _, c := range rootComments {
        isDigg := false
        if uSN != 0 {
            isDigg, _ = digg.IsUserCommentDigg(uSN, c.SN)
        }
        ud := userMap[c.UserSN]
        resp = append(resp, model.ResponseCommentList{
            MODEL:           c.MODEL,
            Content:         c.Content,
            ParentCommentSN: c.ParentCommentSN,
            ArticleSN:       c.ArticleSN,
            DiggCount:       c.DiggCount,
            CommentCount:    c.CommentCount,
            SubComments:     nil,
            IsDigg:          isDigg,
            UserDetail:      &ud,
        })
    }

    return &model.ResponseCommentListWrapper{List: resp, Count: count}, nil
}

func GetRootComments(pcl *model.ParamCommentList, uSN int64, userMap map[int64]model.UserDetail) ([]model.ResponseCommentList, int64, error) {
	// 查找根评论
	var rootComments []model.CommentModel
	var count int64
	db := global.DB.Table("comment_models").
		Where("article_sn = ? AND parent_comment_sn = -1", pcl.ArticleSN)

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (pcl.Page - 1) * pcl.Size
	err := db.Order("create_time ASC").Limit(pcl.Size).Offset(offset).Find(&rootComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, 0, err
	}

	// 构建响应数据
	responseCommentsList := make([]model.ResponseCommentList, 0, len(rootComments))
	for _, comment := range rootComments {
		// 查找用户详情，直接通过评论中的 UserID 获取
		userDetail := userMap[comment.UserSN]
		// 递归获取子评论
		subComments, err := getSubComments(comment.SN, pcl.ArticleSN, uSN, userMap)
		if err != nil {
			return nil, 0, err
		}

		// 检查点赞状态
		isDigg := false
		if uSN != 0 {
			isDigg, _ = digg.IsUserCommentDigg(uSN, comment.SN)
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
			IsDigg:          isDigg,      // 是否点赞
			UserDetail:      &userDetail,
		}

		// 将根评论添加到响应列表
		responseCommentsList = append(responseCommentsList, responseComment)
	}
	return responseCommentsList, count, nil
}

func getSubComments(parentSN int64, articleSN int64, uSN int64, userMap map[int64]model.UserDetail) ([]model.ResponseCommentList, error) {
	var subComments []model.CommentModel
	// 查找子评论
	err := global.DB.Table("comment_models").
		Where("article_sn = ? AND parent_comment_sn = ?", articleSN, parentSN).
		Order("create_time ASC").
		Find(&subComments).Error
	if err != nil {
		return nil, err
	}

	var responseSubComments []model.ResponseCommentList
	for _, subComment := range subComments {
		userDetail := userMap[subComment.UserSN]

		// 递归获取孙子评论
		grandChildren, err := getSubComments(subComment.SN, articleSN, uSN, userMap)
		if err != nil {
			return nil, err
		}

		isDigg := false
		if uSN != 0 {
			isDigg, _ = digg.IsUserCommentDigg(uSN, subComment.SN)
		}

		responseSubComment := model.ResponseCommentList{
			MODEL:           subComment.MODEL,
			Content:         subComment.Content,
			ParentCommentSN: subComment.ParentCommentSN,
			ArticleSN:       subComment.ArticleSN,
			DiggCount:       subComment.DiggCount,
			CommentCount:    subComment.CommentCount,
			SubComments:     grandChildren,
			IsDigg:          isDigg,
			UserDetail:      &userDetail,
		}
		responseSubComments = append(responseSubComments, responseSubComment)
	}
	return responseSubComments, nil
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

	// 删除 Redis 中的评论缓存 (Disabled as per requirement)
	// if err := redis.DeleteArticleComments(uSN, deleteCommentSNList); err != nil {
	// 	return "", err
	// }

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
