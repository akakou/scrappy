import (builtins.fetchTarball {
  url = "https://github.com/NixOS/nixpkgs/archive/dc5d91f840324650bac8c379428c7037a416959a.tar.gz";
  sha256 = "1v403gscmmsaq6v3rrbwin62bvjps737klmpg9v868cwmwkqd9am";
}) {
  config = {
    allowUnfree = true;
    android_sdk.accept_license = true;
  };
  overlays = [];
}
