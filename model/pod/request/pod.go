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

type Resources struct {
	//是否配置容器配额
	Enable bool `json:"enable"`
	//内存
	MemRequest int32 `json:"memRequest"`
	MemLimit   int32 `json:"memLimit"`

	//CPU
	CpuRequest int32 `json:"cpuRequest"`
	CpuLimit   int32 `json:"cpuLimit"`
}

type VolumeMount struct {
	//挂载卷名称
	MountName string `json:"mountName"`
	//挂载卷-> 对应容器内的路径
	MountPath string `json:"mountPath"`
	//是否只读
	ReadOnly bool `json:"readOnly"`
}

type ProbeHttpGet struct {
	//请求协议 http、https
	Scheme string `json:"scheme"`
	//请求host，如果为空就为pod内请求
	Host string `json:"host"`
	//请求路径
	Path string `json:"path"`
	//请求端口
	Port int32 `json:"port"`
	//请求的header
	HttpHeaders []ListMapItem `json:"httpHeaders"`
}

type ProbeCommand struct {
	//cat /test/test.txt
	Command []string `json:"command"`
}

type ProbeTcpSocket struct {
	//请求host，如果为空就为pod内请求
	Host string `json:"host"`
	//请求端口
	Port int32 `json:"port"`
}

type ProbeTime struct {
	//初始化时间 初始化若干秒然后开始探针
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`
	//探测间隔时间，每个若干秒探测
	PeriodSeconds int32 `json:"periodSeconds"`
	//探测等待时间，等待若干秒没有返回就表示失败
	TimeoutSeconds int32 `json:"timeoutSeconds"`
	//探测若干次成功了则表示这次探针成功
	SuccessThreshold int32 `json:"successThreshold"`
	//探测若干次失败才认为这次探针失败
	FailureThreshold int32 `json:"failureThreshold"`
}

type ContainerProbe struct {
	//是否打开探针
	Enable bool `json:"enable"`
	//探针类型
	Type      string         `json:"type"`
	HttpGet   ProbeHttpGet   `json:"httpGet"`
	Command   ProbeCommand   `json:"command"`
	TcpSocket ProbeTcpSocket `json:"tcpSocket"`
	ProbeTime
}

type Container struct {
	//容器名称
	Name string `json:"name"`
	//容器镜像
	Image string `json:"image"`
	//镜像拉去策略
	ImagePullPolicy string `json:"imagePullPolicy"`
	//是否开启伪终端
	Tty bool `json:"tty"`
	//工作目录
	WorkingDir string `json:"workingDir"`
	//执行命令
	Command []string `json:"command"`
	//参数
	Args []string `json:"args"`
	//环境变量
	Envs []ListMapItem `json:"envs"`
	//是否开启模式
	Privileged bool `json:"privileged"`
	//容器申请配额
	Resources Resources `json:"resources"`
	//容器卷挂载
	VolumeMounts []VolumeMount `json:"volumeMounts"`
	//启动探针
	StartupProbe ContainerProbe `json:"startupProbe"`
	//存活探针
	LivenessProbe ContainerProbe `json:"livenessProbe"`
	//就绪探针
	ReadinessProbe ContainerProbe `json:"readinessProbe"`
}

type Pod struct {
	//基础定义信息
	Base Base `json:"base"`
	//卷
	Volumes []Volumes `json:"volumes"`
	//网络相关
	NetWorking NetWorking `json:"netWorking"`
	//container
	Containers []Container `json:"containers"`
	//init containers
	InitContainers []Container `json:"initContainers"`
}
