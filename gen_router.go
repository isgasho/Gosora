// Code generated by. DO NOT EDIT.
/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main

import (
	"log"
	"strings"
	"sync"
	"errors"
	"net/http"

	"./common"
)

var ErrNoRoute = errors.New("That route doesn't exist.")
// TODO: What about the /uploads/ route? x.x
var RouteMap = map[string]interface{}{ 
	"routeAPI": routeAPI,
	"routeOverview": routeOverview,
	"routeCustomPage": routeCustomPage,
	"routeForums": routeForums,
	"routeForum": routeForum,
	"routeChangeTheme": routeChangeTheme,
	"routeShowAttachment": routeShowAttachment,
	"routeReportSubmit": routeReportSubmit,
	"routeTopicCreate": routeTopicCreate,
	"routeTopics": routeTopics,
	"routePanelForums": routePanelForums,
	"routePanelForumsCreateSubmit": routePanelForumsCreateSubmit,
	"routePanelForumsDelete": routePanelForumsDelete,
	"routePanelForumsDeleteSubmit": routePanelForumsDeleteSubmit,
	"routePanelForumsEdit": routePanelForumsEdit,
	"routePanelForumsEditSubmit": routePanelForumsEditSubmit,
	"routePanelForumsEditPermsSubmit": routePanelForumsEditPermsSubmit,
	"routePanelForumsEditPermsAdvance": routePanelForumsEditPermsAdvance,
	"routePanelSettings": routePanelSettings,
	"routePanelSettingEdit": routePanelSettingEdit,
	"routePanelSettingEditSubmit": routePanelSettingEditSubmit,
	"routePanelWordFilters": routePanelWordFilters,
	"routePanelWordFiltersCreate": routePanelWordFiltersCreate,
	"routePanelWordFiltersEdit": routePanelWordFiltersEdit,
	"routePanelWordFiltersEditSubmit": routePanelWordFiltersEditSubmit,
	"routePanelWordFiltersDeleteSubmit": routePanelWordFiltersDeleteSubmit,
	"routePanelThemes": routePanelThemes,
	"routePanelThemesSetDefault": routePanelThemesSetDefault,
	"routePanelPlugins": routePanelPlugins,
	"routePanelPluginsActivate": routePanelPluginsActivate,
	"routePanelPluginsDeactivate": routePanelPluginsDeactivate,
	"routePanelPluginsInstall": routePanelPluginsInstall,
	"routePanelUsers": routePanelUsers,
	"routePanelUsersEdit": routePanelUsersEdit,
	"routePanelUsersEditSubmit": routePanelUsersEditSubmit,
	"routePanelAnalyticsViews": routePanelAnalyticsViews,
	"routePanelAnalyticsRoutes": routePanelAnalyticsRoutes,
	"routePanelAnalyticsRouteViews": routePanelAnalyticsRouteViews,
	"routePanelGroups": routePanelGroups,
	"routePanelGroupsEdit": routePanelGroupsEdit,
	"routePanelGroupsEditPerms": routePanelGroupsEditPerms,
	"routePanelGroupsEditSubmit": routePanelGroupsEditSubmit,
	"routePanelGroupsEditPermsSubmit": routePanelGroupsEditPermsSubmit,
	"routePanelGroupsCreateSubmit": routePanelGroupsCreateSubmit,
	"routePanelBackups": routePanelBackups,
	"routePanelLogsMod": routePanelLogsMod,
	"routePanelDebug": routePanelDebug,
	"routePanel": routePanel,
	"routeAccountEditCritical": routeAccountEditCritical,
	"routeAccountEditCriticalSubmit": routeAccountEditCriticalSubmit,
	"routeAccountEditAvatar": routeAccountEditAvatar,
	"routeAccountEditAvatarSubmit": routeAccountEditAvatarSubmit,
	"routeAccountEditUsername": routeAccountEditUsername,
	"routeAccountEditUsernameSubmit": routeAccountEditUsernameSubmit,
	"routeAccountEditEmail": routeAccountEditEmail,
	"routeAccountEditEmailTokenSubmit": routeAccountEditEmailTokenSubmit,
	"routeProfile": routeProfile,
	"routeBanSubmit": routeBanSubmit,
	"routeUnban": routeUnban,
	"routeActivate": routeActivate,
	"routeIps": routeIps,
	"routeDynamic": routeDynamic,
	"routeUploads": routeUploads,
}

