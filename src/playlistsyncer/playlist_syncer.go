package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"musicserviceclients"
)

const PLAYLIST_ALL = "--all"

//Services
const (
	SPOTIFY           = "spotify"
	GOOGLE_PLAY_MUSIC = "gpm"
)

type CliArguments struct {
	sourceService      string
	destinationService string
	playList           string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatalf("%v", err)
	}
	sourceClient, err := client(args.sourceService)
	if err != nil {
		log.Fatalf("Failed to initialize client for [service=%s, err=%v]", args.sourceService, err)
	}
	err = sourceClient.Login()
	if err != nil {
		log.Fatalf("Failed to login for [service=%s, err=%v]", args.sourceService, err)
	}
	destinationClient, err := client(args.destinationService)
	if err != nil {
		log.Fatalf("Failed to initialize client for [service=%s, err=%v]", args.destinationService, err)
	}
	err = destinationClient.Login()
	if err != nil {
		log.Fatalf("Failed to login for [service=%s, err=%v]", args.destinationService, err)
	}
	if args.playList == PLAYLIST_ALL {
		log.Println("Listing all playlists")
		playlists, err := sourceClient.ListAllPlaylists()
		if err != nil {
			log.Fatalf("Failed to list playlists for [service=%s, err=%v]", args.sourceService, err)
		}
		for _, playlist := range playlists {
			log.Printf("Creating Playlist %s", playlist.Name)
			err = destinationClient.CreatePlaylist(playlist.Name, playlist.Description, playlist.Songs)
			if err != nil {
				log.Printf("Failed to create playlist for [name=%s, service=%s, err=%v]", playlist.Name, args.destinationService, err)
			}
		}
	} else {
		log.Printf("Listing %s playlist", args.playList)
		playlist, err := sourceClient.ListPlaylist(args.playList)
		if err != nil {
			log.Fatalf("Failed to list playlist for [name=%s, service=%s, err=%v]", args.playList, args.sourceService, err)
		}
		log.Printf("Creating Playlist %s", playlist.Name)
		err = destinationClient.CreatePlaylist(playlist.Name, playlist.Description, playlist.Songs)
		if err != nil {
			log.Printf("Failed to create playlist for [name=%s, service=%s, err=%v]", playlist.Name, args.destinationService, err)
		}
	}
}

func client(service string) (musicserviceclients.MediaServiceClient, error) {
	switch service {
	case SPOTIFY:
		return musicserviceclients.NewSpotifyClient()
	case GOOGLE_PLAY_MUSIC:
		return musicserviceclients.NewGooglePlayMusicClient()
	default:
		return nil, fmt.Errorf("Unimplemented service %s", service)
	}
}

func parseArgs() (*CliArguments, error) {
	sourceService := flag.String("source", "", "The source music service")
	destinationService := flag.String("destination", "", "The destination music service")
	playList := flag.String("playlist", "", "The name of the playlist you want to transfer. Use '--all' for moving all playlists")
	flag.Parse()

	var errs []error
	if len(*sourceService) == 0 || !validService(sourceService) {
		errs = append(errs, fmt.Errorf("Invalid source service=%s", *sourceService))
	}

	if len(*destinationService) == 0 || !validService(destinationService) {
		errs = append(errs, fmt.Errorf("Invalid destination service=%s", *destinationService))
	}

	if len(*playList) == 0 {
		errs = append(errs, fmt.Errorf("You need to specify playlist to transfer"))
	}

	var err error
	for _, v := range errs {
		log.Printf("failed to parse args [error=%v]", v)
		err = errors.New("failed to parse args")
	}

	if err != nil {
		flag.PrintDefaults()
		return nil, err
	}

	return &CliArguments{sourceService: *sourceService, destinationService: *destinationService, playList: *playList}, nil
}

func validService(string *string) bool {
	return true
}
