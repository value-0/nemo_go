package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/hanc00l/nemo_go/pkg/comm"
	"github.com/hanc00l/nemo_go/pkg/conf"
	"github.com/hanc00l/nemo_go/pkg/logging"
	minichatConfig "github.com/hanc00l/nemo_go/pkg/minichat/config"
	"github.com/hanc00l/nemo_go/pkg/notify"
	"github.com/hanc00l/nemo_go/pkg/task/ampq"
	"github.com/hanc00l/nemo_go/pkg/task/custom"
	"github.com/hanc00l/nemo_go/pkg/task/domainscan"
	"github.com/hanc00l/nemo_go/pkg/task/onlineapi"
	"github.com/hanc00l/nemo_go/pkg/task/portscan"
	"github.com/hanc00l/nemo_go/pkg/utils"
	"github.com/remeh/sizedwaitgroup"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ConfigController struct {
	BaseController
}

const (
	HoneyPot               string = "honeypot"
	IPLocation             string = "iplocation"
	IPLocationB            string = "iplocationB"
	IPLocationC            string = "iplocationC"
	Service                string = "service"
	Xray                   string = "xray"
	XrayConfig             string = "config.xray"
	XrayPlugin             string = "plugin.xray"
	XrayModule             string = "module.xray"
	BlackDomain            string = "black_domain"
	BlackIP                string = "black_ip"
	TaskWorkspace          string = "task_workspace"
	FOFAFilterKeyword      string = "fofa_filter_keyword"
	FOFAFilterKeywordLocal string = "fofa_filter_keyword_local"
)

type DefaultConfig struct {
	//portscan
	CmdBin string `json:"cmdbin" form:"cmdbin"`
	Port   string `json:"port" form:"port"`
	Rate   int    `json:"rate" form:"rate"`
	Tech   string `json:"tech" form:"tech"`
	IsPing bool   `json:"ping" form:"ping"`
	//task
	IpSliceNumber   int    `json:"ipslicenumber" form:"ipslicenumber"`
	PortSliceNumber int    `json:"portslicenumber" form:"portslicenumber"`
	Version         string `json:"version" form:"version"`
	TaskWorkspace   string `json:"taskworkspace" form:"taskworkspace"`
	//fingerprint
	IsHttpx          bool `json:"httpx" form:"httpx"`
	IsScreenshot     bool `json:"screenshot" form:"screenshot"`
	IsFingerprintHub bool `json:"fingerprinthub" form:"fingerprinthub"`
	IsIconHash       bool `json:"iconhash" form:"iconhash"`
	IsFingerprintx   bool `json:"fingerprintx" form:"fingerprintx"`
	// onlineapi
	IsFofa           bool   `json:"fofa" form:"fofa"`
	IsQuake          bool   `json:"quake" form:"quake"`
	IsHunter         bool   `json:"hunter" form:"hunter"`
	ServerChanToken  string `json:"serverchan" form:"serverchan"`
	DingTalkToken    string `json:"dingtalk" form:"dingtalk"`
	FeishuToken      string `json:"feishu" form:"feishu"`
	FofaToken        string `json:"fofatoken" form:"fofatoken"`
	HunterToken      string `json:"huntertoken" form:"huntertoken"`
	QuakeToken       string `json:"quaketoken" form:"quaketoken"`
	ChinazToken      string `json:"chinaztoken" form:"chinaztoken"`
	SearchPageSize   int    `json:"pagesize" form:"pagesize"`
	SearchLimitCount int    `json:"limitcount" form:"limitcount"`
	// domainscan
	Wordlist           string `json:"wordlist" form:"wordlist"`
	IsSubDomainFinder  bool   `json:"subfinder" form:"subfinder"`
	IsSubDomainBrute   bool   `json:"subdomainbrute" form:"subdomainbrute"`
	IsSubDomainCrawler bool   `json:"subdomaincrawler" form:"subdomaincrawler"`
	IsIgnoreCDN        bool   `json:"ignorecdn" form:"ignorecdn"`
	IsIgnoreOutofChina bool   `json:"ignoreoutofchina" form:"ignoreoutofchina"`
	IsPortscan         bool   `json:"portscan" form:"portscan"`
	IsWhois            bool   `json:"whois" form:"whois"`
	IsICP              bool   `json:"icp" form:"icp"`
	// proxy
	ProxyList string `json:"proxyList" form:"proxyList"`
	// wiki:feishu
	FeishuAppId        string `json:"feishuappid" form:"feishuappid"`
	FeishuAppSecret    string `json:"feishusecret" form:"feishuappsecret"`
	FeishuRefreshToken string `json:"feishurefreshtoken" form:"feishurefreshtoken"`
	// filter
	MaxPortPerIP   int    `json:"maxportperip" form:"maxportperip"`
	MaxDomainPerIP int    `json:"maxdomainperip" form:"maxdomainperip"`
	TitleFilter    string `json:"title" form:"title"`
	// minichat
	Anonymous         bool `json:"anonymous" form:"anonymous"`
	IsNotDelFileDir   bool `json:"notdelfiledir" form:"notdelfiledir"`
	LoadHistory       bool `json:"loadhistory" form:"loadhistory"`
	MaxHistoryMessage int  `json:"maxhistorymessage" form:"maxhistorymessage"`
}

