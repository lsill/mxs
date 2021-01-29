package constcode

const (
	PingHeart = iota+0        // 连接心跳
	PositionMine // 同步周围相关位置信息
	PositionOther  // 同步给他人我的信息
	LoginOn		// 登陆
)