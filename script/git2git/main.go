package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"moul.io/http2curl"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

const GitlabGuest = 10
const GitlabReporter = 20
const GitlabDeveloper = 30
const GitlabMaintainer = 40
const GitlabOwner = 50
const GiteeAdmin = 8
const GiteeCollar = 4
const GiteeVisitor = 2

type GitlabService struct {
	Host          string
	Token         string
	GroupPath     string
	GroupID       int
	ParentID      int
	ProjectPath   string
	ProjectID     int
	ParentPath    string
	GroupFullPath string
	GiteeHost     string
}

type GiteeService struct {
	Token            string
	GroupPath        string
	GroupID          int
	ProjectPath      string
	ProjectID        int
	CodeURLPrefix    string //http://cmb-gitaly.dev.gitee.work/api/code/api/enterprises/cmbchina
	CodeUserQueryURL string //http://code.gitee.work/api/gitlab/users?username=
	PrivateToken     string
	GiteeHost        string
	GroupFullPath    string
}

var GLS GitlabService
var GTS GiteeService
var UserMap map[string]int
var RoleMap map[int]int
var envConfig Config

func (gls *GitlabService) Get(url, queryStrings string) *http.Response {
	url = fmt.Sprintf("%s/%s?access_token=%s", gls.Host, url, gls.Token)
	if queryStrings != "" {
		url = url + "&" + queryStrings
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Printf("gitlab api\n%s\n", command)
	if err != nil {
		fmt.Println("create gitlab service request failed:", err)
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("call gitlab host failed:", err)
		return nil
	}
	return resp
}

type DescendantGroupResponse struct {
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

// TranslateGroupsDescendant 获取组下所有子组
func (gls *GitlabService) TranslateGroupsDescendant() {
	//先同步自己
	result := gls.TranslateGroupsByName()
	if result {
		gls.TranslateProjectsByGroup()
	}
	path := url2.QueryEscape(gls.GroupPath)
	i := 1
	for {
		url := fmt.Sprintf("api/v4/groups/%s/descendant_groups", path)
		resp := gls.Get(url, fmt.Sprintf("page=%d", i))
		ret := make([]DescendantGroupResponse, 0)
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}
		for _, info := range ret {
			gls.GroupPath = info.FullPath
			//gls.GroupFullPath = info.FullPath
			result := gls.TranslateGroupsByName()
			if result {
				gls.TranslateProjectsByGroup()
			}
		}
		if len(ret) < 20 {
			break
		}
		i++
	}
}

type GroupResponse struct {
	ID                             int         `json:"id"`
	WebURL                         string      `json:"web_url"`
	Name                           string      `json:"name"`
	Path                           string      `json:"path"`
	Description                    string      `json:"description"`
	Visibility                     string      `json:"visibility"`
	ShareWithGroupLock             bool        `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool        `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int         `json:"two_factor_grace_period"`
	ProjectCreationLevel           string      `json:"project_creation_level"`
	AutoDevopsEnabled              interface{} `json:"auto_devops_enabled"`
	SubgroupCreationLevel          string      `json:"subgroup_creation_level"`
	EmailsDisabled                 interface{} `json:"emails_disabled"`
	MentionsDisabled               interface{} `json:"mentions_disabled"`
	LfsEnabled                     bool        `json:"lfs_enabled"`
	DefaultBranchProtection        int         `json:"default_branch_protection"`
	AvatarURL                      interface{} `json:"avatar_url"`
	RequestAccessEnabled           bool        `json:"request_access_enabled"`
	FullName                       string      `json:"full_name"`
	FullPath                       string      `json:"full_path"`
	CreatedAt                      time.Time   `json:"created_at"`
	ParentID                       interface{} `json:"parent_id"`
	LdapCn                         interface{} `json:"ldap_cn"`
	LdapAccess                     interface{} `json:"ldap_access"`
}

func (gls *GitlabService) TranslateGroupsByName() bool {
	groupPaths := strings.Split(gls.GroupPath, "/")
	// 组大于5级不迁移
	if len(groupPaths) > 5 {
		return false
	}
	//parentPath := strings.Join(groupPaths[:len(groupPaths)-1], ",")
	//currentPath := groupPaths[len(groupPaths)-1]

	//ret := make([]Group, 0)
	//if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
	//	panic(err)
	//}

	for i, groupPath := range groupPaths {
		gls.GroupPath = groupPath
		fullPath := groupPaths[0]
		if i > 0 {
			for n := 1; n <= i; n++ {
				fullPath = fullPath + "/" + groupPaths[n]
			}
		}
		url := gls.GiteeHost + fmt.Sprintf("/api/gitlab/groups/can_migrate?full_path=%s", fullPath)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("PRIVATE-TOKEN", GTS.PrivateToken)
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Printf("gitee api invoke\n%s\n", command)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		if resp == nil {
			panic("can not check group")
		}

		if resp.StatusCode == http.StatusOK {
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic("can not check group")
			}
			if string(bytes) == "false" {
				fmt.Printf("duplicate path %s", groupPath)
				return false
			}
		} else {
			panic("can not check group")
		}
		gls.GroupFullPath = fullPath
		gls.TranslateGroupByName()
	}
	return true
}

