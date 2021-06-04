package restapi

const adminJwtToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyTmFtZSI6ImFkbWluIn0.ZZX3Z0zbeCWGhjsdtCxrf3O4xTQ4QYc38AED6RLSUG0`

var ServerFlag string

type commandOut struct {
	Logid   int64  `json:"log_id"`
	Content string `json:"content"`
}

type payload struct {
	Logid     string `json:"logID"`
	Status    string `json:"status"`
	Starttime string `json:"startTime"`
	Endtime   string `json:"endTime"`
	Duration  string `json:"duration"`
}

type MatchInfo struct {
	Branch string
	Stream string
}

type NamePairInfo struct {
	// git email
	GitEmail string `json:"gitEmail,omitempty"`
	// git user name
	GitUserName string `json:"gitUserName,omitempty"`
	// snv user name
	SnvUserName string `json:"svnUserName,omitempty"`
}

type Task struct {
	TaskId       int64
	TaskLogId    int64
	CcPassword   string
	CcUser       string
	Component    string
	GitPassword  string
	GitURL       string
	GitUser      string
	GitEmail     string
	Pvob         string
	IncludeEmpty bool
	Matches      []MatchInfo
	Keep         string
	ModelType    string          `json:"modelType,omitempty"`
	NamePair     []*NamePairInfo `json:"namePair"`
	SvnURL       string          `json:"svnUrl,omitempty"`
	Gitignore    string          `json:"gitignore"`
}

type WorkerTaskDelInfo struct {
	TaskId     int64  `json:"task_id"`
	CcPassword string `json:"cc_password"`
	CcUser     string `json:"cc_user"`
	Exception  string `json:"exception"`
	ModelType  string `json:"modelType,omitempty"`
}

type Conf struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	ServerAddr string `yaml:"server_addr"`
}

type Payload struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
