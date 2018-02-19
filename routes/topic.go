package routes

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../common"
	"../common/counters"
	"../query_gen/lib"
)

type TopicStmts struct {
	getReplies *sql.Stmt
}

var topicStmts TopicStmts

// TODO: Move these DbInits into a TopicList abstraction
func init() {
	common.DbInits.Add(func(acc *qgen.Accumulator) error {
		topicStmts = TopicStmts{
			getReplies: acc.SimpleLeftJoin("replies", "users", "replies.rid, replies.content, replies.createdBy, replies.createdAt, replies.lastEdit, replies.lastEditBy, users.avatar, users.name, users.group, users.url_prefix, users.url_name, users.level, replies.ipaddress, replies.likeCount, replies.actionType", "replies.createdBy = users.uid", "replies.tid = ?", "replies.rid ASC", "?,?"),
		}
		return acc.FirstError()
	})
}

var successJSONBytes = []byte(`{"success":"1"}`)

func ViewTopic(w http.ResponseWriter, r *http.Request, user common.User, urlBit string) common.RouteError {
	page, _ := strconv.Atoi(r.FormValue("page"))

	// SEO URLs...
	// TODO: Make a shared function for this
	halves := strings.Split(urlBit, ".")
	if len(halves) < 2 {
		halves = append(halves, halves[0])
	}

	tid, err := strconv.Atoi(halves[1])
	if err != nil {
		return common.PreError("The provided TopicID is not a valid number.", w, r)
	}

	// Get the topic...
	topic, err := common.GetTopicUser(tid)
	if err == sql.ErrNoRows {
		return common.NotFound(w, r, nil) // TODO: Can we add a simplified invocation of headerVars here? This is likely to be an extremely common NotFound
	} else if err != nil {
		return common.InternalError(err, w, r)
	}
	topic.ClassName = ""
	//log.Printf("topic: %+v\n", topic)

	headerVars, ferr := common.ForumUserCheck(w, r, &user, topic.ParentID)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic {
		//log.Printf("user.Perms: %+v\n", user.Perms)
		return common.NoPermissions(w, r, user)
	}
	headerVars.Zone = "view_topic"
	// TODO: Only include these on pages with polls
	headerVars.Stylesheets = append(headerVars.Stylesheets, "chartist/chartist.min.css")
	headerVars.Scripts = append(headerVars.Scripts, "chartist/chartist.min.js")

	topic.ContentHTML = common.ParseMessage(topic.Content, topic.ParentID, "forums")
	topic.ContentLines = strings.Count(topic.Content, "\n")

	// We don't want users posting in locked topics...
	if topic.IsClosed && !user.IsMod {
		user.Perms.CreateReply = false
	}

	postGroup, err := common.Groups.Get(topic.Group)
	if err != nil {
		return common.InternalError(err, w, r)
	}

	topic.Tag = postGroup.Tag
	if postGroup.IsMod || postGroup.IsAdmin {
		topic.ClassName = common.Config.StaffCSS
	}
	topic.RelativeCreatedAt = common.RelativeTime(topic.CreatedAt)

	// TODO: Make a function for this? Build a more sophisticated noavatar handling system?
	topic.Avatar = common.BuildAvatar(topic.CreatedBy, topic.Avatar)

	var poll common.Poll
	if topic.Poll != 0 {
		pPoll, err := common.Polls.Get(topic.Poll)
		if err != nil {
			log.Print("Couldn't find the attached poll for topic " + strconv.Itoa(topic.ID))
			return common.InternalError(err, w, r)
		}
		poll = pPoll.Copy()
	}

	// Calculate the offset
	offset, page, lastPage := common.PageOffset(topic.PostCount, page, common.Config.ItemsPerPage)
	tpage := common.TopicPage{topic.Title, user, headerVars, []common.ReplyUser{}, topic, poll, page, lastPage}

	// Get the replies if we have any...
	if topic.PostCount > 0 {
		rows, err := topicStmts.getReplies.Query(topic.ID, offset, common.Config.ItemsPerPage)
		if err == sql.ErrNoRows {
			return common.LocalError("Bad Page. Some of the posts may have been deleted or you got here by directly typing in the page number.", w, r, user)
		} else if err != nil {
			return common.InternalError(err, w, r)
		}
		defer rows.Close()

		replyItem := common.ReplyUser{ClassName: ""}
		for rows.Next() {
			err := rows.Scan(&replyItem.ID, &replyItem.Content, &replyItem.CreatedBy, &replyItem.CreatedAt, &replyItem.LastEdit, &replyItem.LastEditBy, &replyItem.Avatar, &replyItem.CreatedByName, &replyItem.Group, &replyItem.URLPrefix, &replyItem.URLName, &replyItem.Level, &replyItem.IPAddress, &replyItem.LikeCount, &replyItem.ActionType)
			if err != nil {
				return common.InternalError(err, w, r)
			}

			replyItem.UserLink = common.BuildProfileURL(common.NameToSlug(replyItem.CreatedByName), replyItem.CreatedBy)
			replyItem.ParentID = topic.ID
			replyItem.ContentHtml = common.ParseMessage(replyItem.Content, topic.ParentID, "forums")
			replyItem.ContentLines = strings.Count(replyItem.Content, "\n")

			postGroup, err = common.Groups.Get(replyItem.Group)
			if err != nil {
				return common.InternalError(err, w, r)
			}

			if postGroup.IsMod || postGroup.IsAdmin {
				replyItem.ClassName = common.Config.StaffCSS
			} else {
				replyItem.ClassName = ""
			}

			// TODO: Make a function for this? Build a more sophisticated noavatar handling system? Do bulk user loads and let the common.UserStore initialise this?
			replyItem.Avatar = common.BuildAvatar(replyItem.CreatedBy, replyItem.Avatar)
			replyItem.Tag = postGroup.Tag
			replyItem.RelativeCreatedAt = common.RelativeTime(replyItem.CreatedAt)

			// We really shouldn't have inline HTML, we should do something about this...
			if replyItem.ActionType != "" {
				switch replyItem.ActionType {
				case "lock":
					replyItem.ActionType = "This topic has been locked by <a href='" + replyItem.UserLink + "'>" + replyItem.CreatedByName + "</a>"
					replyItem.ActionIcon = "&#x1F512;&#xFE0E"
				case "unlock":
					replyItem.ActionType = "This topic has been reopened by <a href='" + replyItem.UserLink + "'>" + replyItem.CreatedByName + "</a>"
					replyItem.ActionIcon = "&#x1F513;&#xFE0E"
				case "stick":
					replyItem.ActionType = "This topic has been pinned by <a href='" + replyItem.UserLink + "'>" + replyItem.CreatedByName + "</a>"
					replyItem.ActionIcon = "&#x1F4CC;&#xFE0E"
				case "unstick":
					replyItem.ActionType = "This topic has been unpinned by <a href='" + replyItem.UserLink + "'>" + replyItem.CreatedByName + "</a>"
					replyItem.ActionIcon = "&#x1F4CC;&#xFE0E"
				case "move":
					replyItem.ActionType = "This topic has been moved by <a href='" + replyItem.UserLink + "'>" + replyItem.CreatedByName + "</a>"
				default:
					replyItem.ActionType = replyItem.ActionType + " has happened"
					replyItem.ActionIcon = ""
				}
			}
			replyItem.Liked = false

			common.RunVhook("topic_reply_row_assign", &tpage, &replyItem)
			// TODO: Use a pointer instead to make it easier to abstract this loop? What impact would this have on escape analysis?
			tpage.ItemList = append(tpage.ItemList, replyItem)
		}
		err = rows.Err()
		if err != nil {
			return common.InternalError(err, w, r)
		}
	}

	if common.RunPreRenderHook("pre_render_view_topic", w, r, &user, &tpage) {
		return nil
	}
	err = common.RunThemeTemplate(headerVars.Theme.Name, "topic", tpage, w)
	if err != nil {
		return common.InternalError(err, w, r)
	}
	counters.TopicViewCounter.Bump(topic.ID) // TODO: Move this into the router?
	return nil
}

