package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/maxstanley/tango/handler"
)

func NewPostCanDrink() handler.Handler {
	return &getCanDrinkEvent{}
}

type getCanDrinkEvent struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Country string `path:"country"`
	City    string `path:"city"`
}

func (e *getCanDrinkEvent) UnmarshalBody(data []byte) error {
	return json.Unmarshal(data, e)
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

func (e *getCanDrinkEvent) Handler() (int, string) {
	if e.Name == "" {
		return http.StatusBadRequest, "no name found in body"
	}
	if e.Age == "" {
		return http.StatusBadRequest, "no age found in body"
	}
	if e.City == "" {
		return http.StatusBadRequest, "no city parameter found"
	}
	if e.Country == "" {
		return http.StatusBadRequest, "no country parameter found"
	}

	age, err := strconv.Atoi(e.Age)
	if err != nil {
		return http.StatusBadRequest, "age must be a valid number"
	}

	if age > 17 {
		return http.StatusOK, fmt.Sprintf("%s can have alcohol in %s, %s!", e.Name, e.City, e.Country)
	}

	return http.StatusOK, fmt.Sprintf("%s can have a soft drink!", e.Name)
}
