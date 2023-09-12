package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"

	"go.uber.org/atomic"

	"x-ui/logger"
	"x-ui/xray"
)

type XrayService struct {
	inboundService *InboundService
	settingService *SettingService

	p                 *xray.Process
	lock              sync.Mutex
	isNeedXrayRestart atomic.Bool
	result            string
}

func NewXrayService(inboundService *InboundService, settingService *SettingService) *XrayService {
	return &XrayService{
		inboundService: inboundService,
		settingService: settingService,
	}
}

func (s *XrayService) GetXrayErr() error {
	if s.p == nil {
		return nil
	}
	return s.p.GetErr()
}

func (s *XrayService) GetXrayResult() string {
	if s.result != "" || !s.IsXrayRunning() {
		return s.result
	}
	s.result = s.p.GetResult()
	return s.result
}

func (s *XrayService) IsXrayRunning() bool {
	return s.p != nil && s.p.IsRunning()
}

func (s *XrayService) GetXrayConfig() (*xray.Config, error) {
	templateConfig, err := s.settingService.GetXrayConfigTemplate()
	if err != nil {
		return nil, err
	}

	xrayConfig := &xray.Config{}
	err = json.Unmarshal([]byte(templateConfig), xrayConfig)
	if err != nil {
		return nil, err
	}

	inboundConfigs, err := s.getAllInboundConfigs()
	if err != nil {
		return nil, err
	}

	xrayConfig.InboundConfigs = inboundConfigs

	return xrayConfig, nil
}

func (s *XrayService) getAllInboundConfigs() ([]xray.InboundConfig, error) {
	filePath := "/path/to/config.json" // 替换为实际的配置文件路径

	configBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config xray.Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return config.InboundConfigs, nil
}

func (s *XrayService) RestartXray(isForce bool) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	logger.Debug("restart xray, force:", isForce)

	xrayConfig, err := s.GetXrayConfig()
	if err != nil {
		return err
	}

	if s.p != nil && s.IsXrayRunning() {
		if !isForce && s.p.GetConfig().Equals(xrayConfig) {
			logger.Debug("not need to restart xray")
			return nil
		}
		s.p.Stop()
	}

	s.p = xray.NewProcess(xrayConfig)
	s.result = ""
	return s.p.Start()
}

func (s *XrayService) StopXray() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	logger.Debug("stop xray")
	if s.IsXrayRunning() {
		return s.p.Stop()
	}
	return errors.New("xray is not running")
}

func (s *XrayService) SetToNeedRestart() {
	s.isNeedXrayRestart.Store(true)
}

func (s *XrayService) IsNeedRestartAndSetFalse() bool {
	return s.isNeedXrayRestart.CAS(true, false)
}

func (s *XrayService) GetXrayTraffic() ([]*xray.Traffic, error) {
	if !s.IsXrayRunning() {
		return nil, errors.New("xray is not running")
	}
	return s.p.GetTraffic(true)
}
