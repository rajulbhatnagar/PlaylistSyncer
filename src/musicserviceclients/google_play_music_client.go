package musicserviceclients

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gpsoauth"
	"io"
	"io/ioutil"
	"log"
	"musicserviceclients/models"
	"net/http"
	"net/url"
	"os"
	"strings"
	"uuid"
)

const GPM_OAUTH_SERVICE = "sj"

const MAX_GPM_SEARCH_RESULTS = "10"

const BASE_GPM_URI = "https://mclients.googleapis.com/sj/v2.5/"

const (
	PATH_GPM_CREATE_PLAYLIST       = "playlistbatch"
	PATH_GPM_SEARCH                = "query"
	PATH_GPM_ADD_SONGS_TO_PLAYLIST = "plentriesbatch"
)

type googlePlayMusicClient struct {
	oAuthToken string
	client     *http.Client
}

func NewGooglePlayMusicClient() (MediaServiceClient, error) {
	client := &http.Client{}
	return &googlePlayMusicClient{client: client}, nil
}

func (c *googlePlayMusicClient) Login() error {
	reader := bufio.NewReader(os.Stdin)
	log.Println("If you are using Gmail 2 factor authentication please create a app specific password at https://security.google.com/settings/security/apppasswords and use that.")
	log.Print("Enter gmail id: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to get username %v", err)
	}
	log.Print("Enter gmail password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to get password %v", err)
	}
	oAuthToken, err := gpsoauth.Login(username, password, gpsoauth.GetNode(), GPM_OAUTH_SERVICE)
	if err != nil {
		return fmt.Errorf("failed to get master token %v", err)
	}
	c.oAuthToken = oAuthToken
	return nil
}

func (c *googlePlayMusicClient) ListPlaylist(playListName string) (*Playlist, error) {
	return nil, errors.New("Unimplemented function")
}

func (c *googlePlayMusicClient) ListAllPlaylists() ([]Playlist, error) {
	return nil, errors.New("Unimplemented function")
}

func (c *googlePlayMusicClient) CreatePlaylist(playlistName, playListDescription string, songs []Song) error {
	id, err := c.createNewPlaylist(playlistName, playListDescription)
	if err != nil {
		return fmt.Errorf("failed to create new empty playlist [name=%s][description=%s][err=%v]", playlistName, playListDescription, err)
	}
	err = c.addTracksToPlaylist(id, songs)
	if err != nil {
		return fmt.Errorf("failed to add the following songs %v", err)
	}
	return nil
}

func (c *googlePlayMusicClient) createNewPlaylist(playlistName, playListDescription string) (string, error) {
	request := &models.GpmCreatePlaylistRequestMutations{Mutations: []models.GpmCreatePlaylistRequest{
		{GpmCreatePlaylist: models.GpmCreatePlaylist{
			Name:                  playlistName,
			Description:           playListDescription,
			Deleted:               false,
			CreationTimestamp:     "-1",
			LastModifiedTimestamp: "0",
			PlaylistType:          models.GPM_PLAYLIST_TYPE,
			ShareState:            models.GPM_PLAYLIST_SHARESTATE_PUBLIC}}}}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to create json request [err=%v]", err)
	}
	response, err := c.makeRequest(http.MethodPost, PATH_GPM_CREATE_PLAYLIST, bytes.NewReader(jsonRequest))
	if err != nil {
		return "", fmt.Errorf("failed to create playlist [name=%s][err=%v]", playlistName, err)
	}
	dec := json.NewDecoder(strings.NewReader(response))
	var responseObj models.GpmCreatePlaylistMutationsResponse
	err = dec.Decode(&responseObj)
	if err != nil {
		return "", fmt.Errorf("failed to parse response [response=%s][err=%v]", response, err)
	}
	return responseObj.Response[0].Id, nil
}

