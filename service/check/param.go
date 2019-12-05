// @Contact:    huaxinrui
// @Time:       2019/9/27 上午11:52

package check

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LTE      = "_lte"
	GTE      = "_gte"
	LBracket = "("
	RBracket = ")"
	NONE     = ""
	PERCENT  = "%"
	JORN_OR  = " or "
	JORN_AND = " and "
)

var ParamParse ParamCheck

func init() {
	tmp := new(param)
	ParamParse = tmp
}

type param struct {
}

func (p *param) addPrefix(str, part string) string {
	return part + str + part
}

func (p *param) ParseParamsToWhereArgs(c *gin.Context, allowed []string, using_blurry bool) (string, []interface{}) {
	var where_str []string
	var args []interface{}
	for _, k := range allowed {
		value := c.DefaultQuery(k, NONE)
		if value != "" {
			if strings.Contains(k, LTE) {
				field := strings.Replace(k, LTE, NONE, 1)
				where_str = append(where_str, fmt.Sprintf("`%s` < ?", field))
				args = append(args, value)
				continue
			}
			if strings.Contains(k, GTE) {
				field := strings.Replace(k, GTE, NONE, 1)
				where_str = append(where_str, fmt.Sprintf("`%s` > ?", field))
				args = append(args, value)
				continue
			}
			// the front must obey the conditions
			multi_query := strings.Split(value, ",")
			if using_blurry {
				if len(multi_query) > 1 {
					var tmp_where_str []string
					var tmp_args []interface{}
					for _, i := range multi_query {
						tmp_where_str = append(tmp_where_str, fmt.Sprintf("`%s` like ?", k))
						tmp_args = append(tmp_args, p.addPrefix(i, PERCENT))
					}
					w := strings.Join(tmp_where_str, JORN_OR)
					where_str = append(where_str, LBracket+w+RBracket)
					args = append(args, tmp_args...)
				} else {
					where_str = append(where_str, fmt.Sprintf("`%s` like ?", k))
					args = append(args, p.addPrefix(value, PERCENT))
				}
			} else {
				if len(multi_query) > 1 {
					var tmp_where_str []string
					var tmp_args []interface{}
					for _, i := range multi_query {
						tmp_where_str = append(tmp_where_str, fmt.Sprintf("`%s` = ?", k))
						tmp_args = append(tmp_args, i)
					}
					w := strings.Join(tmp_where_str, JORN_OR)
					where_str = append(where_str, LBracket+w+RBracket)
					args = append(args, tmp_args...)
				} else {
					where_str = append(where_str, fmt.Sprintf("`%s` = ?", k))
					args = append(args, value)
				}
			}
		}
	}
	log.Println(strings.Join(where_str, JORN_AND), args)
	return strings.Join(where_str, JORN_AND), args
}

func (p *param) GetParamsBlurry(c *gin.Context, allowed []string) (string, []interface{}) {
	return p.ParseParamsToWhereArgs(c, allowed, true)
}

func (p *param) GetParamsSpecific(c *gin.Context, allowed []string) (string, []interface{}) {
	return p.ParseParamsToWhereArgs(c, allowed, false)
}

func (p *param) GetPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := page_size * (page - 1)
	return offset, page_size
}

func (p *param) CheckValue(k string, detail map[string]interface{}, t interface{}) interface{} {
	res, ok := detail[k]
	switch t.(type) {
	case int, float64:
		if !ok || res == nil {
			return float64(0)
		}
	case string:
		if !ok || res == nil {
			return ""
		}

	case time.Time:
		fmt.Println(k, res)
		if !ok || res == nil {
			return time.Now()
		}
	}
	return res
}

func (p *param) CheckLevel(l string) int {
	switch l {
	case "Low":
		return 1
	case "Medium":
		return 2
	case "High":
		return 3
	case "Critical":
		return 4
	}
	return 1
}
