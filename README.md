# batsig

batsig polls `upower` for battery percentage and sends notifications through `libnotify`.
Alerts are stored in the XDG config directory as `batsig/alerts` (`$XDG_CONFIG_HOME/batsig/alerts` or `~/.config/batsig/alerts`), with files named by percentage threshold and the alert text stored as the file contents.

## Usage

- `batsig`
  - Start monitoring battery percentage and fire saved alerts.
- `batsig daemon`
  - Detach and run the monitor in the background.
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

## Alert storage

- Alert file path: `alerts/<percentage>`
- Example: `alerts/100`
- File content: notification text to display when the battery reaches that percentage.

## Example

```bash
batsig alert new 100 "Your battery is fully charged. You can safely unplug."
batsig alert clear 100
```

## Behavior

- The monitor reads the battery device from `upower -e`.
- It parses the percentage from `upower -i <device>`.
- When the battery percentage reaches or exceeds a configured threshold while the battery is discharging, `notify-send` is used to show a libnotify notification.
- Alerts are only fired once per threshold until the battery drops below that threshold again.
