tar --owner="arx" --group="arx" \
    --transform 's|app|dicevault|' \
    -cvf dicevault.tar app admin-bot.js