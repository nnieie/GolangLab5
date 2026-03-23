package dal

import (
	"github.com/nnieie/golanglab5/internal/interaction/dal/cache"
	"github.com/nnieie/golanglab5/internal/interaction/dal/db"
)

func Init() {
	db.InitMySQL()
	cache.InitRedis()
}