func (gls *GitlabService) TranslateGroupByName() {
	url := fmt.Sprintf("api/v4/groups/%s", url2.QueryEscape(gls.GroupFullPath))
	//queryStrings := "search=" + gls.GroupFullPath
	resp := gls.Get(url, "")
	if resp == nil {
		panic("no valid group")
	}
	ret := GroupResponse{}
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}
		gls.GroupID = ret.ID
		GTS.CreateGroupWithPath(GLS.GroupPath, ret.Name)
		gls.TranslateMemberPermissionByGroupOrProject("group")
		return
	}

	//for _, info := range ret {
	//	if info.Path == gls.GroupPath {
	//		gls.GroupID = info.ID
	//		GTS.CreateGroupWithPath(GLS.GroupPath, info.Name)
	//		gls.TranslateMemberPermissionByGroupOrProject("group")
	//		return
	//	}
	//}
	panic("no group found")
}

func (gls *GitlabService) TranslateProjectsByGroup() {
	i := 1
	for {
		url := fmt.Sprintf("api/v4/groups/%d/projects", gls.GroupID)
		resp := gls.Get(url, fmt.Sprintf("page=%d", i))
		if resp == nil {
			panic("no valid project in group")
		}
		type ProjectResponse struct {
			ID                int       `json:"id"`
			Description       string    `json:"description"`
			Name              string    `json:"name"`
			NameWithNamespace string    `json:"name_with_namespace"`
			Path              string    `json:"path"`
			PathWithNamespace string    `json:"path_with_namespace"`
			CreatedAt         time.Time `json:"created_at"`
			DefaultBranch     string    `json:"default_branch"`
			SSHURLToRepo      string    `json:"ssh_url_to_repo"`
			HTTPURLToRepo     string    `json:"http_url_to_repo"`
			WebURL            string    `json:"web_url"`
			ForksCount        int       `json:"forks_count"`
			StarCount         int       `json:"star_count"`
			EmptyRepo         bool      `json:"empty_repo"`
			Archived          bool      `json:"archived"`
			Visibility        string    `json:"visibility"`
		}
		ret := make([]ProjectResponse, 0)
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}
		fmt.Printf("get gitlab group's project list: %s", ret)
		for _, info := range ret {
			if gls.ProjectPath == "" || (gls.ProjectPath != "" && info.Path == gls.ProjectPath) {
				gls.ProjectID = info.ID
				result := GTS.CreateProjectWithName(info.Name, info.Path, info.Description)
				if result {
					GLS.TranslateMemberPermissionByGroupOrProject("project")
				}
			}
		}
		if len(ret) < 20 {
			break
		}
		i++
	}
}