// ? - Should we add a new permission or permission zone (like per-forum permissions) specifically for profile comment creation
// ? - Should we allow banned users to make reports? How should we handle report abuse?
// TODO: Add a permission to stop certain users from using custom avatars
// ? - Log username changes and put restrictions on this?
func CreateTopic(w http.ResponseWriter, r *http.Request, user common.User, sfid string) common.RouteError {
	var fid int
	var err error
	if sfid != "" {
		fid, err = strconv.Atoi(sfid)
		if err != nil {
			return common.LocalError("You didn't provide a valid number for the forum ID.", w, r, user)
		}
	}
	if fid == 0 {
		fid = common.Config.DefaultForum
	}

	headerVars, ferr := common.ForumUserCheck(w, r, &user, fid)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.CreateTopic {
		return common.NoPermissions(w, r, user)
	}
	headerVars.Zone = "create_topic"

	// Lock this to the forum being linked?
	// Should we always put it in strictmode when it's linked from another forum? Well, the user might end up changing their mind on what forum they want to post in and it would be a hassle, if they had to switch pages, even if it is a single click for many (exc. mobile)
	var strictmode bool
	common.RunVhook("topic_create_pre_loop", w, r, fid, &headerVars, &user, &strictmode)

	// TODO: Re-add support for plugin_guilds
	var forumList []common.Forum
	var canSee []int
	if user.IsSuperAdmin {
		canSee, err = common.Forums.GetAllVisibleIDs()
		if err != nil {
			return common.InternalError(err, w, r)
		}
	} else {
		group, err := common.Groups.Get(user.Group)
		if err != nil {
			// TODO: Refactor this
			common.LocalError("Something weird happened behind the scenes", w, r, user)
			log.Printf("Group #%d doesn't exist, but it's set on common.User #%d", user.Group, user.ID)
			return nil
		}
		canSee = group.CanSee
	}

	// TODO: plugin_superadmin needs to be able to override this loop. Skip flag on topic_create_pre_loop?
	for _, ffid := range canSee {
		// TODO: Surely, there's a better way of doing this. I've added it in for now to support plugin_guilds, but we really need to clean this up
		if strictmode && ffid != fid {
			continue
		}

		// Do a bulk forum fetch, just in case it's the SqlForumStore?
		forum := common.Forums.DirtyGet(ffid)
		if forum.Name != "" && forum.Active {
			fcopy := forum.Copy()
			// TODO: Abstract this
			if common.Hooks["topic_create_frow_assign"] != nil {
				// TODO: Add the skip feature to all the other row based hooks?
				if common.RunHook("topic_create_frow_assign", &fcopy).(bool) {
					continue
				}
			}
			forumList = append(forumList, fcopy)
		}
	}

	ctpage := common.CreateTopicPage{"Create Topic", user, headerVars, forumList, fid}
	if common.RunPreRenderHook("pre_render_create_topic", w, r, &user, &ctpage) {
		return nil
	}

	err = common.RunThemeTemplate(headerVars.Theme.Name, "create_topic", ctpage, w)
	if err != nil {
		return common.InternalError(err, w, r)
	}
	return nil
}

func CreateTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User) common.RouteError {
	fid, err := strconv.Atoi(r.PostFormValue("topic-board"))
	if err != nil {
		return common.LocalError("The provided ForumID is not a valid number.", w, r, user)
	}

	// TODO: Add hooks to make use of headerLite
	_, ferr := common.SimpleForumUserCheck(w, r, &user, fid)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.CreateTopic {
		return common.NoPermissions(w, r, user)
	}

	topicName := html.EscapeString(strings.Replace(r.PostFormValue("topic-name"), "\n", "", -1))
	content := common.PreparseMessage(r.PostFormValue("topic-content"))
	// TODO: Fully parse the post and store it in the parsed column
	tid, err := common.Topics.Create(fid, topicName, content, user.ID, user.LastIP)
	if err != nil {
		switch err {
		case common.ErrNoRows:
			return common.LocalError("Something went wrong, perhaps the forum got deleted?", w, r, user)
		case common.ErrNoTitle:
			return common.LocalError("This topic doesn't have a title", w, r, user)
		case common.ErrNoBody:
			return common.LocalError("This topic doesn't have a body", w, r, user)
		}
		return common.InternalError(err, w, r)
	}

	topic, err := common.Topics.Get(tid)
	if err != nil {
		return common.LocalError("Unable to load the topic", w, r, user)
	}
	if r.PostFormValue("has_poll") == "1" {
		var maxPollOptions = 10
		var pollInputItems = make(map[int]string)
		for key, values := range r.Form {
			//if common.Dev.SuperDebug {
			log.Print("key: ", key)
			log.Printf("values: %+v\n", values)
			//}
			for _, value := range values {
				if strings.HasPrefix(key, "pollinputitem[") {
					halves := strings.Split(key, "[")
					if len(halves) != 2 {
						return common.LocalError("Malformed pollinputitem", w, r, user)
					}
					halves[1] = strings.TrimSuffix(halves[1], "]")

					index, err := strconv.Atoi(halves[1])
					if err != nil {
						return common.LocalError("Malformed pollinputitem", w, r, user)
					}

					// If there are duplicates, then something has gone horribly wrong, so let's ignore them, this'll likely happen during an attack
					_, exists := pollInputItems[index]
					if !exists && len(html.EscapeString(value)) != 0 {
						pollInputItems[index] = html.EscapeString(value)

						if len(pollInputItems) >= maxPollOptions {
							break
						}
					}
				}
			}
		}

		// Make sure the indices are sequential to avoid out of bounds issues
		var seqPollInputItems = make(map[int]string)
		for i := 0; i < len(pollInputItems); i++ {
			seqPollInputItems[i] = pollInputItems[i]
		}

		pollType := 0 // Basic single choice
		_, err := common.Polls.Create(topic, pollType, seqPollInputItems)
		if err != nil {
			return common.LocalError("Failed to add poll to topic", w, r, user) // TODO: Might need to be an internal error as it could leave phantom polls?
		}
	}

	err = common.Subscriptions.Add(user.ID, tid, "topic")
	if err != nil {
		return common.InternalError(err, w, r)
	}

	err = user.IncreasePostStats(common.WordCount(content), true)
	if err != nil {
		return common.InternalError(err, w, r)
	}

	// Handle the file attachments
	// TODO: Stop duplicating this code
	if user.Perms.UploadFiles {
		files, ok := r.MultipartForm.File["upload_files"]
		if ok {
			if len(files) > 5 {
				return common.LocalError("You can't attach more than five files", w, r, user)
			}

			for _, file := range files {
				common.DebugLog("file.Filename ", file.Filename)
				extarr := strings.Split(file.Filename, ".")
				if len(extarr) < 2 {
					return common.LocalError("Bad file", w, r, user)
				}
				ext := extarr[len(extarr)-1]

				// TODO: Can we do this without a regex?
				reg, err := regexp.Compile("[^A-Za-z0-9]+")
				if err != nil {
					return common.LocalError("Bad file extension", w, r, user)
				}
				ext = strings.ToLower(reg.ReplaceAllString(ext, ""))
				if !common.AllowedFileExts.Contains(ext) {
					return common.LocalError("You're not allowed to upload files with this extension", w, r, user)
				}

				infile, err := file.Open()
				if err != nil {
					return common.LocalError("Upload failed", w, r, user)
				}
				defer infile.Close()

				hasher := sha256.New()
				_, err = io.Copy(hasher, infile)
				if err != nil {
					return common.LocalError("Upload failed [Hashing Failed]", w, r, user)
				}
				infile.Close()

				checksum := hex.EncodeToString(hasher.Sum(nil))
				filename := checksum + "." + ext
				outfile, err := os.Create("." + "/attachs/" + filename)
				if err != nil {
					return common.LocalError("Upload failed [File Creation Failed]", w, r, user)
				}
				defer outfile.Close()

				infile, err = file.Open()
				if err != nil {
					return common.LocalError("Upload failed", w, r, user)
				}
				defer infile.Close()

				_, err = io.Copy(outfile, infile)
				if err != nil {
					return common.LocalError("Upload failed [Copy Failed]", w, r, user)
				}

				err = common.Attachments.Add(fid, "forums", tid, "topics", user.ID, filename)
				if err != nil {
					return common.InternalError(err, w, r)
				}
			}
		}
	}

	counters.PostCounter.Bump()
	counters.TopicCounter.Bump()
	http.Redirect(w, r, "/topic/"+strconv.Itoa(tid), http.StatusSeeOther)
	return nil
}

