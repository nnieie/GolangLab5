package dal

import "github.com/nnieie/golanglab5/internal/social/dal/db"

func Init() {
	db.InitMySQL()
}
