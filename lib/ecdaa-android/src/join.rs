use ecdaa::join::ReqForJoin;

use crate::utils::random;

const JOIN_MSG: &[u8; 7] = b"scrappy";

pub fn join_for_request() -> String {
    let mut rng = random();

    let join_req_ppt = ReqForJoin::random(JOIN_MSG.as_slice(), &mut rng);
    let req_for_join = join_req_ppt.unwrap();

    let sk = serde_json::to_string(&req_for_join.0).unwrap();
    let req = serde_json::to_string(&req_for_join.1).unwrap();

    let result = format!("{}\n{}", sk, req);

    return result;
}
