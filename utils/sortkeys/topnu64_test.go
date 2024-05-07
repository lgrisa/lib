package sortkeys

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestU64TopN(t *testing.T) {
	RegisterTestingT(t)

	topN := NewU64TopN(5)

	topN.Add(1)
	topN.Add(2)
	topN.Add(4)
	topN.Add(3)

	Ω(topN.Array()).Should(BeEquivalentTo([]uint64{
		1,
		2,
		4,
		3,
	}))

	Ω(topN.SortAsc()).Should(BeEquivalentTo([]uint64{
		1,
		2,
		3,
		4,
	}))

	Ω(topN.SortDesc()).Should(BeEquivalentTo([]uint64{
		4,
		3,
		2,
		1,
	}))

	topN.Add(6)

	Ω(topN.Array()).Should(BeEquivalentTo([]uint64{
		1,
		2,
		4,
		3,
		6,
	}))

	topN.Add(5)

	Ω(topN.Array()).Should(BeEquivalentTo([]uint64{
		5,
		2,
		4,
		3,
		6,
	}))

	topN.Add(7)

	Ω(topN.Array()).Should(BeEquivalentTo([]uint64{
		5,
		7,
		4,
		3,
		6,
	}))
}
