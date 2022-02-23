#include <string.h>
#include <ecdaa.h>
#include <ecdaa-tpm.h>
#include <tss2/tss2_tcti.h>
#include <tss2/tss2_sys.h>
#include <tss2/tss2_tcti_device.h>
#include "../../attestation/common/examples_rand.h"
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
    TPM2_HANDLE key_handle;

    if (0 != read_public_key_from_files(serialized_public_key, &key_handle, PUB_KEY_PATH, HANDLE_FILE_PATH))
    {
        printf("Error: error reading in public key files '%s' and '%s'\n", PUB_KEY_PATH, HANDLE_FILE_PATH);
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


int read_public_key_from_files(uint8_t *public_key,
                               TPM2_HANDLE *key_handle,
                               const char *pub_key_filename,
                               const char *handle_filename)
{
    int ret = 0;

    FILE *pub_key_file_ptr = fopen(pub_key_filename, "r");
    if (NULL == pub_key_file_ptr)
        return -1;
    do
    {
        for (unsigned i = 0; i < 65; i++)
        {
            unsigned byt;
            if (fscanf(pub_key_file_ptr, "%02X", &byt) != 1)
            {
                ret = -1;
                break;
            }
            public_key[i] = (uint8_t)byt;
        }
    } while (0);
    (void)fclose(pub_key_file_ptr);
    if (0 != ret)
        return -1;

    FILE *handle_file_ptr = fopen(handle_filename, "r");
    if (NULL == handle_file_ptr)
        return -1;
    do
    {
        for (int i = (sizeof(TPM2_HANDLE) - 1); i >= 0; i--)
        {
            unsigned byt;
            if (fscanf(handle_file_ptr, "%02X", &byt) != 1)
            {
                ret = -1;
                break;
            }
            *key_handle += byt << (i * 8);
        }
        if (0 != ret)
            break;
    } while (0);
    (void)fclose(handle_file_ptr);

    return ret;
}
