package db

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primarykey,autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Email   string `gorm:"index:idx_email,unique" json:"email"`
	EID     string `gorm:"index:idx_email,unique" json:"eid"`
	Name    string `json:"name"`
	Picture string `json:"picture"`

	BinanceApiKey      string  `json:"binance_api_key"`
	BinanceApiSecret   string  `json:"binance_api_secret"`
	BinanceKeyEnable   bool    `json:"binance_key_enable"`
	UsdtOrderLimit     float64 `json:"usdt_order_limit"`
	MaxOpenPositions   float64 `json:"max_open_positions"` // TODO: convert it to uint
	CrossMargin        bool    `json:"cross_margin"`
	MoveStlWhenReachTp bool    `json:"move_stl_when_reach_tp"`
	Banned             bool    `json:"banned"`

	TelegramGroup []TelegramGroup `gorm:"many2many:user_telegram_group;"`
}

type TelegramGroup struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ID         int64  `gorm:"primarykey,unique" json:"id"`
	GroupTitle string `json:"group_title"`
}
