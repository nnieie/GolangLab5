package service

import "github.com/nnieie/golanglab5/cmd/user/dal/db"

func (s *UserService) SearchUserIdsByName(pattern string, page, pageSize int64) ([]string, error) {
	userIds, err := db.SearchUserIdsByName(s.ctx, pattern, page, pageSize)
	if err != nil {
		return nil, err
	}
	return userIds, nil
}
