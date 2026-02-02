package comment

import (
	"errors"

	"github.com/bramble555/blog/dao/mysql/digg"
	"github.com/bramble555/blog/dao/mysql/user"
	dao_redis "github.com/bramble555/blog/dao/redis"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"gorm.io/gorm"
)

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
	// 防止出现查询重复的用户,所以先需要 map 去重,相当于 set
	snSet := make(map[int64]struct{}, len(rootComments))
	for _, c := range rootComments {
		snSet[c.UserSN] = struct{}{}
	}
	userSNs := make([]int64, 0, len(snSet))
	for sn := range snSet {
		userSNs = append(userSNs, sn)
	}
	// 查询用户详情
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
		udCopy := ud
		resp = append(resp, model.ResponseCommentList{
			MODEL:           c.MODEL,
			Content:         c.Content,
			ParentCommentSN: c.ParentCommentSN,
			ArticleSN:       c.ArticleSN,
			UserSN:          c.UserSN,
			DiggCount:       c.DiggCount,
			CommentCount:    c.CommentCount,
			SubComments:     nil,
			IsDigg:          isDigg,
			UserDetail:      &udCopy,
		})
	}

	// 成功获取评论后，合并 Redis 实时计数
	mergeCommentRealtimeCounts(&resp)

	return &model.ResponseCommentListWrapper{List: resp, Count: count}, nil
}
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
	userMap := make(map[int64]model.UserDetail, len(*udl))
	for _, ud := range *udl {
		userMap[ud.SN] = ud
	}

	if pcl.Page <= 0 {
		pcl.Page = 1
	}
	if pcl.Size <= 0 {
		pcl.Size = 10
	}
	if pcl.Order == "" {
		pcl.Order = model.OrderByTime
	}
	// 获取根评论
	rootComments, count, err := getRootComments(pcl, uSN, userMap)
	if err != nil {
		return nil, err
	}

	// 合并 Redis 实时计数
	mergeCommentRealtimeCounts(&rootComments)

	return &model.ResponseCommentListWrapper{
		List:  rootComments,
		Count: count,
	}, nil
}

