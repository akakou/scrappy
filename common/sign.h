#pragma once

#include "init.h"
#include "common.h"

#define ERROR_READING_SECRET 1
#define ERROR_DESRIALIZE_SECRET 2
#define ERROR_READING_CRED 3
#define ERROR_DESRIALIZE_CRED 4
#define ERROR_SIGNING 5
#define ERROR_INIT_TPM 6
#define ERROR_INIT_TPMECDAA 7
#define OK 0


int sign(uint8_t serialized_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH], uint8_t *message, int msg_len, char *basename, int basename_len)
{
    uint8_t tpm_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH];

    int ret;
    size_t size;
    uint8_t buffer[1024];
    uint8_t tcti_buffer[256];
    TPM2_HANDLE handle;

    ECP_FP256BN public_key;
    TSS2_TCTI_CONTEXT *tcti_ctx;

    struct ecdaa_tpm_context ecdaa_ctx;
    struct ecdaa_member_secret_key_FP256BN sk;
    struct ecdaa_credential_FP256BN cred;
    struct ecdaa_signature_FP256BN sig;

    ret = read_key_handle(&handle, HANDLE_FILE_PATH);

    if (0 != ret)
    {
        fprintf(stderr, "Error: read_key_handle failed: 0x%x\n", ret);
        return ERROR_INIT_TPMECDAA;
    }

    ecdaa_init(&ecdaa_ctx, handle, tcti_buffer, sizeof(tcti_buffer));

    if (0 != ret)
    {
        fprintf(stderr ,"Error: ecdaa_tpm_context_init failed: 0x%x\n", ret);
        return ERROR_INIT_TPMECDAA;
    }

    // Read member credential from disk
    ret = ecdaa_credential_FP256BN_deserialize_file(&cred, CRED_PATH);

    if (0 != ret)
    {
        fprintf(stderr, "Error deserializing member credential (%d)\n", ret);
        return ERROR_DESRIALIZE_CRED;
    }

    // Create signature
    ret = ecdaa_signature_TPM_FP256BN_sign(&sig, message, msg_len, basename, basename_len, &cred, examples_rand, &ecdaa_ctx);
    if (ret != 0)
    {
        fprintf(stderr, "Error signing (%d)\n", ret);
        return ERROR_SIGNING;
    }

    ecdaa_signature_FP256BN_serialize(tpm_sig, &sig, HAS_BASENAME);
    memcpy(serialized_sig, tpm_sig, ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH);

    return OK;
}

void print2(uint8_t *buffer, size_t len)
{
    for (int i = 0; i < len; i++)
        printf("[%d] %02x\n", i, buffer[i]);
}
