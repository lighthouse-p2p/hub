<h1 align="center">Lighthouse Hub</h1>
<p align="center">The DHT and BlockChain for Lighthouse p2p</p>

---

## Building

### Requirements

- Golang installed on your system
- A fully configured and secure postgres server
- A fully configured and secure redis server

### How to build

- `git clone` this repository
- Copy `.env.sample` to `.env`
- Change the values in `.env` **(WARNING: DO NOT MODIFY `.env.sample`)**
- Make a directory named `build`
- - For POSIX systems, run `go build -o build/lighthousehub cmd/lighthousehub/main.go`
  - For Windows systems, in a `PowerShell` window, `go build -o .\build\lighthousehub.exe .\cmd\lighthousehub\main.go`

## Running

- Double check all the values in `.env`
- Run
  - `./cmd/lighthousehub` for POSIX systems, or
  - `.\cmd\lighthousehub.exe` for Windows systems in a `PowerShell` window
- Verify that the server is running on the specified `HTTP_ADDR`. If you see `Cannot GET /`, it means the server is successfully configured and running

---

## How it works

Click [here](https://github.com/lighthouse-p2p/docs) to see the documentation about the protocol and the algorithm used.

---

## License

Lighthouse Hub is licensed under the `AGPL-3.0-or-later` license. You can obtain a copy [here](https://www.gnu.org/licenses/agpl-3.0.html).