func (c *googlePlayMusicClient) addTracksToPlaylist(id string, songs []Song) error {
	var errorList []error
	var tracks []models.GpmSearchItem
	//TODO Make concurrent
	for _, song := range songs {
		track, err := c.findBestMatchSong(song)
		if err != nil {
			errorList = append(errorList, err)
		} else {
			tracks = append(tracks, *track)
		}
	}
	if len(tracks) > 0 {
		var addTrackEntries []models.GpmCreateSongEntry
		prevId := ""
		currId := uuid.NewUUID().String()
		nextId := uuid.NewUUID().String()
		for i, track := range tracks {
			log.Printf("Adding Track %s\n", track.Track.Name)
			source := 1
			if strings.HasPrefix(track.Track.Id, "T") {
				source = 2
			}
			entry := models.GpmCreateSongEntry{CreateGpmSongEntry: models.GpmSongEntry{
				CreationTimestamp:     "-1",
				Deleted:               false,
				LastModifiedTimestamp: "0",
				PlayListId:            id,
				SongId:                track.Track.Id,
				Source:                source,
				ClientId:              currId}}

			if i > 0 {
				entry.CreateGpmSongEntry.PreviousEntryId = prevId
			}
			if i < len(tracks)-1 {
				entry.CreateGpmSongEntry.NextEntryId = nextId
			}

			addTrackEntries = append(addTrackEntries, entry)

			prevId = currId
			currId = nextId
			nextId = uuid.NewUUID().String()
		}
		addTracksRequest := models.GpmCreateSongEntryMutations{Mutations: addTrackEntries}
		jsonRequest, err := json.Marshal(addTracksRequest)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to create json request [err=%v]", err))
		} else {
			response, err := c.makeRequest(http.MethodPost, PATH_GPM_ADD_SONGS_TO_PLAYLIST, bytes.NewReader(jsonRequest))
			if err != nil {
				errorList = append(errorList, fmt.Errorf("failed to add songs to playlist [id=%s][err=%v]", id, err))
			}
			dec := json.NewDecoder(strings.NewReader(response))
			var responseObj models.GpmAddTracksMutationsResponse
			err = dec.Decode(&responseObj)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("failed to parse response [response=%s][err=%v]", response, err))
			}
			for _, responseEntry := range responseObj.Response {
				if responseEntry.ResponseCode != "OK" {
					errorList = append(errorList, fmt.Errorf("failed to add track [id=%s][err=%v]", responseEntry.Id, err))
				}
			}
		}
	}
	if len(errorList) == 0 {
		return nil
	} else {
		return flattenErrors(errorList)
	}
}

func (c *googlePlayMusicClient) makeRequest(method, path string, body io.Reader) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", BASE_GPM_URI, path), body)
	if err != nil {
		return "", fmt.Errorf("failed to create http request for [path=%s][err=%v]", path, err)
	}
	q := req.URL.Query()
	q.Add("tier", "aa")
	q.Add("hl", "en_US")
	q.Add("dv", "0")
	q.Add("alt", "json")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("GoogleLogin auth=%s", c.oAuthToken))
	req.Header.Add("Content-type", "application/json")
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
func (c *googlePlayMusicClient) findBestMatchSong(song Song) (*models.GpmSearchItem, error) {
	seachQuery := c.searchQuery(song)
	for {
		query := url.Values{}
		query.Add("q", seachQuery)
		query.Add("max-results", MAX_GPM_SEARCH_RESULTS)
		query.Add("ct", "1")
		response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("%s?%s", PATH_GPM_SEARCH, query.Encode()), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to search request [song=%s][err=%v]", song.Name, err)
		}
		dec := json.NewDecoder(strings.NewReader(response))
		var responseObj models.GpmSearchResponse
		err = dec.Decode(&responseObj)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response [response=%s][err=%v]", response, err)
		}
		if len(responseObj.SuggestedQuery) != 0 {
			seachQuery = responseObj.SuggestedQuery
			continue
		}
		for _, track := range responseObj.Entries {
			if track.ItemType == "1" {
				return &track, nil
			}
		}
		return nil, fmt.Errorf("failed to find a match [song=%s][err=%v]", song.Name, err)
	}

}

func (c *googlePlayMusicClient) searchQuery(song Song) string {
	if len(song.Artists) == 0 {
		return song.Name
	} else {
		return fmt.Sprintf("%s - %s", song.Artists[0].Name, song.Name)
	}
}

func flattenErrors(errorList []error) error {
	var err error = nil
	for _, errorVal := range errorList {
		err = errors.New(fmt.Sprintf("\n[err=%v]", errorVal))
	}
	return err
}
