package call

import (
	"github.com/lgrisa/lib/utils"
)

func AsyncCall(f func()) {
	utils.AsyncCall(f)
}
