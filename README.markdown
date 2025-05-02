# Go Blockchain Project

This project implements a simple blockchain system with a blockchain server and a wallet server, written in Go. The
blockchain server manages the blockchain and transactions, while the wallet server provides a web interface for users to
interact with the blockchain, view their wallet details, check their balance, and send transactions.

## Project Structure

The project is organized into the following directories:

```
goblockchain/
├── block/
│   ├── block.go           # Defines the block structure and methods
│   └── marshal.go         # Handles JSON marshaling for blocks
├── blockchain_server/
│   ├── blockchain_server.go  # Blockchain server implementation
│   └── main.go            # Entry point for the blockchain server
├── cmd/
│   └── main.go            # Main entry point for CLI commands (if any)
├── utils/
│   ├── ecdsa.go           # ECDSA utility functions for cryptography
│   ├── json.go            # JSON utility functions
│   └── neighbor.go        # Neighbor node handling (for consensus)
├── wallet/
│   ├── wallet.go          # Wallet implementation (key generation, etc.)
├── wallet_server/
│    ├── static/        # Static files for the wallet server
│    │   └── wallet.js  # JavaScript logic for the wallet interface
│    ├── templates/     # HTML templates for the wallet server
│    │   └── index.html # Main HTML page for the wallet interface
│    ├── main.go    # Entry point for the wallet server (if separate)
│    └── wallet_server.go  # Wallet server implementation
├── .gitignore            # Git ignore file
└── go.mod                # Go module file for dependency management
```

## Prerequisites

To run this project, you need the following:

- **Go**: Version 1.16 or higher (recommended: latest stable version).
- **Operating System**: Tested blockchain server on WSL (Ubuntu 18.04) with GoLand IDE run configuration, but should
  work on any OS supported by Go (Windows, macOS, Linux).
- **Internet Connection**: Required to fetch Go dependencies and Tailwind CSS CDN (used in `index.html`).

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/vantri1010/goblockchain.git
   cd goblockchain
   ```

2. **Install Dependencies**:
   Ensure you have Go installed, then fetch the required dependencies:
   ```bash
   go mod tidy
   ```
   This will download dependencies specified in `go.mod`, such as standard Go libraries for HTTP, JSON, and
   cryptography.

## Running the Project

The project consists of two main components: the blockchain server and the wallet server. You need to run both servers
to use the application.

### 1. Run the Blockchain Server

The blockchain server manages the blockchain, transactions, mining, and consensus. It exposes APIs like `/wallet`,
`/transactions`, `/mine`, `/amount`, and `/consensus`.

#### Command:

```bash
cd blockchain_server
go run . -port 5000
```

#### Configuration:

- **Port**: `5000` (default). You can change the port using the `-port` flag.
- **Example Configurations to run concurrently**:
    - Port 5000: `go run . -port 5000`
    - Port 5001: `go run . -port 5001`
    - Port 5002: `go run . -port 5002`

#### Notes:

- The blockchain server must be running before the wallet server, as the wallet server depends on its APIs (e.g.,
  `GET /wallet`, `GET /amount`).
- If you use goland ide you can also run it on WSL environment. The program can also find the ips of neighboring nodes
  on both windows and wsl environments.

### 2. Run the Wallet Server

The wallet server provides a web interface to interact with the blockchain. It communicates with the blockchain server
via a gateway (e.g., `http://localhost:5000`).

#### Command:

```bash
cd wallet_server
go run . -port 8080 -gateway http://localhost:5000
```

#### Configuration:

- **Port**: `8080` (default). You can change the port using the `-port` flag.
- **Gateway**: `http://localhost:5000` (default). This should match the port of the blockchain server.
- **Example Configurations to run concurrently**:
    - Port 8080, Gateway 5000: `go run . -port 8080 -gateway http://localhost:5000`
    - Port 8081, Gateway 5001: `go run . -port 8081 -gateway http://localhost:5001`
    - Port 8082, Gateway 5002: `go run . -port 8082 -gateway http://localhost:5002`

#### Notes:

- The wallet server serves static files (`wallet.js`) from `wallet_server/static/` and templates (`index.html`) from
  `wallet_server/templates/`.
- Ensure the blockchain server is running on the specified gateway port before starting the wallet server.

## Accessing the Application

