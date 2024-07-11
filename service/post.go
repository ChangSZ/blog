package service

import (
	"context"
	"errors"
	"html/template"
	"strings"
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/model"
)

func ConsolePostCount(ctx context.Context, limit int, offset int, isTrash bool) (count int64, err error) {
	post := new(model.Posts)
	if isTrash {
		err = conf.SqlServer.WithContext(ctx).Table(post.TableName()).Unscoped().
			Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").
			Order("id desc").Offset(offset).Limit(limit).Count(&count).Error
	} else {
		err = conf.SqlServer.WithContext(ctx).Table(post.TableName()).
			Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
			Order("id desc").Offset(offset).Limit(limit).Count(&count).Error
	}
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return 0, err
	}
	return count, nil
}

func ConsolePostIndex(ctx context.Context, limit int, offset int, isTrash bool) (postListArr []*common.ConsolePostList, err error) {
	posts := make([]model.Posts, 0)
	if isTrash {
		err = conf.SqlServer.WithContext(ctx).Unscoped().
			Where("`deleted_at` IS NOT NULL OR `deleted_at`=?", "0001-01-01 00:00:00").
			Order("id desc").Offset(offset).Limit(limit).Find(&posts).Error
	} else {
		err = conf.SqlServer.WithContext(ctx).
			Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
			Order("id desc").Offset(offset).Limit(limit).Find(&posts).Error
	}
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}

	for _, post := range posts {
		consolePost := common.ConsolePost{
			Id:        post.ID,
			Uid:       post.UID,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		//category
		cates, err := GetPostCateByPostId(ctx, post.ID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		consoleCate := common.ConsoleCate{
			Id:          cates.ID,
			Name:        cates.Name,
			DisplayName: cates.DisplayName,
			SeoDesc:     cates.SeoDesc,
		}

		//tag
		tagIds, err := GetPostTagsByPostId(ctx, post.ID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		tags, err := GetTagsByIds(ctx, tagIds)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		var consoleTags []common.ConsoleTag
		for _, v := range tags {
			consoleTag := common.ConsoleTag{
				Id:          v.ID,
				Name:        v.Name,
				DisplayName: v.DisplayName,
				SeoDesc:     v.SeoDesc,
				Num:         v.Num,
			}
			consoleTags = append(consoleTags, consoleTag)
		}

		//view
		view, err := PostView(ctx, post.ID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		consoleView := common.ConsoleView{
			Num: view.Num,
		}

		//user
		user, err := GetUserById(ctx, post.UserID)
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return nil, err
		}
		consoleUser := common.ConsoleUser{
			Id:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status,
		}

		postList := common.ConsolePostList{
			Post:     consolePost,
			Category: consoleCate,
			Tags:     consoleTags,
			View:     consoleView,
			Author:   consoleUser,
		}
		postListArr = append(postListArr, &postList)
	}

	return postListArr, nil
}

func PostView(ctx context.Context, postId int) (*model.PostViews, error) {
	postV := new(model.PostViews)
	err := conf.SqlServer.WithContext(ctx).Where("post_id = ?", postId).Find(&postV).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
	}
	return postV, nil
}

func PostStore(ctx context.Context, ps common.PostStore, userId int) {
	postCreate := &model.Posts{
		Title:    ps.Title,
		UserID:   userId,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	policy := bluemonday.UGCPolicy()
	policy.AllowStandardURLs()
	policy.AllowAttrs("href").OnElements("a")
	html := policy.SanitizeBytes(unsafe)
	postCreate.Content = string(html)

	var err error
	err = conf.SqlServer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(postCreate).Error; err != nil {
			return err
		}

		if ps.Category > 0 {
			postCateCreate := &model.PostCate{
				PostID: postCreate.ID,
				CateID: ps.Category,
			}
			if err = tx.Create(postCateCreate).Error; err != nil {
				return err
			}
		}

		if len(ps.Tags) > 0 {
			for _, v := range ps.Tags {
				postTagCreate := &model.PostTag{
					PostID: postCreate.ID,
					TagID:  v,
				}
				if err = tx.Create(postTagCreate).Error; err != nil {
					return err
				}

				if err = tx.Table((&model.Tags{}).TableName()).Where("id=?", v).UpdateColumn("num", gorm.Expr("num + ?", 1)).Error; err != nil {
					return err
				}
			}
		}

		postView := &model.PostViews{
			PostID: postCreate.ID,
			Num:    1,
		}
		if err = tx.Create(postView).Error; err != nil {
			return err
		}

		uid, err := conf.ZHashId.Encode([]int{postCreate.ID})
		if err != nil {
			log.WithTrace(ctx).Error(err)
			return err
		}

		newPostCreate := &model.Posts{
			ID:  postCreate.ID,
			UID: uid,
		}
		if err = tx.Updates(newPostCreate).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.WithTrace(ctx).Error(err)
	}
}

func PostDetail(ctx context.Context, postId int) (p *model.Posts, err error) {
	post := new(model.Posts)
	err = conf.SqlServer.WithContext(ctx).Where("id = ?", postId).Find(&post).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return post, err
	}
	return post, nil
}

