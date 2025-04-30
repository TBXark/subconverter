package converter

type ConvertParams struct {
	Target       string `form:"target"`        // 必要参数，配置类型如surge&ver=4
	URL          string `form:"url"`           // 可选，订阅链接URLEncode处理
	Group        string `form:"group"`         // 可选，设置订阅组名
	UploadPath   string `form:"upload_path"`   // 可选，Gist上传文件名URLEncode处理
	Include      string `form:"include"`       // 可选，保留节点正则URLEncode处理
	Exclude      string `form:"exclude"`       // 可选，排除节点正则URLEncode处理
	Config       string `form:"config"`        // 可选，外部配置地址URLEncode处理
	DevID        string `form:"dev_id"`        // 可选，QuantumultX远程设备ID
	Filename     string `form:"filename"`      // 可选，生成订阅文件名
	Interval     int    `form:"interval"`      // 可选，托管配置更新间隔(秒)
	Rename       string `form:"rename"`        // 可选，自定义重命名URLEncode处理
	FilterScript string `form:"filter_script"` // 可选，筛选节点js代码URLEncode处理
	Strict       bool   `form:"strict"`        // 可选，是否强制更新
	Upload       bool   `form:"upload"`        // 可选，是否上传至Gist
	Emoji        bool   `form:"emoji"`         // 可选，节点名是否含Emoji
	AddEmoji     bool   `form:"add_emoji"`     // 可选，节点名前加Emoji
	RemoveEmoji  bool   `form:"remove_emoji"`  // 可选，是否删除原有Emoji
	AppendType   bool   `form:"append_type"`   // 可选，节点名前插入类型
	TFO          bool   `form:"tfo"`           // 可选，开启TCP Fast Open
	UDP          bool   `form:"udp"`           // 可选，开启UDP
	List         bool   `form:"list"`          // 可选，输出节点列表类型
	Sort         bool   `form:"sort"`          // 可选，是否按节点名排序
	SortScript   string `form:"sort_script"`   // 可选，自定义排序js代码URLEncode处理
	Script       bool   `form:"script"`        // 可选，生成Clash Script
	Insert       bool   `form:"insert"`        // 可选，是否插入insert_url
	SCV          bool   `form:"scv"`           // 可选，关闭TLS证书检查
	FDN          bool   `form:"fdn"`           // 可选，过滤不支持节点
	Expand       bool   `form:"expand"`        // 可选，是否处理规则列表
	AppendInfo   bool   `form:"append_info"`   // 可选，输出流量/到期信息
	Prepend      bool   `form:"prepend"`       // 可选，insert_url插入位置
	Classic      bool   `form:"classic"`       // 可选，生成Clash classical rule-provider
	TLS13        bool   `form:"tls13"`         // 可选，增加tls1.3参数
	NewName      bool   `form:"new_name"`      // 可选，启用Clash新组名
}

type ShadowSocks struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Server            string `json:"server"`
	Port              int    `json:"port"`
	Cipher            string `json:"cipher"`
	Password          string `json:"password"`
	Udp               bool   `json:"udp"`
	UdpOverTcp        bool   `json:"udp-over-tcp"`
	UdpOverTcpVersion int    `json:"udp-over-tcp-version"`
	IpVersion         string `json:"ip-version"`
	Plugin            string `json:"plugin"`
	PluginOpts        struct {
		Mode string `json:"mode"`
	} `json:"plugin-opts"`
	Smux struct {
		Enabled bool `json:"enabled"`
	} `json:"smux"`
}
