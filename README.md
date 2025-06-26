# Hyprcircade

Daemon that manages dark/light theme switching for system

## How it works

You specify time when you want to switch your dark or light theme, files that you want to modify when theme switches (e.g. configuration files where you want to switch between dark/light modes) and commands that need to be executed when theme switches (e.g. send notification about theme switching or trigger applications to reload its configuration)

## Installation

```bash
go install github.com/tfk70/hyprcircade/cmd/hyprcircade@v0.0.7
```

## Usage

1. Create `hyprcircade.conf` in `$HOME/.config/hypr` directory:

```hyprlang
general {
  anchor = THEME_SWITCHER_TARGET # hyprcircade will search for this anchor in files for strings replacement

  # 24-hours format specification
  dark-at = 20 # switch to dark theme at 20:00 (8 P.M.)
  light-at = 8 # switch to light theme at 8:00 (8 A.M.)
}

file {
  path = ./configuration.yaml # path to configuration file that needs to be modified
  day-value = light # value that needs to be placed when theme is light
  night-value = dark # value that needs to be placed when theme is dark
  ignore-anchor = false # if you don't want (or can't) specify anchors in configuration file - you can ignore them
}

file { # several files can be specified
  path = ./somefile2.yaml
  day-value = light2
  night-value = dark2
  ignore-anchor = true
}

command {
  day-exec = notify-send -t 3000 "Switching theme to light" # command that will be executed when light theme is applied
}


command {
  day-exec = another command # several commands can be specified
}

command {
  night-exec = notify-send -t 3000 "Switching theme to dark" # command that will be executed when dark theme is applied
}
```

If we have `configuration.yaml` like this:

```yaml
configuration:
    theme: light # THEME_SWITCHER_TARGET
    name: lightning mcqueen
```

When dark theme is applied hyprcircade will update this file like this: (if ignore-anchor is false)

```yaml
configuration:
    theme: dark # THEME_SWITCHER_TARGET
    name: lightning mcqueen # line has substring 'light' but it won't be modified because this line has no anchor
```

If ignore-anchor was set to true then the result will be:

```yaml
configuration:
    theme: dark # THEME_SWITCHER_TARGET - anchor will be ignored
    name: darkning mcqueen # was modified
```

2. Start `hyprcircade` daemon

```bash
hyprcircade
```

This will start daemon in foreground mode and also apply theme based on your current time of day. If you want to omit this functionality just do:

```bash
hyprcircade --apply-on-start=false
```

You can also trigger theme switching manually

```bash
hyprcircade switch light
hyprcircade switch dark
```

## Starting in background

### Hyprland

Recommended option is to start hyprcircade using hyprland's `exec-one` option:

```hyprlang
exec-once = hyprcircade &
```

### From current session

However, if you want to start hyprcircade from already running hyprland session you can do:

```bash
hyprcircade &
disown %1
```

This will start hyprcircade daemon in background and detach it from current session. However, this will print stdout logs to you current terminal session and it is recommended to restart the session after launching hyprcircade this way

You can also use `nohup` to redirect stdout logs to a file:

```bash
nohup hyprcircade &
disown %1
```

### Systemd

If you want more control over your daemons you can launch hyprircade as a systemd service. Create `.config/systemd/hyprcircade.service` file:

```service
[Unit]
Description=Hyprcircade dark/light theme switching daemon

[Service]
Type=simple
ExecStart=/home/user/go/bin/hyprcircade # Path to hyprcircade binary

[Install]
WantedBy=graphical-session.target # We want to start our daemon after graphical session was initialized
```

Now enable it:

```bash
systemctl --user enable --now hyprcircade.service
```

And view it's status:

```bash
systemctl --user status hyprcircade.service
```

## CLI Reference

```bash
NAME:
   hyprcircade - Dark/light theme manager for hyprland

USAGE:
   hyprcircade [global options] [command [command options]]

VERSION:
   v0.0.1

COMMANDS:
   switch
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config string, -c string  Path to hyprcircade configuration file (default: "/home/darius/.config/hypr/hyprcircade.conf") [$HYPRCIRCADE_CONFIGURATION_FILE]
   --debug                     Enable debug logging (default: false) [$HYPRCIRCADE_DEBUG]
   --apply-on-start            Apply theme based on time of day on daemon startup (default: true) [$HYPRCIRCADE_APPLY_ON_START]
   --help, -h                  show help
   --version, -v               print the version

COPYRIGHT:
   (c) 2025 TFK70
```

## Real example

```hyprlang
general {
  anchor = THEME_SWITCHER_TARGET

  dark-at = 20
  light-at = 8
}

# Set palette for wallust form dark to light
file {
  path = /home/user/.config/wallust/wallust.toml
  day-value = softlight
  night-value = dark
  ignore-anchor = false
}

# Set background for terminal
file {
  path = /home/user/.config/alacritty/alacritty.toml
  day-value = 0xffffff
  night-value = 0x000000
  ignore-anchor = false
}

# Set colors for rofi theme
file {
  path = /home/user/.config/rofi/themes/custom/launcher.rasi
  day-value = white
  night-value = black
  ignore-anchor = false
}

# Update taskwarrior theme
file {
  path = /home/user/.taskrc
  day-value = light-256
  night-value = dark-blue-256
  ignore-anchor = false
}

# Update neovim colorscheme (I use plugin to read colorscheme from this file and update it in real time)
file {
  path = /home/user/.config/nvim/colorscheme/current_colorscheme
  day-value = github_light
  night-value = tokyonight-moon
  ignore-anchor = true
}

# Send notification about theme switching
command {
  day-exec = notify-send -t 3000 "Switching theme to light"
}

# Execute script that updates wallpaper, regenerates wallust theme and reloads waybar, swaync, etc.
command {
  day-exec = /home/user/.config/hypr/compositions/wallpaper.sh day-wallpaper.png
}

# Set system-wide color scheme to light so applications like Google Chrome (and websites) could synchronize its theme with OS
command {
  day-exec = gsettings set org.gnome.desktop.interface color-scheme 'prefer-light'
}

# Set GTK theme
command {
  day-exec = gsettings set org.gnome.desktop.interface gtk-theme 'Flat-Remix-GTK-Blue-Light'
}

## Following commands do the same for dark theme

command {
  night-exec = notify-send -t 3000 "Switching theme to dark"
}

command {
  night-exec = /home/user/.config/hypr/compositions/wallpaper.sh night-wallpaper.png
}

command {
  night-exec = gsettings set org.gnome.desktop.interface gtk-theme 'Flat-Remix-GTK-Blue-Dark'
}

command {
  night-exec = gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'
}

```

## Relation to hyprland

Tool was called **hypr**circade because you use hyprlang for hyprcircade configuration. However, I plan to add support for yaml configuration file for people who want to use this tool in setups without hyprland (because hyprland is not required).

## Known issues and limitations

- No modules or variables support (like in hyprland configuration file)
    - The reason is that there is not official hyprlang parser for golang

## Todo

- [ ] Implement yaml configuration file as an alternative