// ! NEVER RELY ON THESE REMAINING THE SAME BETWEEN COMMITS
var routeMapEnum = map[string]int{ 
	"routeAPI": 0,
	"routeOverview": 1,
	"routeCustomPage": 2,
	"routeForums": 3,
	"routeForum": 4,
	"routeChangeTheme": 5,
	"routeShowAttachment": 6,
	"routeReportSubmit": 7,
	"routeTopicCreate": 8,
	"routeTopics": 9,
	"routePanelForums": 10,
	"routePanelForumsCreateSubmit": 11,
	"routePanelForumsDelete": 12,
	"routePanelForumsDeleteSubmit": 13,
	"routePanelForumsEdit": 14,
	"routePanelForumsEditSubmit": 15,
	"routePanelForumsEditPermsSubmit": 16,
	"routePanelForumsEditPermsAdvance": 17,
	"routePanelSettings": 18,
	"routePanelSettingEdit": 19,
	"routePanelSettingEditSubmit": 20,
	"routePanelWordFilters": 21,
	"routePanelWordFiltersCreate": 22,
	"routePanelWordFiltersEdit": 23,
	"routePanelWordFiltersEditSubmit": 24,
	"routePanelWordFiltersDeleteSubmit": 25,
	"routePanelThemes": 26,
	"routePanelThemesSetDefault": 27,
	"routePanelPlugins": 28,
	"routePanelPluginsActivate": 29,
	"routePanelPluginsDeactivate": 30,
	"routePanelPluginsInstall": 31,
	"routePanelUsers": 32,
	"routePanelUsersEdit": 33,
	"routePanelUsersEditSubmit": 34,
	"routePanelAnalyticsViews": 35,
	"routePanelAnalyticsRoutes": 36,
	"routePanelAnalyticsRouteViews": 37,
	"routePanelGroups": 38,
	"routePanelGroupsEdit": 39,
	"routePanelGroupsEditPerms": 40,
	"routePanelGroupsEditSubmit": 41,
	"routePanelGroupsEditPermsSubmit": 42,
	"routePanelGroupsCreateSubmit": 43,
	"routePanelBackups": 44,
	"routePanelLogsMod": 45,
	"routePanelDebug": 46,
	"routePanel": 47,
	"routeAccountEditCritical": 48,
	"routeAccountEditCriticalSubmit": 49,
	"routeAccountEditAvatar": 50,
	"routeAccountEditAvatarSubmit": 51,
	"routeAccountEditUsername": 52,
	"routeAccountEditUsernameSubmit": 53,
	"routeAccountEditEmail": 54,
	"routeAccountEditEmailTokenSubmit": 55,
	"routeProfile": 56,
	"routeBanSubmit": 57,
	"routeUnban": 58,
	"routeActivate": 59,
	"routeIps": 60,
	"routeDynamic": 61,
	"routeUploads": 62,
}
var reverseRouteMapEnum = map[int]string{ 
	0: "routeAPI",
	1: "routeOverview",
	2: "routeCustomPage",
	3: "routeForums",
	4: "routeForum",
	5: "routeChangeTheme",
	6: "routeShowAttachment",
	7: "routeReportSubmit",
	8: "routeTopicCreate",
	9: "routeTopics",
	10: "routePanelForums",
	11: "routePanelForumsCreateSubmit",
	12: "routePanelForumsDelete",
	13: "routePanelForumsDeleteSubmit",
	14: "routePanelForumsEdit",
	15: "routePanelForumsEditSubmit",
	16: "routePanelForumsEditPermsSubmit",
	17: "routePanelForumsEditPermsAdvance",
	18: "routePanelSettings",
	19: "routePanelSettingEdit",
	20: "routePanelSettingEditSubmit",
	21: "routePanelWordFilters",
	22: "routePanelWordFiltersCreate",
	23: "routePanelWordFiltersEdit",
	24: "routePanelWordFiltersEditSubmit",
	25: "routePanelWordFiltersDeleteSubmit",
	26: "routePanelThemes",
	27: "routePanelThemesSetDefault",
	28: "routePanelPlugins",
	29: "routePanelPluginsActivate",
	30: "routePanelPluginsDeactivate",
	31: "routePanelPluginsInstall",
	32: "routePanelUsers",
	33: "routePanelUsersEdit",
	34: "routePanelUsersEditSubmit",
	35: "routePanelAnalyticsViews",
	36: "routePanelAnalyticsRoutes",
	37: "routePanelAnalyticsRouteViews",
	38: "routePanelGroups",
	39: "routePanelGroupsEdit",
	40: "routePanelGroupsEditPerms",
	41: "routePanelGroupsEditSubmit",
	42: "routePanelGroupsEditPermsSubmit",
	43: "routePanelGroupsCreateSubmit",
	44: "routePanelBackups",
	45: "routePanelLogsMod",
	46: "routePanelDebug",
	47: "routePanel",
	48: "routeAccountEditCritical",
	49: "routeAccountEditCriticalSubmit",
	50: "routeAccountEditAvatar",
	51: "routeAccountEditAvatarSubmit",
	52: "routeAccountEditUsername",
	53: "routeAccountEditUsernameSubmit",
	54: "routeAccountEditEmail",
	55: "routeAccountEditEmailTokenSubmit",
	56: "routeProfile",
	57: "routeBanSubmit",
	58: "routeUnban",
	59: "routeActivate",
	60: "routeIps",
	61: "routeDynamic",
	62: "routeUploads",
}

