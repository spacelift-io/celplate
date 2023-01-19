package macros

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

var anyListType = reflect.TypeOf([]any{})

func GetJoinFunction() cel.EnvOption {
	joinFunction := cel.Function(
		"join",
		cel.MemberOverload("join_list", []*cel.Type{cel.AnyType, cel.StringType}, cel.StringType,
			cel.BinaryBinding(func(list, delim ref.Val) ref.Val {
				d := delim.(types.String)

				if anyList, err := list.ConvertToNative(anyListType); err == nil {
					var strList []string
					for _, a := range anyList.([]any) {
						strList = append(strList, fmt.Sprint(a))
					}
					return types.String(strings.Join(strList, string(d)))
				}

				return types.NewErr("unsupported type for join_list: %v", list.Type())
			})))

	return joinFunction
}
