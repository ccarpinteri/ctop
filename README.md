<p align="center"><img width="200px" src="/_docs/img/logo.png" alt="ctop"/></p>

#

![release][release]

Top-like interface for container metrics

`ctop` provides a concise and condensed overview of real-time metrics for multiple containers:
<p align="center"><img src="_docs/img/grid.gif" alt="ctop"/></p>

as well as a [single container view][single_view] for inspecting a specific container.

`ctop` comes with built-in support for Docker and runC; connectors for other container and cluster systems are planned for future releases.

## Install

Fetch the [latest release](https://github.com/ccarpinteri/ctop/releases) for your platform:


#### Linux (Generic)

```bash
sudo wget https://github.com/ccarpinteri/ctop/releases/download/v0.2.0/ctop-0.2.0-linux-amd64.tar.xz
sudo tar -xJf ctop-0.2.0-linux-amd64.tar.xz
sudo mv ctop /usr/local/bin/ctop
sudo chmod +x /usr/local/bin/ctop
```

#### OS X

```bash
curl -Lo ctop.tar.xz https://github.com/ccarpinteri/ctop/releases/download/v0.2.0/ctop-0.2.0-apple-darwin-arm64.tar.xz
tar -xJf ctop.tar.xz
sudo mv ctop /usr/local/bin/ctop
sudo chmod +x /usr/local/bin/ctop
```

#### Windows

##### Currently only available using the binary release.

Work in progress to provide a distribution via Scoop. or winget.

'```powershell
ctop-windows-amd64.exe
```

#### Docker

```bash
docker run --rm -ti \
  --name=ctop \
  --volume /var/run/docker.sock:/var/run/docker.sock:ro \
  ghcr.io/ccarpinteri/ctop:latest
```

## Building

Build steps can be found [here][build].

## Usage

`ctop` requires no arguments and uses Docker host variables by default. See [connectors][connectors] for further configuration options.

### Config file

While running, use `S` to save the current filters, sort field, and other options to a default config path (`~/.config/ctop/config` on XDG systems, else `~/.ctop`).

Config file values will be loaded and applied the next time `ctop` is started.

### Options

Option | Description
--- | ---
`-a`	| show active containers only
`--context <string>` | use named Docker context
`-f <string>` | set an initial filter string
`-h`	| display help dialog
`--host <string>` | Docker daemon endpoint (e.g., tcp://192.168.1.100:2376)
`-i`  | invert default colors
`-r`	| reverse container sort order
`-s`  | select initial container sort field
`-v`	| output version information and exit

### Keybindings

|           Key            | Action                                                     |
| :----------------------: | ---------------------------------------------------------- |
| <kbd>&lt;ENTER&gt;</kbd> | Open container menu                                        |
|       <kbd>a</kbd>       | Toggle display of all (running and non-running) containers |
|       <kbd>f</kbd>       | Filter displayed containers (`esc` to clear when open)     |
|       <kbd>H</kbd>       | Toggle ctop header                                         |
|       <kbd>h</kbd>       | Open help dialog                                           |
|       <kbd>s</kbd>       | Select container sort field                                |
|       <kbd>r</kbd>       | Reverse container sort order                               |
|       <kbd>o</kbd>       | Open single view                                           |
|       <kbd>l</kbd>       | View container logs (`t` to toggle timestamp when open)    |
|       <kbd>e</kbd>       | Exec Shell                                                 |
|       <kbd>c</kbd>       | Configure columns                                          |
|       <kbd>S</kbd>       | Save current configuration to file                         |
|       <kbd>q</kbd>       | Quit ctop                                                  |

[build]: _docs/build.md
[connectors]: _docs/connectors.md
[single_view]: _docs/single.md
[release]: https://img.shields.io/github/release/lordoverlord/ctop.svg "ctop"

## Alternatives

See [Awesome Docker list](https://github.com/veggiemonk/awesome-docker/blob/master/README.md#terminal) for similar tools to work with Docker.
