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
		if event.EventTypeID == 5 && event.RiskScore > risk {
			p.handleLoginPossibleThreat(event)
		} else {
			p.API.LogInfo("Not implemented yet")
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

func addFields(fields []*model.SlackAttachmentField, title, msg string, short bool) []*model.SlackAttachmentField {
	return append(fields, &model.SlackAttachmentField{
		Title: title,
		Value: msg,
		Short: short,
	})
}
