package unsafeid

func GetUnsafeId(heroId int64) uint32 {

	// slot(8) + incre(40) + heroIdx(4)
	// => incre(24) + slot(8)
	u := uint64(heroId)
	slot := u >> 44
	incre := (u >> 4) & (1<<40 - 1)
	unsafeId := incre<<8 | slot
	return uint32(unsafeId)
}
