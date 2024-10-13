package common

type BinanceAccount struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type CoinInfo struct {
	Coin              string    `json:"coin"`
	DepositAllEnable  bool      `json:"depositAllEnable"`
	Free              string    `json:"free"`
	Freeze            string    `json:"freeze"`
	Ipoable           string    `json:"ipoable"`
	Ipoing            string    `json:"ipoing"`
	IsLegalMoney      bool      `json:"isLegalMoney"`
	Locked            string    `json:"locked"`
	Name              string    `json:"name"`
	NetworkList       []Network `json:"networkList"`
	Storage           string    `json:"storage"`
	Trading           bool      `json:"trading"`
	WithdrawAllEnable bool      `json:"withdrawAllEnable"`
	Withdrawing       string    `json:"withdrawing"`
}

type Network struct {
	AddressRegex            string `json:"addressRegex"`
	Coin                    string `json:"coin"`
	DepositDesc             string `json:"depositDesc,omitempty"` // 仅在充值关闭时返回
	DepositEnable           bool   `json:"depositEnable"`
	IsDefault               bool   `json:"isDefault"`
	MemoRegex               string `json:"memoRegex"`
	MinConfirm              int    `json:"minConfirm"` // 上账所需的最小确认数
	Name                    string `json:"name"`
	Network                 string `json:"network"`
	ResetAddressStatus      bool   `json:"resetAddressStatus"`
	SpecialTips             string `json:"specialTips"`
	UnLockConfirm           int    `json:"unLockConfirm"`          // 解锁需要的确认数
	WithdrawDesc            string `json:"withdrawDesc,omitempty"` // 仅在提现关闭时返回
	WithdrawEnable          bool   `json:"withdrawEnable"`
	WithdrawFee             string `json:"withdrawFee"`
	WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`
	WithdrawMax             string `json:"withdrawMax"`
	WithdrawMin             string `json:"withdrawMin"`
	SameAddress             bool   `json:"sameAddress"` // 是否需要memo
	ContractAddress         string `json:"contractAddress"`
}
