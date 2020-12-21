# T mint
Interface for managing Tmux sessions, windows, and panes.
- created with [tview](https://github.com/rivo/tview) 

## Status

Working beta with more features being added.

## Getting Started

Download and unzip the file to the '/usr/bin' directory
1. Download the zip file for your sysytem
https://github.com/gschnall/tmint/releases
for macOs use `tmint_0.1.0_Darwin_x86_64.tar.gz`
3. `sudo tar -xf ~/Downloads/tmint_0.1.0_Darwin_x86_64.tar.gz -C /usr/bin`
3. the `tmint` command should now be available in your terminal
4. once tmint has been started, use the `?` key to bring up the help menu 

## Args
`-p` prevents tmint from zooming the current pane (`tmint -p`)

## Example Keybindings

### vim
Within ~/.vmrc
`map <C-t> :!tmint<CR>`
Tmint can now be launched in vim with Ctrl-t

## Issues

Find a bug or want a new feature? Feel free to create an issue 

## Contributions

Create a new branch and submit a Pull request with screenshots and a description of changes.

## Licensing
MIT - see LICENSE
