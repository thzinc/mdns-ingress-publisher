# mdns-ingress-publisher

Watches for Kubernetes Ingress object events and publishes mDNS records accordingly

## Quickstart

### Run as needed

Ensure your `.kube/config` is configured to connect to your preferred Kubernetes cluster, then run `mdns-ingress-publisher`.

More information is available by running `mdns-ingress-publisher --help`

### Run as systemd service using microk8s

1. Build the `microk8s-systemd` target
   ```bash
   make microk8s-systemd
   ```
2. Copy the appropriate binary to `/usr/local/bin`
   ```bash
   cp artifacts/linux-amd64/mdns-ingress-publisher /usr/local/bin
   ```
3. Run `install.sh`
   ```bash
   cd artifacts/microk8s-systemd
   sudo ./install.sh
   ```

## Building

```bash
make build
```

## Code of Conduct

We are committed to fostering an open and welcoming environment. Please read our [code of conduct](CODE_OF_CONDUCT.md) before participating in or contributing to this project.

## Contributing

We welcome contributions and collaboration on this project. Please read our [contributor's guide](CONTRIBUTING.md) to understand how best to work with us.

## License and Authors

[![Daniel James logo](https://secure.gravatar.com/avatar/eaeac922b9f3cc9fd18cb9629b9e79f6.png?size=16) Daniel James](https://github.com/thzinc)

[![license](https://img.shields.io/github/license/go-sensors/rpii2c.svg)](https://github.com/go-sensors/rpii2c/blob/master/LICENSE)
[![GitHub contributors](https://img.shields.io/github/contributors/go-sensors/rpii2c.svg)](https://github.com/go-sensors/rpii2c/graphs/contributors)

This software is made available by Daniel James under the MIT license.