func (gls *GitlabService) MigrateProjectAllBranches(projectID int) {
	url := fmt.Sprintf("api/v4/projects/%d/repository/branches", projectID)
	resp := gls.Get(url, "")
	if resp == nil {
		panic("no valid branch info")
	}
	type ProjectBranchResponse struct {
		Name   string `json:"name"`
		Commit struct {
			ID             string    `json:"id"`
			ShortID        string    `json:"short_id"`
			CreatedAt      time.Time `json:"created_at"`
			Title          string    `json:"title"`
			Message        string    `json:"message"`
			AuthorName     string    `json:"author_name"`
			AuthorEmail    string    `json:"author_email"`
			AuthoredDate   time.Time `json:"authored_date"`
			CommitterName  string    `json:"committer_name"`
			CommitterEmail string    `json:"committer_email"`
			CommittedDate  time.Time `json:"committed_date"`
			WebURL         string    `json:"web_url"`
		} `json:"commit"`
		Merged    bool `json:"merged"`
		Protected bool `json:"protected"`
	}
	ret := make([]ProjectBranchResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		panic(err)
	}
	for _, info := range ret {
		//TODO create branch by gitee api
		fmt.Println(info)
	}
}

func (gls *GitlabService) TranslateMemberPermissionByGroupOrProject(targetType string) {
	url := ""
	if targetType == "group" {
		url = fmt.Sprintf("api/v4/groups/%d/members/all", gls.GroupID)
	} else {
		url = fmt.Sprintf("api/v4/projects/%d/members/all", gls.ProjectID)
	}
	resp := gls.Get(url, "")
	if resp == nil {
		panic("no valid project in group")
	}
	type UserPermissionResponse struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Username    string    `json:"username"`
		State       string    `json:"state"`
		AvatarURL   string    `json:"avatar_url"`
		WebURL      string    `json:"web_url"`
		AccessLevel int       `json:"access_level"`
		CreatedAt   time.Time `json:"created_at"`
		//ExpiresAt   *time.Time `json:"expires_at"`
	}
	ret := make([]UserPermissionResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		panic(err)
	}
	for _, info := range ret {
		if !GTS.GetGiteeUserInfo(info.Username) {
			fmt.Println("no valid user found " + info.Username)
			continue
		}
		GTS.CreateGroupOrProjectMember(info.Username, targetType, info.AccessLevel)
	}
}

func (gts *GiteeService) PostOrGet(url, method string, body io.Reader) *http.Response {
	url = gts.CodeURLPrefix + url
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("create gitee request failed", err)
		return nil
	}
	cookie := http.Cookie{Name: "PRE-GW-SESSION", Value: gts.Token, Path: "/"}
	req.AddCookie(&cookie)
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Printf("gitee api invoke\n%s\n", command)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("connect gitee failed", err)
		return nil
	}
	return r
}

//
//func translateWord(words string) string {
//	args := pinyin.Args{}
//	retKeywordsPinyin := ""
//	for _, w := range words {
//		if unicode.Is(unicode.Han, w) {
//			pinyinKeywords := pinyin.Pinyin(string(w), args)
//			if len(pinyinKeywords) > 0 {
//				retKeywordsPinyin += pinyinKeywords[0][0]
//			}
//		} else {
//			retKeywordsPinyin += string(w)
//		}
//	}
//	return retKeywordsPinyin
//}

