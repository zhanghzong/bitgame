package apollo

import (
	"fmt"
	"github.com/zouyx/agollo/v4/storage"
)

type changeListener struct {
}

// 配置被修改
func (c *changeListener) OnChange(changeEvent *storage.ChangeEvent) {
	for key, value := range changeEvent.Changes {
		fmt.Println("change key : ", key, ", value :", value)
	}

	fmt.Println(changeEvent.Namespace)
}

// 新增配置不处理
func (c *changeListener) OnNewestChange(event *storage.FullChangeEvent) {
	//write your code here
}
