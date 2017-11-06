// +build pgsql

// This file was generated by Gosora's Query Generator. Please try to avoid modifying this file, as it might change at any time.
package main

import "log"
import "database/sql"

// nolint
type Stmts struct {
	addRepliesToTopic *sql.Stmt
	removeRepliesFromTopic *sql.Stmt
	addLikesToTopic *sql.Stmt
	addLikesToReply *sql.Stmt
	editTopic *sql.Stmt
	editReply *sql.Stmt
	stickTopic *sql.Stmt
	unstickTopic *sql.Stmt
	lockTopic *sql.Stmt
	unlockTopic *sql.Stmt
	updateLastIP *sql.Stmt
	updateSession *sql.Stmt
	setPassword *sql.Stmt
	setAvatar *sql.Stmt
	setUsername *sql.Stmt
	changeGroup *sql.Stmt
	activateUser *sql.Stmt
	updateUserLevel *sql.Stmt
	incrementUserScore *sql.Stmt
	incrementUserPosts *sql.Stmt
	incrementUserBigposts *sql.Stmt
	incrementUserMegaposts *sql.Stmt
	incrementUserTopics *sql.Stmt
	editProfileReply *sql.Stmt
	updateForum *sql.Stmt
	updateSetting *sql.Stmt
	updatePlugin *sql.Stmt
	updatePluginInstall *sql.Stmt
	updateTheme *sql.Stmt
	updateUser *sql.Stmt
	updateUserGroup *sql.Stmt
	updateGroupPerms *sql.Stmt
	updateGroupRank *sql.Stmt
	updateGroup *sql.Stmt
	updateEmail *sql.Stmt
	verifyEmail *sql.Stmt
	setTempGroup *sql.Stmt
	updateWordFilter *sql.Stmt
	bumpSync *sql.Stmt

	getActivityFeedByWatcher *sql.Stmt
	getActivityCountByWatcher *sql.Stmt
	todaysPostCount *sql.Stmt
	todaysTopicCount *sql.Stmt
	todaysReportCount *sql.Stmt
	todaysNewUserCount *sql.Stmt
	findUsersByIPUsers *sql.Stmt
	findUsersByIPTopics *sql.Stmt
	findUsersByIPReplies *sql.Stmt

	Mocks bool
}