func (gts *GiteeService) GetGroupInfoByPath(groupPath string) bool {
	//if gts.GroupPath != "" {
	//	groupPath = gts.GroupPath
	//}
	url := "/repo_groups/" + groupPath
	resp := gts.PostOrGet(url, http.MethodGet, nil)
	if resp == nil || resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusNotFound {
		return false
	} else {
		ret := struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Path        string `json:"path"`
			Belong      string `json:"belong"`
			Description string `json:"description"`
			Creator     struct {
				Name      string `json:"name"`
				AvatarURL string `json:"avatar_url"`
			} `json:"creator"`
			CreatedAt        time.Time `json:"created_at"`
			CanModify        bool      `json:"can_modify"`
			ChildExists      bool      `json:"child_exists"`
			Depth            int       `json:"depth"`
			GroupUUID        string    `json:"group_uuid"`
			PublicIdent      string    `json:"public_ident"`
			PullRequestCount int       `json:"pull_request_count"`
			IssueCount       int       `json:"issue_count"`
			PermissionPolicy string    `json:"permission_policy"`
			UserPermission   struct {
				Code       int    `json:"code"`
				Permission string `json:"permission"`
			} `json:"user_permission"`
			IconColor string `json:"icon_color"`
			IconPath  string `json:"icon_path"`
			IconURL   string `json:"icon_url"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
			panic(err)
		}
		gts.GroupID = ret.ID
		gts.GroupPath = ret.Path
		return true
	}
}

func (gts *GiteeService) CreateGroupWithPath(groupPath string, groupName string) {
	if gts.GetGroupInfoByPath(groupPath) {
		return
	}
	fmt.Println("create new group", groupPath)
	url := "/repo_groups"
	postBody := struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		ParentID    int    `json:"parent_id"`
		Description string `json:"description"`
	}{
		Name: groupName,
		Path: groupPath,
		// 如果此时GTS已有GroupID，说明本次迁移过程中已经有创建好了的Group且Path不等于当前需要创建的，那么这个已有的Group应该为当前创建的Group的Parent
		ParentID: gts.GroupID,
	}
	postBodyByte, _ := json.Marshal(postBody)
	resp := gts.PostOrGet(url, http.MethodPost, bytes.NewBuffer(postBodyByte))
	if resp == nil {
		panic("create group failed " + groupName)
	}
	if !gts.GetGroupInfoByPath(groupPath) {
		panic("create group failed" + groupName)
	}
}

func (gts *GiteeService) GetProjectByPath(path string) bool {
	url := fmt.Sprintf("/programs/%s/projects/%s", gts.GroupPath, path)
	resp := gts.PostOrGet(url, http.MethodGet, nil)
	if resp == nil {
		panic("no project found " + path)
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
		return false
	}
	ret := struct {
		ID   int    `json:"id"`
		Path string `json:"path"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		panic(err)
	}
	gts.ProjectID = ret.ID
	gts.ProjectPath = ret.Path
	return true
}

func (gts *GiteeService) CreateProjectWithName(name string, path string, description string) bool {
	if gts.GetProjectByPath(path) {
		return false
	}
	fmt.Println("create new project", path)
	url := "/inner_source/projects"
	postBody := struct {
		Readme          bool   `json:"readme"`
		Name            string `json:"name"`
		Path            string `json:"path"`
		CanCreateBranch bool   `json:"can_create_branch"`
		BranchModelName string `json:"branch_model_name"`
		Public          bool   `json:"public"`
		ParentID        int    `json:"parent_id"`
		Description     string `json:"description"`
	}{
		Readme:          false,
		Name:            name,
		ParentID:        gts.GroupID,
		Path:            path,
		Public:          false,
		CanCreateBranch: true,
		BranchModelName: "single_branch_model",
		Description:     description,
	}
	postBodyByte, _ := json.Marshal(postBody)
	fmt.Printf("create project post body: %s", postBody)
	resp := gts.PostOrGet(url, http.MethodPost, bytes.NewBuffer(postBodyByte))
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("create project response: %s", string(bytes))
	if resp == nil || resp.StatusCode != http.StatusCreated {
		fmt.Println("create project failed " + name)
		return false
	}
	if !gts.GetProjectByPath(path) {
		fmt.Println("create project failed " + name)
		return false
	}
	return true
}

