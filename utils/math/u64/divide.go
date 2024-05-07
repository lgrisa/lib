package u64

func DivideTimes(x, y uint64) uint64 {
	if x <= 0 || y <= 0 {
		return 0
	}

	return (x + y - 1) / y
}

// 找最佳使用方案，比如武将升级，需要10001经验，存在4种经验丹，分别是100，500，1000，2000经验值的
// 寻找最优使用方案
func GetDividePlan(totalAmount uint64, amountArray, ownCountArray []uint64) (ok bool, plan []uint64) {

	// amountArray 是从小到大的值
	n := len(amountArray)
	if n != len(ownCountArray) {
		return false, nil
	}
	plan = make([]uint64, n)

	if totalAmount <= 0 {
		return true, plan
	}

	// 从大到小开始扣，如果刚刚好扣完，那么就退出，如果超出，则少扣一个，找小的来补
	// 找完一圈都不够，那么开始第二轮，从小到大的找，找第一个有富余的，补进来，应该就超出了
	// 这个时候还要再看下，是否有小的溢出，如果溢出，那么可以移除小的
	var curAmt uint64
	hasEnoughCount := false
	for i := 0; i < n; i++ {
		idx := n - 1 - i
		amt := amountArray[idx]
		if amt <= 0 {
			continue
		}

		remain := Sub(totalAmount, curAmt)
		needCount := remain / amt
		ownCount := ownCountArray[idx]
		useCount := Min(needCount, ownCount)

		plan[idx] = useCount
		if useCount >= needCount && remain%amt == 0 {
			// 如果这次扣的数量足够，并且刚好扣完，那么扣除结束
			return true, plan
		}

		curAmt += amt * useCount

		if ownCount > needCount {
			// 如果拥有的个数超过，所需的个数，那么可以肯定数量是足够的
			hasEnoughCount = true
		}
	}

	if !hasEnoughCount {
		// 没有足够的数量扣
		return false, nil
	}

	if curAmt >= totalAmount {
		// 防御性
		return false, nil
	}

	// 第二轮，从小到大的找，找第一个有富余的，补进来，应该就超出了
	for i := 0; i < n; i++ {
		amt := amountArray[i]
		if amt <= 0 {
			continue
		}

		ownCount := ownCountArray[i]
		if plan[i] < ownCount {
			// 补一个进来
			plan[i]++
			curAmt += amt
			break
		}
	}

	if curAmt <= totalAmount {
		// 防御性
		return false, nil
	}

	// 第三轮，是否有小的溢出，如果溢出，那么可以移除小的
	overflow := Sub(curAmt, totalAmount)
	for i := 0; i < n; i++ {
		amt := amountArray[i]
		if amt <= 0 {
			continue
		}

		overflowCount := overflow / amt
		if overflowCount <= 0 {
			break
		}

		backCount := Min(overflowCount, plan[i])
		if backCount > 0 {
			plan[i] = Sub(plan[i], backCount)
			overflow = Sub(overflow, amt*backCount)
		}
	}

	return true, plan
}
