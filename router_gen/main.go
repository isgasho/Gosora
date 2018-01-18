/* WIP Under Construction */
package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"text/template"
)

var routeList []*RouteImpl
var routeGroups []*RouteGroup

type TmplVars struct {
	RouteList     []*RouteImpl
	RouteGroups   []*RouteGroup
	AllRouteNames []string
	AllRouteMap   map[string]int
	AllAgentNames []string
	AllAgentMap   map[string]int
}

func main() {
	log.Println("Generating the router...")

	// Load all the routes...
	routes()

	var tmplVars = TmplVars{
		RouteList:   routeList,
		RouteGroups: routeGroups,
	}
	var allRouteNames []string
	var allRouteMap = make(map[string]int)

	var out string
	var mapIt = func(name string) {
		allRouteNames = append(allRouteNames, name)
		allRouteMap[name] = len(allRouteNames) - 1
	}
	var countToIndents = func(indent int) (indentor string) {
		for i := 0; i < indent; i++ {
			indentor += "\t"
		}
		return indentor
	}
	var runBefore = func(runnables []Runnable, indent int) (out string) {
		var indentor = countToIndents(indent)
		if len(runnables) > 0 {
			for _, runnable := range runnables {
				if runnable.Literal {
					out += "\n\t" + indentor + runnable.Contents
				} else {
					out += "\n" + indentor + "err = common." + runnable.Contents + "(w,req,user)\n" +
						indentor + "if err != nil {\n" +
						indentor + "\trouter.handleError(err,w,req,user)\n" +
						indentor + "\treturn\n" +
						indentor + "}\n" + indentor
				}
			}
		}
		return out
	}

	for _, route := range routeList {
		mapIt(route.Name)
		var end = len(route.Path) - 1
		out += "\n\t\tcase \"" + route.Path[0:end] + "\":"
		out += runBefore(route.RunBefore, 4)
		out += "\n\t\t\tcommon.RouteViewCounter.Bump(" + strconv.Itoa(allRouteMap[route.Name]) + ")"
		out += "\n\t\t\terr = " + route.Name + "(w,req,user"
		for _, item := range route.Vars {
			out += "," + item
		}
		out += `)
			if err != nil {
				router.handleError(err,w,req,user)
			}`
	}

	for _, group := range routeGroups {
		var end = len(group.Path) - 1
		out += "\n\t\tcase \"" + group.Path[0:end] + "\":"
		out += runBefore(group.RunBefore, 3)
		out += "\n\t\t\tswitch(req.URL.Path) {"

		var defaultRoute = blankRoute()
		for _, route := range group.RouteList {
			if group.Path == route.Path {
				defaultRoute = route
				continue
			}
			mapIt(route.Name)

			out += "\n\t\t\t\tcase \"" + route.Path + "\":"
			if len(route.RunBefore) > 0 {
			skipRunnable:
				for _, runnable := range route.RunBefore {
					for _, gRunnable := range group.RunBefore {
						if gRunnable.Contents == runnable.Contents {
							continue
						}
						// TODO: Stop hard-coding these
						if gRunnable.Contents == "AdminOnly" && runnable.Contents == "MemberOnly" {
							continue skipRunnable
						}
						if gRunnable.Contents == "AdminOnly" && runnable.Contents == "SuperModOnly" {
							continue skipRunnable
						}
						if gRunnable.Contents == "SuperModOnly" && runnable.Contents == "MemberOnly" {
							continue skipRunnable
						}
					}

					if runnable.Literal {
						out += "\n\t\t\t\t\t" + runnable.Contents
					} else {
						out += `
					err = common.` + runnable.Contents + `(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					`
					}
				}
			}
			out += "\n\t\t\t\t\tcommon.RouteViewCounter.Bump(" + strconv.Itoa(allRouteMap[route.Name]) + ")"
			out += "\n\t\t\t\t\terr = " + route.Name + "(w,req,user"
			for _, item := range route.Vars {
				out += "," + item
			}
			out += ")"
		}

		if defaultRoute.Name != "" {
			mapIt(defaultRoute.Name)
			out += "\n\t\t\t\tdefault:"
			out += runBefore(defaultRoute.RunBefore, 4)
			out += "\n\t\t\t\t\tcommon.RouteViewCounter.Bump(" + strconv.Itoa(allRouteMap[defaultRoute.Name]) + ")"
			out += "\n\t\t\t\t\terr = " + defaultRoute.Name + "(w,req,user"
			for _, item := range defaultRoute.Vars {
				out += ", " + item
			}
			out += ")"
		}
		out += `
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}`
	}

	// Stubs for us to refer to these routes through
	mapIt("routeDynamic")
	mapIt("routeUploads")
	tmplVars.AllRouteNames = allRouteNames
	tmplVars.AllRouteMap = allRouteMap
	tmplVars.AllAgentNames = []string{
		"unknown",
		"firefox",
		"chrome",
		"opera",
		"safari",
		"edge",
		"internetexplorer",

		"googlebot",
		"yandex",
		"bing",
		"baidu",
		"duckduckgo",
		"discord",
		"cloudflarealwayson",
		"lynx",
		"blank",
	}

	tmplVars.AllAgentMap = make(map[string]int)
	for id, agent := range tmplVars.AllAgentNames {
		tmplVars.AllAgentMap[agent] = id
	}

	var fileData = `// Code generated by. DO NOT EDIT.
/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main

import (
	"log"
	"strings"
	"sync"
	"errors"
	"net/http"

	"./common"
	"./routes"
)

var ErrNoRoute = errors.New("That route doesn't exist.")
// TODO: What about the /uploads/ route? x.x
var RouteMap = map[string]interface{}{ {{range .AllRouteNames}}
	"{{.}}": {{.}},{{end}}
}

// ! NEVER RELY ON THESE REMAINING THE SAME BETWEEN COMMITS
var routeMapEnum = map[string]int{ {{range $index, $element := .AllRouteNames}}
	"{{$element}}": {{$index}},{{end}}
}
var reverseRouteMapEnum = map[int]string{ {{range $index, $element := .AllRouteNames}}
	{{$index}}: "{{$element}}",{{end}}
}
var agentMapEnum = map[string]int{ {{range $index, $element := .AllAgentNames}}
	"{{$element}}": {{$index}},{{end}}
}
var reverseAgentMapEnum = map[int]string{ {{range $index, $element := .AllAgentNames}}
	{{$index}}: "{{$element}}",{{end}}
}

// TODO: Stop spilling these into the package scope?
func init() {
	common.SetRouteMapEnum(routeMapEnum)
	common.SetReverseRouteMapEnum(reverseRouteMapEnum)
	common.SetAgentMapEnum(agentMapEnum)
	common.SetReverseAgentMapEnum(reverseAgentMapEnum)
}

type GenRouter struct {
	UploadHandler func(http.ResponseWriter, *http.Request)
	extraRoutes map[string]func(http.ResponseWriter, *http.Request, common.User) common.RouteError
	
	sync.RWMutex
}

func NewGenRouter(uploads http.Handler) *GenRouter {
	return &GenRouter{
		UploadHandler: http.StripPrefix("/uploads/",uploads).ServeHTTP,
		extraRoutes: make(map[string]func(http.ResponseWriter, *http.Request, common.User) common.RouteError),
	}
}

func (router *GenRouter) handleError(err common.RouteError, w http.ResponseWriter, r *http.Request, user common.User) {
	if err.Handled() {
		return
	}
	
	if err.Type() == "system" {
		common.InternalErrorJSQ(err, w, r, err.JSON())
		return
	}
	common.LocalErrorJSQ(err.Error(), w, r, user,err.JSON())
}

func (router *GenRouter) Handle(_ string, _ http.Handler) {
}

func (router *GenRouter) HandleFunc(pattern string, handle func(http.ResponseWriter, *http.Request, common.User) common.RouteError) {
	router.Lock()
	defer router.Unlock()
	router.extraRoutes[pattern] = handle
}

func (router *GenRouter) RemoveFunc(pattern string) error {
	router.Lock()
	defer router.Unlock()
	_, ok := router.extraRoutes[pattern]
	if !ok {
		return ErrNoRoute
	}
	delete(router.extraRoutes, pattern)
	return nil
}

// TODO: Pass the default route or config struct to the router rather than accessing it via a package global
// TODO: SetDefaultRoute
// TODO: GetDefaultRoute

func (router *GenRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) == 0 || req.URL.Path[0] != '/' {
		w.WriteHeader(405)
		w.Write([]byte(""))
		return
	}
	
	var prefix, extraData string
	prefix = req.URL.Path[0:strings.IndexByte(req.URL.Path[1:],'/') + 1]
	if req.URL.Path[len(req.URL.Path) - 1] != '/' {
		// TODO: Cover more suspicious strings and at a lower layer than this
		for _, char := range req.URL.Path {
			if char != '&' && !(char > 44 && char < 58) && char != '=' && char != '?' && !(char > 64 && char < 91) && char != '\\' && char != '_' && !(char > 96 && char < 123) {
				log.Print("Suspicious UA: ", req.UserAgent())
				log.Print("Method: ", req.Method)
				for key, value := range req.Header {
					for _, vvalue := range value {
						log.Print("Header '" + key + "': " + vvalue + "!!")
					}
				}
				log.Print("req.URL.Path: ", req.URL.Path)
				log.Print("req.Referer(): ", req.Referer())
				log.Print("req.RemoteAddr: ", req.RemoteAddr)
			}
		}
		if strings.Contains(req.URL.Path,"..") || strings.Contains(req.URL.Path,"--") {
			log.Print("Suspicious UA: ", req.UserAgent())
			log.Print("Method: ", req.Method)
			for key, value := range req.Header {
				for _, vvalue := range value {
					log.Print("Header '" + key + "': " + vvalue + "!!")
				}
			}
			log.Print("req.URL.Path: ", req.URL.Path)
			log.Print("req.Referer(): ", req.Referer())
			log.Print("req.RemoteAddr: ", req.RemoteAddr)
		}
		extraData = req.URL.Path[strings.LastIndexByte(req.URL.Path,'/') + 1:]
		req.URL.Path = req.URL.Path[:strings.LastIndexByte(req.URL.Path,'/') + 1]
	}
	
	if common.Dev.SuperDebug {
		log.Print("before routeStatic")
		log.Print("Method: ", req.Method)
		for key, value := range req.Header {
			for _, vvalue := range value {
				log.Print("Header '" + key + "': " + vvalue + "!!")
			}
		}
		log.Print("prefix: ", prefix)
		log.Print("req.URL.Path: ", req.URL.Path)
		log.Print("extraData: ", extraData)
		log.Print("req.Referer(): ", req.Referer())
		log.Print("req.RemoteAddr: ", req.RemoteAddr)
	}
	
	if prefix == "/static" {
		req.URL.Path += extraData
		routeStatic(w, req)
		return
	}
	if common.Dev.SuperDebug {
		log.Print("before PreRoute")
	}

	// Increment the global view counter
	common.GlobalViewCounter.Bump()

	// Track the user agents. Unfortunately, everyone pretends to be Mozilla, so this'll be a little less efficient than I would like.
	// TODO: Add a setting to disable this?
	// TODO: Use a more efficient detector instead of smashing every possible combination in
	ua := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(req.UserAgent(),"Mozilla/5.0 ")," Safari/537.36")) // Noise, no one's going to be running this and it complicates implementing an efficient UA parser, particularly the more efficient right-to-left one I have in mind
	switch {
	case strings.Contains(ua,"Google"):
		common.AgentViewCounter.Bump({{.AllAgentMap.googlebot}})
	case strings.Contains(ua,"Yandex"):
		common.AgentViewCounter.Bump({{.AllAgentMap.yandex}})
	case strings.Contains(ua,"bingbot"), strings.Contains(ua,"BingPreview"):
		common.AgentViewCounter.Bump({{.AllAgentMap.bing}})
	case strings.Contains(ua,"OPR"): // Pretends to be Chrome, needs to run before that
		common.AgentViewCounter.Bump({{.AllAgentMap.opera}})
	case strings.Contains(ua,"Edge"):
		common.AgentViewCounter.Bump({{.AllAgentMap.edge}})
	case strings.Contains(ua,"Chrome"):
		common.AgentViewCounter.Bump({{.AllAgentMap.chrome}})
	case strings.Contains(ua,"Firefox"):
		common.AgentViewCounter.Bump({{.AllAgentMap.firefox}})
	case strings.Contains(ua,"Safari"):
		common.AgentViewCounter.Bump({{.AllAgentMap.safari}})
	case strings.Contains(ua,"MSIE"):
		common.AgentViewCounter.Bump({{.AllAgentMap.internetexplorer}})
	case strings.Contains(ua,"Baiduspider"):
		common.AgentViewCounter.Bump({{.AllAgentMap.baidu}})
	case strings.Contains(ua,"DuckDuckBot"):
		common.AgentViewCounter.Bump({{.AllAgentMap.duckduckgo}})
	case strings.Contains(ua,"Discordbot"):
		common.AgentViewCounter.Bump({{.AllAgentMap.discord}})
	case strings.Contains(ua,"Lynx"):
		common.AgentViewCounter.Bump({{.AllAgentMap.lynx}})
	case strings.Contains(ua,"CloudFlare-AlwaysOnline"):
		common.AgentViewCounter.Bump({{.AllAgentMap.cloudflarealwayson}})
	case ua == "":
		common.AgentViewCounter.Bump({{.AllAgentMap.blank}})
		if common.Dev.DebugMode {
			log.Print("Blank UA: ", req.UserAgent())
			log.Print("Method: ", req.Method)

			for key, value := range req.Header {
				for _, vvalue := range value {
					log.Print("Header '" + key + "': " + vvalue + "!!")
				}
			}
			log.Print("prefix: ", prefix)
			log.Print("req.URL.Path: ", req.URL.Path)
			log.Print("extraData: ", extraData)
			log.Print("req.Referer(): ", req.Referer())
			log.Print("req.RemoteAddr: ", req.RemoteAddr)
		}
	default:
		common.AgentViewCounter.Bump({{.AllAgentMap.unknown}})
		if common.Dev.DebugMode {
			log.Print("Unknown UA: ", req.UserAgent())
			log.Print("Method: ", req.Method)
			for key, value := range req.Header {
				for _, vvalue := range value {
					log.Print("Header '" + key + "': " + vvalue + "!!")
				}
			}
			log.Print("prefix: ", prefix)
			log.Print("req.URL.Path: ", req.URL.Path)
			log.Print("extraData: ", extraData)
			log.Print("req.Referer(): ", req.Referer())
			log.Print("req.RemoteAddr: ", req.RemoteAddr)
		}
	}
	
	// Deal with the session stuff, etc.
	user, ok := common.PreRoute(w, req)
	if !ok {
		return
	}
	if common.Dev.SuperDebug {
		log.Print("after PreRoute")
		log.Print("routeMapEnum: ", routeMapEnum)
	}
	
	var err common.RouteError
	switch(prefix) {` + out + `
		/*case "/sitemaps": // TODO: Count these views
			req.URL.Path += extraData
			err = sitemapSwitch(w,req)
			if err != nil {
				router.handleError(err,w,req,user)
			}*/
		case "/uploads":
			if extraData == "" {
				common.NotFound(w,req)
				return
			}
			common.RouteViewCounter.Bump({{.AllRouteMap.routeUploads}})
			req.URL.Path += extraData
			// TODO: Find a way to propagate errors up from this?
			router.UploadHandler(w,req) // TODO: Count these views
		case "":
			// Stop the favicons, robots.txt file, etc. resolving to the topics list
			// TODO: Add support for favicons and robots.txt files
			switch(extraData) {
				case "robots.txt":
					err = routeRobotsTxt(w,req) // TODO: Count these views
					if err != nil {
						router.handleError(err,w,req,user)
					}
					return
				/*case "sitemap.xml":
					err = routeSitemapXml(w,req) // TODO: Count these views
					if err != nil {
						router.handleError(err,w,req,user)
					}
					return*/
			}
			
			if extraData != "" {
				common.NotFound(w,req)
				return
			}

			handle, ok := RouteMap[common.Config.DefaultRoute]
			if !ok {
				// TODO: Make this a startup error not a runtime one
				log.Print("Unable to find the default route")
				common.NotFound(w,req)
				return
			}
			common.RouteViewCounter.Bump(routeMapEnum[common.Config.DefaultRoute])

			handle.(func(http.ResponseWriter, *http.Request, common.User) common.RouteError)(w,req,user)
		default:
			// A fallback for the routes which haven't been converted to the new router yet or plugins
			router.RLock()
			handle, ok := router.extraRoutes[req.URL.Path]
			router.RUnlock()
			
			if ok {
				common.RouteViewCounter.Bump({{.AllRouteMap.routeDynamic}}) // TODO: Be more specific about *which* dynamic route it is
				req.URL.Path += extraData
				err = handle(w,req,user)
				if err != nil {
					router.handleError(err,w,req,user)
				}
				return
			}
			common.NotFound(w,req) // TODO: Collect all the error view counts so we can add a replacement for GlobalViewCounter by adding up the view counts of every route? Complex and may be inaccurate, collecting it globally and locally would at-least help find places we aren't capturing views
	}
}
`
	var tmpl = template.Must(template.New("router").Parse(fileData))
	var b bytes.Buffer
	err := tmpl.Execute(&b, tmplVars)
	if err != nil {
		log.Fatal(err)
	}

	writeFile("./gen_router.go", string(b.Bytes()))
	log.Println("Successfully generated the router")
}

func writeFile(name string, content string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
	f.Close()
}