func (c *ConfigController) IndexServerAction() {
	c.Layout = "base.html"
	c.TplName = "config-server.html"
}

func (c *ConfigController) IndexWorkerAction() {
	c.Layout = "base.html"
	c.TplName = "config-worker.html"
}

func (c *ConfigController) CustomAction() {
	c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, true)

	c.Layout = "base.html"
	c.TplName = "custom.html"
}

// LoadDefaultConfigAction 获取默认的worker扫描使用的配置参数
func (c *ConfigController) LoadDefaultConfigAction() {
	defer c.ServeJSON()

	err := conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	portscan := conf.GlobalWorkerConfig().Portscan
	fingerprint := conf.GlobalWorkerConfig().Fingerprint
	onlineAPI := conf.GlobalWorkerConfig().OnlineAPI
	domainscan := conf.GlobalWorkerConfig().Domainscan

	data := DefaultConfig{
		CmdBin: portscan.Cmdbin,
		Port:   portscan.Port,
		Rate:   portscan.Rate,
		Tech:   portscan.Tech,
		IsPing: portscan.IsPing,
		//
		IsHttpx:          fingerprint.IsHttpx,
		IsScreenshot:     fingerprint.IsScreenshot,
		IsFingerprintHub: fingerprint.IsFingerprintHub,
		IsIconHash:       fingerprint.IsIconHash,
		IsFingerprintx:   fingerprint.IsFingerprintx,
		//
		IsSubDomainFinder:  domainscan.IsSubDomainFinder,
		IsSubDomainBrute:   domainscan.IsSubDomainBrute,
		IsSubDomainCrawler: domainscan.IsSubdomainCrawler,
		IsIgnoreCDN:        domainscan.IsIgnoreCDN,
		IsIgnoreOutofChina: domainscan.IsIgnoreOutofChina,
		IsPortscan:         domainscan.IsPortScan,
		IsWhois:            domainscan.IsWhois,
		IsICP:              domainscan.IsICP,
		//
		IsFofa:   onlineAPI.IsFofa,
		IsHunter: onlineAPI.IsHunter,
		IsQuake:  onlineAPI.IsQuake,
	}

	if fileContent, err1 := os.ReadFile(filepath.Join(conf.GetRootPath(), "version.txt")); err1 == nil {
		data.Version = strings.TrimSpace(string(fileContent))
	}
	c.Data["json"] = data
}

