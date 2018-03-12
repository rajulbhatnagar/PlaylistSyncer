package models

const (
	GPM_PLAYLIST_TYPE              = "USER_GENERATED"
	GPM_PLAYLIST_SHARESTATE_PUBLIC = "PUBLIC"
)

type GpmCreatePlaylist struct {
	CreationTimestamp     string `json:"creationTimestamp"`
	Deleted               bool   `json:"deleted"`
	LastModifiedTimestamp string `json:"lastModifiedTimestamp"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	PlaylistType          string `json:"type"`
	ShareState            string `json:"shareState"`
}

type GpmSongEntry struct {
	CreationTimestamp     string `json:"creationTimestamp"`
	Deleted               bool   `json:"deleted"`
	LastModifiedTimestamp string `json:"lastModifiedTimestamp"`
	PreviousEntryId       string `json:"precedingEntryId"`
	NextEntryId           string `json:"followingEntryId"`
	ClientId              string `json:"clientId"`
	Source                int    `json:"source"`
	PlayListId            string `json:"playlistId"`
	SongId                string `json:"trackId"`
}

type GpmCreateSongEntry struct {
	CreateGpmSongEntry GpmSongEntry `json:"create"`
}

type GpmCreateSongEntryMutations struct {
	Mutations []GpmCreateSongEntry `json:"mutations"`
}

type GpmAddTracksResponse struct {
	Id           string `json:"id"`
	ResponseCode string `json:"response_code"`
}

type GpmAddTracksMutationsResponse struct {
	Response []GpmAddTracksResponse `json:"mutate_response"`
}

type GpmCreatePlaylistRequest struct {
	GpmCreatePlaylist GpmCreatePlaylist `json:"create"`
}

type GpmCreatePlaylistRequestMutations struct {
	Mutations []GpmCreatePlaylistRequest `json:"mutations"`
}

type GpmCreatePlaylistResponse struct {
	Id           string `json:"id"`
	ResponseCode string `json:"response_code"`
}

type GpmCreatePlaylistMutationsResponse struct {
	Response []GpmCreatePlaylistResponse `json:"mutate_response"`
}

type TrackItem struct {
	Name   string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Id     string `json:"storeId"`
}

type GpmSearchItem struct {
	ItemType string    `json:"type"`
	Track    TrackItem `json:"track"`
}

type GpmSearchResponse struct {
	SuggestedQuery string          `json:"suggestedQuery"`
	Entries        []GpmSearchItem `json:"entries"`
}
