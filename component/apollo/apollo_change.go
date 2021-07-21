package apollo

import (
	"github.com/sirupsen/logrus"
	"github.com/zouyx/agollo/v4/storage"
)

type changeListener struct {
}

// OnChange 配置被修改
func (c *changeListener) OnChange(changeEvent *storage.ChangeEvent) {
	for key, value := range changeEvent.Changes {
		logrus.Infof("Apollo 配置文件发生变化, namespace:%+v, key:%s, value:%+v", changeEvent.Namespace, key, value)
	}
}

// OnNewestChange 新增配置不处理
func (c *changeListener) OnNewestChange(event *storage.FullChangeEvent) {
}
