package dal

import "github.com/nnieie/golanglab5/cmd/chat/dal/db"

func Init() {
	db.InitMySQL()
}
