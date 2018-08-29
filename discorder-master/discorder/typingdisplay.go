package discorder

import (
	"github.com/jonas747/discorder/common"
	"github.com/jonas747/discorder/ui"
	"github.com/jonas747/discordgo"
)

// Shows who's typing
type TypingDisplay struct {
	*ui.BaseEntity

	App *App

	text *ui.Text
}

func NewTypingDisplay(app *App) *TypingDisplay {
	td := &TypingDisplay{
		BaseEntity: &ui.BaseEntity{},
		App:        app,
		text:       ui.NewText(),
	}

	td.text.Transform.AnchorMax = common.NewVector2I(1, 1)
	app.ApplyThemeToText(td.text, "typing_bar")
	td.Transform.AddChildren(td.text)

	return td
}

func (t *TypingDisplay) Update() {
	if t.App.session == nil || !t.App.session.DataReady {
		return
	}
	tab := t.App.ViewManager.ActiveTab

	if tab == nil {
		return
	}

	channels := make([]string, len(tab.MessageView.Channels))
	copy(channels, tab.MessageView.Channels)

	if tab.MessageView.ShowAllPrivate {
		t.App.session.State.RLock()
		for _, pChan := range t.App.session.State.PrivateChannels {
			found := false
			for _, subChan := range tab.MessageView.Channels {
				if subChan == pChan.ID {
					found = true
					break
				}
			}
			if !found {
				channels = append(channels, pChan.ID)
			}
		}
		t.App.session.State.RUnlock()
	}

	typing := t.App.typingRoutine.GetTyping(channels)

	if len(typing) > 0 {
		t.text.Disabled = false

		typingStr := "Typing: "

		for _, v := range typing {
			channel, err := t.App.session.State.Channel(v.ChannelID)
			if err != nil {
				continue
			}
			if !(channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM) {
				member, err := t.App.session.State.Member(channel.GuildID, v.UserID)
				if err != nil {
					continue
				}
				typingStr += "#" + GetChannelNameOrRecipient(channel) + ":" + member.User.Username + ", "
			} else {
				typingStr += "#DM:" + channel.Recipients[0].Username + ", "
			}
		}
		// Remove trailing ","
		typingStr = typingStr[:len(typingStr)-2]
		t.text.Text = typingStr
	} else {
		t.text.Disabled = true
	}
}

func (t *TypingDisplay) Destroy() { t.DestroyChildren() }

func (t *TypingDisplay) GetRequiredSize() common.Vector2F {
	rect := t.text.Transform.GetRect()
	//log.Println(float32(t.text.HeightRequired()), t.text.Text)
	return common.Vector2F{rect.W, float32(t.text.HeightRequired())}
}

func (t *TypingDisplay) IsLayoutDynamic() bool {
	return false
}