// LoadWorkerConfigAction 获取worker的参数
func (c *ConfigController) LoadWorkerConfigAction() {
	if !c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) {
		c.LoadDefaultConfigAction()
		return
	}
	defer c.ServeJSON()

	err := conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	portscan := conf.GlobalWorkerConfig().Portscan
	fingerprint := conf.GlobalWorkerConfig().Fingerprint
	apiConfig := conf.GlobalWorkerConfig().API
	onlineAPI := conf.GlobalWorkerConfig().OnlineAPI
	domainscan := conf.GlobalWorkerConfig().Domainscan
	proxy := conf.GlobalWorkerConfig().Proxy
	filter := conf.GlobalWorkerConfig().Filter

	data := DefaultConfig{
		CmdBin: portscan.Cmdbin,
		Port:   portscan.Port,
		Rate:   portscan.Rate,
		Tech:   portscan.Tech,
		IsPing: portscan.IsPing,
		//
		IsHttpx:          fingerprint.IsHttpx,
		IsScreenshot:     fingerprint.IsScreenshot,
		IsFingerprintHub: fingerprint.IsFingerprintHub,
		IsIconHash:       fingerprint.IsIconHash,
		IsFingerprintx:   fingerprint.IsFingerprintx,
		//
		FofaToken:   apiConfig.Fofa.Key,
		HunterToken: apiConfig.Hunter.Key,
		QuakeToken:  apiConfig.Quake.Key,
		ChinazToken: apiConfig.ICP.Key,
		//
		Wordlist:           domainscan.Wordlist,
		IsSubDomainFinder:  domainscan.IsSubDomainFinder,
		IsSubDomainBrute:   domainscan.IsSubDomainBrute,
		IsSubDomainCrawler: domainscan.IsSubdomainCrawler,
		IsIgnoreCDN:        domainscan.IsIgnoreCDN,
		IsIgnoreOutofChina: domainscan.IsIgnoreOutofChina,
		IsPortscan:         domainscan.IsPortScan,
		IsWhois:            domainscan.IsWhois,
		IsICP:              domainscan.IsICP,
		//onlineAPI:
		IsFofa:           onlineAPI.IsFofa,
		IsHunter:         onlineAPI.IsHunter,
		IsQuake:          onlineAPI.IsQuake,
		SearchPageSize:   apiConfig.SearchPageSize,
		SearchLimitCount: apiConfig.SearchLimitCount,
		//
		MaxPortPerIP:   filter.MaxPortPerIp,
		MaxDomainPerIP: filter.MaxDomainPerIp,
		TitleFilter:    filter.Title,
		//
		ProxyList: strings.Join(proxy.Host, "\n"),
	}

	if fileContent, err1 := os.ReadFile(filepath.Join(conf.GetRootPath(), "version.txt")); err1 == nil {
		data.Version = strings.TrimSpace(string(fileContent))
	}
	c.Data["json"] = data
}

// LoadServerConfigAction 获取Server的管理的参数
func (c *ConfigController) LoadServerConfigAction() {
	if !c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) {
		return
	}
	defer c.ServeJSON()

	err := conf.GlobalServerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	task := conf.GlobalServerConfig().Task
	notifyToken := conf.GlobalServerConfig().Notify
	feishu := conf.GlobalServerConfig().Wiki.Feishu

	data := DefaultConfig{
		IpSliceNumber:   task.IpSliceNumber,
		PortSliceNumber: task.PortSliceNumber,
		//
		ServerChanToken: notifyToken["serverchan"].Token,
		DingTalkToken:   notifyToken["dingtalk"].Token,
		FeishuToken:     notifyToken["feishu"].Token,
		//
		FeishuAppId:        feishu.AppId,
		FeishuAppSecret:    feishu.AppSecret,
		FeishuRefreshToken: feishu.UserAccessRefreshToken,
		//
		Anonymous:         minichatConfig.EnableAnonymous,
		IsNotDelFileDir:   minichatConfig.IsNotDelFileDir,
		LoadHistory:       minichatConfig.LoadHistory,
		MaxHistoryMessage: minichatConfig.MaxHistoryMessage,
	}
	if fileContent, err1 := os.ReadFile(filepath.Join(conf.GetRootPath(), "version.txt")); err1 == nil {
		data.Version = strings.TrimSpace(string(fileContent))
	}
	c.Data["json"] = data
}

// ChangePasswordAction 修改密码
func (c *ConfigController) ChangePasswordAction() {
	defer c.ServeJSON()

	op := c.GetString("oldpass", "")
	np := c.GetString("newpass", "")
	if op == "" || np == "" {
		c.FailedStatus("参数为空！")
		return
	}
	// 密码使用rsa进行解密：
	oldPassEncrypt, err1 := base64.StdEncoding.DecodeString(op)
	newPassEncrypt, err2 := base64.StdEncoding.DecodeString(np)
	if err1 != nil || err2 != nil || len(oldPassEncrypt) == 0 || len(newPassEncrypt) == 0 {
		c.FailedStatus("Base64密码解密出错！")
		return
	}
	oldPassDecrypted, err1 := utils.RSADecryptFromPemText(oldPassEncrypt, comm.RsaPrivateKeyText)
	newPassDecrypted, err2 := utils.RSADecryptFromPemText(newPassEncrypt, comm.RsaPrivateKeyText)
	if err1 != nil || err2 != nil || len(oldPassDecrypted) == 0 || len(newPassDecrypted) == 0 {
		c.FailedStatus("RSA密码解密出错！")
		return
	}
	oldPass := string(oldPassDecrypted)
	newPass := string(newPassDecrypted)
	if oldPass == "" || newPass == "" {
		c.FailedStatus("密码为空！")
		return
	}
	userName := c.GetCurrentUser()
	if len(userName) == 0 {
		c.FailedStatus("修改密码失败！")
		return
	}
	if UpdatePassword(userName, oldPass, newPass) {
		c.SucceededStatus("OK！")
	} else {
		c.FailedStatus("修改密码失败！")
	}
}

