# Troubleshooting

## Script not found: no such file or directory

If you receive error `no such file or directory`,
please check the `"script"` attribute in
your configuration file. You can
specify an absolute path, or relative path
related to directory of configuration file.

## Command not found

If your script use subprocesses, it is possible
it will unable to find corresponding binaries.

It is because `cnbot` doesn't pass its
environment to subprocesses.
Long story short, you can specify full
absolute paths or export `PATH` environment variable.

Why does `cnbot` behave is that way?
Where are two crucial reasons:

- Repeatable behavior. You script has to work in the same environment on your laptop,
on the dedicated server, from  ordinary user and from root.
- Security reasons. You have to specify `PATH`
explicitly in your script.

Long story. If you use some kind of `shell`
for scripting you may be surprised to know
that `PATH` is a magic variable. It exists
even if it doesn't. At least it is true for
[bash](http://git.savannah.gnu.org/cgit/bash.git/tree/variables.c#n486)
([DEFAULT_PATH_VALUE](http://git.savannah.gnu.org/cgit/bash.git/tree/config-top.h#n66)
defined on compilation time) and
[zsh](https://github.com/zsh-users/zsh/blob/master/Src/init.c#L978).

You may want to set variables like
`LD_PRELOAD`, `LANG`, `LC_*`,
`POSIXLY_CORRECT`, `PYTHONPATH`, `PERL5LIB`
along with `PATH`.

## HTTP port is already in use

If you receive log message like that
```
listen tcp :9090: bind: address already in use
```
it means what it means. If you don't need Web server
for asynchronous messages, just remove `server` section
from config. Otherwise, change ports to eliminate
conflict.

## Two bot drivers at the same time

It is impossible to run more than one bot driver at the
same time. Otherwise, you
will receive `HTTP 409 Conflict` error with payload like this:
```json
{"ok":false,"error_code":409,"description":"Conflict: terminated by other getUpdates request; make sure that only one bot instance is running"}
```
You are to have only one bot driver.

## Web hook

If your bot has been drowned by another tool/library/SDK and
this bot starts to receive `HTTP 409 Conflict` error with
following payload:
```
{"ok":false,"error_code":409,"description":"Conflict: can't use getUpdates method while webhook is active; use deleteWebhook to delete the webhook first"}
```
It means you set up web hook. You can check it by command like that: 
```sh
curl 'https://api.telegram.org/bot$TOKEN/getWebhookInfo'
```
```json
{"ok":true,"result":{"url":"https://your.url/","has_custom_certificate":false,"pending_update_count":0,"max_connections":40,"allowed_updates":["message"]}}
```
You can easily remove web hook manually:
```sh
curl 'https://api.telegram.org/bot$TOKEN/deleteWebhook'
```
```json
{"ok":true,"result":true,"description":"Webhook was deleted"}
```

## Telegram API is unreachable (use proxy)

Use `HTTPS_PROXY` environment
```sh
HTTPS_PROXY=socks5://localhost:8888 ./cnbot -c config.yaml
```
