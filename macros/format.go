package macros

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func GetFormatFunction() cel.EnvOption {
	formatFn := cel.Function(
		"format",
		cel.MemberOverload("timestamp_format", []*cel.Type{cel.TimestampType, cel.StringType}, cel.StringType,
			cel.BinaryBinding(func(timestamp, format ref.Val) ref.Val {
				ts := timestamp.(types.Timestamp)
				dateFmt := format.(types.String)
				return types.String(ts.Time.Format(string(dateFmt)))
			})))

	return formatFn
}
