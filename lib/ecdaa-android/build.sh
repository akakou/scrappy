cargo build --release --lib --target aarch64-linux-android
cargo build --release --lib --target armv7-linux-androideabi
cargo build --release --lib --target i686-linux-android

mkdir /out/x86
mkdir /out/armeabi-v7a
mkdir /out/arm64-v8a

mv ./target/i686-linux-android/release/libecdaa_android.so /out/x86
mv ./target/armv7-linux-androideabi/release/libecdaa_android.so /out/armeabi-v7a 
mv ./target/aarch64-linux-android/release/libecdaa_android.so /out/arm64-v8a
