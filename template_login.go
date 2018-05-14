// +build !no_templategen

// Code generated by Gosora. More below:
/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main
import "io"
import "./common"

var login_tmpl_phrase_id int

// nolint
func init() {
	common.Template_login_handle = Template_login
	common.Ctemplates = append(common.Ctemplates,"login")
	common.TmplPtrMap["login"] = &common.Template_login_handle
	common.TmplPtrMap["o_login"] = Template_login
	login_tmpl_phrase_id = common.RegisterTmplPhraseNames([]string{
		"menu_hamburger_tooltip",
		"login_head",
		"login_account_name",
		"login_account_name",
		"login_account_password",
		"login_submit_button",
		"login_no_account",
		"footer_powered_by",
		"footer_made_with_love",
		"footer_theme_selector_aria",
	})
}

// nolint
func Template_login(tmpl_login_vars common.Page, w io.Writer) error {
var phrases = common.GetTmplPhrasesBytes(login_tmpl_phrase_id)
w.Write(header_frags[0])
w.Write([]byte(tmpl_login_vars.Title))
w.Write(header_frags[1])
w.Write([]byte(tmpl_login_vars.Header.Site.Name))
w.Write(header_frags[2])
w.Write([]byte(tmpl_login_vars.Header.Theme.Name))
w.Write(header_frags[3])
if len(tmpl_login_vars.Header.Stylesheets) != 0 {
for _, item := range tmpl_login_vars.Header.Stylesheets {
w.Write(header_frags[4])
w.Write([]byte(item))
w.Write(header_frags[5])
}
}
w.Write(header_frags[6])
if len(tmpl_login_vars.Header.Scripts) != 0 {
for _, item := range tmpl_login_vars.Header.Scripts {
w.Write(header_frags[7])
w.Write([]byte(item))
w.Write(header_frags[8])
}
}
w.Write(header_frags[9])
w.Write([]byte(tmpl_login_vars.CurrentUser.Session))
w.Write(header_frags[10])
w.Write([]byte(tmpl_login_vars.Header.Site.URL))
w.Write(header_frags[11])
if tmpl_login_vars.Header.MetaDesc != "" {
w.Write(header_frags[12])
w.Write([]byte(tmpl_login_vars.Header.MetaDesc))
w.Write(header_frags[13])
}
w.Write(header_frags[14])
if !tmpl_login_vars.CurrentUser.IsSuperMod {
w.Write(header_frags[15])
}
w.Write(header_frags[16])
w.Write([]byte(common.BuildWidget("leftOfNav",tmpl_login_vars.Header)))
w.Write(header_frags[17])
w.Write([]byte(tmpl_login_vars.Header.Site.ShortName))
w.Write(header_frags[18])
w.Write([]byte(common.BuildWidget("topMenu",tmpl_login_vars.Header)))
w.Write(header_frags[19])
w.Write(phrases[0])
w.Write(header_frags[20])
w.Write([]byte(common.BuildWidget("rightOfNav",tmpl_login_vars.Header)))
w.Write(header_frags[21])
if tmpl_login_vars.Header.Widgets.RightSidebar != "" {
w.Write(header_frags[22])
}
w.Write(header_frags[23])
if len(tmpl_login_vars.Header.NoticeList) != 0 {
for _, item := range tmpl_login_vars.Header.NoticeList {
w.Write(notice_frags[0])
w.Write([]byte(item))
w.Write(notice_frags[1])
}
}
w.Write(header_frags[24])
w.Write(login_frags[0])
w.Write(phrases[1])
w.Write(login_frags[1])
w.Write(phrases[2])
w.Write(login_frags[2])
w.Write(phrases[3])
w.Write(login_frags[3])
w.Write(phrases[4])
w.Write(login_frags[4])
w.Write(phrases[5])
w.Write(login_frags[5])
w.Write(phrases[6])
w.Write(login_frags[6])
w.Write(footer_frags[0])
w.Write([]byte(common.BuildWidget("footer",tmpl_login_vars.Header)))
w.Write(footer_frags[1])
w.Write(phrases[7])
w.Write(footer_frags[2])
w.Write(phrases[8])
w.Write(footer_frags[3])
w.Write(phrases[9])
w.Write(footer_frags[4])
if len(tmpl_login_vars.Header.Themes) != 0 {
for _, item := range tmpl_login_vars.Header.Themes {
if !item.HideFromThemes {
w.Write(footer_frags[5])
w.Write([]byte(item.Name))
w.Write(footer_frags[6])
if tmpl_login_vars.Header.Theme.Name == item.Name {
w.Write(footer_frags[7])
}
w.Write(footer_frags[8])
w.Write([]byte(item.FriendlyName))
w.Write(footer_frags[9])
}
}
}
w.Write(footer_frags[10])
w.Write([]byte(common.BuildWidget("rightSidebar",tmpl_login_vars.Header)))
w.Write(footer_frags[11])
return nil
}
