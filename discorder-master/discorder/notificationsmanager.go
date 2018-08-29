package discorder

import (
	"fmt"
	"github.com/jonas747/discorder/common"
	"github.com/jonas747/discorder/ui"
	"github.com/jonas747/discordgo"
	"log"
)

type NotificationsManager struct {
	*ui.BaseEntity
	App  *App
	text *ui.Text
}

func NewNotificationsManager(app *App) *NotificationsManager {
	nm := &NotificationsManager{
		BaseEntity: &ui.BaseEntity{},
		App:        app,
		text:       ui.NewText(),
	}

	nm.text.Transform.AnchorMax = common.NewVector2I(1, 1)
	nm.text.Layer = 2
	app.ApplyThemeToText(nm.text, "notifications_bar")
	nm.Transform.AddChildren(nm.text)
	return nm
}

func (nm *NotificationsManager) Update() {
	str := ""

	if nm.App.session == nil {
		return
	}

	nm.App.session.State.RLock()
	defer nm.App.session.State.RUnlock()

	readStates := nm.App.session.State.ReadState

	if len(readStates) > 0 {
		total := 0

		for _, v := range readStates {
			if v.MentionCount == 0 {
				continue
			}

			total += v.MentionCount
			nm.App.session.State.RUnlock() // 2 read locks in same goroutines can cause deadlock
			channel, err := nm.App.session.State.Channel(v.ID)
			nm.App.session.State.RLock()
			if err != nil {
				log.Println("Error getting channel:", err)
				continue
			}
			if (channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM) {
				str += fmt.Sprintf("@%s: %d, ", channel.Recipients[0].Username, v.MentionCount)
				continue
			}

			nm.App.session.State.RUnlock() // 2 read locks in same goroutines can cause deadlock
			guild, err := nm.App.session.State.Guild(channel.GuildID)
			nm.App.session.State.RLock()
			if err != nil {
				log.Println("Error getting guild:", err)
				continue
			}

			str += fmt.Sprintf("%s/%s: %d, ", guild.Name, GetChannelNameOrRecipient(channel), v.MentionCount)
		}
		if str != "" {
			str = str[:len(str)-2]
			str = fmt.Sprintf("Mentions: %d (%s)", total, str)
		}
	}
	nm.text.Text = str
}

func (nm *NotificationsManager) AddMention(msg *discordgo.Message) {
	found := false

	state := nm.App.session.State
	state.Lock()
	defer state.Unlock()

	for _, v := range state.ReadState {
		if v.ID == msg.ChannelID {
			v.MentionCount += 1
			found = true
			break
		}
	}

	if !found {
		state.ReadState = append(state.ReadState, &discordgo.ReadState{
			ID:            msg.ChannelID,
			MentionCount:  1,
			LastMessageID: "",
		})
	}
}

func (nm *NotificationsManager) HandleAck(a *discordgo.MessageAck) {
	// du di da

	state := nm.App.session.State
	state.Lock()
	defer state.Unlock()

	var rs *discordgo.ReadState
	for _, v := range state.ReadState {
		if v.ID == a.ChannelID {
			rs = v
			break
		}
	}

	if rs == nil {
		rs = &discordgo.ReadState{
			ID:            a.ChannelID,
			LastMessageID: a.MessageID,
		}
		state.ReadState = append(state.ReadState, rs)
	}
	state.Unlock()
	channel, err := state.Channel(rs.ID)
	state.Lock()
	if err != nil {
		log.Println("Failed getting channel in HandleAck.")
		return
	}

	if channel.LastMessageID == a.MessageID {
		rs.MentionCount = 0
	}
}

func (nm *NotificationsManager) Destroy() { nm.DestroyChildren() }

func (nm *NotificationsManager) GetRequiredSize() common.Vector2F {
	return common.NewVector2F(0, float32(nm.text.HeightRequired()))
}

func (nm *NotificationsManager) IsLayoutDynamic() bool {
	return false
}

type NotificationSource struct {
	ChannelId string
	LastRead  string
	Count     int
}
