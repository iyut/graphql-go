package helper

import (
	"reflect"
	"strconv"

	"github.com/graph-gophers/graphql-go"
)

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func IntToGraphqlID(id int64) graphql.ID {

	idStr := strconv.FormatInt(id, 10)
	return graphql.ID(idStr)

}
