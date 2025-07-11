# PlaylistSyncer

PlaylistSyncer is a command-line tool designed to synchronize your music playlists between different music streaming services. The goal is to allow you to transfer playlists from a source service to a destination service.

## Supported Services

Currently, the tool aims to support:

*   Google Play Music (GPM)
*   Spotify

**Important Note on Current Status:**

*   **Google Play Music is a deprecated service.** While the client for GPM is included, its long-term functionality is not guaranteed.
*   **Reading playlists from Google Play Music is NOT YET IMPLEMENTED.**
*   **Writing playlists to Spotify is NOT YET IMPLEMENTED.**

Due to these limitations, the only currently functional synchronization path is from **Spotify (as source) to Google Play Music (as destination).**

## Features

*   Transfer a specific playlist from a source service to a destination service.
*   Transfer all playlists from a source service to a destination service.
*   Command-line interface for easy operation.

## Prerequisites

*   Go (Golang) installed on your system. You can find installation instructions at [https://golang.org/doc/install](https://golang.org/doc/install).

## Installation & Build

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <repository-url>
    cd PlaylistSyncer
    ```
    *(Assuming the project root is named PlaylistSyncer. This might need adjustment based on the actual cloned directory name.)*

2.  **Build the executable:**
    The main application is `playlistsyncer`. Navigate to its directory and build it:
    ```bash
    cd src/playlistsyncer
    go build
    ```
    This will create an executable named `playlistsyncer` (or `playlistsyncer.exe` on Windows) in the current directory (`src/playlistsyncer`). You might want to move this executable to a directory in your system's PATH for easier access, or run it directly from this location.

## Getting Started & Usage

This section will guide you on how to configure and use PlaylistSyncer.

**Important Reminder:** As noted above, the only currently functional synchronization is from **Spotify (source) to Google Play Music (destination)**. The instructions below reflect this current capability.

### 1. Authentication

Before you can sync playlists, you need to authenticate with both services. The tool will prompt you for credentials when you first run it for a service.

*   **Spotify (Source Service):**
    *   You will need a Spotify OAuth Token.
    *   The tool will prompt: `Enter Spotify OAuth Token.`
    *   You can retrieve this token from the [Spotify Web API Console](https://developer.spotify.com/console/get-playlist/).
    *   When obtaining the token, ensure you select the following scopes:
        *   `playlist-read-private`
        *   `playlist-read-collaborative`
        *   `playlist-modify-public` (though playlist creation on Spotify is not yet implemented)
        *   `playlist-modify-private` (though playlist creation on Spotify is not yet implemented)
        *   `user-read-private`
    *   Paste the obtained token when prompted.

*   **Google Play Music (Destination Service):**
    *   You will be prompted for your Gmail ID and password.
        *   `Enter gmail id:`
        *   `Enter gmail password:`
    *   **2-Factor Authentication:** If you use 2-factor authentication (2FA) with your Google account, you **must** create an "App Password" specifically for PlaylistSyncer. You can generate one at [Google App Passwords](https://myaccount.google.com/apppasswords). Use this app password instead of your regular Google account password when prompted by the tool.

### 2. Running the Syncer

Once you have built the `playlistsyncer` executable (as described in the "Installation & Build" section), you can run it from your terminal.

The command structure is:

```bash
./playlistsyncer -source <source_service> -destination <destination_service> -playlist "<playlist_name>"
```

Or, if you moved it to a directory in your PATH:

```bash
playlistsyncer -source <source_service> -destination <destination_service> -playlist "<playlist_name>"
```

**Arguments:**

*   `-source`: The music service to copy playlists from. Currently, this **must be `spotify`**.
*   `-destination`: The music service to copy playlists to. Currently, this **must be `gpm`** (for Google Play Music).
*   `-playlist`:
    *   To transfer a specific playlist, provide its name in quotes (e.g., `"My Awesome Mix"`). The name must match exactly how it appears in Spotify.
    *   To transfer all your playlists from the source, use `"--all"` (including the quotes).

**Example (Syncing a single playlist from Spotify to GPM):**

Let's say you want to sync your Spotify playlist named "Chill Vibes" to Google Play Music.

1.  Navigate to the directory where `playlistsyncer` was built (e.g., `src/playlistsyncer`):
    ```bash
    cd path/to/your/PlaylistSyncer/src/playlistsyncer
    ```
2.  Run the command:
    ```bash
    ./playlistsyncer -source spotify -destination gpm -playlist "Chill Vibes"
    ```

**Example (Syncing all playlists from Spotify to GPM):**

1.  Navigate to the directory where `playlistsyncer` was built.
2.  Run the command:
    ```bash
    ./playlistsyncer -source spotify -destination gpm -playlist "--all"
    ```

### 3. What to Expect

*   The tool will first attempt to log in to the source service (Spotify). If it's your first time or the token is invalid, it will prompt you for the OAuth token.
*   Then, it will attempt to log in to the destination service (Google Play Music). If it's your first time, it will prompt you for your Gmail ID and password (or App Password).
*   If you chose a specific playlist:
    *   It will search for that playlist on Spotify.
    *   If found, it will then attempt to create a new playlist with the same name and description on Google Play Music.
    *   It will then search for each song from the Spotify playlist on Google Play Music and add any matches to the new GPM playlist.
*   If you chose `--all`:
    *   It will list all your playlists on Spotify.
    *   For each Spotify playlist, it will attempt to create a corresponding playlist on Google Play Music and transfer the songs, similar to the single playlist process.
*   Log messages will be printed to the console indicating the progress and any errors encountered. Songs that cannot be found on GPM will be skipped for that playlist.

### Important Considerations:

*   **Song Matching:** Song matching between services is not always perfect. The tool typically searches by song title and artist. Discrepancies in metadata or availability on the destination service can lead to some songs not being found or being incorrectly matched.
*   **API Limitations & Errors:** Music service APIs can have rate limits or other restrictions. If you are transferring a very large number of playlists or songs, you might encounter errors. The tool attempts to log these.
*   **Google Play Music Deprecation:** As Google Play Music is deprecated, its API may become unavailable or stop working correctly at any point. This tool's GPM functionality is dependent on the continued operation of those APIs.

## How It Works

The tool works by:

1.  Authenticating with both the source and destination music services.
2.  Fetching the specified playlist(s) data (including track information) from the source service.
3.  For each track in the source playlist, it attempts to find a matching track on the destination service.
4.  It then creates a new playlist on the destination service and adds the matched tracks to it.

## Contributing

Contributions are welcome! If you'd like to help improve PlaylistSyncer, please consider:

*   Implementing the missing functionalities (GPM playlist reading, Spotify playlist writing).
*   Adding support for other music services.
*   Improving error handling and track matching logic.
*   Reporting bugs or suggesting new features by opening an issue.

*(Further details on development setup and contribution guidelines could be added here).*

## License

This project is licensed under the MIT License - see the LICENSE file for details.
*(Corrected: The LICENSE file is MIT).*