func IndexPostDetailDao(ctx context.Context, postId int) (postDetail common.IndexPostDetail, err error) {
	post := new(model.Posts)
	err = conf.SqlServer.WithContext(ctx).Where("id = ?", postId).Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Find(&post).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	if post.ID <= 0 {
		return postDetail, errors.New("post do not exists ")
	}
	Post := common.IndexPost{
		Id:        post.ID,
		Uid:       post.UID,
		Title:     post.Title,
		Summary:   post.Summary,
		Original:  post.Original,
		Content:   template.HTML(post.Content),
		Password:  post.Password,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	tags, err := PostIdTags(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	var Tags []common.ConsoleTag
	for _, v := range tags {
		consoleTag := common.ConsoleTag{
			Id:          v.ID,
			Name:        v.Name,
			DisplayName: v.DisplayName,
			SeoDesc:     v.SeoDesc,
			Num:         v.Num,
		}
		Tags = append(Tags, consoleTag)
	}

	cate, err := PostCates(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	Cate := common.ConsoleCate{
		Id:          cate.ID,
		Name:        cate.Name,
		DisplayName: cate.DisplayName,
		SeoDesc:     cate.SeoDesc,
	}

	//view
	view, err := PostView(ctx, post.ID)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	View := common.ConsoleView{
		Num: view.Num,
	}

	//user
	user, err := GetUserById(ctx, post.UserID)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	Author := common.ConsoleUser{
		Id:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}

	// last post
	lastPost, err := LastPost(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	// next post
	nextPost, err := NextPost(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	postDetail = common.IndexPostDetail{
		Post:     Post,
		Category: Cate,
		Tags:     Tags,
		View:     View,
		Author:   Author,
		LastPost: lastPost,
		NextPost: nextPost,
	}

	return postDetail, nil
}

func LastPost(ctx context.Context, postId int) (post *model.Posts, err error) {
	post = new(model.Posts)
	err = conf.SqlServer.WithContext(ctx).Where("id < ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Order("id desc").Find(&post).Error
	return
}

func NextPost(ctx context.Context, postId int) (post *model.Posts, err error) {
	post = new(model.Posts)
	err = conf.SqlServer.WithContext(ctx).Where("id > ?", postId).
		Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").
		Order("id asc").Find(&post).Error
	return
}

func PostIdTags(ctx context.Context, postId int) (tags []*model.Tags, err error) {
	tagIds, err := PostIdTag(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	err = conf.SqlServer.WithContext(ctx).Where("id in ?", tagIds).Find(&tags).Error
	return
}

func PostIdTag(ctx context.Context, postId int) (tagIds []int, err error) {
	postTag := make([]model.PostTag, 0)
	err = conf.SqlServer.WithContext(ctx).Where("post_id = ?", postId).Find(&postTag).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}

	for _, v := range postTag {
		tagIds = append(tagIds, v.TagID)
	}
	return tagIds, nil
}

func PostCates(ctx context.Context, postId int) (cate *model.Categories, err error) {
	cateId, err := PostCate(ctx, postId)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	cate = new(model.Categories)
	err = conf.SqlServer.WithContext(ctx).Where("id=?", cateId).Find(&cate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	return
}

func PostCate(ctx context.Context, postId int) (int, error) {
	postCate := new(model.PostCate)
	err := conf.SqlServer.WithContext(ctx).Where("post_id = ?", postId).Find(&postCate).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return 0, err
	}
	return postCate.CateID, nil
}

func PostUpdate(ctx context.Context, postId int, ps common.PostStore) {
	postUpdate := &model.Posts{
		Title:    ps.Title,
		UserID:   1,
		Summary:  ps.Summary,
		Original: ps.Content,
	}

	unsafe := blackfriday.Run([]byte(ps.Content))
	policy := bluemonday.UGCPolicy()
	policy.AllowStandardURLs()
	policy.AllowAttrs("href").OnElements("a")
	html := policy.SanitizeBytes(unsafe)
	postUpdate.Content = strings.Replace(string(html), "&amp;", "&", -1) // 将 &amp; 转换为 &
	postUpdate.ID = postId

	var err error
	err = conf.SqlServer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(postUpdate).Error; err != nil {
			return err
		}

		postCate := new(model.PostCate)
		if err = tx.Where("post_id = ?", postId).Delete(&postCate).Error; err != nil {
			return err
		}

		if ps.Category > 0 {
			postCateCreate := &model.PostCate{
				PostID: postId,
				CateID: ps.Category,
			}

			if err = tx.Create(postCateCreate).Error; err != nil {
				return err
			}
		}

		postTag := make([]model.PostTag, 0)
		if err = tx.Where("post_id = ?", postId).Find(&postTag).Error; err != nil {
			return err
		}
		if len(postTag) > 0 {
			for _, v := range postTag {
				err = tx.Table((&model.Tags{}).TableName()).Where("id=?", v.TagID).UpdateColumn("num", gorm.Expr("num - ?", 1)).Error
			}

			pt := new(model.PostTag)
			if err = tx.Where("post_id = ?", postId).Delete(&pt).Error; err != nil {
				return err
			}
		}

		if len(ps.Tags) > 0 {
			for _, v := range ps.Tags {
				postTagCreate := &model.PostTag{
					PostID: postId,
					TagID:  v,
				}
				if err = tx.Create(postTagCreate).Error; err != nil {
					return err
				}

				if err = tx.Table((&model.Tags{}).TableName()).Where("id=?", v).UpdateColumn("num", gorm.Expr("num + ?", 1)).Error; err != nil {
					return err
				}

			}
		}

		return nil
	})
	if err != nil {
		log.WithTrace(ctx).Error(err)
	}
}

func PostDestroy(ctx context.Context, postId int) (bool, error) {
	post := new(model.Posts)
	toBeCharge := time.Now().Format(time.RFC3339)
	timeLayout := time.RFC3339
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	if err != nil {
		return false, err
	}
	post.DeletedAt = &theTime
	post.ID = postId
	err = conf.SqlServer.WithContext(ctx).Updates(post).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return false, err
	}
	return true, nil
}

func PostUnTrash(ctx context.Context, postId int) (bool, error) {
	post := new(model.Posts)
	theTime, _ := time.Parse("2006-01-02 15:04:05", "")
	post.DeletedAt = &theTime
	post.ID = postId
	err := conf.SqlServer.WithContext(ctx).Updates(post).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return false, err
	}
	return true, nil
}

func PostCnt(ctx context.Context) (cnt int64, err error) {
	post := new(model.Posts)
	err = conf.SqlServer.WithContext(ctx).Table(post.TableName()).Count(&cnt).Error
	return
}

func PostTagListCount(ctx context.Context, tagId int, limit int, offset int) (count int64, err error) {
	postTag := new(model.PostTag)
	err = conf.SqlServer.WithContext(ctx).Table(postTag.TableName()).
		Where("tag_id = ?", tagId).Order("id desc").Offset(offset).Limit(limit).Count(&count).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return 0, err
	}
	return
}

func PostTagList(ctx context.Context, tagId int, limit int, offset int) (postListArr []*common.ConsolePostList, err error) {
	postTags := make([]model.PostTag, 0)
	err = conf.SqlServer.WithContext(ctx).Table("post_tag").Where("tag_id = ?", tagId).Order("id desc").Offset(offset).Limit(limit).Find(&postTags).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}
	for _, postTag := range postTags {
		post := new(model.Posts)
		err = conf.SqlServer.WithContext(ctx).Where(postTag.PostID).Find(&post).Error
		if err != nil {
			return nil, err
		}
		if post.ID == 0 {
			log.WithTrace(ctx).Warnf("not find post by id(%v)", postTag.PostID)
			continue
		}

		consolePost := common.ConsolePost{
			Id:        post.ID,
			Uid:       post.UID,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := common.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}
	return postListArr, nil
}

func PostCateListCount(ctx context.Context, cateId int, limit int, offset int) (count int64, err error) {
	postCate := new(model.PostCate)
	err = conf.SqlServer.WithContext(ctx).Table(postCate.TableName()).
		Where("cate_id = ?", cateId).
		Order("id desc").Offset(offset).Limit(limit).Count(&count).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return 0, err
	}
	return
}

func PostCateList(ctx context.Context, cateId int, limit int, offset int) (postListArr []*common.ConsolePostList, err error) {
	postCates := make([]model.PostCate, 0)
	err = conf.SqlServer.WithContext(ctx).Where("cate_id = ?", cateId).Order("id desc").Offset(offset).Limit(limit).Find(&postCates).Error
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return nil, err
	}

	for _, postCate := range postCates {
		post := new(model.Posts)
		err = conf.SqlServer.WithContext(ctx).Where("id=?", postCate.PostID).Find(&post).Error
		if err != nil {
			return nil, err
		}
		if post.ID == 0 {
			log.WithTrace(ctx).Warnf("not find post by id(%v)", postCate.PostID)
			continue
		}

		consolePost := common.ConsolePost{
			Id:        post.ID,
			Uid:       post.UID,
			Title:     post.Title,
			Summary:   post.Summary,
			Original:  post.Original,
			Content:   post.Content,
			Password:  post.Password,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		postList := common.ConsolePostList{
			Post: consolePost,
		}
		postListArr = append(postListArr, &postList)
	}
	return postListArr, nil
}