// TODO: Update the stats after edits so that we don't under or over decrement stats during deletes
// TODO: Disable stat updates in posts handled by plugin_guilds
func EditTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User, stid string) common.RouteError {
	isJs := (r.PostFormValue("js") == "1")

	tid, err := strconv.Atoi(stid)
	if err != nil {
		return common.PreErrorJSQ("The provided TopicID is not a valid number.", w, r, isJs)
	}

	topic, err := common.Topics.Get(tid)
	if err == sql.ErrNoRows {
		return common.PreErrorJSQ("The topic you tried to edit doesn't exist.", w, r, isJs)
	} else if err != nil {
		return common.InternalErrorJSQ(err, w, r, isJs)
	}

	// TODO: Add hooks to make use of headerLite
	_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.EditTopic {
		return common.NoPermissionsJSQ(w, r, user, isJs)
	}

	err = topic.Update(r.PostFormValue("topic_name"), r.PostFormValue("topic_content"))
	if err != nil {
		return common.InternalErrorJSQ(err, w, r, isJs)
	}

	err = common.Forums.UpdateLastTopic(topic.ID, user.ID, topic.ParentID)
	if err != nil && err != sql.ErrNoRows {
		return common.InternalErrorJSQ(err, w, r, isJs)
	}

	if !isJs {
		http.Redirect(w, r, "/topic/"+strconv.Itoa(tid), http.StatusSeeOther)
	} else {
		_, _ = w.Write(successJSONBytes)
	}
	return nil
}

