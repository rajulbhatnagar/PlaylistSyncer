package musicserviceclients

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"musicserviceclients/models"
	"net/http"
	"os"
	"strings"
)

const BASE_SPOTIFY_URI = "https://api.spotify.com/v1/"

const (
	PATH_SPOTIFY_USER            = "me"
	PATH_SPOTIFY_LIST_PLAYLISTS  = "users/%s/playlists?limit=%d&offset=%d"
	PATH_SPOTIFY_LIST_PLAYLIST   = "users/%s/playlists/%s/tracks?limit=%d&offset=%d"
	PATH_SPOTIFY_SEARCH          = ""
	PATH_SPOTIFY_CREATE_PLAYLIST = ""
	PATH_SPOTIFY_ADD_TRACK       = ""
)

type spotifyClient struct {
	oAuthToken string
	userId     string
	client     *http.Client
}

func NewSpotifyClient() (MediaServiceClient, error) {
	client := &http.Client{}
	return &spotifyClient{client: client}, nil
}

func (c *spotifyClient) Login() error {
	log.Println("\nEnter Spotify OAuth Token.\nYou can retrieve the token at https://developer.spotify.com/web-api/console/get-playlist.\nSelect Scopes[playlist-read-private, playlist-read-collaborative, playlist-modify-public, playlist-modify-collaborative, user-read-private]")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		c.oAuthToken = scanner.Text()
	}
	if scanner.Err() != nil {
		return fmt.Errorf("failed to fetch OAuth token [err=%v]", scanner.Err())
	}
	user, err := c.getCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user profile [err=%v]", err)
	}
	log.Printf("Welcome to Spotify %s", user.Name)
	c.userId = user.Id
	return nil
}

func (c *spotifyClient) ListPlaylist(playListName string) (*Playlist, error) {
	limit := 50
	offset := 0
	for {
		response, err := c.makeRequest(http.MethodGet, fmt.Sprintf(PATH_SPOTIFY_LIST_PLAYLISTS, c.userId, limit, offset), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list of user playlists [err=%v]", err)
		}
		dec := json.NewDecoder(strings.NewReader(response))
		var spotifyPlaylists models.SpotifyPlaylists
		err = dec.Decode(&spotifyPlaylists)
		if err != nil {
			return nil, err
		}
		for i, spotifyPlaylist := range spotifyPlaylists.Playlists {
			if spotifyPlaylist.Name == strings.TrimSpace(playListName) {
				playlist, err := c.getPlaylist(spotifyPlaylist)
				if err != nil {
					return nil, fmt.Errorf("Failed to retrieve playlist info at [offset=%d][name=%s][err=%v]", offset+i, playListName, err)
				}
				return playlist, nil
			}
		}
		offset += limit
		if spotifyPlaylists.Total < offset {
			break
		}
	}
	return nil, fmt.Errorf("failed to find playlist with [name=%s]", playListName)
}

func (c *spotifyClient) ListAllPlaylists() ([]Playlist, error) {
	limit := 50
	offset := 0
	var playlists []Playlist

	for {
		response, err := c.makeRequest(http.MethodGet, fmt.Sprintf(PATH_SPOTIFY_LIST_PLAYLISTS, c.userId, limit, offset), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list of user playlists [err=%v]", err)
		}
		dec := json.NewDecoder(strings.NewReader(response))
		var spotifyPlaylists models.SpotifyPlaylists
		err = dec.Decode(&spotifyPlaylists)
		if err != nil {
			return playlists, err
		}
		for i, spotifyPlaylist := range spotifyPlaylists.Playlists {
			playlist, err := c.getPlaylist(spotifyPlaylist)
			if err != nil {
				return playlists, fmt.Errorf("Failed to retrieve playlist info at [offset=%d][err=%v]", offset+i, err)
			}
			playlists = append(playlists, *playlist)
		}
		offset += limit
		if spotifyPlaylists.Total < offset {
			break
		}
	}
	return playlists, nil
}

func (c *spotifyClient) CreatePlaylist(playlistName, playListDescription string, songs []Song) error {
	return errors.New("Unimplemented function")
}

func (c *spotifyClient) makeRequest(method, path string, body io.Reader) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", BASE_SPOTIFY_URI, path), body)
	if err != nil {
		return "", fmt.Errorf("failed to create http request for [path=%s][err=%v]", path, err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.oAuthToken))
	response, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create http request for [path=%s][err=%v]", path, err)
	}
	result, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to create http request for [path=%s][err=%v]", path, err)
	}
	if response.StatusCode == http.StatusOK {
		return string(result), nil
	} else {
		return "", fmt.Errorf("failed to make http request for [path=%s][httpstatus=%d][err=%s]", path, response.StatusCode, string(result))
	}
}

func (c *spotifyClient) getCurrentUser() (*models.SpotifyUser, error) {
	response, err := c.makeRequest(http.MethodGet, PATH_SPOTIFY_USER, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currentUserId [err=%v]", err)
	}
	dec := json.NewDecoder(strings.NewReader(response))
	var user models.SpotifyUser
	err = dec.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *spotifyClient) getPlaylist(playlist models.SpotifyPlaylist) (*Playlist, error) {
	limit := 100
	offset := 0
	var mergedSpotifyPlaylistTracks *models.SpotifyPlaylistTracks
	for {
		response, err := c.makeRequest(http.MethodGet, fmt.Sprintf(PATH_SPOTIFY_LIST_PLAYLIST, playlist.Owner.Id, playlist.Id, limit, offset), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist [name=%s][id=%s][err=%v]", playlist.Name, playlist.Id, err)
		}
		dec := json.NewDecoder(strings.NewReader(response))
		var spotifyPlaylistTracks models.SpotifyPlaylistTracks
		err = dec.Decode(&spotifyPlaylistTracks)
		if err != nil {
			return nil, err
		}
		if mergedSpotifyPlaylistTracks == nil {
			mergedSpotifyPlaylistTracks = &spotifyPlaylistTracks
		} else {
			mergedSpotifyPlaylistTracks.Tracks = append(mergedSpotifyPlaylistTracks.Tracks, spotifyPlaylistTracks.Tracks...)
		}
		offset += limit
		if mergedSpotifyPlaylistTracks.Total < offset {
			break
		}
	}
	return mediaPlaylist(playlist, *mergedSpotifyPlaylistTracks), nil
}

func mediaPlaylist(spotifyPlaylist models.SpotifyPlaylist, spotifyPlaylistTracks models.SpotifyPlaylistTracks) *Playlist {
	return &Playlist{Name: spotifyPlaylist.Name, Description: spotifyPlaylist.Description, Id: spotifyPlaylist.Id, Songs: mediaSongs(spotifyPlaylistTracks.Tracks)}
}
func mediaSongs(tracks []models.SpotifyTrackWrapper) []Song {
	var Songs []Song
	for _, track := range tracks {
		Songs = append(Songs, mediaSong(track.Track))
	}
	return Songs
}
func mediaSong(track models.SpotifyTrack) Song {
	return Song{Name: track.Name, Album: mediaAlbum(track.Album), Artists: mediaArtists(track.Artists)}
}
func mediaArtists(artists []models.SpotifyArtist) []Artist {
	var Artists []Artist
	for _, artist := range artists {
		Artists = append(Artists, mediaArtist(artist))
	}
	return Artists
}
func mediaArtist(artist models.SpotifyArtist) Artist {
	return Artist{Name: artist.Name}
}

func mediaAlbum(album models.SpotifyAlbum) Album {
	return Album{Name: album.Name}
}
