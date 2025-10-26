package dal

import (
	"github.com/nnieie/golanglab5/cmd/chat/dal/cache"
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
)

var syncWorker *cache.SyncWorker

func Init() {
	db.InitMySQL()
	cache.InitRedis()

	// 启动同步工作器
	syncWorker = cache.NewSyncWorker()
	go syncWorker.Start()
}
