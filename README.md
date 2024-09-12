# EmuWatcher
<div align="center">
  <img src="assets/Icon.png" alt="Icon" width="500"/>
</div>

**EmuWatcher** is a powerful tool designed to monitor and manage the running state of games on specific emulators. It ensures continuous operation by automatically restarting games if they stop running, providing a seamless gaming experience.

[![Discord](https://img.shields.io/discord/1273553809724932128?color=7289DA&logo=discord&logoColor=white)](https://discord.gg/QfpGHB87jK)

## Key Features

- 🎮 **Advanced Roblox Monitoring**
  - Real-time status tracking
  - Crash detection and recovery
  - Loading screen hang prevention
  - In-game user activity monitoring via Lua scripts
- 🔐 **Seamless Emulator Integration**
  - Automatic login to running emulators
  - Push autoexec scripts to emulators
- 📊 **Enhanced Management**
  - Auto-arrange emulator windows
  - Command-line interface (CLI) for advanced control
  - Automatic emulator startup
  - Launch on system startup
- 📡 **System Monitoring**
  - Send system information and screenshots to webhook
- 🔄 **Upcoming: Emulator Cloning** (with identical HWID)

## Quick Start

1. Download the latest release from the [Releases](https://github.com/Irilith/EmuWatcher/releases) page.
2. Run `EmuWatcher.exe` to start the application.
3. For advanced usage, run `EmuWatcher.exe -h` in the command line.

For detailed instructions, please refer to our [Usage Guide](https://github.com/Irilith/EmuWatcher/blob/main/USAGE.md).

## Disclaimer

> [!IMPORTANT]
> **EmuWatcher** is provided "as-is" without any warranties or guarantees. Use at your own risk.

- **No Warranty**: We do not guarantee that the program will meet your requirements, operate without interruption, or be error-free.
- **Use at Your Own Risk**: The developer is not responsible for any damage or loss resulting from the use of this software.
- **Compatibility**: EmuWatcher may not be compatible with all emulators or games. Users are responsible for ensuring compatibility with their setup.
- **Support**: Support and updates are provided at the developer's discretion.

### Third-Party Executors

- This program may interact with third-party executors not created by the EmuWatcher developer.
- Users are solely responsible for ensuring their use complies with all relevant laws and regulations.
- The developer assumes no responsibility for misuse, illegal activities, or consequences resulting from the use of this program or associated executors.

## License

EmuWatcher is open-source software licensed under the [GNU GPL-3.0 License](https://opensource.org/license/gpl-3-0). For more details, see the [LICENSE](https://github.com/Irilith/EmuWatcher/blob/main/LICENSE) file.

## Building from Source

To build EmuWatcher from source:

1. **Clone the Repository**
   ```
   git clone https://github.com/Irilith/EmuWatcher.git
   cd EmuWatcher
   ```

2. **Install Dependencies**
   - [Go 1.22](https://golang.org/dl/)
   - [Make](https://www.gnu.org/software/make/)
   - [rsrc](https://github.com/akavel/rsrc) (Install with `go install github.com/akavel/rsrc@latest`)

3. **Build the Program**
   ```bash
   make
   ```

4. **Run EmuWatcher**
   ```bash
   ./EmuWatcher
   ```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more information on how to get involved.

## Support

For questions, feature requests, or bug reports, please [open an issue](https://github.com/Irilith/EmuWatcher/issues) on GitHub or join our [Discord community](https://discord.gg/QfpGHB87jK).