func getRootComments(pcl *model.ParamCommentList, uSN int64, userMap map[int64]model.UserDetail) ([]model.ResponseCommentList, int64, error) {
	// 查找根评论
	var rootComments []model.CommentModel
	var count int64
	db := global.DB.Table("comment_models").
		Where("article_sn = ? AND parent_comment_sn = -1", pcl.ArticleSN)

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (pcl.Page - 1) * pcl.Size
	err := db.Order(pcl.Order).Limit(pcl.Size).Offset(offset).Find(&rootComments).Error
	if err != nil {
		global.Log.Errorf("comment_models Find err: %s\n", err.Error())
		return nil, 0, err
	}

	// 构建响应数据
	responseCommentsList := make([]model.ResponseCommentList, 0, len(rootComments))
	for _, comment := range rootComments {
		// 查找用户详情，直接通过评论中的 UserID 获取
		userDetail := userMap[comment.UserSN]
		userDetailCopy := userDetail
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
			UserSN:          comment.UserSN,
			DiggCount:       comment.DiggCount,
			CommentCount:    comment.CommentCount,
			SubComments:     subComments, // 添加子评论
			IsDigg:          isDigg,      // 是否点赞
			UserDetail:      &userDetailCopy,
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

	return buildSubComments(subComments, 0, articleSN, uSN, userMap, make([]model.ResponseCommentList, 0, len(subComments)))
}

func buildSubComments(subComments []model.CommentModel, idx int, articleSN int64, uSN int64, userMap map[int64]model.UserDetail, acc []model.ResponseCommentList) ([]model.ResponseCommentList, error) {
	if idx >= len(subComments) {
		return acc, nil
	}

	subComment := subComments[idx]
	userDetail := userMap[subComment.UserSN]
	userDetailCopy := userDetail

	grandChildren, err := getSubComments(subComment.SN, articleSN, uSN, userMap)
	if err != nil {
		return nil, err
	}

	isDigg := false
	if uSN != 0 {
		isDigg, _ = digg.IsUserCommentDigg(uSN, subComment.SN)
	}

	acc = append(acc, model.ResponseCommentList{
		MODEL:           subComment.MODEL,
		Content:         subComment.Content,
		ParentCommentSN: subComment.ParentCommentSN,
		ArticleSN:       subComment.ArticleSN,
		UserSN:          subComment.UserSN,
		DiggCount:       subComment.DiggCount,
		CommentCount:    subComment.CommentCount,
		SubComments:     grandChildren,
		IsDigg:          isDigg,
		UserDetail:      &userDetailCopy,
	})

	// 每次递归后，合并 Redis 实时点赞数 (不仅根评论，子评论也要最新)
	mergeCommentRealtimeCounts(&acc)

	return buildSubComments(subComments, idx+1, articleSN, uSN, userMap, acc)
}

// mergeCommentRealtimeCounts 从 Redis 批量获取评论计数并合并 (大厂标准)
func mergeCommentRealtimeCounts(comments *[]model.ResponseCommentList) {
	if len(*comments) == 0 {
		return
	}

	sns := make([]int64, len(*comments))
	for i, c := range *comments {
		sns[i] = c.SN
	}

	// 批量获取评论点赞数
	diggCounts, _ := dao_redis.GetRedisArticlesCounts(sns, dao_redis.FieldCommentDigg)

	// 合并
	for i := range *comments {
		sn := (*comments)[i].SN
		if val, ok := diggCounts[sn]; ok && val != 0 {
			(*comments)[i].DiggCount += val
		}
		// 如果有子评论且不为空，递归合并子评论计数
		if len((*comments)[i].SubComments) > 0 {
			mergeCommentRealtimeCounts(&(*comments)[i].SubComments)
		}
	}
}

// DeleteArticleComments 删除文章下的评论（包括子评论）
func DeleteArticleComments(uSN int64, role int64, sn int64, articleSN int64) (string, error) {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := deleteCommentTree(tx, uSN, role, sn, articleSN); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit().Error; err != nil {
		return "", err
	}
	return "删除成功", nil
}

// deleteCommentTree 删除评论树，包括子评论和子评论的子评论
// 函数里面的参数 articleSN 为 0，删除某个评论列表,从 comment 表中获取 articleSN
func DeleteCommentsList(uSN int64, role int64, snList []int64) (string, error) {
	if role != int64(ctype.PermissionAdmin) {
		return "", errors.New("无权批量删除评论")
	}
	if len(snList) == 0 {
		return "删除成功", nil
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, sn := range snList {
		if err := deleteCommentTree(tx, uSN, role, sn, 0); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}
	return "删除成功", nil
}

// deleteCommentTree 删除评论树，包括子评论和子评论的子评论
// 如果 articleSN 为 0，则说明的是删除某个评论列表,从 comment 表中获取 articleSN(管理员删除)
// 如果 articleSN 不为 0，则说明的是删除某个文章下的评论(用户删除)
func deleteCommentTree(tx *gorm.DB, uSN int64, role int64, commentSN int64, articleSN int64) error {
	var root model.CommentModel
	if err := tx.Table("comment_models").Where("sn = ?", commentSN).Take(&root).Error; err != nil {
		return err
	}

	if role != int64(ctype.PermissionAdmin) && root.UserSN != uSN {
		return errors.New("无权删除该评论")
	}

	if articleSN == 0 {
		articleSN = root.ArticleSN
	}

	deleteCommentSNList := make([]int64, 0, 8)
	if err := iterativeCollectComments(tx, commentSN, &deleteCommentSNList); err != nil {
		return err
	}
	// 评论树中的所有评论 SN 都收集到 deleteCommentSNList 中了
	deleteCommentSNList = append(deleteCommentSNList, commentSN)
	deleteCount := len(deleteCommentSNList)

	// 减少评论的祖先评论数量
	if err := decrementAncestorCommentCount(tx, root.ParentCommentSN, deleteCount); err != nil {
		return err
	}
	// 减少文章的评论数量
	if err := decrementArticleCommentCount(tx, articleSN, deleteCount); err != nil {
		return err
	}

	// 删除评论点赞表中的点赞记录
	if err := tx.Table("user_comment_digg_models").Where("comment_sn IN ?", deleteCommentSNList).Delete(&model.UserCommentDiggModel{}).Error; err != nil {
		return err
	}

	// 删除评论表中的评论记录
	if err := tx.Table("comment_models").Where("sn IN ?", deleteCommentSNList).Delete(&model.CommentModel{}).Error; err != nil {
		return err
	}

	return nil
}

// iterativeCollectComments 递归收集评论树中的所有评论 SN
// 包括当前评论和所有子评论
func iterativeCollectComments(tx *gorm.DB, commentSN int64, deleteCommentSNList *[]int64) error {
	// 用一个队列来模拟递归,用队列不用递归
	var queue []int64
	queue = append(queue, commentSN)

	// 不断从队列中取出评论，查找其子评论
	for len(queue) > 0 {
		// 取出队列中的第一个评论
		currentSN := queue[0]
		queue = queue[1:]

		// 查询当前评论的所有子评论的 SN
		var subCommentSNList []int64
		if err := tx.Table("comment_models").
			Select("sn").
			Where("parent_comment_sn = ?", currentSN).
			Find(&subCommentSNList).Error; err != nil {
			return err
		}

		// 将查询到的子评论 SN 加入到 deleteCommentSNList
		*deleteCommentSNList = append(*deleteCommentSNList, subCommentSNList...)

		// 将查询到的子评论 SN 添加到队列中
		for _, subSN := range subCommentSNList {
			queue = append(queue, subSN)
		}
	}

	return nil
}

// decrementAncestorCommentCount 减少评论的祖先评论数量
// 包括当前评论和所有父评论
func decrementAncestorCommentCount(tx *gorm.DB, parentSN int64, deleteCount int) error {
	if parentSN == -1 {
		return nil
	}

	// 更新当前评论的 comment_count
	if err := tx.Table("comment_models").
		Where("sn = ?", parentSN).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", deleteCount)).Error; err != nil {
		return err
	}

	// 直接查出父评论 SN
	var newParentSN int64
	if err := tx.Table("comment_models").
		Select("parent_comment_sn").
		Where("sn = ?", parentSN).
		Scan(&newParentSN).Error; err != nil {
		return err
	}

	// 递归更新
	return decrementAncestorCommentCount(tx, newParentSN, deleteCount)
}

// decrementArticleCommentCount 减少文章的评论数量
func decrementArticleCommentCount(tx *gorm.DB, articleSN int64, count int) error {
	return tx.Table("article_models").
		Where("sn = ?", articleSN).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", count)).Error
}
