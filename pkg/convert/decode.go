package convert

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToTimeHookFunc .
func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) && t != reflect.TypeOf(timestamppb.Timestamp{}) {
			return data, nil
		}

		// timestamppb.Timestamp
		if t == reflect.TypeOf(timestamppb.Timestamp{}) {
			if f == reflect.TypeOf(time.Time{}) {
				v := data.(time.Time)
				return ptypes.TimestampProto(v)
			} else if f == reflect.TypeOf(&time.Time{}) {
				v := data.(*time.Time)
				return ptypes.TimestampProto(*v)
			} else if f == reflect.TypeOf(map[string]interface{}{}) {
				// @todo 检查time.time转为空map的问题，暂时设为当前时间
				v := time.Now()
				return ptypes.TimestampProto(v)
			}
		}

		// timestamppb.Timestamp
		if t == reflect.TypeOf(time.Time{}) || t == reflect.TypeOf(&time.Time{}) {
			if f == reflect.TypeOf("") {
				return cast.ToTime(data), nil
			}
		}

		return data, nil
	}
}

// Decode 类型转换
func Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(ToTimeHookFunc()),
		Result:     result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