// nolint
func _gen_pgsql() (err error) {
	if dev.DebugMode {
		log.Print("Building the generated statements")
	}
	
	log.Print("Preparing addRepliesToTopic statement.")
	stmts.addRepliesToTopic, err = db.Prepare("UPDATE `topics` SET `postCount` = `postCount` + ?,`lastReplyBy` = ?,`lastReplyAt` = LOCALTIMESTAMP() WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing removeRepliesFromTopic statement.")
	stmts.removeRepliesFromTopic, err = db.Prepare("UPDATE `topics` SET `postCount` = `postCount` - ? WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing addLikesToTopic statement.")
	stmts.addLikesToTopic, err = db.Prepare("UPDATE `topics` SET `likeCount` = `likeCount` + ? WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing addLikesToReply statement.")
	stmts.addLikesToReply, err = db.Prepare("UPDATE `replies` SET `likeCount` = `likeCount` + ? WHERE `rid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing editTopic statement.")
	stmts.editTopic, err = db.Prepare("UPDATE `topics` SET `title` = ?,`content` = ?,`parsed_content` = ? WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing editReply statement.")
	stmts.editReply, err = db.Prepare("UPDATE `replies` SET `content` = ?,`parsed_content` = ? WHERE `rid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing stickTopic statement.")
	stmts.stickTopic, err = db.Prepare("UPDATE `topics` SET `sticky` = 1 WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing unstickTopic statement.")
	stmts.unstickTopic, err = db.Prepare("UPDATE `topics` SET `sticky` = 0 WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing lockTopic statement.")
	stmts.lockTopic, err = db.Prepare("UPDATE `topics` SET `is_closed` = 1 WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing unlockTopic statement.")
	stmts.unlockTopic, err = db.Prepare("UPDATE `topics` SET `is_closed` = 0 WHERE `tid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateLastIP statement.")
	stmts.updateLastIP, err = db.Prepare("UPDATE `users` SET `last_ip` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateSession statement.")
	stmts.updateSession, err = db.Prepare("UPDATE `users` SET `session` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing setPassword statement.")
	stmts.setPassword, err = db.Prepare("UPDATE `users` SET `password` = ?,`salt` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing setAvatar statement.")
	stmts.setAvatar, err = db.Prepare("UPDATE `users` SET `avatar` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing setUsername statement.")
	stmts.setUsername, err = db.Prepare("UPDATE `users` SET `name` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing changeGroup statement.")
	stmts.changeGroup, err = db.Prepare("UPDATE `users` SET `group` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing activateUser statement.")
	stmts.activateUser, err = db.Prepare("UPDATE `users` SET `active` = 1 WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateUserLevel statement.")
	stmts.updateUserLevel, err = db.Prepare("UPDATE `users` SET `level` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing incrementUserScore statement.")
	stmts.incrementUserScore, err = db.Prepare("UPDATE `users` SET `score` = `score` + ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing incrementUserPosts statement.")
	stmts.incrementUserPosts, err = db.Prepare("UPDATE `users` SET `posts` = `posts` + ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing incrementUserBigposts statement.")
	stmts.incrementUserBigposts, err = db.Prepare("UPDATE `users` SET `posts` = `posts` + ?,`bigposts` = `bigposts` + ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing incrementUserMegaposts statement.")
	stmts.incrementUserMegaposts, err = db.Prepare("UPDATE `users` SET `posts` = `posts` + ?,`bigposts` = `bigposts` + ?,`megaposts` = `megaposts` + ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing incrementUserTopics statement.")
	stmts.incrementUserTopics, err = db.Prepare("UPDATE `users` SET `topics` = `topics` + ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing editProfileReply statement.")
	stmts.editProfileReply, err = db.Prepare("UPDATE `users_replies` SET `content` = ?,`parsed_content` = ? WHERE `rid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateForum statement.")
	stmts.updateForum, err = db.Prepare("UPDATE `forums` SET `name` = ?,`desc` = ?,`active` = ?,`preset` = ? WHERE `fid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateSetting statement.")
	stmts.updateSetting, err = db.Prepare("UPDATE `settings` SET `content` = ? WHERE `name` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updatePlugin statement.")
	stmts.updatePlugin, err = db.Prepare("UPDATE `plugins` SET `active` = ? WHERE `uname` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updatePluginInstall statement.")
	stmts.updatePluginInstall, err = db.Prepare("UPDATE `plugins` SET `installed` = ? WHERE `uname` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateTheme statement.")
	stmts.updateTheme, err = db.Prepare("UPDATE `themes` SET `default` = ? WHERE `uname` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateUser statement.")
	stmts.updateUser, err = db.Prepare("UPDATE `users` SET `name` = ?,`email` = ?,`group` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateUserGroup statement.")
	stmts.updateUserGroup, err = db.Prepare("UPDATE `users` SET `group` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateGroupPerms statement.")
	stmts.updateGroupPerms, err = db.Prepare("UPDATE `users_groups` SET `permissions` = ? WHERE `gid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateGroupRank statement.")
	stmts.updateGroupRank, err = db.Prepare("UPDATE `users_groups` SET `is_admin` = ?,`is_mod` = ?,`is_banned` = ? WHERE `gid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateGroup statement.")
	stmts.updateGroup, err = db.Prepare("UPDATE `users_groups` SET `name` = ?,`tag` = ? WHERE `gid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateEmail statement.")
	stmts.updateEmail, err = db.Prepare("UPDATE `emails` SET `email` = ?,`uid` = ?,`validated` = ?,`token` = ? WHERE `email` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing verifyEmail statement.")
	stmts.verifyEmail, err = db.Prepare("UPDATE `emails` SET `validated` = 1,`token` = '' WHERE `email` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing setTempGroup statement.")
	stmts.setTempGroup, err = db.Prepare("UPDATE `users` SET `temp_group` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing updateWordFilter statement.")
	stmts.updateWordFilter, err = db.Prepare("UPDATE `word_filters` SET `find` = ?,`replacement` = ? WHERE `wfid` = ?")
	if err != nil {
		return err
	}
		
	log.Print("Preparing bumpSync statement.")
	stmts.bumpSync, err = db.Prepare("UPDATE `sync` SET `last_update` = LOCALTIMESTAMP()")
	if err != nil {
		return err
	}
	
	return nil
}
