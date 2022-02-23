#include <string.h>
#include <ecdaa.h>
#include <ecdaa-tpm.h>
#include <tss2/tss2_tcti.h>
#include <tss2/tss2_sys.h>
#include <tss2/tss2_tcti_device.h>
#include "../../attestation/common/examples_rand.h"
#include "../../attestation/thirdparty/ecdaa/build/common/amcl-extensions/ecp_FP256BN.h"
#define HANDLE_FILE_PATH "/attestation/ignored_workspace/handle.txt"
#define PUB_KEY_PATH "/attestation/ignored_workspace/pub_key.txt"
#define ISSUER_PRIV_KEY_PATH "/attestation/ignored_workspace/issuer_private.bin"
#define CRED_PATH "/attestation/ignored_workspace/member_credential.bin"
#define TPM_PATH "/dev/tpm0"
#define KEY_HANDLE 81010000
#define PUBKEY_NUF_LEN 65

void print_bin(u_int8_t *bytes, size_t size);

int read_public_key_from_files(uint8_t *public_key,
                               TPM2_HANDLE *key_handle,
                               const char *pub_key_filename,
                               const char *handle_filename);

int main() {
    struct ecdaa_issuer_secret_key_FP256BN isk;
    struct ecdaa_credential_FP256BN cred;
    struct ecdaa_credential_FP256BN_signature cred_sig;
    struct ecdaa_member_public_key_FP256BN pk;
    TSS2_TCTI_CONTEXT *tcti_ctx;
    struct ecdaa_tpm_context ecdaa_ctx;

    uint8_t tcti_buffer[256];
    uint8_t serialized_public_key[PUBKEY_NUF_LEN];

    const char *mssim_conf = "host=localhost,port=2321";
    const char *device_conf = "/dev/tpm0";

    memset(tcti_buffer, 0, sizeof(tcti_buffer));

    int ret = 0;

    TPM2_HANDLE key_handle = 0;

    if (0 != read_public_key_from_files(serialized_public_key, &key_handle, PUB_KEY_PATH, HANDLE_FILE_PATH))
    {
        printf("Error: error reading in public key files '%s' and '%s'\n", PUB_KEY_PATH, HANDLE_FILE_PATH);
        return -1;
    }

    // if (0 != ecp_ZZZ_deserialize(&ctx->public_key, (uint8_t *)ctx->serialized_public_key))
    // {
    //     printf("Error: error public key to point\n");
    //     return -1;
    // }

    tcti_ctx = (TSS2_TCTI_CONTEXT *)tcti_buffer;
#ifdef USE_TCP_TPM
    (void)device_conf;
    size_t size;
    ret = Tss2_Tcti_Mssim_Init(NULL, &size, mssim_conf);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Failed to get allocation size for tcti context\n");
        return -1;
    }
    if (size > sizeof(ctx->tcti_buffer))
    {
        printf("Error: socket TCTI context size larger than pre-allocated buffer\n");
        return -1;
    }
    ret = Tss2_Tcti_Mssim_Init(ctx->tcti_context, &size, mssim_conf);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Error: Unable to initialize socket TCTI context\n");
        return -1;
    }
#else
    (void)mssim_conf;
    size_t size;
    ret = Tss2_Tcti_Device_Init(NULL, &size, device_conf);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Failed to get allocation size for tcti context\n");
        return -1;
    }
    if (size > sizeof(tcti_buffer))
    {
        printf("Error: device TCTI context size larger than pre-allocated buffer\n");
        return -1;
    }
    ret = Tss2_Tcti_Device_Init(tcti_ctx, &size, device_conf);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Error: Unable to initialize device TCTI context\n");
        return -1;
    }
#endif

    ret = ecdaa_tpm_context_init(&ecdaa_ctx, key_handle, NULL, 0, tcti_ctx);
    if (0 != ret)
    {
        printf("Error: ecdaa_tpm_context_init failed: 0x%x\n", ret);
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
