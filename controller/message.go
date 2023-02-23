package controller

import (
	"github.com/SCUTKing/service"
	"github.com/SCUTKing/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageListResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	// 参数绑定

	toUserID, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 64)
	//actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)   默认是1
	content := c.Query("content")

	// 判断是否登录
	var (
		isLogin  bool
		viewerID uint64 //发消息的人
	)
	// 判断传入的token是否合法，用户是否存在
	if token := c.Query("token"); token != "" {
		claims, err := util.ParseToken(token)
		if err == nil {
			viewerID = claims.UserID
			isLogin = true
		}
	}

	//发送消息给那个人  （更新数据库）

	//判断是否登录
	if isLogin {
		IsSendSuccess, err := service.SendMessage(viewerID, toUserID, content)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送失败！"})
			return
		}
		if IsSendSuccess {
			// 返回成功并生成响应 json
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "发送成功！"})
		}
	}

}

// MessageList 获取消息列表
func MessageList(c *gin.Context) {
	//参数绑定
	toUserID, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 64)
	// 判断是否登录
	var (
		isLogin  bool
		viewerID uint64 //发消息的人
	)
	// 判断传入的token是否合法，用户是否存在
	if token := c.Query("token"); token != "" {
		claims, err := util.ParseToken(token)
		if err == nil {
			viewerID = claims.UserID
			isLogin = true
		}
	}

	//判断是否登录
	if isLogin {
		var messageList []Message
		//查到了最新的消息
		messageGetList, err := service.GetMessageList(viewerID, toUserID)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取失败！"})
			return
		}
		//直接删除最新的消息  防止它重复出现
		service.DeleteMessageList(viewerID, toUserID)

		for _, each := range messageGetList {
			var message Message
			message.ID = each.MessageId
			message.Content = each.Content
			message.CreateDate = each.CreatedAt.Unix()
			messageList = append(messageList, message)
		}

		// 返回成功并生成响应 json
		c.JSON(http.StatusOK, MessageListResponse{
			Response:    Response{StatusCode: 0, StatusMsg: "OK"},
			MessageList: messageList,
		})

	}

}
