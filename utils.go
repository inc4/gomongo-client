package mongoclient

import (
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func isStructOrPtrToStruct(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func isMongoOperator(v any) bool {
	m, ok := v.(bson.M)
	if !ok {
		return false
	}
	for k := range m {
		if len(k) > 0 && k[0] == '$' {
			return true
		}
	}
	return false
}
