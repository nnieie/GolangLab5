package dal

import "github.com/nnieie/golanglab5/internal/video/dal/db"

func Init() {
	db.InitMySQL()
}
