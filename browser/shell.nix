{ pkgs ? import ../nix/pkgs.nix }:

let
  base = import ../shell.nix { inherit pkgs; };
in
pkgs.mkShell {
  inputsFrom = [ base ];
  packages = with pkgs; [ chromium jq ];
  inherit (base) CGO_ENABLED GOTOOLCHAIN GOFLAGS TPM2TOOLS_TCTI;

  shellHook = ''
    scrappy-browser() (
      bash "$SCRAPPY_ROOT/browser/run.sh" "$@"
    )
  '';
}
