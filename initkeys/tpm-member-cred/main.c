#include <string.h>
#include <ecdaa.h>
#include <ecdaa-tpm.h>
#include "../../attestation/common/common.h"
#include "../../attestation/common/init.h"

int main() {    
    struct ecdaa_issuer_secret_key_FP256BN isk;
    struct ecdaa_credential_FP256BN cred;
    struct ecdaa_credential_FP256BN_signature cred_sig;
    struct ecdaa_member_public_key_FP256BN pk;
    struct ecdaa_tpm_context ecdaa_ctx;

    uint8_t tcti_buffer[256];
    uint8_t serialized_public_key[PUBKEY_NUF_LEN];

    int ret = 0;
    TPM2_HANDLE key_handle = 0;

    if (0 != read_public_key(serialized_public_key, PUB_KEY_PATH))
    {
        printf("Error: error reading in public key file '%s'\n", PUB_KEY_PATH);
        return 1;
    }

    if (0 != read_key_handle(&key_handle, HANDLE_FILE_PATH))
    {
        printf("Error: error reading in handle file '%s'\n", HANDLE_FILE_PATH);
        return 1;
    }

    ret = ecdaa_init(&ecdaa_ctx, key_handle, tcti_buffer, sizeof(tcti_buffer));

    if (ret != 0) {
        printf("Error: ecdaa_init failed: 0x%x\n", ret);
        return -1;
    }

    uint8_t *nonce = (uint8_t *)"nonce";
    uint32_t nonce_len = 5;

    ret = ecdaa_member_key_pair_TPM_FP256BN_generate(&pk, serialized_public_key, nonce, nonce_len, &ecdaa_ctx);

    if (ret != 0) {
        printf("Error: ecdaa_member_key_pair_generate failed: 0x%x\n", ret);
        return -1;
    }

    ret = ecdaa_issuer_secret_key_FP256BN_deserialize_file(&isk, ISSUER_PRIV_KEY_PATH);

    if (ret != 0)
    {
        printf("Error: ecdaa_issuer_secret_key_FP256BN_deserialize_file failed: 0x%x\n", ret);
        return -1;
    }

    ret = ecdaa_credential_FP256BN_generate(&cred, &cred_sig, &isk, &pk, examples_rand);

    if (ret != 0)
    {
        printf("Error: ecdaa_credential_FP256BN_generate failed: 0x%x\n", ret);
        return -1;
    }

    ret = ecdaa_credential_FP256BN_serialize_file(CRED_PATH, &cred);

    if (ret != 0)
    {
        printf("Error: ecdaa_credential_FP256BN_serialize_file failed: 0x%x\n", ret);
        return -1;
    }

    return 0;
}


