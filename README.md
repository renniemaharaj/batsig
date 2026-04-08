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
- `batsig alert -c <message>`
  - Set the charging state alert message for when wall power is connected.
- `batsig alert -d <message>`
  - Set the discharging state alert message for when wall power is disconnected.
- `batsig alert new [-c] <percentage> <message>`
  - Create or update an alert in the discharging folder by default, or in the charging folder when `-c` is used.
- `batsig alert set [-c] <percentage> <message>`
  - Create a new alert or override an existing one.
- `batsig alert clear [-c] <percentage>`
  - Remove the saved alert from the selected folder.
- `batsig alert clear -c`
  - Remove the charging state alert.
- `batsig alert clear -d`
  - Remove the discharging state alert.
- `batsig alert test [-c] <percentage>`
  - Trigger a specific alert for testing.
- `batsig alert test -c`
  - Trigger the charging state alert for testing.
- `batsig alert test -d`
  - Trigger the discharging state alert for testing.
- `batsig alert mv <percentage> <charging|discharging>`
  - Move an alert between the charging and discharging folders.
- `batsig alert list`
  - List all configured alerts and their messages.

## Audio notifications

batsig can play an accompanying audio notification when alerts are shown.

- Set `BATSIG_NOTIFY_AUDIO=1` to play the default system notification sound.
- Set `BATSIG_NOTIFY_SOUND=<event>` to play a libcanberra event ID, e.g. `message`.
- Set `BATSIG_NOTIFY_SOUND=/path/to/sound.ogg` to play a local audio file.
- If `notification.wav` exists in your config directory, batsig will use that file automatically.
  Audio playback requires `canberra-gtk-play` for sound events, or `paplay` / `aplay` for audio files.

## Alert storage

- Discharging alerts: `alerts/<percentage>`
- Charging alerts: `alerts/charging/<percentage>`
- Charging state alert: `alerts/state/charging`
- Discharging state alert: `alerts/state/discharging`
- File content: notification text to display when the battery reaches that percentage exactly, or a state event occurs.

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
