# scrappy-android
An Android app acting a signer of the scrappy protocol.

Build the Go library from the repository root:

```sh
nix-shell android/shell.nix --run scrappy-android
```

The shell provides Go with CGO, JDK 11, Android platform 33, build tools 33.0.0,
NDK 23.1.7779620 (matching the former Docker environment), and gomobile from the
shared nixpkgs pin. It sets `JAVA_HOME`, `ANDROID_HOME`, `ANDROID_SDK_ROOT`,
`ANDROID_NDK_HOME`, and `ANDROID_NDK_ROOT`. The SDK license is accepted by the
Nix configuration. No emulator or system images are included.

`core/android/go.mod` pins `golang.org/x/mobile` to the same revision as gomobile
in nixpkgs. Update both together when changing the nixpkgs pin. This version of
the Android module requires Go 1.22 or later, supplied by the shell.

The output is `core/android/scrappy_crypto.aar`, which the Gradle application
loads relative to the checkout. Open `android/scrappy/` in Android Studio to
build the app; the Nix command above builds the Go library only.
