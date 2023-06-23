package runtime

import (
	"fmt"
	"math"
)

var NumberType = &DataType{Name: "Number"}

type NumberImpl struct {
	Value float64
}

func (n NumberImpl) Repr() string {
	if math.Mod(n.Value, 1.0) == 0 {
		return fmt.Sprintf("%.0f", n.Value)
	} else {
		return fmt.Sprintf("%f", n.Value)
	}
}
