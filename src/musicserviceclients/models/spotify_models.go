package models

type SpotifyUser struct {
	Name string `json:"display_name"`
	Id   string `json:"id"`
}

type SpotifyAlbum struct {
	Name string `json:"name"`
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyTrackWrapper struct {
	Track SpotifyTrack `json:"track"`
}

type SpotifyTrack struct {
	Name    string          `json:"name"`
	Album   SpotifyAlbum    `json:"album"`
	Artists []SpotifyArtist `json:"artists"`
}

type SpotifyPlaylistTracks struct {
	Tracks []SpotifyTrackWrapper `json:"items"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Total  int                   `json:"total"`
}

type SpotifyPlaylistsOwner struct {
	Name string `json:"display_name"`
	Id   string `json:"id"`
}

type SpotifyPlaylists struct {
	Playlists []SpotifyPlaylist `json:"items"`
	Limit     int               `json:"limit"`
	Offset    int               `json:"offset"`
	Total     int               `json:"total"`
}

type SpotifyPlaylist struct {
	Name        string                `json:"name"`
	Id          string                `json:"id"`
	Description string                `json:"description"`
	Tracks      SpotifyPlaylistTracks `json:"tracks"`
	Owner       SpotifyPlaylistsOwner `json:"owner"`
}
