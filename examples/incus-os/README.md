# Incus OS

Example for installing [Incus OS](https://linuxcontainers.org/incus-os) using an ephemeral Alpine VM to generate a installation ISO.

## How it works

An ephemeral Alpine VM uses [`flasher-tool`](https://github.com/lxc/incus-os/tree/main/incus-osd/cmd/flasher-tool) to produce a customised Incus OS ISO that has seed data baked in. The seed data includes the host's client certificate so the resulting Incus OS installation automatically trusts the host.

The generated ISO is written to the host's `/tmp` directory and is then used to perform the Incus OS installation.

## Usage

```sh
incus-apply incus-os.yaml
```

## Examples

| File            | Description                                                     |
| --------------- | --------------------------------------------------------------- |
| `incus-os.yaml` | Downloads installation ISO and begins installation for Incus OS |