// LoadCustomConfigAction 加载一个自定义文件
func (c *ConfigController) LoadCustomConfigAction() {
	defer c.ServeJSON()

	customType := c.GetString("type", "")
	if customType == "" {
		c.FailedStatus("未指定类型")
		return
	}
	customFile := getCustomFilename(customType)
	if customFile == "" {
		logging.RuntimeLog.Errorf("error custom file:%s", customType)
		c.FailedStatus("错误的类型")
		return
	}
	content, err := os.ReadFile(filepath.Join(conf.GetRootPath(), "thirdparty", customFile))
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	c.SucceededStatus(string(content))
}

// SaveCustomConfigAction 保存一个自定义文件
func (c *ConfigController) SaveCustomConfigAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	customType := c.GetString("type", "")
	customContent := c.GetString("content", "")
	if customType == "" {
		c.FailedStatus("未指定类型")
		return
	}
	customFile := getCustomFilename(customType)
	if customFile == "" {
		logging.RuntimeLog.Errorf("get custom file type:%s fail", customType)
		c.FailedStatus("错误的类型")
		return
	}
	err := os.WriteFile(filepath.Join(conf.GetRootPath(), "thirdparty", customFile), []byte(customContent), 0666)
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	c.SucceededStatus("保存配置成功")
}

// SaveCustomTaskWorkspaceConfigAction 保存自定义任务的GUID，然后server重新加载
func (c *ConfigController) SaveCustomTaskWorkspaceConfigAction() {
	c.SaveCustomConfigAction()
	ampq.CustomTaskWorkspaceMap = custom.LoadCustomTaskWorkspace()
}

