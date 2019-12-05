// @Time:       2019/11/15 下午4:36

package check

import (
	"code.byted.org/byte_security/platform_api/utils"
)

type Resource struct {
	Result map[string][]string
}

func (p *Resource) HasPerm(typ , perm string) bool {
	res, ok := p.Result[typ]
	if !ok {
		return false
	}
	if utils.ContainsString(res, "all") {
		return true
	}
	return utils.ContainsString(res, perm)
}


