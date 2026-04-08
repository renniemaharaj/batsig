# batsig

batsig polls `upower` for battery percentage and sends notifications through `libnotify`.
Alerts are stored in the XDG config directory as `batsig/alerts` (`$XDG_CONFIG_HOME/batsig/alerts` or `~/.config/batsig/alerts`), with files named by percentage threshold and the alert text stored as the file contents.

## Usage

- `batsig`
  - Start monitoring battery percentage and fire saved alerts.
- `batsig daemon`
  - Detach and run the monitor in the background.
- `batsig install`
  - Install the `batsig` binary to your user bin directory, copy `notification.wav` to your config, and pre-install the local `alerts` folder.
- `batsig alert new <percentage> <message>`
  - Create or update an alert for the given percentage.
- `batsig alert set <percentage> <message>`
  - Create a new alert or override an existing one.
- `batsig alert clear <percentage>`
  - Remove the saved alert for the given percentage.
- `batsig alert test <percentage>`
  - Trigger the alert notification immediately for testing.
- `batsig alert list`
  - List all configured alerts and their messages.

## Audio notifications

batsig can play an accompanying audio notification when alerts are shown.

- Set `BATSIG_NOTIFY_AUDIO=1` to play the default system notification sound.
- Set `BATSIG_NOTIFY_SOUND=<event>` to play a libcanberra event ID, e.g. `message`.
- Set `BATSIG_NOTIFY_SOUND=/path/to/sound.ogg` to play a local audio file.- If `notification.wav` exists in your config directory, batsig will use that file automatically.
  Audio playback requires `canberra-gtk-play` for sound events, or `paplay` / `aplay` for audio files.

## Alert storage

- Alert file path: `alerts/<percentage>`
- Example: `alerts/100`
- File content: notification text to display when the battery reaches that percentage.

## Example

```bash
batsig alert new 50 "Your battery halfway discharged. You may want to plug in your device."
batsig daemon
batsig alert clear 100
```

## Behavior

- The monitor reads the battery device from `upower -e`.
- It parses the percentage from `upower -i <device>`.
- When the battery percentage reaches or exceeds a configured threshold while the battery is discharging, `notify-send` is used to show a libnotify notification.
- If `BATSIG_NOTIFY_AUDIO=1` or `BATSIG_NOTIFY_SOUND` is configured, batsig also plays an audio notification.
- Alerts are only fired once per threshold until the battery drops below that threshold again.
