package big

import (
	"math/big"
)

type Accuracy = big.Accuracy

const (
	Below Accuracy = -1
	Exact Accuracy = 0
	Above Accuracy = +1
)
