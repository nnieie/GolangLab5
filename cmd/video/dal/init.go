package dal

import "github.com/nnieie/golanglab5/cmd/video/dal/db"

func Init() {
	db.InitMySQL()
}
