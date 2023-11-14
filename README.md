# Clickup TUI

![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

**clickup-tui** is a Golang-based project that leverages the Charm library, including components like lipgloss, bubbles, bubbletea, and glamour. This project aims to provide a terminal user interface (TUI) for easy access to ClickUp, allowing users to interact with ClickUp features conveniently from the command line.

## Features

- **Intuitive TUI:** Enjoy a user-friendly terminal interface for managing your ClickUp tasks and projects.
- **Efficient Navigation:** Navigate seamlessly through your ClickUp workspace using keyboard shortcuts.
- **Task Management:** View, create, update, and delete tasks without leaving the terminal.
- **Project Overview:** Get a quick overview of your ClickUp projects and their statuses.

## Installation

To install **clickup-tui**, you can use the following steps:

```bash
go get -u github.com/prgres/clickup-tui
```

Usage
Once installed, you can run `clickup-tui` from the terminal:

```bash
Copy code
clickup-tui
```

Use the arrow keys, enter, and other relevant keyboard shortcuts to navigate through the TUI and interact with ClickUp.

## Configuration
Before using the tool, set up your ClickUp API key and configure any necessary settings. You can do this by creating a configuration file or using environment variables. Please take a look at the documentation for details on how to set up your configuration.

## Dependencies

- [Charm](https://github.com/charmbracelet/charm): A collection of terminal user interface components.
- [lipgloss](https://github.com/charmbracelet/lipgloss): Styling for your terminal interfaces.
- [bubbles](https://github.com/charmbracelet/bubbles): A delightful way to render terminal tables.
- [bubbletea](https://github.com/charmbracelet/bubbletea): A functional framework for building terminal applications.

## Contributing
Contributions are welcome! If you find any bugs or have suggestions for improvement, please open an issue or submit a pull request.

### TODO
- [ ] async fetch data - sync in schedule
- [ ] config in file
- [ ] oauth login
- [ ] edit status and assigned
- [ ] polish space/folder/list flow
- [ ] add team setting
- [ ] store in the memory or file last selected (team, space, etc) 
- [ ] add team handling

License
This project is licensed under the MIT License - see the LICENSE file for details.