// SaveTaskSliceNumberAction 保存任务切分设置
func (c *ConfigController) SaveTaskSliceNumberAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	ipSliceNumber, err1 := c.GetInt("ipslicenumber", utils.DefaultIpSliceNumber)
	portSliceNumber, err2 := c.GetInt("portslicenumber", utils.DefaultPortSliceNumber)
	if err1 != nil || err2 != nil {
		c.FailedStatus("数量错误")
		return
	}
	err := conf.GlobalServerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	conf.GlobalServerConfig().Task.IpSliceNumber = ipSliceNumber
	conf.GlobalServerConfig().Task.PortSliceNumber = portSliceNumber
	err = conf.GlobalServerConfig().WriteConfig()
	if err != nil {
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// SaveWikiFeishuAction 保存知识库的设置
func (c *ConfigController) SaveWikiFeishuAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	appId := c.GetString("feishuappid", "")
	appSecret := c.GetString("feishusecret", "")
	refreshToken := c.GetString("feishurefreshtoken", "")

	err := conf.GlobalServerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	conf.GlobalServerConfig().Wiki.Feishu.AppId = appId
	conf.GlobalServerConfig().Wiki.Feishu.AppSecret = appSecret
	conf.GlobalServerConfig().Wiki.Feishu.UserAccessRefreshToken = refreshToken

	err = conf.GlobalServerConfig().WriteConfig()
	if err != nil {
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// SaveTaskNotifyAction 保存任务通知的Token设置
func (c *ConfigController) SaveTaskNotifyAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	serverChanToken := c.GetString("token_serverchan", "")
	dingtalkToken := c.GetString("token_dingtalk", "")
	feishuToken := c.GetString("token_feishu", "")

	err := conf.GlobalServerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	if conf.GlobalServerConfig().Notify == nil {
		conf.GlobalServerConfig().Notify = make(map[string]conf.Notify)
	}
	conf.GlobalServerConfig().Notify["serverchan"] = conf.Notify{Token: serverChanToken}
	conf.GlobalServerConfig().Notify["dingtalk"] = conf.Notify{Token: dingtalkToken}
	conf.GlobalServerConfig().Notify["feishu"] = conf.Notify{Token: feishuToken}

	err = conf.GlobalServerConfig().WriteConfig()
	if err != nil {
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// TestTaskNotifyAction 测试任务通知
func (c *ConfigController) TestTaskNotifyAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	message := "这是一个测试消息，来自Nemo的配置管理！"
	notify.Send(message)

	c.SucceededStatus("已发送测试通知，请确认消息是否正确！")
}

// SaveAPITokenAction 保存API的Token
func (c *ConfigController) SaveAPITokenAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	data := DefaultConfig{}
	err := c.ParseForm(&data)
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	err = conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		logging.RuntimeLog.Error("read config file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	//onlineapi
	conf.GlobalWorkerConfig().OnlineAPI.IsFofa = data.IsFofa
	conf.GlobalWorkerConfig().OnlineAPI.IsQuake = data.IsQuake
	conf.GlobalWorkerConfig().OnlineAPI.IsHunter = data.IsHunter
	conf.GlobalWorkerConfig().API.Fofa.Key = data.FofaToken
	conf.GlobalWorkerConfig().API.Hunter.Key = data.HunterToken
	conf.GlobalWorkerConfig().API.Quake.Key = data.QuakeToken
	conf.GlobalWorkerConfig().API.ICP.Key = data.ChinazToken
	conf.GlobalWorkerConfig().API.SearchLimitCount = data.SearchLimitCount
	conf.GlobalWorkerConfig().API.SearchPageSize = data.SearchPageSize
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		logging.RuntimeLog.Error("save config file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	c.SucceededStatus("保存配置成功")
}

// TestOnlineAPIKeyAction 在线测试API的key是否可用
func (c *ConfigController) TestOnlineAPIKeyAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}
	sb := strings.Builder{}
	swg := sizedwaitgroup.New(4)
	msgChan := make(chan string)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-msgChan:
				sb.WriteString(msg)
			case <-done:
				return
			}
		}
	}()

	apiKeys := conf.GlobalWorkerConfig().API
	if len(apiKeys.Fofa.Key) > 0 {
		swg.Add()
		go testOnineAPI("fofa", &swg, msgChan)
	}
	if len(apiKeys.Hunter.Key) > 0 {
		swg.Add()
		go testOnineAPI("hunter", &swg, msgChan)
	}
	if len(apiKeys.Quake.Key) > 0 {
		swg.Add()
		go testOnineAPI("quake", &swg, msgChan)
	}
	if len(apiKeys.ICP.Key) > 0 {
		swg.Add()
		go func(swg *sizedwaitgroup.SizedWaitGroup, testMsgChan chan string) {
			defer swg.Done()

			icp := onlineapi.NewICPQuery(onlineapi.ICPQueryConfig{})
			if icp.RunICPQuery("10086.cn") != nil {
				testMsgChan <- "icp: OK!\n"
			} else {
				testMsgChan <- "icp: fail\n"
			}
		}(&swg, msgChan)
	}
	swg.Wait()
	done <- struct{}{}
	if sb.Len() > 0 {
		c.SucceededStatus(sb.String())
	} else {
		c.FailedStatus("api接口没有可用的key！")
	}
}