// TODO: Stop spilling these into the package scope?
func init() {
	common.SetRouteMapEnum(routeMapEnum)
	common.SetReverseRouteMapEnum(reverseRouteMapEnum)
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
		extraData = req.URL.Path[strings.LastIndexByte(req.URL.Path,'/') + 1:]
		req.URL.Path = req.URL.Path[:strings.LastIndexByte(req.URL.Path,'/') + 1]
	}
	
	if common.Dev.SuperDebug {
		log.Print("before routeStatic")
		log.Print("prefix: ", prefix)
		log.Print("req.URL.Path: ", req.URL.Path)
		log.Print("extraData: ", extraData)
		log.Print("req.Referer(): ", req.Referer())
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
	switch(prefix) {
		case "/api":
			common.RouteViewCounter.Bump(0)
			err = routeAPI(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/overview":
			common.RouteViewCounter.Bump(1)
			err = routeOverview(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/pages":
			common.RouteViewCounter.Bump(2)
			err = routeCustomPage(w,req,user,extraData)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/forums":
			common.RouteViewCounter.Bump(3)
			err = routeForums(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/forum":
			common.RouteViewCounter.Bump(4)
			err = routeForum(w,req,user,extraData)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/theme":
				err = common.ParseForm(w,req,user)
				if err != nil {
					router.handleError(err,w,req,user)
					return
				}
				
			common.RouteViewCounter.Bump(5)
			err = routeChangeTheme(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/attachs":
				err = common.ParseForm(w,req,user)
				if err != nil {
					router.handleError(err,w,req,user)
					return
				}
				
			common.RouteViewCounter.Bump(6)
			err = routeShowAttachment(w,req,user,extraData)
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/report":
			err = common.NoBanned(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
				return
			}
			
			switch(req.URL.Path) {
				case "/report/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(7)
					err = routeReportSubmit(w,req,user,extraData)
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/topics":
			switch(req.URL.Path) {
				case "/topics/create/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(8)
					err = routeTopicCreate(w,req,user,extraData)
				default:
					common.RouteViewCounter.Bump(9)
					err = routeTopics(w,req,user)
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/panel":
			err = common.SuperModOnly(w,req,user)
			if err != nil {
				router.handleError(err,w,req,user)
				return
			}
			
			switch(req.URL.Path) {
				case "/panel/forums/":
					common.RouteViewCounter.Bump(10)
					err = routePanelForums(w,req,user)
				case "/panel/forums/create/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(11)
					err = routePanelForumsCreateSubmit(w,req,user)
				case "/panel/forums/delete/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(12)
					err = routePanelForumsDelete(w,req,user,extraData)
				case "/panel/forums/delete/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(13)
					err = routePanelForumsDeleteSubmit(w,req,user,extraData)
				case "/panel/forums/edit/":
					common.RouteViewCounter.Bump(14)
					err = routePanelForumsEdit(w,req,user,extraData)
				case "/panel/forums/edit/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(15)
					err = routePanelForumsEditSubmit(w,req,user,extraData)
				case "/panel/forums/edit/perms/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(16)
					err = routePanelForumsEditPermsSubmit(w,req,user,extraData)
				case "/panel/forums/edit/perms/":
					common.RouteViewCounter.Bump(17)
					err = routePanelForumsEditPermsAdvance(w,req,user,extraData)
				case "/panel/settings/":
					common.RouteViewCounter.Bump(18)
					err = routePanelSettings(w,req,user)
				case "/panel/settings/edit/":
					common.RouteViewCounter.Bump(19)
					err = routePanelSettingEdit(w,req,user,extraData)
				case "/panel/settings/edit/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(20)
					err = routePanelSettingEditSubmit(w,req,user,extraData)
				case "/panel/settings/word-filters/":
					common.RouteViewCounter.Bump(21)
					err = routePanelWordFilters(w,req,user)
				case "/panel/settings/word-filters/create/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(22)
					err = routePanelWordFiltersCreate(w,req,user)
				case "/panel/settings/word-filters/edit/":
					common.RouteViewCounter.Bump(23)
					err = routePanelWordFiltersEdit(w,req,user,extraData)
				case "/panel/settings/word-filters/edit/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(24)
					err = routePanelWordFiltersEditSubmit(w,req,user,extraData)
				case "/panel/settings/word-filters/delete/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(25)
					err = routePanelWordFiltersDeleteSubmit(w,req,user,extraData)
				case "/panel/themes/":
					common.RouteViewCounter.Bump(26)
					err = routePanelThemes(w,req,user)
				case "/panel/themes/default/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(27)
					err = routePanelThemesSetDefault(w,req,user,extraData)
				case "/panel/plugins/":
					common.RouteViewCounter.Bump(28)
					err = routePanelPlugins(w,req,user)
				case "/panel/plugins/activate/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(29)
					err = routePanelPluginsActivate(w,req,user,extraData)
				case "/panel/plugins/deactivate/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(30)
					err = routePanelPluginsDeactivate(w,req,user,extraData)
				case "/panel/plugins/install/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(31)
					err = routePanelPluginsInstall(w,req,user,extraData)
				case "/panel/users/":
					common.RouteViewCounter.Bump(32)
					err = routePanelUsers(w,req,user)
				case "/panel/users/edit/":
					common.RouteViewCounter.Bump(33)
					err = routePanelUsersEdit(w,req,user,extraData)
				case "/panel/users/edit/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(34)
					err = routePanelUsersEditSubmit(w,req,user,extraData)
				case "/panel/analytics/views/":
					err = common.ParseForm(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(35)
					err = routePanelAnalyticsViews(w,req,user)
				case "/panel/analytics/routes/":
					common.RouteViewCounter.Bump(36)
					err = routePanelAnalyticsRoutes(w,req,user)
				case "/panel/analytics/route/":
					common.RouteViewCounter.Bump(37)
					err = routePanelAnalyticsRouteViews(w,req,user,extraData)
				case "/panel/groups/":
					common.RouteViewCounter.Bump(38)
					err = routePanelGroups(w,req,user)
				case "/panel/groups/edit/":
					common.RouteViewCounter.Bump(39)
					err = routePanelGroupsEdit(w,req,user,extraData)
				case "/panel/groups/edit/perms/":
					common.RouteViewCounter.Bump(40)
					err = routePanelGroupsEditPerms(w,req,user,extraData)
				case "/panel/groups/edit/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(41)
					err = routePanelGroupsEditSubmit(w,req,user,extraData)
				case "/panel/groups/edit/perms/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(42)
					err = routePanelGroupsEditPermsSubmit(w,req,user,extraData)
				case "/panel/groups/create/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(43)
					err = routePanelGroupsCreateSubmit(w,req,user)
				case "/panel/backups/":
					err = common.SuperAdminOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(44)
					err = routePanelBackups(w,req,user,extraData)
				case "/panel/logs/mod/":
					common.RouteViewCounter.Bump(45)
					err = routePanelLogsMod(w,req,user)
				case "/panel/debug/":
					err = common.AdminOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(46)
					err = routePanelDebug(w,req,user)
				default:
					common.RouteViewCounter.Bump(47)
					err = routePanel(w,req,user)
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/user":
			switch(req.URL.Path) {
				case "/user/edit/critical/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(48)
					err = routeAccountEditCritical(w,req,user)
				case "/user/edit/critical/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(49)
					err = routeAccountEditCriticalSubmit(w,req,user)
				case "/user/edit/avatar/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(50)
					err = routeAccountEditAvatar(w,req,user)
				case "/user/edit/avatar/submit/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(51)
					err = routeAccountEditAvatarSubmit(w,req,user)
				case "/user/edit/username/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(52)
					err = routeAccountEditUsername(w,req,user)
				case "/user/edit/username/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(53)
					err = routeAccountEditUsernameSubmit(w,req,user)
				case "/user/edit/email/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(54)
					err = routeAccountEditEmail(w,req,user)
				case "/user/edit/token/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(55)
					err = routeAccountEditEmailTokenSubmit(w,req,user,extraData)
				default:
					req.URL.Path += extraData
					common.RouteViewCounter.Bump(56)
					err = routeProfile(w,req,user)
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}
		case "/users":
			switch(req.URL.Path) {
				case "/users/ban/submit/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(57)
					err = routeBanSubmit(w,req,user,extraData)
				case "/users/unban/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(58)
					err = routeUnban(w,req,user,extraData)
				case "/users/activate/":
					err = common.NoSessionMismatch(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(59)
					err = routeActivate(w,req,user,extraData)
				case "/users/ips/":
					err = common.MemberOnly(w,req,user)
					if err != nil {
						router.handleError(err,w,req,user)
						return
					}
					
					common.RouteViewCounter.Bump(60)
					err = routeIps(w,req,user)
			}
			if err != nil {
				router.handleError(err,w,req,user)
			}
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
			common.RouteViewCounter.Bump(62)
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
				common.RouteViewCounter.Bump(61) // TODO: Be more specific about *which* dynamic route it is
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
