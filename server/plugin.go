package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	ONELOGIN_ICON_URL = "https://www.onelogin.com/assets/img/press/presskit/downloads/Onelogin_Logomark/Screen/png/Onelogin_Mark_black_RGB.png"
	ONELOGIN_USERNAME = "OneLogin Bot"
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	TeamID    string
	ChannelID string
	BotUserID string
}

func (p *Plugin) OnActivate() error {
	configuration := p.getConfiguration()

	if err := p.IsValid(configuration); err != nil {
		return err
	}

	split := strings.Split(p.configuration.TeamChannel, ",")
	teamSplit := split[0]
	channelSplit := split[1]

	team, err := p.API.GetTeamByName(teamSplit)
	if err != nil {
		return err
	}
	p.TeamID = team.Id

	user, err := p.API.GetUserByUsername(p.configuration.UserName)
	if err != nil {
		p.API.LogError(err.Error())
		return fmt.Errorf("Unable to find user with configured username: %v", p.configuration.UserName)
	}
	p.BotUserID = user.Id

	channel, err := p.API.GetChannelByName(team.Id, channelSplit, false)
	if err != nil && err.StatusCode == http.StatusNotFound {
		channelToCreate := &model.Channel{
			Name:        channelSplit,
			DisplayName: channelSplit,
			Type:        model.CHANNEL_OPEN,
			TeamId:      p.TeamID,
			CreatorId:   p.BotUserID,
		}

		newChannel, errChannel := p.API.CreateChannel(channelToCreate)
		if err != nil {
			return errChannel
		}
		p.ChannelID = newChannel.Id
	} else if err != nil {
		return err
	} else {
		p.ChannelID = channel.Id
	}

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if err := p.checkHeaderToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		p.API.LogError("Onelogin TOKEN INVALID")
		return
	}

	var oneLoginEvents []OneLogin
	if err := json.NewDecoder(r.Body).Decode(&oneLoginEvents); err != nil {
		return
	}

	for _, event := range oneLoginEvents {
		risk, _ := strconv.Atoi(p.configuration.RiskThreshold)
		if event.EventTypeID == 5 && event.RiskScore > risk { // USER_LOGGED_INTO_ONELOGIN
			p.handleLoginPossibleThreat(event)
		} else if event.EventTypeID == 12 { // UNLOCKED_USER
			p.handleLoginUserUnlocked(event)
		} else if event.EventTypeID == 13 { // CREATED_USER
			p.handleLoginUserCreation(event)
		} else if event.EventTypeID == 15 { // DEACTIVATED_USER
			p.handleLoginUserDeactived(event)
		} else if event.EventTypeID == 17 { // DELETED_USER
			p.handleLoginUserDeleted(event)
		} else if event.EventTypeID == 19 { // USER_LOCKED
			p.handleLoginUserLocked(event)
		} else if event.EventTypeID == 24 { // USER_REMOVED_OTP_DEVICE
			p.handleLoginUserRemovedOTP(event)
		} else if event.EventTypeID == 69 { // USER_REJECTED_BY_RADIUS
			p.handleRejectedRadius(event)
		} else {
			p.API.LogInfo("Not implemented yet", "Event Type", event.EventTypeID)
		}
	}

	return

}

func (p *Plugin) checkHeaderToken(r *http.Request) error {
	headerToken := r.Header.Get("X-OneLogin-Token")
	if headerToken == "" || strings.Compare(headerToken, p.configuration.Token) != 0 {
		return fmt.Errorf("Invalid or missing token")
	}
	return nil
}

func (p *Plugin) handleLoginPossibleThreat(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "Risk Score", strconv.Itoa(event.RiskScore), false)
	fields = addFields(fields, "Risk Reasons", event.RiskReasons, false)
	fields = addFields(fields, "Notes", event.Notes, false)
	fields = addFields(fields, "User Agent", event.UserAgent, true)
	fields = addFields(fields, "IP Address", event.Ipaddr, true)

	title := fmt.Sprintf("%s had a risky login to OneLogin.", event.UserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#ff0000",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserCreation(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "User Name", event.UserName, false)

	title := fmt.Sprintf("%s was created by %s", event.UserName, event.ActorUserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#008000",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserDeactived(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "User Name", event.UserName, false)
	fields = addFields(fields, "Login Name", event.LoginName, false)
	fields = addFields(fields, "Notes", event.Notes, false)

	title := fmt.Sprintf("%s was deactived by %s", event.UserName, event.ActorUserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#0000FF",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserDeleted(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "User Name", event.UserName, false)

	title := fmt.Sprintf("%s was deleted by %s", event.UserName, event.ActorUserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#0000FF",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserUnlocked(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "Login Name", event.LoginName, true)
	fields = addFields(fields, "IP Address", event.Ipaddr, true)

	title := fmt.Sprintf("%s unlocked %s", event.ActorUserName, event.UserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#008000",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserLocked(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "Login Name", event.LoginName, true)
	fields = addFields(fields, "IP Address", event.Ipaddr, true)
	fields = addFields(fields, "Notes", event.Notes, false)

	title := fmt.Sprintf("%s locked", event.UserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#FFA500",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleLoginUserRemovedOTP(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "User Name", event.UserName, false)
	fields = addFields(fields, "Notes", event.Notes, false)
	fields = addFields(fields, "OTP Device Name", event.OtpDeviceName, true)
	fields = addFields(fields, "OTP Device ID", strconv.Itoa(event.OtpDeviceID), true)

	title := fmt.Sprintf("%s deregistered for %s", event.OtpDeviceName, event.UserName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#FFA500",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func (p *Plugin) handleRejectedRadius(event OneLogin) {
	var fields []*model.SlackAttachmentField
	fields = addFields(fields, "Login Name", event.LoginName, true)
	fields = addFields(fields, "IP Address", event.Ipaddr, true)

	title := fmt.Sprintf("%s rejected by %s", event.ActorUserName, event.RadiusConfigName)
	attachment := &model.SlackAttachment{
		Title:  title,
		Fields: fields,
		Color:  "#008000",
	}

	post := &model.Post{
		ChannelId: p.ChannelID,
		UserId:    p.BotUserID,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": ONELOGIN_USERNAME,
			"override_icon_url": ONELOGIN_ICON_URL,
		},
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		return
	}
	return
}

func addFields(fields []*model.SlackAttachmentField, title, msg string, short bool) []*model.SlackAttachmentField {
	return append(fields, &model.SlackAttachmentField{
		Title: title,
		Value: msg,
		Short: short,
	})
}
