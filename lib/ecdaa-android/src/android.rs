use crate::join::join_for_request;
use jni::objects::JObject;
use jni::sys::jstring;
use jni::JNIEnv;

#[no_mangle]
pub unsafe extern "C" fn Java_com_github_akakou_scrappy_Scrappy_joinForReqest(
    env: JNIEnv,
    _this: JObject,
) -> jstring {
    let result = join_for_request();

    env.new_string(result)
        .expect("Couldn't create Java string!")
        .into_inner()
}
