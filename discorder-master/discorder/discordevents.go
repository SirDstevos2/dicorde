package discorder

import (
	"github.com/0xAX/notificator"
	"github.com/jonas747/discorder/ui"
	"github.com/jonas747/discordgo"
	"log"
)

func (app *App) Ready(s *discordgo.Session, r *discordgo.Ready) {
	app.Lock()
	defer app.Unlock()

	log.Println("Received ready from Discord!")

	app.settings = r.Settings
	app.guildSettings = r.UserGuildSettings

	if app.firstReady {
		return
	}
	app.firstReady = true

	app.requestRoutineRunning = true
	go app.requestRoutine.Run()

	app.ViewManager.OnReady()
	app.PrintWelcome()
}

func (app *App) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	app.Lock()
	defer app.Unlock()

	// Have no idea how this happens but it has so... just gonna leave this here to be sure
	if app.session.State.User == nil {
		return
	}

	var settings *ChannelNotificationSettings
	if !app.session.State.User.Bot {
		settings = app.GetNotificationSettingsForChannel(m.ChannelID)
	} else {
		settings = &ChannelNotificationSettings{
			Notifications:    2,
			Muted:            true,
			SurpressEveryone: true,
		}
	}

	if m.Author == nil {
		// I believe this only happens in mesasge edits
		// but to be sure i just put this check here, will prob get removed in the future
		log.Println("!MESSAGE HAS NO AUTHOR!")
		return
	}

	author := m.Author.Username
	authorId := m.Author.ID

	if app.typingRoutine != nil {
		app.typingRoutine.msgEvtIn <- authorId
	}

	// Check if we should do a notification
	if authorId != app.session.State.User.ID && !s.State.User.Bot {

		shouldNotify := false

		if !settings.Muted && settings.Notifications == MessageNotificationsAll {
			shouldNotify = true
		} else if !settings.Muted && settings.Notifications == MessageNotificationsMentions {
			for _, v := range m.Mentions {
				if v.ID == s.State.User.ID {
					shouldNotify = true
					break
				}
			}
		} else if !settings.SurpressEveryone && m.MentionEveryone {
			shouldNotify = true
		}

		if shouldNotify {
			if app.notifications != nil {
				app.notifications.Push(author, m.ContentWithMentionsReplaced(), "", notificator.UR_NORMAL)
			}
			app.ViewManager.notificationsManager.AddMention(m.Message)
		}
	}

	// Update last message
	channel, err := app.session.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel", err)
	} else {
		channel.LastMessageID = m.ID
	}

	// Emit event
	ui.RunFunc(app, func(e ui.Entity) {
		cast, ok := e.(MessageCreateHandler)
		if ok {
			cast.HandleMessageCreate(m.Message)
		}
	})
}

func (app *App) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// Emit event
	app.Lock()
	defer app.Unlock()

	ui.RunFunc(app, func(e ui.Entity) {
		cast, ok := e.(MessageUpdateHandler)
		if ok {
			cast.HandleMessageUpdate(m.Message)
		}
	})
}

func (app *App) messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	app.Lock()
	defer app.Unlock()

	// Emit event
	ui.RunFunc(app, func(e ui.Entity) {
		cast, ok := e.(MessageRemoveHandler)
		if ok {
			cast.HandleMessageRemove(m.Message)
		}
	})
}

func (app *App) messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	if app.options.DebugEnabled {
		log.Println("Received ack!")
	}
	app.ViewManager.notificationsManager.HandleAck(a)
}

func (app *App) guildSettingsUpdated(s *discordgo.Session, a *discordgo.UserGuildSettingsUpdate) {
	app.Lock()
	defer app.Unlock()

	set := false
	for k, settings := range app.guildSettings {
		if settings.GuildID == a.GuildID {
			app.guildSettings[k] = a.UserGuildSettings
			set = true
			break
		}
	}

	if !set {
		app.guildSettings = append(app.guildSettings, a.UserGuildSettings)
	}
}

func (app *App) userSettingsUpdated(s *discordgo.Session, u *discordgo.UserSettingsUpdate) {
}

func (app *App) typingStart(s *discordgo.Session, t *discordgo.TypingStart) {
	app.typingRoutine.typingEvtIn <- t
}
