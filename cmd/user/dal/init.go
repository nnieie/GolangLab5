package dal

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/cache"
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
)

func Init() {
	db.InitMySQL()
	cache.InitRedis()
}
