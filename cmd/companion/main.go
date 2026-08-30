package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"rocketleaguecompanion/internal/config"
	"rocketleaguecompanion/internal/state"
	"rocketleaguecompanion/internal/stats"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	settingsPath := filepath.Join(homeDir(), ".rocketleague-companion", "settings.json")
	defaultStatsAPIPath, userStatsAPIPath := rocketLeagueConfigPaths()
	settings, err := config.Load(settingsPath)
	if err != nil {
		fmt.Printf("failed to load settings: %v\n", err)
		settings = config.DefaultSettings()
	}

	tracker := state.NewSessionStats()
	ui := newCompanionUI(settings, settingsPath, defaultStatsAPIPath, userStatsAPIPath, tracker)
	ui.show()
}

func rocketLeagueConfigPaths() (string, string) {
	executable, err := os.Executable()
	if err != nil {
		executable = "."
	}
	defaultPath := filepath.Join(filepath.Dir(executable), "..", "DefaultStatsAPI.ini")
	userPath := filepath.Join(homeDir(), "My Games", "Rocket League", "TAGame", "Config", "TAStatsAPI.ini")
	return filepath.Clean(defaultPath), userPath
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

type companionUI struct {
	settings            config.Settings
	settingsPath        string
	defaultStatsAPIPath string
	userStatsAPIPath    string
	session             *state.SessionStats
	app                 fyne.App
	window              fyne.Window
	status              *widget.Label
	match               *widget.Label
	score               *widget.Label
	trackerWins         *widget.Label
	trackerLoss         *widget.Label
	rateSlider          *widget.Slider
}

func newCompanionUI(settings config.Settings, settingsPath, defaultStatsAPIPath, userStatsAPIPath string, session *state.SessionStats) *companionUI {
	return &companionUI{
		settings:            settings,
		settingsPath:        settingsPath,
		defaultStatsAPIPath: defaultStatsAPIPath,
		userStatsAPIPath:    userStatsAPIPath,
		session:             session,
	}
}

func (u *companionUI) show() {
	u.app = app.New()
	u.window = u.app.NewWindow("Rocket League Companion")
	u.window.Resize(fyne.NewSize(760, 540))

	u.status = widget.NewLabel("Waiting for stats stream...")
	u.match = widget.NewLabel("Arena: --\nClock: --\nPlaylist: --")
	u.score = widget.NewLabel("Score: 0 - 0")
	u.trackerWins = widget.NewLabel("0")
	u.trackerLoss = widget.NewLabel("0")

	u.rateSlider = widget.NewSlider(5, 120)
	u.rateSlider.Step = 5
	u.rateSlider.Value = u.settings.PacketSendRate
	rateLabel := widget.NewLabel(fmt.Sprintf("Current: %.0f Hz", u.settings.PacketSendRate))
	u.rateSlider.OnChanged = func(value float64) {
		u.settings.PacketSendRate = value
		rateLabel.SetText(fmt.Sprintf("Current: %.0f Hz", value))
	}

	settingsForm := container.NewVBox(
		widget.NewLabel("PacketSendRate"),
		rateLabel,
		u.rateSlider,
		widget.NewLabel("Recommended setting: 30 Hz for most systems; increase for smoother tracking on faster hardware."),
		widget.NewLabel("Set this before launching Rocket League; changes require a client restart."),
		widget.NewButton("Save settings", func() {
			if err := config.Save(u.settingsPath, u.settings); err != nil {
				u.status.SetText("Failed to save settings: " + err.Error())
				return
			}
			if err := config.SaveRocketLeague(u.defaultStatsAPIPath, u.settings); err != nil {
				u.status.SetText("Failed to save Rocket League defaults: " + err.Error())
				return
			}
			if err := config.SaveRocketLeague(u.userStatsAPIPath, u.settings); err != nil {
				u.status.SetText("Failed to save Rocket League user config: " + err.Error())
				return
			}
			u.status.SetText("Settings saved. Restart Rocket League for changes to take effect.")
		}),
	)

	winButton := widget.NewButton("+ Win", func() {
		u.session.RecordWin()
		u.refreshSessionTracker()
	})
	lossButton := widget.NewButton("+ Loss", func() {
		u.session.RecordLoss()
		u.refreshSessionTracker()
	})
	resetButton := widget.NewButton("Reset session", func() {
		u.session.Reset()
		u.refreshSessionTracker()
	})

	sessionRow := container.NewHBox(
		widget.NewLabel("Session record:"),
		widget.NewLabel("Wins:"),
		u.trackerWins,
		widget.NewLabel("Losses:"),
		u.trackerLoss,
		winButton,
		lossButton,
		resetButton,
	)

	statsCard := container.NewVBox(
		widget.NewLabel("Live Match"),
		u.score,
		u.match,
	)

	content := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(
			widget.NewLabel("Rocket League Companion"),
			u.status,
			sessionRow,
			statsCard,
			settingsForm,
		),
	)

	u.window.SetContent(content)
	u.refreshSessionTracker()

	client := stats.NewClient(u.settings.Port)
	go func() {
		if err := client.Start(); err != nil {
			u.status.SetText("Connection error: " + err.Error())
			return
		}
		for snapshot := range client.Updates() {
			u.updateMatch(snapshot)
		}
	}()

	u.window.ShowAndRun()
}

func (u *companionUI) refreshSessionTracker() {
	u.trackerWins.SetText(strconv.Itoa(u.session.Wins))
	u.trackerLoss.SetText(strconv.Itoa(u.session.Losses))
}

func (u *companionUI) updateMatch(snapshot stats.Snapshot) {
	u.status.SetText("Connected to match data stream")
	u.match.SetText(fmt.Sprintf("Arena: %s\nClock: %ds\nPlaylist: %s", snapshot.Arena, snapshot.TimeSeconds, snapshot.PlaylistName))
	u.score.SetText(fmt.Sprintf("Score: %d - %d", snapshot.ScoreBlue, snapshot.ScoreOrange))
}
