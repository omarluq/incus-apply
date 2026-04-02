# Configuration Reference

This page contains the full field reference for `incus-apply` resource documents.

## Common Fields

| Field         | Type   | Description                                    |
| ------------- | ------ | ---------------------------------------------- |
| `type`        | string | **Required.** Resource type                    |
| `name`        | string | **Required.** Resource name                    |
| `project`     | string | Incus project (overridden by `--project` flag) |
| `config`      | map    | Resource configuration options                 |
| `devices`     | map    | Device configurations                          |
| `description` | string | Resource description                           |

## Instance Fields

| Field      | Type   | Description                                   |
| ---------- | ------ | --------------------------------------------- |
| `image`    | string | Image to use (for example `images:debian/12`) |
| `vm`       | bool   | Create a VM instead of a container            |
| `empty`    | bool   | Create an empty instance                      |
| `profiles` | list   | Profiles to apply                             |
| `storage`  | string | Storage pool for the root disk                |
| `network`  | string | Network to attach                             |
| `target`   | string | Cluster member target                         |

## Storage Pool Fields

| Field    | Type   | Description                                 |
| -------- | ------ | ------------------------------------------- |
| `driver` | string | Storage driver (dir, zfs, btrfs, lvm, ceph) |
| `source` | string | Source path or device                       |

## Storage Volume And Bucket Fields

| Field  | Type   | Description                     |
| ------ | ------ | ------------------------------- |
| `pool` | string | **Required.** Storage pool name |

## Network Fields

| Field         | Type   | Description                                          |
| ------------- | ------ | ---------------------------------------------------- |
| `networkType` | string | Network type (bridge, ovn, macvlan, sriov, physical) |

## Network ACL Fields

| Field     | Type | Description            |
| --------- | ---- | ---------------------- |
| `ingress` | list | Ingress firewall rules |
| `egress`  | list | Egress firewall rules  |

## Variables

Variables are declared with a `type: vars` document and referenced from resource documents with `$VAR` or `${VAR}`.

### Example

```yaml
---
type: vars
vars:
  DB_NAME: myapp
  DB_USER: appuser
  DB_PASS: ${MYSQL_PASSWORD}
files:
  - secrets.env
---
type: instance
name: db
image: docker:mysql
config:
  environment.MYSQL_DATABASE: $DB_NAME
  environment.MYSQL_USER: $DB_USER
  environment.MYSQL_PASSWORD: $DB_PASS
```

### Scoping

- Variables are file-scoped by default.
- Use `global: true` in a `vars` document to share variables across files.
- File-scoped variables override global variables with the same name.

### Shell Environment

- Shell environment variables can be referenced only inside the `vars` document.
- Resource documents expand only variables declared through `type: vars`.
- If a referenced variable is not declared in `type: vars`, it is left unchanged in the resource document.

Example:

```yaml
---
type: vars
vars:
  APP_NAME: myapp
---
type: instance
name: web
image: images:debian/12
config:
  environment.APP_NAME: $APP_NAME
  environment.HOME_DIR: $HOME
```

In this example, `environment.APP_NAME` becomes `myapp`, while `environment.HOME_DIR` remains `$HOME` because `HOME` was not declared in `type: vars`.

### Preview Redaction

- Instance preview output redacts `config.environment.*` values by default.
- This affects preview rendering only. Interpolation and apply behavior still use the actual values.
- Use `--show-env` when you want preview output to show the real environment variable values.

### Syntax

- `$VAR` uses the value of `VAR`.
- `${VAR}` uses the value of `VAR`.
- `${VAR:-default}` uses `default` when `VAR` is unset or empty.
- `$$` produces a literal `$`.

## Notes

- Configuration files may contain multiple YAML documents separated by `---`.
- Variables are declared with `type: vars` and are not resource documents.
- See [../README.md](../README.md) for quick start and common usage.