// TODO: Add support for soft-deletion and add a permission for hard delete in addition to the usual
// TODO: Disable stat updates in posts handled by plugin_guilds
func DeleteTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User) common.RouteError {
	// TODO: Move this to some sort of middleware
	var tids []int
	var isJs = false
	if common.ReqIsJson(r) {
		if r.Body == nil {
			return common.PreErrorJS("No request body", w, r)
		}
		err := json.NewDecoder(r.Body).Decode(&tids)
		if err != nil {
			return common.PreErrorJS("We weren't able to parse your data", w, r)
		}
		isJs = true
	} else {
		tid, err := strconv.Atoi(r.URL.Path[len("/topic/delete/submit/"):])
		if err != nil {
			return common.PreError("The provided TopicID is not a valid number.", w, r)
		}
		tids = append(tids, tid)
	}
	if len(tids) == 0 {
		return common.LocalErrorJSQ("You haven't provided any IDs", w, r, user, isJs)
	}

	for _, tid := range tids {
		topic, err := common.Topics.Get(tid)
		if err == sql.ErrNoRows {
			return common.PreErrorJSQ("The topic you tried to delete doesn't exist.", w, r, isJs)
		} else if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}

		// TODO: Add hooks to make use of headerLite
		_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
		if ferr != nil {
			return ferr
		}
		if !user.Perms.ViewTopic || !user.Perms.DeleteTopic {
			return common.NoPermissionsJSQ(w, r, user, isJs)
		}

		// We might be able to handle this err better
		err = topic.Delete()
		if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}

		err = common.ModLogs.Create("delete", tid, "topic", user.LastIP, user.ID)
		if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}

		// ? - We might need to add soft-delete before we can do an action reply for this
		/*_, err = stmts.createActionReply.Exec(tid,"delete",ipaddress,user.ID)
		if err != nil {
			return common.InternalErrorJSQ(err,w,r,isJs)
		}*/

		log.Printf("Topic #%d was deleted by common.User #%d", tid, user.ID)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func StickTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User, stid string) common.RouteError {
	tid, err := strconv.Atoi(stid)
	if err != nil {
		return common.PreError("The provided TopicID is not a valid number.", w, r)
	}

	topic, err := common.Topics.Get(tid)
	if err == sql.ErrNoRows {
		return common.PreError("The topic you tried to pin doesn't exist.", w, r)
	} else if err != nil {
		return common.InternalError(err, w, r)
	}

	// TODO: Add hooks to make use of headerLite
	_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.PinTopic {
		return common.NoPermissions(w, r, user)
	}

	err = topic.Stick()
	if err != nil {
		return common.InternalError(err, w, r)
	}

	err = addTopicAction("stick", topic, user)
	if err != nil {
		return common.InternalError(err, w, r)
	}
	http.Redirect(w, r, "/topic/"+strconv.Itoa(tid), http.StatusSeeOther)
	return nil
}

func UnstickTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User, stid string) common.RouteError {
	tid, err := strconv.Atoi(stid)
	if err != nil {
		return common.PreError("The provided TopicID is not a valid number.", w, r)
	}

	topic, err := common.Topics.Get(tid)
	if err == sql.ErrNoRows {
		return common.PreError("The topic you tried to unpin doesn't exist.", w, r)
	} else if err != nil {
		return common.InternalError(err, w, r)
	}

	// TODO: Add hooks to make use of headerLite
	_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.PinTopic {
		return common.NoPermissions(w, r, user)
	}

	err = topic.Unstick()
	if err != nil {
		return common.InternalError(err, w, r)
	}

	err = addTopicAction("unstick", topic, user)
	if err != nil {
		return common.InternalError(err, w, r)
	}

	http.Redirect(w, r, "/topic/"+strconv.Itoa(tid), http.StatusSeeOther)
	return nil
}

func LockTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User) common.RouteError {
	// TODO: Move this to some sort of middleware
	var tids []int
	var isJs = false
	if common.ReqIsJson(r) {
		if r.Body == nil {
			return common.PreErrorJS("No request body", w, r)
		}
		err := json.NewDecoder(r.Body).Decode(&tids)
		if err != nil {
			return common.PreErrorJS("We weren't able to parse your data", w, r)
		}
		isJs = true
	} else {
		tid, err := strconv.Atoi(r.URL.Path[len("/topic/lock/submit/"):])
		if err != nil {
			return common.PreError("The provided TopicID is not a valid number.", w, r)
		}
		tids = append(tids, tid)
	}
	if len(tids) == 0 {
		return common.LocalErrorJSQ("You haven't provided any IDs", w, r, user, isJs)
	}

	for _, tid := range tids {
		topic, err := common.Topics.Get(tid)
		if err == sql.ErrNoRows {
			return common.PreErrorJSQ("The topic you tried to lock doesn't exist.", w, r, isJs)
		} else if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}

		// TODO: Add hooks to make use of headerLite
		_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
		if ferr != nil {
			return ferr
		}
		if !user.Perms.ViewTopic || !user.Perms.CloseTopic {
			return common.NoPermissionsJSQ(w, r, user, isJs)
		}

		err = topic.Lock()
		if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}

		err = addTopicAction("lock", topic, user)
		if err != nil {
			return common.InternalErrorJSQ(err, w, r, isJs)
		}
	}

	if len(tids) == 1 {
		http.Redirect(w, r, "/topic/"+strconv.Itoa(tids[0]), http.StatusSeeOther)
	}
	return nil
}

func UnlockTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User, stid string) common.RouteError {
	tid, err := strconv.Atoi(stid)
	if err != nil {
		return common.PreError("The provided TopicID is not a valid number.", w, r)
	}

	topic, err := common.Topics.Get(tid)
	if err == sql.ErrNoRows {
		return common.PreError("The topic you tried to unlock doesn't exist.", w, r)
	} else if err != nil {
		return common.InternalError(err, w, r)
	}

	// TODO: Add hooks to make use of headerLite
	_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
	if ferr != nil {
		return ferr
	}
	if !user.Perms.ViewTopic || !user.Perms.CloseTopic {
		return common.NoPermissions(w, r, user)
	}

	err = topic.Unlock()
	if err != nil {
		return common.InternalError(err, w, r)
	}

	err = addTopicAction("unlock", topic, user)
	if err != nil {
		return common.InternalError(err, w, r)
	}

	http.Redirect(w, r, "/topic/"+strconv.Itoa(tid), http.StatusSeeOther)
	return nil
}

// ! JS only route
// TODO: Figure a way to get this route to work without JS
func MoveTopicSubmit(w http.ResponseWriter, r *http.Request, user common.User, sfid string) common.RouteError {
	fid, err := strconv.Atoi(sfid)
	if err != nil {
		return common.PreErrorJS("The provided Forum ID is not a valid number.", w, r)
	}

	// TODO: Move this to some sort of middleware
	var tids []int
	if r.Body == nil {
		return common.PreErrorJS("No request body", w, r)
	}
	err = json.NewDecoder(r.Body).Decode(&tids)
	if err != nil {
		return common.PreErrorJS("We weren't able to parse your data", w, r)
	}
	if len(tids) == 0 {
		return common.LocalErrorJS("You haven't provided any IDs", w, r)
	}

	for _, tid := range tids {
		topic, err := common.Topics.Get(tid)
		if err == sql.ErrNoRows {
			return common.PreErrorJS("The topic you tried to move doesn't exist.", w, r)
		} else if err != nil {
			return common.InternalErrorJS(err, w, r)
		}

		// TODO: Add hooks to make use of headerLite
		_, ferr := common.SimpleForumUserCheck(w, r, &user, topic.ParentID)
		if ferr != nil {
			return ferr
		}
		if !user.Perms.ViewTopic || !user.Perms.MoveTopic {
			return common.NoPermissionsJS(w, r, user)
		}
		_, ferr = common.SimpleForumUserCheck(w, r, &user, fid)
		if ferr != nil {
			return ferr
		}
		if !user.Perms.ViewTopic || !user.Perms.MoveTopic {
			return common.NoPermissionsJS(w, r, user)
		}

		err = topic.MoveTo(fid)
		if err != nil {
			return common.InternalErrorJS(err, w, r)
		}

		// TODO: Log more data so we can list the destination forum in the action post?
		err = addTopicAction("move", topic, user)
		if err != nil {
			return common.InternalErrorJS(err, w, r)
		}
	}

	if len(tids) == 1 {
		http.Redirect(w, r, "/topic/"+strconv.Itoa(tids[0]), http.StatusSeeOther)
	}
	return nil
}

func addTopicAction(action string, topic *common.Topic, user common.User) error {
	err := common.ModLogs.Create(action, topic.ID, "topic", user.LastIP, user.ID)
	if err != nil {
		return err
	}
	return topic.CreateActionReply(action, user.LastIP, user)
}
