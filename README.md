# T mint
Interface for managing Tmux sessions, windows, and panes.
- created with [tview](https://github.com/rivo/tview) 

## Status

Working beta with more features being added.

## Getting Started

1. Download and unzip the release for your system
https://github.com/gschnall/tmint/releases
for macOs use `tmint_0.1.0_Darwin_x86_64.tar.gz`
2. save the executable somewhere like `~/.tmux/tmint`
3. run it `~/path-to-my-executable/tmint`
4. use the `?` key to bring up the help menu

## Keybindings

### bashrc or zshrc alias
`alias tmint=~/path-to-my-executable/tmint`
example:
`alias tmint=~/.tmux/tmint`

### vim keybinding
Within ~/.vmrc
`map <C-t> :!~/.tmux/tmint<CR>`
Tmint can now be launched in vim with Ctrl-t

### tmux keybinding
Work in progress.

## Issues

Find a bug or want a new feature? Feel free to create an issue 

## Contributions

Create a new branch and submit a Pull request with screenshots and text describing changes.

## Licensing
MIT - see LICENSE
