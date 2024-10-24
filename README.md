
# Project README

This project is a Go-based application that interacts with Taproot Assets. It includes various components and tests to ensure the correct functionality of the wallet service and other features.

## Prerequisites

To set up and run this project, you need the following:

1. Go Environment: Ensure you have Go `1.22+` installed on your system. The project is written in Go and requires a working Go environment.

2. Taproot Assets Daemon: Clone and run the Taproot Assets Daemon from the following fork:

-   [habibitcoin/taproot-assets](https://github.com/habibitcoin/taproot-assets/tree/tajfi-fork)
-  This fork has custom logic in place that allows signatures to be communicated between `tapd` and `tajfi-server`

3. Configuration: Set the TaprootSigsDir to match the directory where the  tapd binary is run from. This ensures that the application can correctly locate and interact with the Taproot Assets Daemon.

## Setup Instructions

1.  Clone the Repository: Clone this repository to your local machine.

2. Install Dependencies: Navigate to the project directory and run:

	`go  mod  tidy`

	This command will install all necessary dependencies.

3.  Configure Environment: Ensure that the  TaprootSigsDir environment variable is set to the directory where the  tapd binary is located.
4. Begin the backend server by running `go run cmd/server.go`. Ensure the `tapd` daemon is also already running.