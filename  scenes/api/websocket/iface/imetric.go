package iface

type Metric_Typ int

const (
	Metric_Typ_Send Metric_Typ = iota
	Metric_Typ_Rev
)

type MerticFun func(typ Metric_Typ, sessionId uint64, size int)