package musicserviceclients

type Album struct {
	Name string
}

type Artist struct {
	Name string
}

type Song struct {
	Name    string
	Album   Album
	Artists []Artist
}

type Playlist struct {
	Name        string
	Description string
	Id          string
	Songs       []Song
}

type MediaServiceClient interface {
	Login() error
	ListPlaylist(string) (*Playlist, error)
	ListAllPlaylists() ([]Playlist, error)
	CreatePlaylist(string, string, []Song) error
}
