package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_get_IP(t *testing.T) {
	if ip, err := utils.GetLocalIP(); err != nil {
		t.Log(err)
	} else {
		id := utils.GetIDbyIP(ip)
		t.Log(ip)
		t.Log(id)
	}

}

func Test_removeHTMLTags(t *testing.T) {
	desc := `+9 <lol-uikit-tooltipped-keyword key='LinkTooltip_Description_Adaptive'><font color='#48C4B7'>适应之力</font></lol-uikit-tooltipped-keyword>`
	conf := global.ChaConf.AmqpConf
	t.Log(conf.User)
	t.Log(utils.RemoveHTMLTags(desc))
}
