package scenc

type Unit struct {
	*Entity
	content string	// 气泡消息
}

func NewUnit() *Unit {
	// 生成一个实体id
	en := NewEntity()
	un := &Unit{
		content: "",
		Entity: en,
	}
	return un
}