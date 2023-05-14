use ecdaa::fp256bn_amcl::rand::RAND;

use rand::{thread_rng, RngCore};

pub fn random() -> RAND {
    let mut th_rng = thread_rng();
    let mut seed = [0; 32].to_vec();
    th_rng.fill_bytes(&mut seed);

    let mut rng = RAND::new();
    rng.seed(32, &seed);

    rng
}
