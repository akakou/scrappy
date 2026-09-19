#!/usr/bin/env bash
set -euo pipefail

scrappy_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)
profile="$scrappy_root/.cache/chromium"
host="$scrappy_root/.cache/bin/scrappy-host"

# The ID of an unpacked extension depends on its location. Register the ID
# displayed by chrome://extensions, rather than assuming a container path.
if [[ $# -ne 1 || ! $1 =~ ^[a-p]{32}$ ]]; then
    echo "Usage: scrappy-browser EXTENSION_ID (from chrome://extensions)" >&2
    echo "First load $scrappy_root/browser/extension in:" >&2
    printf 'chromium --user-data-dir=%q chrome://extensions\n' "$profile" >&2
    exit 2
fi

mkdir -p "$profile/NativeMessagingHosts" "$(dirname -- "$host")"
(
    cd "$scrappy_root/browser/app"
    go build -o "$host" .
)

# Chromium starts native hosts in their executable's directory. The wrapper
# supplies the same working directory as the server for configs and databases.
wrapper="$scrappy_root/.cache/bin/scrappy-native"
{
    printf '#!%s\n' "$(command -v bash)"
    printf 'cd -- %q || exit\n' "$scrappy_root/example"
    printf 'exec %q "$@"\n' "$host"
} > "$wrapper"
chmod +x "$wrapper"

jq --arg path "$wrapper" --arg origin "chrome-extension://$1/" \
    '.path = $path | .allowed_origins = [$origin]' \
    "$scrappy_root/browser/app/scrappy.json" \
    > "$profile/NativeMessagingHosts/com.akakou.scrappy.json"

exec chromium --user-data-dir="$profile" http://localhost:8081/
