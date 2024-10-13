package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
	l  *zap.SugaredLogger
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{
		db: db,
		l:  zap.S(),
	}
}

func (p *Postgres) GetUsers() ([]User, error) {
	var users = []User{}

	p.db.Select("email", "name", "binance_key_enable",
		"usdt_order_limit", "max_open_positions", "move_stl_when_reach_tp", "banned").
		Find(&users)
	return users, nil
}

func (p *Postgres) GetEnableBinanceUsers() ([]User, error) {
	var users = []User{}
	p.db.Where("binance_key_enable = ?", true).Find(&users)
	return users, nil
}

func (p *Postgres) GetUserWithPermission(eid, email string, permisison string) (User, error) {
	var user = User{}
	result := p.db.Preload("Permissions").
		Joins("inner join user_permission up on up.user_id = users.id ").
		Joins("inner join Permissions p on p.id= up.permission_id ").
		Where("p.desc = ?", permisison).
		Where("e_id = ? and email = ?", eid, email).First(&user)
	return user, result.Error
}
