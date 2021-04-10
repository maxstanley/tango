package handler

import "google.golang.org/protobuf/proto"

type Handler interface {
	Handler() (int, proto.Message)

	UnmarshalJSONBody([]byte) error
	UnmarshalProtoBody([]byte) error
	UnmarshalPath(map[string]string)
}