func (gts *GiteeService) GetGiteeUserInfo(username string) bool {
	if _, ok := UserMap[username]; ok {
		return true
	}
	url := gts.CodeUserQueryURL + username
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("PRIVATE-TOKEN", gts.PrivateToken)
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Printf("gitee api invoke\n%s\n", command)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	type UserResponse struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
		State     string `json:"state"`
	}
	ret := make([]UserResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		bytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("GetGiteeUserInfo resp: %s", string(bytes))
		println(err)
		return false
		//panic(err)
	}
	if len(ret) == 1 {
		if ret[0].Username == username {
			UserMap[username] = ret[0].ID
			return true
		}
	}

	return false
}

func (gts *GiteeService) CreateGroupOrProjectMember(username, targetType string, accessLevel int) {
	if userID, ok := UserMap[username]; !ok {
		panic("no user found" + username)
	} else {
		if (gts.GroupID == 0 && targetType == "group") || (gts.ProjectID == 0 && targetType == "project") {
			panic(fmt.Sprintf("no valid group or project: %d %d %s", gts.GroupID, gts.ProjectID, targetType))
		}
		codePermission := GiteeVisitor
		if accessLevel == GitlabOwner || accessLevel == GitlabMaintainer {
			codePermission = GiteeAdmin
		} else if accessLevel == GitlabDeveloper {
			codePermission = GiteeCollar
		}
		if targetType == "project" {
			postBody := struct {
				RoleID  int   `json:"role_id"`
				UserIDs []int `json:"user_ids"`
			}{
				RoleID:  RoleMap[codePermission],
				UserIDs: []int{userID},
			}
			postBodyByte, _ := json.Marshal(postBody)
			url := fmt.Sprintf("/programs/%s/projects/%s/members", gts.GroupPath, gts.ProjectPath)
			resp := gts.PostOrGet(url, http.MethodPost, bytes.NewBuffer(postBodyByte))
			if resp == nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
				// todo 迁移时会报错
				fmt.Println("create project member failed " + gts.ProjectPath + " " + username)
				return
				//panic("create project member failed " + gts.ProjectPath + " " + username)
			}
		} else {
			postBody := struct {
				CurrentID int    `json:"current_id"`
				RoleID    int    `json:"role_id"`
				Type      string `json:"type"`
				UserIDs   []int  `json:"user_ids"`
			}{
				CurrentID: gts.GroupID,
				RoleID:    RoleMap[codePermission],
				Type:      "repo_group",
				UserIDs:   []int{userID},
			}
			postBodyByte, _ := json.Marshal(postBody)
			url := "/repo_groups/members"
			resp := gts.PostOrGet(url, http.MethodPost, bytes.NewBuffer(postBodyByte))
			if resp == nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
				// todo 迁移时会panic
				fmt.Println("create group member failed " + gts.GroupPath + " " + username)
				return
				//panic("create group member failed " + gts.GroupPath + " " + username)
			}
		}
	}
}

func (gts *GiteeService) GetRoles() {
	url := "/members/member_roles?is_query=true"
	resp := gts.PostOrGet(url, http.MethodGet, nil)
	if resp == nil || resp.StatusCode != http.StatusOK {
		panic("no valid role")
	}
	type RoleResponse struct {
		ID         int `json:"id"`
		Permission int `json:"permission"`
	}
	ret := make([]RoleResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		panic(err)
	}
	for _, info := range ret {
		RoleMap[info.Permission] = info.ID
	}
}

type Config struct {
	GitlabHost            string
	GitlabToken           string
	GitlabProtocol        string
	GiteeToken            string
	GiteeCodeURLPrefix    string
	GiteeCodeUserQueryURL string
	GiteePrivateToken     string
	GiteeHost             string
}

