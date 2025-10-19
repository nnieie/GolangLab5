package pack

import (
	apiBase "github.com/nnieie/golanglab5/cmd/api/biz/model/base"
	kitBase "github.com/nnieie/golanglab5/kitex_gen/base"
)

func BaseRespRPCToBaseResp(base *kitBase.BaseResp) *apiBase.BaseResp {
	if base == nil {
		return nil
	}
	return &apiBase.BaseResp{
		Code: base.Code,
		Msg:  base.Msg,
	}
}

func UserRPCToUser(user *kitBase.User) *apiBase.User {
	if user == nil {
		return nil
	}
	return &apiBase.User{
		ID:        user.Id,
		Username:  user.Username,
		AvatarURL: user.AvatarUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func UsersRPCToUsers(users []*kitBase.User) []*apiBase.User {
	res := make([]*apiBase.User, 0, len(users))
	for _, user := range users {
		res = append(res, UserRPCToUser(user))
	}
	return res
}

func VideoRPCToVideo(video *kitBase.Video) *apiBase.Video {
	if video == nil {
		return nil
	}
	return &apiBase.Video{
		ID:           video.Id,
		UserID:       video.UserId,
		VideoURL:     video.VideoUrl,
		CoverURL:     video.CoverUrl,
		Title:        video.Title,
		Description:  video.Description,
		VisitCount:   video.VisitCount,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
		CreatedAt:    video.CreatedAt,
		UpdatedAt:    video.UpdatedAt,
		DeletedAt:    video.DeletedAt,
	}
}

func VideosRPCToVideos(videos []*kitBase.Video) []*apiBase.Video {
	res := make([]*apiBase.Video, 0, len(videos))
	for _, video := range videos {
		res = append(res, VideoRPCToVideo(video))
	}
	return res
}

func CommentRPCToComment(comment *kitBase.Comment) *apiBase.Comment {
	if comment == nil {
		return nil
	}
	return &apiBase.Comment{
		ID:         comment.Id,
		UserID:     comment.UserId,
		VideoID:    comment.VideoId,
		ParentID:   comment.ParentId,
		ChildCount: comment.ChildCount,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
		DeletedAt:  comment.DeletedAt,
	}
}

func CommentsRPCToComments(comments []*kitBase.Comment) []*apiBase.Comment {
	res := make([]*apiBase.Comment, 0, len(comments))
	for _, comment := range comments {
		res = append(res, CommentRPCToComment(comment))
	}
	return res
}

func PrivateMessageRPCToPrivateMessage(msg *kitBase.PrivateMessage) *apiBase.PrivateMessage {
	if msg == nil {
		return nil
	}
	return &apiBase.PrivateMessage{
		FromUserID: msg.FromUserId,
		ToUserID:   msg.ToUserId,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}
}

func PrivateMessagesRPCToPrivateMessages(msgs []*kitBase.PrivateMessage) []*apiBase.PrivateMessage {
	res := make([]*apiBase.PrivateMessage, 0, len(msgs))
	for _, msg := range msgs {
		res = append(res, PrivateMessageRPCToPrivateMessage(msg))
	}
	return res
}

func GroupMessageRPCToGroupMessage(msg *kitBase.GroupMessage) *apiBase.GroupMessage {
	if msg == nil {
		return nil
	}
	return &apiBase.GroupMessage{
		FromUserID: msg.FromUserId,
		GroupID:    msg.GroupId,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}
}

func GroupMessagesRPCToGroupMessages(msgs []*kitBase.GroupMessage) []*apiBase.GroupMessage {
	res := make([]*apiBase.GroupMessage, 0, len(msgs))
	for _, msg := range msgs {
		res = append(res, GroupMessageRPCToGroupMessage(msg))
	}
	return res
}
