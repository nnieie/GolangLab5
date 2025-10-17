package dal

import "github.com/nnieie/golanglab5/cmd/social/dal/db"

func Init() {
	db.InitMySQL()
}
