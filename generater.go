// Title：组件合成器
//
// Description:
//
// Author:black
//
// Createtime:2013-08-08 09:52
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-08 09:52 black 创建文档
package lessgo

import (
	"encoding/xml"
	"net/http"
)

//运行期间用来缓存组件内存对象的容器
var runtimeComponentContain = make(map[string]interface{})

//根据用户定义的view文件和数据库内容生成数据流
func Generate(viewContent []byte, terminal, packageName string, r *http.Request) []byte {

	var viewport Viewport

	err := xml.Unmarshal(viewContent, &viewport)

	if err != nil {
		Log.Error(err)
		return []byte{}
	}

	return viewport.GenerateViewport(terminal, packageName, r)
}
