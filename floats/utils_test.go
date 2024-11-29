package floats

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/puellanivis/math/bits"
)

func fromAny[SPEC spec[D], D datum](t *testing.T, val any) (D, error) {
	t.Helper()

	var z D

	switch val := val.(type) {
	case uint16:
		return convert[binary16, SPEC](val, RoundTiesToEven{}), nil
	case uint32:
		return convert[binary32, SPEC](val, RoundTiesToEven{}), nil
	case uint64:
		return convert[binary64, SPEC](val, RoundTiesToEven{}), nil
	case bits.Uint128:
		return convert[binary128, SPEC](val, RoundTiesToEven{}), nil
	case float32:
		return convert[binary32, SPEC](math.Float32bits(val), RoundTiesToEven{}), nil
	case float64:
		return convert[binary64, SPEC](math.Float64bits(val), RoundTiesToEven{}), nil

	case string:
		var spec SPEC

		f := new(big.Float).SetPrec(uint(spec.width()))
		if _, _, err := f.Parse(val, 0); err != nil {
			return z, err
		}

		return fromBigFloat[SPEC](f, RoundTiesToEven{}), nil
	}

	return z, fmt.Errorf("unsupported type: %T", val)
}
