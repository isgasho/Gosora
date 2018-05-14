package common

import (
	"html/template"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"./alerts"
	"./templates"
)

var Ctemplates []string
var Templates = template.New("")
var PrebuildTmplList []func(User, *Header) CTmpl

type CTmpl struct {
	Name       string
	Filename   string
	Path       string
	StructName string
	Data       interface{}
	Imports    []string
}

// TODO: Stop duplicating these bits of code
// nolint
func interpretedTopicTemplate(pi TopicPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["topic"]
	if !ok {
		mapping = "topic"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_topic_handle = interpretedTopicTemplate
var Template_topic_alt_handle = interpretedTopicTemplate

// nolint
var Template_topics_handle = func(pi TopicListPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["topics"]
	if !ok {
		mapping = "topics"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_forum_handle = func(pi ForumPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["forum"]
	if !ok {
		mapping = "forum"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_forums_handle = func(pi ForumsPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["forums"]
	if !ok {
		mapping = "forums"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_profile_handle = func(pi ProfilePage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["profile"]
	if !ok {
		mapping = "profile"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_create_topic_handle = func(pi CreateTopicPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["create_topic"]
	if !ok {
		mapping = "create_topic"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_login_handle = func(pi Page, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["login"]
	if !ok {
		mapping = "login"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_register_handle = func(pi Page, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["register"]
	if !ok {
		mapping = "register"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_error_handle = func(pi Page, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["error"]
	if !ok {
		mapping = "error"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// nolint
var Template_ip_search_handle = func(pi IPSearchPage, w io.Writer) error {
	mapping, ok := Themes[DefaultThemeBox.Load().(string)].TemplatesMap["ip_search"]
	if !ok {
		mapping = "ip_search"
	}
	return Templates.ExecuteTemplate(w, mapping+".html", pi)
}

// ? - Add template hooks?
func CompileTemplates() error {
	var config tmpl.CTemplateConfig
	config.Minify = Config.MinifyTemplates
	config.SuperDebug = Dev.TemplateDebug

	c := tmpl.NewCTemplateSet()
	c.SetConfig(config)
	c.SetBaseImportMap(map[string]string{
		"io":       "io",
		"./common": "./common",
	})

	// Schemas to train the template compiler on what to expect
	// TODO: Add support for interface{}s
	user := User{62, BuildProfileURL("fake-user", 62), "Fake User", "compiler@localhost", 0, false, false, false, false, false, false, GuestPerms, make(map[string]bool), "", false, BuildAvatar(62, ""), "", "", "", "", 0, 0, 0, "0.0.0.0.0", 0}
	// TODO: Do a more accurate level calculation for this?
	user2 := User{1, BuildProfileURL("admin-alice", 1), "Admin Alice", "alice@localhost", 1, true, true, true, true, false, false, AllPerms, make(map[string]bool), "", true, BuildAvatar(1, ""), "", "", "", "", 58, 1000, 0, "127.0.0.1", 0}
	user3 := User{2, BuildProfileURL("admin-fred", 62), "Admin Fred", "fred@localhost", 1, true, true, true, true, false, false, AllPerms, make(map[string]bool), "", true, BuildAvatar(2, ""), "", "", "", "", 42, 900, 0, "::1", 0}
	header := &Header{
		Site:        Site,
		Settings:    SettingBox.Load().(SettingMap),
		Themes:      Themes,
		Theme:       Themes[DefaultThemeBox.Load().(string)],
		CurrentUser: user,
		NoticeList:  []string{"test"},
		Stylesheets: []string{"panel"},
		Scripts:     []string{"whatever"},
		Widgets: PageWidgets{
			LeftSidebar: template.HTML("lalala"),
		},
	}

	log.Print("Compiling the templates")

	var now = time.Now()
	poll := Poll{ID: 1, Type: 0, Options: map[int]string{0: "Nothing", 1: "Something"}, Results: map[int]int{0: 5, 1: 2}, QuickOptions: []PollOption{
		PollOption{0, "Nothing"},
		PollOption{1, "Something"},
	}, VoteCount: 7}
	topic := TopicUser{1, "blah", "Blah", "Hey there!", 0, false, false, now, RelativeTime(now), now, RelativeTime(now), 0, "", "127.0.0.1", 0, 1, "classname", poll.ID, "weird-data", BuildProfileURL("fake-user", 62), "Fake User", Config.DefaultGroup, "", 0, "", "", "", "", "", 58, false}
	var replyList []ReplyUser
	replyList = append(replyList, ReplyUser{0, 0, "Yo!", "Yo!", 0, "alice", "Alice", Config.DefaultGroup, now, RelativeTime(now), 0, 0, "", "", 0, "", "", "", "", 0, "127.0.0.1", false, 1, "", ""})

	var varList = make(map[string]tmpl.VarItem)
	tpage := TopicPage{"Title", user, header, replyList, topic, poll, 1, 1}
	topicIDTmpl, err := c.Compile("topic.html", "templates/", "common.TopicPage", tpage, varList)
	if err != nil {
		return err
	}
	topicIDAltTmpl, err := c.Compile("topic_alt.html", "templates/", "common.TopicPage", tpage, varList)
	if err != nil {
		return err
	}

	varList = make(map[string]tmpl.VarItem)
	ppage := ProfilePage{"User 526", user, header, replyList, user}
	profileTmpl, err := c.Compile("profile.html", "templates/", "common.ProfilePage", ppage, varList)
	if err != nil {
		return err
	}

	// TODO: Use a dummy forum list to avoid o(n) problems
	var forumList []Forum
	forums, err := Forums.GetAll()
	if err != nil {
		return err
	}

	for _, forum := range forums {
		forumList = append(forumList, *forum)
	}
	varList = make(map[string]tmpl.VarItem)
	forumsPage := ForumsPage{"Forum List", user, header, forumList}
	forumsTmpl, err := c.Compile("forums.html", "templates/", "common.ForumsPage", forumsPage, varList)
	if err != nil {
		return err
	}

	var topicsList []*TopicsRow
	topicsList = append(topicsList, &TopicsRow{1, "topic-title", "Topic Title", "The topic content.", 1, false, false, "Date", time.Now(), "Date", user3.ID, 1, "", "127.0.0.1", 0, 1, "classname", "", &user2, "", 0, &user3, "General", "/forum/general.2"})
	header.Title = "Topic List"
	topicListPage := TopicListPage{header, topicsList, forumList, Config.DefaultForum, Paginator{[]int{1}, 1, 1}}
	topicListTmpl, err := c.Compile("topics.html", "templates/", "common.TopicListPage", topicListPage, varList)
	if err != nil {
		return err
	}

	//var topicList []TopicUser
	//topicList = append(topicList,TopicUser{1,"topic-title","Topic Title","The topic content.",1,false,false,"Date","Date",1,"","127.0.0.1",0,1,"classname","","admin-fred","Admin Fred",config.DefaultGroup,"",0,"","","","",58,false})
	forumItem := BlankForum(1, "general-forum.1", "General Forum", "Where the general stuff happens", true, "all", 0, "", 0)
	forumPage := ForumPage{"General Forum", user, header, topicsList, forumItem, Paginator{[]int{1}, 1, 1}}
	forumTmpl, err := c.Compile("forum.html", "templates/", "common.ForumPage", forumPage, varList)
	if err != nil {
		return err
	}

	loginPage := Page{"Login Page", user, header, tList, nil}
	loginTmpl, err := c.Compile("login.html", "templates/", "common.Page", loginPage, varList)
	if err != nil {
		return err
	}

	registerPage := Page{"Registration Page", user, header, tList, nil}
	registerTmpl, err := c.Compile("register.html", "templates/", "common.Page", registerPage, varList)
	if err != nil {
		return err
	}

	errorPage := Page{"Error", user, header, tList, "A problem has occurred in the system."}
	errorTmpl, err := c.Compile("error.html", "templates/", "common.Page", errorPage, varList)
	if err != nil {
		return err
	}

	var ipUserList = make(map[int]*User)
	ipUserList[1] = &user2
	ipSearchPage := IPSearchPage{"IP Search", user2, header, ipUserList, "::1"}
	ipSearchTmpl, err := c.Compile("ip_search.html", "templates/", "common.IPSearchPage", ipSearchPage, varList)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var writeTemplate = func(name string, content string) {
		log.Print("Writing template '" + name + "'")
		if content == "" {
			log.Fatal("No content body")
		}

		wg.Add(1)
		go func() {
			err := writeFile("./template_"+name+".go", content)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}

	// Let plugins register their own templates
	DebugLog("Registering the templates for the plugins")
	config = c.GetConfig()
	config.SkipHandles = true
	c.SetConfig(config)
	for _, tmplfunc := range PrebuildTmplList {
		tmplItem := tmplfunc(user, header)
		varList = make(map[string]tmpl.VarItem)
		compiledTmpl, err := c.Compile(tmplItem.Filename, tmplItem.Path, tmplItem.StructName, tmplItem.Data, varList, tmplItem.Imports...)
		if err != nil {
			return err
		}
		writeTemplate(tmplItem.Name, compiledTmpl)
	}

	log.Print("Writing the templates")
	writeTemplate("topic", topicIDTmpl)
	writeTemplate("topic_alt", topicIDAltTmpl)
	writeTemplate("profile", profileTmpl)
	writeTemplate("forums", forumsTmpl)
	writeTemplate("topics", topicListTmpl)
	writeTemplate("forum", forumTmpl)
	writeTemplate("login", loginTmpl)
	writeTemplate("register", registerTmpl)
	writeTemplate("ip_search", ipSearchTmpl)
	writeTemplate("error", errorTmpl)
	writeTemplateList(c, &wg, "./")
	return nil
}

func CompileJSTemplates() error {
	log.Print("Compiling the JS templates")
	var config tmpl.CTemplateConfig
	config.Minify = Config.MinifyTemplates
	config.SuperDebug = Dev.TemplateDebug
	config.SkipHandles = true
	config.SkipInitBlock = true
	config.PackageName = "tmpl"

	c := tmpl.NewCTemplateSet()
	c.SetConfig(config)
	c.SetBaseImportMap(map[string]string{
		"io":               "io",
		"../common/alerts": "../common/alerts",
	})
	var varList = make(map[string]tmpl.VarItem)

	// TODO: Check what sort of path is sent exactly and use it here
	alertItem := alerts.AlertItem{Avatar: "", ASID: 1, Path: "/", Message: "uh oh, something happened"}
	alertTmpl, err := c.Compile("alert.html", "templates/", "alerts.AlertItem", alertItem, varList)
	if err != nil {
		return err
	}

	var dirPrefix = "./tmpl_gen/"
	var wg sync.WaitGroup
	var writeTemplate = func(name string, content string) {
		log.Print("Writing template '" + name + "'")
		if content == "" {
			log.Fatal("No content body")
		}

		wg.Add(1)
		go func() {
			err := writeFile(dirPrefix+"template_"+name+".go", content)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	writeTemplate("alert", alertTmpl)
	writeTemplateList(c, &wg, dirPrefix)
	return nil
}

func writeTemplateList(c *tmpl.CTemplateSet, wg *sync.WaitGroup, prefix string) {
	log.Print("Writing template list")
	wg.Add(1)
	go func() {
		out := "package " + c.GetConfig().PackageName + "\n\n"
		for templateName, count := range c.TemplateFragmentCount {
			out += "var " + templateName + "_frags = make([][]byte," + strconv.Itoa(count) + ")\n"
			out += "\n// nolint\nfunc Get_" + templateName + "_frags() [][]byte {\nreturn " + templateName + "_frags\n}\n"
		}
		out += "\n// nolint\nfunc init() {\n" + c.FragOut + "}\n"
		err := writeFile(prefix+"template_list.go", out)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()
}

func arithToInt64(in interface{}) (out int64) {
	switch in := in.(type) {
	case int64:
		out = in
	case int32:
		out = int64(in)
	case int:
		out = int64(in)
	case uint32:
		out = int64(in)
	case uint16:
		out = int64(in)
	case uint8:
		out = int64(in)
	case uint:
		out = int64(in)
	}
	return out
}

func arithDuoToInt64(left interface{}, right interface{}) (leftInt int64, rightInt int64) {
	return arithToInt64(left), arithToInt64(right)
}

func InitTemplates() error {
	DebugLog("Initialising the template system")
	// TODO: Add support for floats
	fmap := make(map[string]interface{})
	fmap["add"] = func(left interface{}, right interface{}) interface{} {
		leftInt, rightInt := arithDuoToInt64(left, right)
		return leftInt + rightInt
	}

	fmap["subtract"] = func(left interface{}, right interface{}) interface{} {
		leftInt, rightInt := arithDuoToInt64(left, right)
		return leftInt - rightInt
	}

	fmap["multiply"] = func(left interface{}, right interface{}) interface{} {
		leftInt, rightInt := arithDuoToInt64(left, right)
		return leftInt * rightInt
	}

	fmap["divide"] = func(left interface{}, right interface{}) interface{} {
		leftInt, rightInt := arithDuoToInt64(left, right)
		if leftInt == 0 || rightInt == 0 {
			return 0
		}
		return leftInt / rightInt
	}

	fmap["dock"] = func(dock interface{}, headerInt interface{}) interface{} {
		return template.HTML(BuildWidget(dock.(string), headerInt.(*Header)))
	}

	fmap["lang"] = func(phraseNameInt interface{}) interface{} {
		phraseName, ok := phraseNameInt.(string)
		if !ok {
			panic("phraseNameInt is not a string")
		}
		return GetTmplPhrase(phraseName) // TODO: Log non-existent phrases?
	}

	fmap["scope"] = func(name interface{}) interface{} {
		return ""
	}

	// The interpreted templates...
	DebugLog("Loading the template files...")
	Templates.Funcs(fmap)
	templateFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		return err
	}

	var templateFileMap = make(map[string]int)
	for index, path := range templateFiles {
		path = strings.Replace(path, "\\", "/", -1)
		log.Print("templateFile: ", path)
		templateFileMap[path] = index
	}

	overrideFiles, err := filepath.Glob("templates/overrides/*.html")
	if err != nil {
		return err
	}
	for _, path := range overrideFiles {
		path = strings.Replace(path, "\\", "/", -1)
		log.Print("overrideFile: ", path)
		index, ok := templateFileMap["templates/"+strings.TrimPrefix(path, "templates/overrides/")]
		if !ok {
			log.Print("not ok: templates/" + strings.TrimPrefix(path, "templates/overrides/"))
			templateFiles = append(templateFiles, path)
			continue
		}
		templateFiles[index] = path
	}
	template.Must(Templates.ParseFiles(templateFiles...))
	template.Must(Templates.ParseGlob("pages/*"))

	return nil
}
