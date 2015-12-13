# TODO

- Unit tests
- Write documentation

# Features

- Deploy a golang app (go binary with static files) easily with ssh and rsync

# Usage

## Example config
This is an example for a go app (in subfolder `backed` with a JSPM-based frontend in subfolder `webui`). Lets assume you saved this file to `/path/to/simple-cd/go_app1.yml`

```
ssh_private_key_path: /root/.ssh/id_rsa
remote_user: usrdeploy
remote_host: "123.456.789.000"
remote_port: 22
base_remote_dir: /var/subdomains/go_app1
static_dirs:
    - backend/sql_migrations/
    - webui/dist/
    - webui/jspm_packages/
static_files:
    - webui/config.js
    - webui/favicon.ico
    - webui/index.html
    - webui/index.js
exec_files:
    - backend/go_app1_binary
stop_cmd: "/usr/sbin/service go_app1 stop"
start_cmd: "/usr/sbin/service go_app1 start"
```

## Example command-line call with the above YAML file.

`simple-continuous-deployment "/path/to/simple-cd/go_app1.yml" "/base/local/directory"`
