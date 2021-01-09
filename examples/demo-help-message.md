*Simple text messaging*
• \/start — just say hello

*Simple command execution*
• \/date — show current date
• \/uname — show `uname -a` output
• \/uptime — execute `uptime`

*Special cases*
• \/noout — what happens if script does not produce output
• \/nothing — how to say to `cnbot` that we don't want to reply

*Looking at scripts input closely*
• \/args — show command line arguments \(try to say `args a b c` for example\)
• \/env — show environment variables \(your can play with it like `env something`\)

*Preformatted text*
• \/cal — show calendar

*More scripting \(markdown; working with raw messages\)*
• \/calc — calculator \(just for demo; do not use it in production\)

*Simple images*
• \/gologo — just throw image to stdout to show it
• \/du — show your disk utilization in pie chart

*Asynchronous interaction*
• \/async — messaging using HTTP API
• \/cap — image with caption

*Inline keyboards*
• \/btn — simple inline keyboard to invoke ordinary commands
• \/update — simple mutable message with inline keyboard
• \/notify — show different kinds of notifications

*This help message \(one more markdown\)*
• \/help

*And besides, the bot accepts*
• contacts and
• forwarded messages
to figure out user/chat/channel ID

*You can also share you location*
• location
• or live location
