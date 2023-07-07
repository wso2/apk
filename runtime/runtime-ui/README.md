# APK Config Language Support

A Visual Studio Code extension for providing language support for APK configuration YAML files.

## Features
- Syntax highlighting for APK configuration YAML files.
![screencast](https://raw.githubusercontent.com/wso2/apk/main/runtime/runtime-ui/images/demo1.gif)

- Auto-completion for APK configuration YAML properties and values.
- Validation and error checking for APK configuration YAML files.
- Create new apk configuration files with provided templates.
![screencast](https://raw.githubusercontent.com/wso2/apk/main/runtime/runtime-ui/images/demo2.gif)


## Requirements

- Visual Studio Code version 1.63.0 or newer.

## Installation

1. Launch Visual Studio Code.
2. Go to the Extensions view by clicking on the square icon in the sidebar or pressing `Ctrl+Shift+X`.
3. Search for "APK Config Language Support".
4. Click on the "Install" button for the extension published by "APK Config Language Support".
5. The extension will be installed and activated automatically. ( This extension is depending on the YAML Language Support by Red Hat extension). If it's not installed, this will install it automatically.

## Usage

1. Open a APK Configuration file (.apk-conf).
2. The extension will automatically provide syntax highlighting and code completion for APK configuration properties and values.
3. Errors and warnings will be displayed if there are any issues with the YAML structure or values.
4. Use the provided code snippets to quickly insert common APK configuration patterns.
5. Create new APK Configuration will provide you with set of templates to get started.

## Contributing

Contributions are welcome! Here's how you can get involved:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your changes to your fork.
5. Submit a pull request describing your changes.

### Extension Packaging Guide

This guide provides instructions on how to package the extension.

#### Prerequisites

- Node.js 16+ (https://nodejs.org)
- Yarn (https://yarnpkg.com)

#### Installation

1. Install Node.js 16+ from the official website.
2. Install Yarn globally by following the instructions on the official Yarn website.

#### Usage

1. Open your terminal or command prompt and navigate to the project's root directory.
2. Run the following command to install the project dependencies:

    ```shell
    yarn
    ```
3. Once the dependencies are installed, run the following command to build the extension:

    ``` shell
    yarn run build
    ```
    This command triggers the build process, which compiles the source code, optimizes assets, and generates the final packaged extension.

    After running the build command, you should find the packaged extension files in the designated build output directory.


## License

This extension is licensed under the [Apache 2.0 License](LICENSE).
