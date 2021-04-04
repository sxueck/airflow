package realnode

type UserInfo struct {
	Data struct {
		Email             string      `json:"email"`
		TransferEnable    int64       `json:"transfer_enable"`
		LastLoginAt       interface{} `json:"last_login_at"`
		CreatedAt         string      `json:"created_at"`
		Banned            int         `json:"banned"`
		RemindExpire      string      `json:"remind_expire"`
		RemindTraffic     string      `json:"remind_traffic"`
		ExpiredAt         int         `json:"expired_at"`
		Balance           string      `json:"balance"`
		CommissionBalance int         `json:"commission_balance"`
		PlanID            string      `json:"plan_id"`
		Discount          interface{} `json:"discount"`
		CommissionRate    interface{} `json:"commission_rate"`
		TelegramID        interface{} `json:"telegram_id"`
		AvatarURL         string      `json:"avatar_url"`
	} `json:"data"`
}
