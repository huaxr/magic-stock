// @Contact:    huaxinrui
// @Time:       2019/10/25 下午5:12

package upload

import (
	"fmt"
	"testing"
)

func TestUpload(t *testing.T) {
	GlobalUpload.UploadByFilePath("/Users/huaxinrui/go/src/code.byted.org/byte_security/platform_api/conf/loc.json")
	GlobalUpload.UploadByFilePath("/Users/huaxinrui/go/src/code.byted.org/byte_security/platform_api/conf/tce.json")
}

func TestUploads(t *testing.T) {
	x := GlobalUpload.genHashName("a.yx")
	fmt.Println(x)
}
