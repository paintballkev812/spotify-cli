package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charlesyu108/spotify-cli/spotify"
	"github.com/urfave/cli/v2"
)

// ConfigFile defines which config JSON file to load
const ConfigFile = "config.json"

func main() {
	app := &cli.App{
		Name:                 "spotify-cli",
		Usage:                "Use Spotify from the Command Line.",
		EnableBashCompletion: true,
	}

	app.Commands = []*cli.Command{
		// Define Playback category commands.
		{
			Name:     "play",
			Category: "Playback",
			Usage:    "Play/Resume playback.",
			Action:   handlePlay,
			Aliases:  []string{"pl"},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "device", Aliases: []string{"d"}, Usage: "Play music on a device. Use any partial identifier i.e. 'mbp', '064a', 'smartphone', etc."},
				&cli.StringFlag{Name: "track", Aliases: []string{"t"}, Usage: "A track to play."},
				&cli.StringFlag{Name: "album", Aliases: []string{"m"}, Usage: "An album to play."},
				&cli.StringFlag{Name: "artist", Aliases: []string{"r"}, Usage: "An artist to play."},
				&cli.StringFlag{Name: "playlist", Aliases: []string{"l"}, Usage: "An playlist to play."},
			},
		},
		{
			Name:     "pause",
			Category: "Playback",
			Usage:    "Pause playback.",
			Aliases:  []string{"ps"},
			Action:   handlePause,
		},
		{
			Name:     "next",
			Category: "Playback",
			Usage:    "Skip to next track.",
			Aliases:  []string{"nx"},
			Action:   handleNextTrack,
		},
		{
			Name:     "prev",
			Category: "Playback",
			Usage:    "Skip to last track.",
			Aliases:  []string{"pv"},
			Action:   handlePrevTrack,
		},
		// Define Info category commands.
		{
			Name:     "devices",
			Category: "Info",
			Usage:    "Show playable devices.",
			Aliases:  []string{"d", "dv"},
			Action:   handleDevices,
		},
		{
			Name:     "info",
			Category: "Info",
			Usage:    "Show what's currently playing and any track info.",
			Aliases:  []string{"i"},
			Action:   handleInfo,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePlay(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()

	device := c.String("device")
	track := c.String("track")
	album := c.String("album")
	artist := c.String("artist")
	playlist := c.String("playlist")

	switch true {
	case device != "":
		search := strings.ToLower(device)
		for _, d := range Spotify.GetDevices() {

			id, name, t := strings.ToLower(d.ID), strings.ToLower(d.Name), strings.ToLower(d.Type)
			if strings.Contains(id, search) ||
				strings.Contains(name, search) ||
				strings.Contains(t, search) {

				Spotify.PlayOnDevice(d)
				return nil
			}
		}
		fmt.Printf("Could not find any devices '%s'.\n", device)
		os.Exit(1)

	case track != "":
		uri := Spotify.SimpleSearch(track, "track")
		Spotify.PlayURI(uri)

	case album != "":
		uri := Spotify.SimpleSearch(album, "album")
		Spotify.PlayURI(uri)

	case artist != "":
		uri := Spotify.SimpleSearch(artist, "artist")
		Spotify.PlayURI(uri)

	case playlist != "":
		uri := Spotify.SimpleSearch(playlist, "playlist")
		Spotify.PlayURI(uri)

	default:
		Spotify.Play()
	}

	return nil
}

func handlePause(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()
	Spotify.Pause()
	return nil
}

func handleNextTrack(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()
	Spotify.NextTrack()
	return nil
}

func handlePrevTrack(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()
	Spotify.PreviousTrack()
	return nil
}

func handleDevices(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()
	fmt.Printf("[DeviceID]\t\t\t\t\tDeviceType\tName\n")
	for _, d := range Spotify.GetDevices() {
		fmt.Printf("[%s]\t%s\t%s\n", d.ID, d.Type, d.Name)
	}
	return nil
}

func handleInfo(c *cli.Context) error {
	Spotify := spotify.New(ConfigFile)
	Spotify.Authorize()
	state := Spotify.CurrentState()
	isPlayingDesc := "Stopped/Paused"
	if state.IsPlaying {
		isPlayingDesc = "Playing"
	}
	fmt.Printf("[Playback State]: %s\n", isPlayingDesc)
	fmt.Printf("%s\t\t\t%s\t\t%s\t\t%s\n", isPlayingDesc, state.Track.Name, state.Track.Album, state.Track.Artists)
	return nil
}
