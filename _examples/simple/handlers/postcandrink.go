package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"simple/protobuf"

	"github.com/maxstanley/tango/handler"
	"google.golang.org/protobuf/proto"
)

func NewPostCanDrink() handler.Handler {
	return &getCanDrinkEvent{}
}

type getCanDrinkEvent struct {
	protobuf.DrinkEvent

	Country string `path:"country"`
	City    string `path:"city"`
}

func (e *getCanDrinkEvent) UnmarshalJSONBody(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *getCanDrinkEvent) UnmarshalProtoBody(data []byte) error {
	return proto.Unmarshal(data, e)
}

func (e *getCanDrinkEvent) UnmarshalPath(pathMap map[string]string) {
	eventType := reflect.TypeOf(*e)
	if eventType.Kind() != reflect.Struct {
		panic("type should be struct")
	}

	for i := 0; i < eventType.NumField(); i++ {
		field := eventType.Field(i)
		pathValue := field.Tag.Get("path")

		if pathValue == "" {
			continue
		}

		reflect.Indirect(reflect.ValueOf(e)).FieldByName(field.Name).SetString(pathMap[pathValue])
	}
}

func newStringResponse(status int32, message string) *protobuf.StringResponse {
	return &protobuf.StringResponse{
		Status:  status,
		Message: message,
	}
}

func (e *getCanDrinkEvent) Handler() (int, proto.Message) {
	if e.Name == "" {
		return http.StatusBadRequest, newStringResponse(http.StatusBadRequest, "no name found in body")
	}
	if e.Age == "" {
		return http.StatusBadRequest, newStringResponse(http.StatusBadRequest, "no age found in body")
	}
	if e.City == "" {
		return http.StatusBadRequest, newStringResponse(http.StatusBadRequest, "no city parameter found")
	}
	if e.Country == "" {
		return http.StatusBadRequest, newStringResponse(http.StatusBadRequest, "no country parameter found")
	}

	age, err := strconv.Atoi(e.Age)
	if err != nil {
		return http.StatusBadRequest, newStringResponse(http.StatusBadRequest, "age must be a valid number")
	}

	if age > 17 {
		return http.StatusOK, newStringResponse(http.StatusOK, fmt.Sprintf("%s can have alcohol in %s, %s!", e.Name, e.City, e.Country))
	}

	return http.StatusOK, newStringResponse(
		http.StatusOK,
		fmt.Sprintf("%s can have a soft drink!", e.Name),
	)
}
