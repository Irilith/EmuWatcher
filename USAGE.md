# EmuWatcher Usage Guide

## Downloading and Installing EmuWatcher

1. **Access the Releases Page**: Navigate to the [EmuWatcher Releases](https://github.com/Irilith/EmuWatcher/releases) page on GitHub.
2. **Download the Latest Version**: Locate and download the latest `Updater.exe` file.
3. **Installation**: Run the `Updater.exe` file to install EmuWatcher in your preferred location.

## Essential Prerequisites

EmuWatcher requires two key components:

1. **ADB (Android Debug Bridge)**
2. **Tesseract OCR**

### Setting up ADB

1. **Obtain ADB**: Download ADB from the [Android developer website](https://developer.android.com/studio/releases/platform-tools).
2. **Installation**:
    - In your EmuWatcher directory, create a `tools` folder.
    - Within `tools`, create an `adb` subfolder.
    - Place the `adb.exe` file in the `tools/adb/` directory.

    Your folder structure should resemble:
    ```
    EmuWatcher/
    └── tools/
        └── adb/
            └── adb.exe
    ```

### Installing Tesseract OCR

1. **Download Tesseract**: Get Tesseract OCR from the [official repository](https://github.com/tesseract-ocr/tesseract).
2. **Setup**:
    - In your EmuWatcher directory, create a `tools/ocr` folder.
    - Place `tesseract.exe` in `tools/ocr/`.
    - Download [vie.traineddata](https://github.com/tesseract-ocr/tessdata/blob/main/vie.traineddata).
    - Create a `tessdata` folder within `tools/ocr/` and place `vie.traineddata` there.

    Your folder structure should look like:
    ```
    EmuWatcher/
    └── tools/
        └── ocr/
            ├── tesseract.exe
            └── tessdata/
                └── vie.traineddata
    ```

## Configuring EmuWatcher

Ensure EmuWatcher is correctly configured:

- Set **ADB Path** to: `tools/adb/adb.exe`
- Set **Tesseract OCR Path** to: `tools/ocr/tesseract.exe`

Verify your complete folder structure matches:

    EmuWatcher/
    ├── EmuWatcher.exe
    ├── config.json
    ├── data/
    │   ├── cookies.txt
    │   └── autoexec/
    │       └── EmuWatcher.lua
    ├── assets/
    │   └── datasets/
    │       ├── Cookies
    │       └── appStorage.json
    └── tools/
        ├── ocr/
        │   ├── tesseract.exe
        │   └── tessdata/
        │       └── vie.traineddata
        └── adb/
            └── adb.exe

## Launching and Using EmuWatcher

1. **Open Command Prompt**:
    - Press `Win + R`, type `cmd`, and press Enter.
    - Navigate to your EmuWatcher directory:
    ```
    cd path\to\your\EmuWatcher
    ```

2. **Start EmuWatcher**:
    - Type `EmuWatcher.exe` and press Enter.

3. **Configure and Run**:
    - When the menu appears, press `2` to adjust your settings.
    - After configuration, press `1` to begin monitoring.

## Troubleshooting Tips

- **Path Issues**: Double-check that ADB and Tesseract OCR paths are correctly set in your configuration.
- **Dependency Problems**: Ensure ADB and Tesseract OCR are properly installed and functioning.
- **Permission Errors**: Run EmuWatcher as an administrator if you encounter permission issues.
- **ADB Connection**: Make sure your emulator or device is properly connected and recognized by ADB.
- **Didn't Work**: Make sure Root and ADB debug are enabled in your emulator's settings.


If problems persist, consult the EmuWatcher documentation or seek help in the project's GitHub issues section.
