FROM debian:trixie-slim

RUN apt-get update \
&&  echo 'Installing packages...' \
&&  apt-get install -y \
        git golang \
        shared-mime-info ca-certificates \
        curl libcap2-bin imagemagick jq procps \
        less mc vim ne \
        pipx \
&&  echo 'Installing mitmproxy...' \
&&  pipx install --global mitmproxy \
&&  echo 'Adding ordinary user...' \
&&  groupadd -r cnbotgroup \
&&  useradd -r -g cnbotgroup cnbot \
&&  echo 'Adding directories...' \
&&  mkdir \
        /app \
        /app/src \
        /app/logs \
&&  echo 'Building cnbot...' \
&&  cd /app/src \
&&  git clone --depth 1 --branch master https://github.com/michurin/cnbot.git \
&&  cd cnbot \
&&  go mod download \
&&  go mod verify \
&&  go build ./cmd/... \
&&  mv cnbot /app \
&&  cp \
    demo/bot.sh \
    demo/bot_long.sh \
    demo/debugging_wrapper.sh \
    /app \
&&  cd /app \
&&  ln -s debugging_wrapper.sh bot_debug.sh \
&&  ln -s debugging_wrapper.sh bot_long_debug.sh \
&&  chown cnbot:cnbotgroup \
    logs \
    bot.sh \
    bot_long.sh \
    debugging_wrapper.sh \
&&  chown -h cnbot:cnbotgroup \
    bot_debug.sh \
    bot_long_debug.sh \
&&  echo 'RUN successfully finished.'

USER cnbot
WORKDIR /app

ENV TB_TOKEN='' \
    TB_CTRL_ADDR=:9999 \
    TB_SCRIPT=/app/bot.sh \
    TB_LONG_RUNNING_SCRIPT=/app/bot_long.sh

CMD ["/app/cnbot"]
