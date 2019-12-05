// @Contact:    huaxinrui
// @Time:       2019/10/25 下午3:28

package check

import (
	"errors"
	"fmt"
	"strings"

	"magic/stock/utils"

	"golang.org/x/net/html"
)

var (
	Security             SecurityIF
	allowed_suffix       = []string{".pdf", ".jpg", ".jpeg", ".doc", ".docx", ".png"}
	not_allowed_prefixes = []string{"file://", "gopher://", "dict://"}
	sql                  = "'\";%="
)

func init() {
	tmp := new(SecCheck)
	Security = tmp
}

type SecCheck struct {
}

func (s *SecCheck) PreventXSS(content string) string {
	return html.EscapeString(content)
}

func (s *SecCheck) PreventSSRF(url string) error {
	// easy check the illegal protocol
	for _, illegal := range not_allowed_prefixes {
		if strings.HasPrefix(url, illegal) {
			return errors.New("非法请求协议")
		}
	}
	return nil
}

func (s *SecCheck) PreventXXE(url string) error {
	return nil
}

func (s *SecCheck) PreventSystemCommand(exec string) error {
	return nil
}

func (s *SecCheck) PreventFileAnyUpload(suffix string) error {
	if !utils.ContainsString(allowed_suffix, strings.ToLower(suffix)) {
		return errors.New(fmt.Sprintf("%s 后缀不被允许", suffix))
	}
	return nil
}

func (s *SecCheck) PreventSQLI(field ...string) error {
	for _, i := range field {
		valid_filed := strings.ContainsAny(i, sql)
		if valid_filed {
			return errors.New("Field sql injection")
		}
	}
	return nil
}
