<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/branding/fabrico_new.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/branding/fabrico_light.svg">
  <img alt="Fabrico logo" src="https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/branding/fabrico.png">
</picture>

Fabrico Ledger is a project for prototyping a byzantine-fault tolerant distributed ledger system for secure and decentralized distribution of licenses and production data for computer-aided manufacturing systems, such as additive manufacturing / 3D printing and CNC machines.

An integrated digital twin simulating G-Code manufacturing provides the ability to control tolerances and provide quality assurance throughout the complete process.

Interfaces and functions will be developed for three parties, running a full node each:

- Original manufacturer (data provider)
- Data logistics (providing large file storage nodes)
- Contracted manufacturers

Project work for DHBW Lörrach.

## Used key technologies

- SmartBFT Consensus Algorithm / Library
- gRPC (with TLS) & Protobuf
- State of the art encryption / signatures
  - EdDSA with Ed25519
  - X25519
  - ECIES with ChaCha20-Poly1305
- mDNS & DNS-SD for automatic Node Discovery

## User Interface

The user interface for interacting with the distributed ledger is written in HTML5 / JavaScript using the Bootstrap Framework with Chart.js for future visualization of system status.

**System information**

![System information](https://github.com/fabian-z/fabrico-ledger/raw/main/res/screenshots/system-information.png)

**Upload source data**

![Upload](https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/screenshots/upload-action.png)

**Fabricate data file**

![Fabricate](https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/screenshots/fabricate.png)

**Monitoring example**

![Fabricate](https://raw.githubusercontent.com/fabian-z/fabrico-ledger/main/res/screenshots/monitor-demo.png)

## License

This work is licensed under Apache License.
