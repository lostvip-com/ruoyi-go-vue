package service

import (
	"errors"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"strings"
	"things/internal/emqx/vo_emqx"
	"things/internal/iot_product/service"
)

type ClientService struct {
}

func (s ClientService) LoginClient(req *vo_emqx.EmqxLoginVO) error {
	innerPrefix := lv_global.Config().GetValueStr("mqtt.client.prefix.inner")
	devicePrefix := lv_global.Config().GetValueStr("mqtt.client.prefix.device")
	if !strings.HasPrefix(req.ClientId, devicePrefix) && !strings.HasPrefix(req.ClientId, devicePrefix) {
		return errors.New("非法的clientId")
	}
	if strings.HasPrefix(req.ClientId, innerPrefix) { //内部客户端暂时不鉴权
		return nil
	}
	//设备端要鉴权（目前只需要clientId 与 password 即可，username暂时不使用 ）
	arrClientId := strings.Split(req.ClientId, "#")
	deviceId := cast.ToInt(arrClientId[0])
	//获取设备鉴权信息
	dev, err := service.GetDeviceService().FindById(deviceId)
	lv_err.HasErrAndPanic(err)
	if dev.Secret != req.Password {
		return errors.New("设备授权码错误！")
	}
	lv_log.Infof("====> loginClient ok : %v", dev)
	return nil
}

var clientService *ClientService

func GetClientService() *ClientService {
	if clientService == nil {
		clientService = &ClientService{}
	}
	return clientService
}
