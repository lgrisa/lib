package u64

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestDivide(t *testing.T) {
	RegisterTestingT(t)

	Ω(DivideTimes(0, 0)).Should(BeEquivalentTo(0))
	Ω(DivideTimes(1, 0)).Should(BeEquivalentTo(0))
	Ω(DivideTimes(0, 1)).Should(BeEquivalentTo(0))

	Ω(DivideTimes(1, 1)).Should(BeEquivalentTo(1))
	Ω(DivideTimes(5, 6)).Should(BeEquivalentTo(1))
	Ω(DivideTimes(6, 6)).Should(BeEquivalentTo(1))
	Ω(DivideTimes(7, 6)).Should(BeEquivalentTo(2))
}

func TestGetDividePlan(t *testing.T) {
	RegisterTestingT(t)

	datas := []struct {
		totalAmount uint64
		amountArray []uint64
		ownArray    []uint64
		ok          bool
		plan        []uint64
	}{
		// 数据错误
		{100, []uint64{1}, []uint64{}, false, nil},
		{100, []uint64{1}, []uint64{1, 2}, false, nil},

		{9901, []uint64{100, 1000, 10000, 30000}, []uint64{9, 9, 0, 0}, false, nil},
		{9900, []uint64{100, 1000, 10000, 30000}, []uint64{9, 9, 0, 0}, true, []uint64{9, 9, 0, 0}},
		{9899, []uint64{100, 1000, 10000, 30000}, []uint64{9, 9, 0, 0}, true, []uint64{9, 9, 0, 0}},

		{12501, []uint64{100, 1000, 10000, 30000}, []uint64{5, 2, 1, 0}, false, nil},
		{12500, []uint64{100, 1000, 10000, 30000}, []uint64{5, 2, 1, 0}, true, []uint64{5, 2, 1, 0}},
		{10500, []uint64{100, 1000, 10000, 30000}, []uint64{5, 2, 1, 0}, true, []uint64{5, 0, 1, 0}},
		{11500, []uint64{100, 1000, 10000, 30000}, []uint64{5, 2, 1, 0}, true, []uint64{5, 1, 1, 0}},
		{11600, []uint64{100, 1000, 10000, 30000}, []uint64{5, 2, 1, 0}, true, []uint64{0, 2, 1, 0}},

		{99, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{1, 0, 0, 0}},
		{199, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{2, 0, 0, 0}},
		{299, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{0, 1, 0, 0}},

		{1099, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{1, 1, 0, 0}},
		{1199, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{2, 1, 0, 0}},
		{1299, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{0, 2, 0, 0}},

		{2099, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{1, 2, 0, 0}},
		{2199, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{2, 2, 0, 0}},
		{2299, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{0, 0, 1, 0}},

		{82200, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, true, []uint64{2, 2, 2, 2}},
		{82201, []uint64{100, 1000, 10000, 30000}, []uint64{2, 2, 2, 2}, false, nil},
	}

	for _, d := range datas {
		ok, plan := GetDividePlan(d.totalAmount, d.amountArray, d.ownArray)
		Ω(ok).Should(Equal(d.ok))
		Ω(plan).Should(Equal(d.plan))
	}
}
