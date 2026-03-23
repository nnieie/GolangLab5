package dal

import (
	"github.com/nnieie/golanglab5/internal/user/dal/cache"
	"github.com/nnieie/golanglab5/internal/user/dal/db"
)

func Init() {
	db.InitMySQL()
	cache.InitRedis()
}
