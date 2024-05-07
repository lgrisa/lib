package decimal

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func Yuan2Fen(yuan string) (int64, error) {

	if fen, err := Decimal(yuan, 100); err != nil {
		return 0, errors.Wrapf(err, "将人民币元转化成分出错")
	} else {
		return fen, nil
	}
}

func Decimal(s string, multi float64) (int64, error) {

	dec, err := decimal.NewFromString(s)
	if err != nil {
		return 0, errors.Wrapf(err, "Decimal转换出错, %v, %v", s, multi)
	}

	fenDecimal := dec.Mul(decimal.NewFromFloat(multi))

	return fenDecimal.IntPart(), nil
}