func ParseConfig() *Config {
	config := &Config{
		GitlabHost:            os.Getenv("GITLAB_HOST"),
		GitlabToken:           os.Getenv("GITLAB_TOKEN"),
		GiteeToken:            os.Getenv("GITEE_TOKEN"),
		GiteeCodeURLPrefix:    os.Getenv("GITEE_CODE_URL_PREFIX"),
		GiteeCodeUserQueryURL: os.Getenv("GITEE_CODE_USER_QUERY_URL"),
		GiteePrivateToken:     os.Getenv("GITEE_PRIVATE_TOKEN"),
		GiteeHost:             os.Getenv("GITEE_HOST"),
	}
	//config := &Config{
	//	GitlabHost:            "http://192.168.48.60:18787",
	//	GitlabToken:           "yiWxjprfmdBtsZE3tfZk",
	//	GiteeToken:            "79e3b5b3ff3e481c912ccfce15cc5d33",
	//	GiteeCodeURLPrefix:    "/api/code/api/enterprises/osc",
	//	GiteeCodeUserQueryURL: "/api/gitlab/users?username=",
	//	GiteePrivateToken:     "c4ca4238a0b923820dcc509a6f75849b",
	//	GiteeHost:             "http://code.gitee.work",
	//}
	return config
}

func main() {
	var gitlabGroupPath string
	var gitlabProjectPath string
	var gitlabToken string
	var gitlabHost string
	var giteeGroupPath string
	var giteeToken string
	flag.StringVar(&gitlabGroupPath, "gitlab_group_path", "", "")
	flag.StringVar(&gitlabProjectPath, "gitlab_project_path", "", "")
	flag.StringVar(&gitlabToken, "gitlab_token", "", "")
	flag.StringVar(&gitlabHost, "gitlab_host", "", "")
	flag.StringVar(&giteeGroupPath, "gitee_group_path", "", "")
	flag.StringVar(&giteeToken, "gitee_token", "", "")
	flag.Parse()

	//gitlabGroupPath = "0919c"

	UserMap = make(map[string]int, 0)
	RoleMap = make(map[int]int, 0)
	config := ParseConfig()
	if gitlabToken == "" {
		gitlabToken = config.GitlabToken
	}
	if gitlabHost == "" {
		gitlabHost = config.GitlabHost
	}
	if giteeToken == "" {
		giteeToken = config.GiteeToken
	}

	gitlabGroupPath = strings.TrimSuffix(strings.TrimPrefix(gitlabGroupPath, "/"), "/")
	gitlabProjectPath = strings.TrimSuffix(strings.TrimPrefix(gitlabProjectPath, "/"), "/")
	gitlabProjectPaths := strings.Split(gitlabProjectPath, "/")
	lenGitlabProjectPaths := len(gitlabProjectPaths)
	if lenGitlabProjectPaths > 1 {
		tmp := strings.Join(gitlabProjectPaths[:lenGitlabProjectPaths-1], "/")
		if len(gitlabGroupPath) > 0 {
			gitlabGroupPath = gitlabGroupPath + "/" + tmp
		} else {
			gitlabGroupPath = tmp
		}
		gitlabProjectPath = gitlabProjectPaths[lenGitlabProjectPaths-1]
	}

	GLS = GitlabService{
		Host:        gitlabHost,
		Token:       gitlabToken,
		GroupPath:   gitlabGroupPath,
		GroupID:     0,
		ProjectPath: gitlabProjectPath,
		ProjectID:   0,
		GiteeHost:   config.GiteeHost,
	}

	GTS = GiteeService{
		Token:            giteeToken,
		GroupPath:        giteeGroupPath,
		GroupID:          0,
		ProjectPath:      "",
		ProjectID:        0,
		CodeURLPrefix:    config.GiteeHost + config.GiteeCodeURLPrefix,
		CodeUserQueryURL: config.GiteeHost + config.GiteeCodeUserQueryURL,
		GiteeHost:        config.GiteeHost,
	}

	GTS.GetRoles()
	GLS.TranslateGroupsDescendant()
	//GLS.TranslateProjectsByGroup()
}
