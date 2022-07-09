# tmint

![tmint screenshot](./readme_screenshots/tmint_logo.png?raw=true "tmint logo")

Interface for managing Tmux sessions, windows, and panes.

- created with [tview](https://github.com/rivo/tview) | https://github.com/rivo/tview

![tmint screenshot](./readme_screenshots/tmint_2.png?raw=true "tmint example")

## Status

Work in progress with more features being added.

## Getting Started

Place the executable in your `/usr/local/bin` directory

### Linux

```shell
wget -O tmint.tar.gz https://github.com/gschnall/tmint/releases/download/v0.3.1/tmint_0.3.1_Linux_x86_64.tar.gz && tar xvf tmint.tar.gz && mv tmint /usr/local/bin
```

### Mac & Windows

1. Download the zip file for your sysytem
   https://github.com/gschnall/tmint/releases  
   for macOs use `Darwin_x86_64.tar.gz`
2. Unzip the contents
3. move the `tmint` executable to your `/usr/local/bin` directory
4. once tmint has been started, use the `?` key to bring up a help menu

## Features

- Visualize, navigate, and search target sessions, windows, and panes
- Create, kill, and rename targets
- Detach and attach sessions
- View and save pane scrollback history
- Resize Panes

## Args

`-r` Use resize pane utility to modify pane dimensions (`tmint -r`)  
`-p` prevents tmint from zooming the current pane (`tmint -p`)  

- these args only work while a tmux session is attached

## Issues

Find a bug or want a new feature? Feel free to create an issue

## Contributions

Create a new branch and submit a Pull request with screenshots and a description of changes.

## Licensing

MIT - see LICENSE