1. **Open the Wallet Interface**:
    - After starting blockchain servers (port 5000, 5001, 5002), open your browser and navigate to:
      ```
      http://localhost:8080
      http://localhost:8081
      http://localhost:8082
      ```
    - This will load the wallet interface (`index.html`).

2. **Features**:
    - **Wallet Details**: The interface automatically loads the miner's wallet details (`public_key`, `private_key`,
      `blockchain_address`) by calling `GET /wallet`.
    - **Balance**: The balance is updated every second by calling `GET /wallet/amount` with the current
      `blockchain_address`.
    - **Send Money**: Enter a recipient blockchain address and amount, then click "Send" to create a transaction via
      `POST /transaction`.

3. **Example Workflow**:
    - Start 3 blockchain servers on port 5000, 5001, 5002 on WSL
    - Start the wallet server on port 8080 with gateway `http://localhost:5000`.
    - Start the wallet server on port 8081 with gateway `http://localhost:5001`.
   - Open `http://localhost:8080` and `http://localhost:8081` in your browser.
   - View your wallet details, check your balance, copy current wallet address and try to send money from another
     wallet.
   - You can also view transaction pool after a transaction is created and before block mining via
     `http://localhost:<500x>/transactions`
    - You can also view transactions history including mining and deposit transactions via
      `http://localhost:<500x>/chain`

## Testing the Application

1. **Verify Wallet Details**:
    - When the page loads, ensure the `public_key`, `private_key`, and `blockchain_address` fields are populated.
    - Compare these values with the logs from the blockchain server (printed during initialization).
    - Refresh the page multiple times to confirm the same wallet is returned (not a new one each time).

2. **Check Balance Updates**:
    - Ensure the balance (displayed in `#wallet_amount`) updates every second.
    - Edit the `blockchain_address` field and confirm the balance updates accordingly.

3. **Send a Transaction**:
    - Enter a valid `recipient_blockchain_address` (e.g., from another wallet instance) and a `send_amount`.
    - Click "Send" and confirm:
        - A confirmation dialog (`Are you sure to send?`) appears.
        - A success (`Send success`) or failure (`Send fail`) message is displayed.
    - Check the wallet server logs for the `POST /transaction` request and the blockchain server logs for the
      corresponding transaction processing.

4. **Inspect Network Requests**:
    - Open the browser's DevTools (F12) and check the Network tab:
        - `GET /wallet`: Should return the wallet JSON.
        - `GET /wallet/amount`: Should return the balance JSON.
        - `POST /transaction`: Should return a success/failure message.
    - Confirm `wallet.js` is loaded from `http://localhost:8080/static/wallet.js`.

5. **Check Responsiveness**:
    - Test the interface on different screen sizes (e.g., mobile view in DevTools) to ensure Tailwind CSS responsive
      classes work correctly.

## Troubleshooting

- **Wallet Details Not Loaded**:
    - Check the blockchain server logs to see if the wallet information is printed out and to ensure `GET /wallet` is
      being called.
    - Verify the blockchain server is running on the correct port (e.g., 5000) and the wallet server's gateway matches.
    - Check the browser console for JavaScript errors (`Failed to fetch`, etc.).

- **Balance Not Updating**:
    - Ensure `GET /wallet/amount` requests are sent with the correct `blockchain_address`.
    - Check wallet server logs for errors in handling `GET /wallet/amount`.
  - Check IDE console log to see if neighbor IPs are logged (found) for example :
    `Find neighbors:  [172.16.0.2:5000 172.16.0.2:5002]`
  - Try `http://localhost:<500x>/consensus` to update the transaction history of all nodes

- **Transaction Fails**:
    - Verify the `recipient_blockchain_address` and `send_amount` are valid.
    - Check the browser console and server logs for errors in `POST /transaction`.

## Additional Notes

- **Technology Stack**:
    - Backend: Go (net/http for servers, custom blockchain implementation).
    - Frontend: HTML with Tailwind CSS for styling, vanilla JavaScript (`wallet.js`) for interactivity.

- **Scalability**:
    - The blockchain server supports consensus (`/consensus` endpoint) for distributed nodes.
    - To scale, run multiple blockchain servers on different ports (e.g., 5000, 5001, 5002) and wallet servers with
      corresponding gateways.
## License

This project is for educational purposes. Feel free to modify and distribute as needed.