# PortForwarder
 **PortForwarder is a lightweight, single-purpose utility for forwarding local ports to destination ports without installation.**

## Developer Guide (GUI)
> For Windows users, it will only support **Win10+** operating systems.

### Prerequisites
* GO: 1.22+ (For building the Go-based backend)
* NODEJS: 20+ (For managing the Node.js-based scripts and dependencies)

### Installing Dependencies
Navigate to the **[gui](./gui)** directory and run the following command to install all dependencies:
```
npm install
```

### Scripts
```
# Start the backend application
npm run start:app

# Start the frontend application
npm run start:gui

# Build the releasable application
npm run build
```

## Developer Guide (CLI)
> TODO: develop a command-line client for the convenience of more users who can utilize it on the server.

## License
PortForwarder is licensed under the Apache License, Version 2.0. See the [LICENSE](./LICENSE) file for more details.

## Contributions
Contributions to PortForwarder are welcome! Please feel free to submit pull requests, report issues, or provide feedback. We appreciate any help you can offer to improve this tool.

Thank you for considering PortForwarder!
