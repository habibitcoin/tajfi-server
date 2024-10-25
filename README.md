
# Project README

This project is a Go-based application that interacts with Taproot Assets to allow operating a pocket universe. It includes various components and tests to ensure the correct functionality of the wallet service and other features. A demo server is available at [https://demo.tajfi.com](https://demo.tajfi.com) for testing purposes and is configured to automatically fund receive invoices with an amount of 10 units.

It is intended to be run in tandem with the frontend web-app, [tajfi-web](https://github.com/topether21/tajfi-web), which enables users to authenticate with their Nostr pubkey.

This allows Tajfi users to have complete custody over their Taproot Assets virtual UTXOs, while the Tajfi-server handles control of the underlying Bitcoin UTXOs that the vUTXOs are anchored to. This leads to greater privacy than utilizing an onchain ERC20 token, and greater custody than traditional centralized Fintech apps of today.

In the future, the Musig2 will be implemented to allow the Bitcoin UTXOs to be 2-of-2, removing any ability for the server to accidentally burn any Taproot Assets.

This opens up opportunities for a pocket universe operator to provide additional value-added services, such as Taproot Asset channel management, and PSBT marketplaces.

## Prerequisites

To set up and run this project, you need the following:

1. Go Environment: Ensure you have Go `1.22+` installed on your system. The project is written in Go and requires a working Go environment.

2. Taproot Assets Daemon: Clone and run the Taproot Assets Daemon from the following fork:

-   [habibitcoin/taproot-assets](https://github.com/habibitcoin/taproot-assets/tree/tajfi-fork)
-  This fork has custom logic in place that allows signatures to be communicated between `tapd` and `tajfi-server`

3. Configuration: Set the TaprootSigsDir to match the directory where the  tapd binary is run from. This ensures that the application can correctly locate and interact with the Taproot Assets Daemon.

- Optionally configure `DemoMode` to true and configure an external `DemoTapdNode` to fund all receive invoices equal to `DemoAmount`.

## Setup Instructions

1.  Clone the Repository: Clone this repository to your local machine.

2. Install Dependencies: Navigate to the project directory and run:

	`go  mod  tidy`

	This command will install all necessary dependencies.

3.  Configure Environment: Ensure that the  TaprootSigsDir environment variable is set to the directory where the  tapd binary is located.
4. Begin the backend server by running `go run cmd/server.go`. Ensure the `tapd` daemon is also already running.
