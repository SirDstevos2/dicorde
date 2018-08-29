package discorder

import (
	"github.com/jonas747/discorder/common"
	"github.com/jonas747/discorder/ui"
	"log"
)

func OpenLoginWindow(app *App) *CommandWindow {
	cmdWindow := NewCommandWindow(app, 5, nil, "Because of Discord API restrictions Discorder does not support 2 Factor Authentication. Discord has made it very clear that only official clients can use that API endpoint.\n\nTo grab a token from the official client, open Developer Tools with Ctrl-I (Ctrl-Alt-I on Mac), click the Application tab and navigate to the Local Storage section, and find 'token' in the list. Copy the token without quotes below.\n")
	cmdWindow.menuWindow.SearchInput.SetActive(false)
	cmdWindow.commands = LoginCommands
	cmdWindow.categories = []*CommandCategory{}

	return cmdWindow
}

var LoginCommands = []Command{
	&SimpleCommand{
		Name:           "Token login",
		Description:    "Login using a token that was peviously aquired from somewhere else (like the official client)",
		CustomExecText: "Login",
		Args: []*ArgumentDef{
			{Name: "token", Description: "The token", Datatype: ui.DataTypeString},
		},
		RunFunc: func(app *App, args Arguments) {
			token, _ := args.String("token")
			if token == "" {
				log.Println("No token provided.")
				return
			} else {
				log.Println("Trying to log in using a token. If this takes too long, either the API is having problems or the token is invalid.")
				err := app.Login("", "", token)
				if err != nil {
					log.Println("Failed logging in with token:", err)
				}
			}
		},
	},
	&SimpleCommand{
		Name:           "Email and password login (no 2fa)",
		Description:    "Login using email and password",
		CustomExecText: "Login",
		Args: []*ArgumentDef{
			{Name: "Email", Description: "Your email", Datatype: ui.DataTypeString},
			{Name: "Password", Description: "Your email", Datatype: ui.DataTypePassword},
		},
		RunFunc: func(app *App, args Arguments) {
			email, _ := args.String("Email")
			pw, _ := args.String("Password")
			if email == "" {
				log.Println("Email empty")
				return
			}

			if pw == "" {
				log.Println("Password empty")
				return
			}
			err := app.Login(email, pw, "")
			log.Println("Trying to log in using email and password")
			if err != nil {
				log.Println("Error logging in", err)
			} else {
				log.Println("Sucessfully logged in!")
			}
		},
	},
}

type WaitingForLogin struct {
	*ui.BaseEntity

	app    *App
	window *ui.Window
}

func NewWaitingForLogin(app *App, layer int) *WaitingForLogin {
	window := ui.NewWindow(app.ViewManager.UIManager)

	loginWindow := &WaitingForLogin{
		BaseEntity: &ui.BaseEntity{},
		window:     window,
		app:        app,
	}

	loginWindow.Transform.AddChildren(window)
	window.Transform.AnchorMax = common.NewVector2I(1, 1)

	text := ui.NewText()

	text.Text = "Logging in...\n\nIf you logged in by token and this is taking a long time, either Discord is having API problems or the token is invalid. Reset the token using the -r switch."
	app.ApplyThemeToText(text, "text_window_normal")
	window.Transform.AddChildren(text)
	text.Transform.AnchorMax = common.NewVector2I(1, 1)

	app.ViewManager.UIManager.AddWindow(loginWindow)

	loginWindow.Transform.AnchorMax = common.NewVector2I(1, 1)
	loginWindow.Transform.Right = 2
	loginWindow.Transform.Left = 1

	return loginWindow
}

func (w *WaitingForLogin) Update() {
	if w.app.firstReady {
		w.app.ViewManager.RemoveWindow(w)
	}
}

func (w *WaitingForLogin) Destroy() {
	w.app.ViewManager.UIManager.RemoveWindow(w)
	w.DestroyChildren()
}