// SaveFingerprintAction 保存默认指纹设置
func (c *ConfigController) SaveFingerprintAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}
	data := DefaultConfig{}
	err := c.ParseForm(&data)
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	err = conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		logging.RuntimeLog.Error("read config file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	//fingerprint
	conf.GlobalWorkerConfig().Fingerprint.IsHttpx = data.IsHttpx
	conf.GlobalWorkerConfig().Fingerprint.IsFingerprintHub = data.IsFingerprintHub
	conf.GlobalWorkerConfig().Fingerprint.IsScreenshot = data.IsScreenshot
	conf.GlobalWorkerConfig().Fingerprint.IsIconHash = data.IsIconHash
	conf.GlobalWorkerConfig().Fingerprint.IsFingerprintx = data.IsFingerprintx
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		logging.RuntimeLog.Error("save config file error:", err)
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// SaveWorkerProxyAction 保存worker代理设置
func (c *ConfigController) SaveWorkerProxyAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	data := DefaultConfig{}
	err := c.ParseForm(&data)
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	err = conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		logging.RuntimeLog.Error("read config file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	var hostList []string
	for _, line := range strings.Split(data.ProxyList, "\n") {
		host := strings.TrimSpace(line)
		if host != "" {
			hostList = append(hostList, host)
		}
	}
	conf.GlobalWorkerConfig().Proxy.Host = hostList
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		logging.RuntimeLog.Error("save config file error:", err)
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// SaveDomainscanAction 保存默认域名任务的设置
func (c *ConfigController) SaveDomainscanAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	data := DefaultConfig{}
	err := c.ParseForm(&data)
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	err = conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		logging.RuntimeLog.Error("read config file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	conf.GlobalWorkerConfig().Domainscan.Wordlist = data.Wordlist
	conf.GlobalWorkerConfig().Domainscan.IsSubDomainFinder = data.IsSubDomainFinder
	conf.GlobalWorkerConfig().Domainscan.IsSubDomainBrute = data.IsSubDomainBrute
	conf.GlobalWorkerConfig().Domainscan.IsSubdomainCrawler = data.IsSubDomainCrawler
	conf.GlobalWorkerConfig().Domainscan.IsIgnoreCDN = data.IsIgnoreCDN
	conf.GlobalWorkerConfig().Domainscan.IsIgnoreOutofChina = data.IsIgnoreOutofChina
	conf.GlobalWorkerConfig().Domainscan.IsPortScan = data.IsPortscan
	conf.GlobalWorkerConfig().Domainscan.IsICP = data.IsICP
	conf.GlobalWorkerConfig().Domainscan.IsWhois = data.IsWhois
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		logging.RuntimeLog.Error("save config file error:", err)
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// SavePortscanAction 保存默认端口扫描设置
func (c *ConfigController) SavePortscanAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	cmdbin := c.GetString("cmdbin", "masscan")
	port := c.GetString("port", "--top-ports 1000")
	rate, err1 := c.GetInt("rate", 1000)
	tech := c.GetString("tech", "-sS")
	ping, err2 := c.GetBool("ping", false)
	if err1 != nil || err2 != nil {
		c.FailedStatus("配置参数错误！")
		return
	}
	err := conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		logging.RuntimeLog.Error("read config file error:", err)
		c.FailedStatus(err.Error())
		return
	}

	conf.GlobalWorkerConfig().Portscan.Cmdbin = "masscan"
	if cmdbin == "nmap" {
		conf.GlobalWorkerConfig().Portscan.Cmdbin = "nmap"
	} else if cmdbin == "gogo" {
		conf.GlobalWorkerConfig().Portscan.Cmdbin = "gogo"
	}
	conf.GlobalWorkerConfig().Portscan.Port = port
	conf.GlobalWorkerConfig().Portscan.Rate = rate
	conf.GlobalWorkerConfig().Portscan.Tech = tech
	conf.GlobalWorkerConfig().Portscan.IsPing = ping
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		logging.RuntimeLog.Error("save config file error:", err)
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// UploadPocAction xraypoc的上传
func (c *ConfigController) UploadPocAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	pocType := c.GetString("type", "")
	if len(pocType) == 0 {
		logging.RuntimeLog.Error("get poc file type error")
		return
	}
	// 获取上传信息
	f, h, err := c.GetFile("file")
	if err != nil {
		logging.RuntimeLog.Error("get upload file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	defer f.Close()
	// 检查文件后缀
	fileExt := path.Ext(h.Filename)
	if fileExt != ".yml" && fileExt != ".yaml" {
		logging.RuntimeLog.Warning("invalid file type")
		c.FailedStatus("invalid file type!")
		return
	}
	// 保存到poc目录下
	var pocSavedPathName string
	if pocType == "xray" {
		pocSavedPathName = filepath.Join(conf.GetRootPath(), conf.GlobalWorkerConfig().Pocscan.Xray.PocPath, h.Filename)
	} else if pocType == "nuclei" {
		pocSavedPathName = filepath.Join(conf.GetRootPath(), conf.GlobalWorkerConfig().Pocscan.Nuclei.PocPath, h.Filename)
	}
	if pocSavedPathName == "" {
		logging.RuntimeLog.Error("get poc save path from config file error")
		c.FailedStatus("get poc save path from config file error")
		return
	}
	err = c.SaveToFile("file", pocSavedPathName)
	if err != nil {
		logging.RuntimeLog.Error("save poc file error:", err)
		c.FailedStatus(err.Error())
		return
	}
	c.SucceededStatus("上传成功")
}

// SaveWorkerFilterAction 保存任务过滤的设置
func (c *ConfigController) SaveWorkerFilterAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}

	maxPortPerIp, err1 := c.GetInt("maxportperip", portscan.IpOpenedPortFilterNumber)
	maxDomainPerIp, err2 := c.GetInt("maxdomainperip", domainscan.SameIpToDomainFilterMax)
	titleFilter := c.GetString("title", "")
	if err1 != nil || err2 != nil {
		c.FailedStatus("数量错误")
		return
	}
	err := conf.GlobalWorkerConfig().ReloadConfig()
	if err != nil {
		c.FailedStatus(err.Error())
		return
	}
	conf.GlobalWorkerConfig().Filter.MaxPortPerIp = maxPortPerIp
	conf.GlobalWorkerConfig().Filter.MaxDomainPerIp = maxDomainPerIp
	conf.GlobalWorkerConfig().Filter.Title = titleFilter
	err = conf.GlobalWorkerConfig().WriteConfig()
	if err != nil {
		c.FailedStatus(err.Error())
	}
	c.SucceededStatus("保存配置成功")
}

// UpdateMinichatConfigAction 更新minichat的设置
func (c *ConfigController) UpdateMinichatConfigAction() {
	defer c.ServeJSON()
	if c.CheckMultiAccessRequest([]RequestRole{SuperAdmin, Admin}, false) == false {
		c.FailedStatus("当前用户权限不允许！")
		return
	}
	anonymous, err0 := c.GetBool("anonymous", false)
	notDeleteDir, err1 := c.GetBool("notdelfiledir", false)
	loadHistory, err2 := c.GetBool("loadhistory", true)
	maxHistoryMessage, err3 := c.GetInt("maxhistorymessage", 1000)
	if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
		c.FailedStatus("参数错误")
		return
	}
	minichatConfig.EnableAnonymous = anonymous
	minichatConfig.IsNotDelFileDir = notDeleteDir
	minichatConfig.LoadHistory = loadHistory
	minichatConfig.MaxHistoryMessage = maxHistoryMessage

	c.SucceededStatus("更新配置成功")
}

