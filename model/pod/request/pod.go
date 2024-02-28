package request

type ListMapItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Base struct {
	//名字
	Name string `json:"name"`
	//标签
	Labels []ListMapItem `json:"labels"`
	//命名空间
	Namespace string `json:"namespace"`
	//重启策略 Always|Never|On-Failure
	RestartPolicy string `json:"restartPolicy"`
}

type Volumes struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DnsConfig struct {
	Nameservers []string `json:"nameservers"`
}

type NetWorking struct {
	HostNetwork bool          `json:"hostNetwork"`
	HostName    string        `json:"hostname"`
	DnsPolicy   string        `json:"dnsPolicy"`
	DnsConfig   DnsConfig     `json:"dnsConfig"`
	HostAliases []ListMapItem `json:"hostAliases"`
}

type Pod struct {
	//基础定义信息
	Base string `json:"base"`
	//卷
	Volumes []string `json:"volumes"`
	//网络相关
	NetWorking NetWorking `json:"netWorking"`
	//container
	Containers []string `json:"containers"`
	//init containers
	InitContainers []string `json:"initContainers"`
}
