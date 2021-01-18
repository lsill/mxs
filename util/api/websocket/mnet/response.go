package mnet

type Response struct {
	data []byte
	typ int32
}

func NewResponse(typ int32, data []byte) *Response{
	res := &Response{
		data: nil,
		typ:  0,
	}
	res.SetData(data)
	res.SetPkTyp(typ)
	return res
}

func(s *Response) SetData(data []byte) {
	s.data = data
}

func (s *Response) SetPkTyp(typ int32) {
	s.typ = typ
}