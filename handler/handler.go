package handler

type Handler interface {
	Handler() (int, string)

	UnmarshalBody([]byte) error
	UnmarshalPath(map[string]string)
}
