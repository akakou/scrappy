{ pkgs ? import ./nix/pkgs.nix }:

pkgs.mkShell {
  name = "scrappy";
  packages = with pkgs; [ go pkg-config sqlite tpm2-tools ];

  CGO_ENABLED = "1";
  GOTOOLCHAIN = "local";
  GOFLAGS = "-buildvcs=false";
  TPM2TOOLS_TCTI = "device:/dev/tpm0";

  shellHook = ''
    export SCRAPPY_ROOT=${pkgs.lib.escapeShellArg (toString ./.)}
    scrappy-server() (
      cd "$SCRAPPY_ROOT/example" || exit
      go build -o main . && ./main
    )
  '';
}
