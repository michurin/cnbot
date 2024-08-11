# System administration topics

## Build

You can specify build version

```sh
go build -ldflags "-s -w -X github.com/michurin/cnbot/pkg/bot.Build=`date +%F`-`git rev-parse --short HEAD`" ./cmd/...
```

## Monitoring and health checking

You can turn on alive handler, mentioning it in configuration file like that:

```yaml
bot: ...
alive:
  bind_address: ":9999"
```

After that you can fetch current state of process from this handler: version, start time, memory utilisation details, current number of goroutines...
All this information is available in json format, so you can use tools like `jq` to manage it. For example, it is easy to make input string for
simple RRDtools based monitoring:

```sh
curl -s localhost:9999 | jq -r '"N:\(.num_goroutine):\(.mem_status.Alloc):\(.mem_status.HeapObjects)"'
```

It will return string like `N:16:1365296:10822` suitable for `rrdtool update`.

## Startup: modern `systemd` fashion

As usual create service file like this `/etc/systemd/system/cnbot.service`:

```ini
[Unit]
Description=Telegram bot (cnbot) service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=1
User=nobody
ExecStart=/usr/bin/cnbot -c /etc/cnbot-config.yaml

[Install]
WantedBy=multi-user.target
```

Don't forget to set your actual username after `User=` and the proper path to your script
and configuration file in `ExecStart=`. Keep in mind that relative paths in configuration file
evaluate with respect of configuration file directory.

```
[root]# systemctl daemon-reload
[root]# systemctl status cnbot
● cnbot.service - Telegram bot (cnbot) service
     Loaded: loaded (/etc/systemd/system/cnbot.service; disabled; vendor preset: disabled)
     Active: inactive (dead)
[root]# systemctl start cnbot
[root]# systemctl status cnbot # you will see that bot is active
```

If everything is ok, you can enable service (`systemd enable`).
If something goes wrong, you can inspect logs:

```sh
journalctl -u cnbot
```

You can use `-f` with `journalctl` to read live-tail logs.

## Startup: old fashioned `rc.d` scripts

### Daemonize

`cnbot` doesn't have any daemonization abilities.
If you really wish to daemonize it, you can use tools like `nohup`.
However, if you use `sistemd` you don't need for daemonization.

### Log management

`cnbot` just throws log messages to `stdout`.
There is a lot of tools designed with the primary purpose
of maintaining size-capped, automatically rotated, log file sets.
Personally, I prefer `multilog` from [daemontools](http://cr.yp.to/daemontools.html).

Further reading: [Don't use logrotate or newsyslog in this century](http://jdebp.eu/FGA/do-not-use-logrotate.html).
