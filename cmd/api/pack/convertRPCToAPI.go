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
		Avatar:    user.AvatarUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func VideoRPCToVideo(video *kitBase.Video) *apiBase.Video {
	if video == nil {
		return nil
	}
	return &apiBase.Video{
		ID:          video.Id,
		Title:       video.Title,
		Description: video.Description,
		CreatedAt:   video.CreatedAt,
		UpdatedAt:   video.UpdatedAt,
		DeletedAt:   video.DeletedAt,
	}
}

func VideosRPCToVideos(videos []*kitBase.Video) []*apiBase.Video {
	res := make([]*apiBase.Video, 0, len(videos))
	for _, video := range videos {
		res = append(res, VideoRPCToVideo(video))
	}
	return res
}
