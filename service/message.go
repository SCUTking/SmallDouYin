package service

import (
	"github.com/SCUTKing/global"
	"github.com/SCUTKing/model"
)

func SendMessage(from uint64, to uint64, content string) (bool, error) {

	//查询数据库和缓存看看是否有数据   有的话直接更新
	if isExist, err := GetMessageStatusForUpdate(from, to); err == nil {
		if isExist {
			//如果存在
			if err := global.DB.Model(&model.Message{}).Where("to_user_id = ? and from_user_id = ?", to, from).
				Update("content", content).Error; err != nil {
				return false, err
			}
		} else {
			//如果不存在
			//没有的话  直接写入
			var message model.Message
			// 数据库没有记录，写入数据库
			message.MessageId, _ = global.ID_GENERATOR.NextID()
			message.ToUserId = to
			message.FromUserId = from
			message.Content = content
			if err := global.DB.Create(&message).Error; err != nil {
				return false, err
			}
		}
		//更新缓存

		return true, nil
	} else {
		return false, err
	}
}

// GetMessageStatusForUpdate 判断数据库和缓存中是否有数据
func GetMessageStatusForUpdate(from, to uint64) (bool, error) {

	//查询缓存

	// 缓存不存在，查询数据库
	var count int64
	global.DB.Model(&model.Message{}).Where("to_user_id = ? and from_user_id = ?", to, from).Count(&count)
	if count > 0 {
		return true, nil
	}
	// 更新缓存

	//假设没有  就会一直在数据库中生成记录
	return false, nil
}

func GetMessageList(from, to uint64) ([]model.Message, error) {
	var messageList []model.Message
	//查询缓存

	//缓存没有  查数据库

	//获取消息的时候发送者与接收者应该调换
	result := global.DB.Where("to_user_id = ? and from_user_id = ?", from, to).Find(&messageList)
	if result.Error != nil {
		return nil, result.Error
	}
	//更新redis缓存

	return messageList, nil

}

// DeleteMessageList 为了避免客户端重发消息(但是用户再次加入聊天框  什么都没有了)
func DeleteMessageList(viewerID, toUserID uint64) {
	//直接删除数据库中的消息
	//获取消息的时候发送者与接收者应该调换
	global.DB.Model(&model.Message{}).Where("to_user_id = ? and from_user_id = ?", viewerID, toUserID).Delete(&model.Message{})
}
