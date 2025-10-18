package dal

import "github.com/nnieie/golanglab5/cmd/interaction/dal/db"

func Init() {
	db.InitMySQL()
}