// getCustomFilename  根据类型返回自定义文件名
func getCustomFilename(customType string) (customFile string) {
	switch customType {
	case HoneyPot:
		customFile = "custom/honeypot.txt"
	case IPLocation:
		customFile = "custom/iplocation-custom.txt"
	case IPLocationB:
		customFile = "custom/iplocation-custom-B.txt"
	case IPLocationC:
		customFile = "custom/iplocation-custom-C.txt"
	case Service:
		customFile = "custom/services-custom.txt"
	case Xray:
		customFile = "xray/xray.yaml"
	case XrayConfig:
		customFile = "xray/config.yaml"
	case XrayModule:
		customFile = "xray/module.xray.yaml"
	case XrayPlugin:
		customFile = "xray/plugin.xray.yaml"
	case BlackDomain:
		customFile = "custom/black_domain.txt"
	case BlackIP:
		customFile = "custom/black_ip.txt"
	case TaskWorkspace:
		customFile = "custom/task_workspace.txt"
	case FOFAFilterKeyword:
		customFile = "custom/onlineapi_filter_keyword.txt"
	case FOFAFilterKeywordLocal:
		customFile = "custom/onlineapi_filter_keyword_local.txt"
	}

	return
}

// testOnineAPI 多线程方式测试在线接口的可用性
func testOnineAPI(apiName string, swg *sizedwaitgroup.SizedWaitGroup, testMsgChan chan string) {
	defer swg.Done()

	s := onlineapi.NewOnlineAPISearch(onlineapi.OnlineAPIConfig{Target: "fofa.info"}, apiName)
	s.Config.SearchPageSize = 100
	s.Config.SearchLimitCount = 100
	s.Do()

	if len(s.DomainResult.DomainResult) > 0 || len(s.IpResult.IPResult) > 0 {
		testMsgChan <- fmt.Sprintf("%s: OK!\n", apiName)
	} else {
		testMsgChan <- fmt.Sprintf("%s: fail\n", apiName)
	}

	return
}
