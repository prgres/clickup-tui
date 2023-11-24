# Clickup TUI

![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

**clickup-tui** is a Golang-based project that leverages the Charm library, including components like lipgloss, bubbles, bubbletea, and glamour. This project aims to provide a terminal user interface (TUI) for easy access to ClickUp, allowing users to interact with ClickUp features conveniently from the command line.

## Features

- **Intuitive TUI:** Enjoy a user-friendly terminal interface for managing your ClickUp tasks and projects.
- **Efficient Navigation:** Navigate seamlessly through your ClickUp workspace using keyboard shortcuts.
- **Task Management:** View, create, update, and delete tasks without leaving the terminal.
- **Project Overview:** Get a quick overview of your ClickUp projects and their statuses.

## Installation

> [!CAUTION]
> This will not work OOTB until the repository is private. Cloning the repo is a way to go for now. But if you want to `get` it, follow the https://go.dev/doc/faq#git_https

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
### Clonig the repository
To run without building the binary, simply just clone the repo, set config, and run `go run .` in the root.
## Flags
For flags help simply exec:
```
clickup-tui -h

clickup-tui - A terminal user interface for ClickUp
Usage:
  clickup-tui [flags]
Flags:
      --cache-path string   The path to the cache directory (default "./cache")
      --clean-cache         Cleans cache data
      --clean-cache-only    Cleans cache data and exits
  -c, --config string       A config filename (default "config.yaml")
      --debug               Enable debug mode
      --debug-deep          Enable deep debug mode
  -h, --help                Show help
  -v, --version             Show version
```

## Configuration
Before using the tool, set up your ClickUp API key and configure any necessary settings. You can do this by creating a configuration file or using environment variables. Please take a look at the documentation for details on how to set up your configuration.
The app looks for a config file in paths:
- .
- etc/clickup-tui
- home/user/clickup-tui
- $HOME/.config/clickup-tui
For now, you have to manually create that (this will be addressed) - just copy the [`config.yaml.example`](config.yaml.example) file, remove the example suffix, and fill properties (only token is required). In the future, these settings will be manipulated within the app.
### How to obtain a Clickup token
Follow the steps: [ClickUp API docs: Generate your personal API token](https://clickup.com/api/developer-portal/authentication/#generate-your-personal-api-token)
## Dependencies

- [Charm](https://github.com/charmbracelet/charm): A collection of terminal user interface components.
- [lipgloss](https://github.com/charmbracelet/lipgloss): Styling for your terminal interfaces.
- [bubbles](https://github.com/charmbracelet/bubbles): A delightful way to render terminal tables.
- [bubbletea](https://github.com/charmbracelet/bubbletea): A functional framework for building terminal applications.
- [bubble-tabble](https://github.com/Evertras/bubble-table): A customizable, interactive table component for the Bubble Tea framework

## Contributing
Contributions are welcome! If you find any bugs or have suggestions for improvement, please open an issue or submit a pull request.

# License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
