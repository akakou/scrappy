{ pkgs ? import ../nix/pkgs.nix }:

let
  base = import ../shell.nix { inherit pkgs; };
  android = pkgs.androidenv.composeAndroidPackages {
    platformVersions = [ "33" ];
    buildToolsVersions = [ "33.0.0" ];
    includeNDK = true;
    ndkVersions = [ "23.1.7779620" ];
    includeCmake = false;
    includeEmulator = false;
    includeSystemImages = false;
    includeSources = false;
  };
  sdk = "${android.androidsdk}/libexec/android-sdk";
  gomobile = pkgs.gomobile.override { androidPkgs = android; };
in
pkgs.mkShell {
  inputsFrom = [ base ];
  packages = [ pkgs.jdk11 android.androidsdk gomobile ];
  inherit (base) CGO_ENABLED GOTOOLCHAIN GOFLAGS TPM2TOOLS_TCTI;

  JAVA_HOME = pkgs.jdk11.home;
  ANDROID_HOME = sdk;
  ANDROID_SDK_ROOT = sdk;
  ANDROID_NDK_HOME = "${sdk}/ndk/23.1.7779620";
  ANDROID_NDK_ROOT = "${sdk}/ndk/23.1.7779620";

  shellHook = ''
    scrappy-android() (
      cd "$SCRAPPY_ROOT/core/android" || exit
      # nixpkgs' gomobile uses NIX_BUILD_TOP for its working directory.
      scrappy_work=$(mktemp -d) || exit
      trap 'rm -rf -- "$scrappy_work"' EXIT
      export NIX_BUILD_TOP="$scrappy_work"
      gomobile bind -v -o scrappy_crypto.aar -target=android -androidapi=33
    )
  '';
}
