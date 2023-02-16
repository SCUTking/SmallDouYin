package controller

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/goldenBill/douyin-fighting/global"
	"github.com/goldenBill/douyin-fighting/model"
	"github.com/goldenBill/douyin-fighting/service"
	"github.com/goldenBill/douyin-fighting/util"
)

// CommentActionRequest 璇勮鎿嶄綔鐨勮姹�
type CommentActionRequest struct {
	UserID      uint64 `form:"user_id" json:"user_id"` // apk骞舵病鏈変紶user_id杩欎釜鍙傛暟
	Token       string `form:"token" json:"token"`
	VideoID     uint64 `form:"video_id" json:"video_id"`
	ActionType  uint   `form:"action_type" json:"action_type"`
	CommentText string `form:"comment_text" json:"comment_text"`
	CommentID   uint64 `form:"comment_id" json:"comment_id"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentListRequest 璇勮鍒楄〃鐨勮姹�
type CommentListRequest struct {
	UserID  uint64 `form:"user_id" json:"user_id"`
	Token   string `form:"token" json:"token"`
	VideoID uint64 `form:"video_id" json:"video_id"`
}

// CommentListResponse 璇勮鍒楄〃鐨勫搷搴�
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction 璇勮鎿嶄綔鎺ュ彛
// 1. 纭繚鎿嶄綔绫诲瀷姝ｇ‘ 2. 纭繚褰撳墠鐢ㄦ埛鏈夋潈闄愬垹闄�
func CommentAction(c *gin.Context) {
	// 鍙傛暟缁戝畾
	var r CommentActionRequest
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "bind error"})
		return
	}

	// 鍒ゆ柇 action_type 鏄惁姝ｇ‘
	if r.ActionType != 1 && r.ActionType != 2 {
		// action_type 涓嶅悎娉�
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "action type error"})
		return
	}

	// 鑾峰彇 userID
	r.UserID = c.GetUint64("UserID")

	// 璇勮鎿嶄綔
	if r.ActionType == 1 {
		// 鍒ゆ柇comment鏄惁鍚堟硶
		if utf8.RuneCountInString(r.CommentText) > global.MAX_COMMENT_LENGTH ||
			utf8.RuneCountInString(r.CommentText) <= 0 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "闈炴硶璇勮"})
			return
		}
		// 娣诲姞璇勮
		commentID, err := global.ID_GENERATOR.NextID()
		if err != nil {
			// 鐢熸垚ID澶辫触
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		commentModel := model.Comment{
			CommentID: commentID,
			VideoID:   r.VideoID,
			UserID:    r.UserID,
			Content:   r.CommentText,
		}
		if err = service.AddComment(&commentModel); err != nil {
			// 璇勮澶辫触
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "comment failed"})
			return
		}
		userModel, err := service.UserInfoByUserID(commentModel.UserID)
		if err != nil {
			// 鏈壘鍒拌瘎璁虹殑鐢ㄦ埛
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "comment failed"})
			return
		}
		// 鎵归噺鍒ゆ柇鐢ㄦ埛鏄惁鍏虫敞
		isFollow, err := service.GetFollowStatus(r.UserID, userModel.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		// 杩斿洖JSON
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment: Comment{
				ID: commentModel.CommentID,
				User: User{
					ID:             userModel.UserID,
					Name:           userModel.Name,
					FollowCount:    userModel.FollowCount,
					FollowerCount:  userModel.FollowerCount,
					TotalFavorited: userModel.TotalFavorited,
					FavoriteCount:  userModel.FavoriteCount,
					IsFollow:       isFollow,
				},
				Content:    commentModel.Content,
				CreateDate: commentModel.CreatedAt.Format("2006-01-02 15:04"),
			},
		})
		return
	}

	// 鍒犻櫎璇勮
	if err := service.DeleteComment(r.UserID, r.VideoID, r.CommentID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "comment failed"})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// CommentList 璇勮鍒楄〃鎺ュ彛
func CommentList(c *gin.Context) {
	// 鍙傛暟缁戝畾
	var r CommentListRequest
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "bind error"})
		return
	}

	var commentModelList []model.Comment
	var userModelList []model.User
	// 鑾峰彇璇勮鍒楄〃浠ュ強瀵瑰簲鐨勪綔鑰�
	if err := service.GetCommentListAndUserListRedis(r.VideoID, &commentModelList, &userModelList); err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	var (
		isFollowList []bool
		isLogged     = false // 鐢ㄦ埛鏄惁浼犲叆浜嗗悎娉曟湁鏁堢殑token锛堟槸鍚︾櫥褰曪級
		isFollow     bool
		err          error
	)

	var userID uint64
	// 鍒ゆ柇浼犲叆鐨則oken鏄惁鍚堟硶锛岀敤鎴锋槸鍚﹀瓨鍦�
	if token := c.Query("token"); token != "" {
		claims, err := util.ParseToken(token)
		if err == nil {
			// token鍚堟硶
			userID = claims.UserID
			isLogged = true
		}
	}

	if isLogged {
		// 褰撶敤鎴风櫥褰曟椂 涓€娆℃€ц幏鍙栫敤鎴锋槸鍚︾偣璧炰簡鍒楄〃涓殑瑙嗛浠ュ強鏄惁鍏虫敞浜嗚瘎璁虹殑浣滆€�
		authorIDList := make([]uint64, len(commentModelList))
		for i, user_ := range userModelList {
			authorIDList[i] = user_.UserID
		}
		// 鎵归噺鍒ゆ柇鐢ㄦ埛鏄惁鍏虫敞璇勮鐨勪綔鑰�
		isFollowList, err = service.GetFollowStatusList(userID, authorIDList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
	}

	var (
		commentJsonList = make([]Comment, 0, len(commentModelList))
		commentJson     Comment
		userJson        User
		user            model.User
	)

	for i, comment := range commentModelList {
		// 鏈櫥褰曟椂榛樿涓烘湭鍏虫敞鏈偣璧�
		isFollow = false
		if isLogged {
			// 褰撶敤鎴风櫥褰曟椂锛屽垽鏂槸鍚﹀叧娉ㄥ綋鍓嶄綔鑰�
			isFollow = isFollowList[i]
		}
		user = userModelList[i]
		userJson.ID = user.UserID
		userJson.Name = user.Name
		userJson.FollowCount = user.FollowCount
		userJson.FollowerCount = user.FollowerCount
		userJson.TotalFavorited = user.TotalFavorited
		userJson.FavoriteCount = user.FavoriteCount
		userJson.IsFollow = isFollow

		commentJson.ID = comment.CommentID
		commentJson.User = userJson
		commentJson.Content = comment.Content
		commentJson.CreateDate = comment.CreatedAt.Format("2006-01-02 15:04")

		commentJsonList = append(commentJsonList, commentJson)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentJsonList,
	})
}
