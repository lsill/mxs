package iface

type Endpoint interface {
	StartReading()
	// 终止
	Terminate()
	Output() chan []byte
	Send([]byte)bool
}

func PipeEndPoints(e1, e2 Endpoint) {
	e1.StartReading()
	e2.StartReading()
	defer e1.Terminate()
	defer e2.Terminate()
	for {
		select {
		case msgOne, ok := <- e1.Output():
			if !ok || !e2.Send(msgOne) {
				return
			}
		case msgTwo, ok := <- e2.Output():
			if !ok || !e1.Send(msgTwo) {
				return
			}
		}
	}
}
