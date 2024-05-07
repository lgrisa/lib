package msg

import "github.com/lgrisa/lib/utils/pbutil"

type ErrMsg interface {
	error
	ErrMsg() pbutil.Buffer
}
