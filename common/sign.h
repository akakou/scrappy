#include "examples_rand.h"
#include "file_utils.h"
#include "init.h"

#include <ecdaa.h>
#include <ecdaa-tpm.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <tss2/tss2_tcti.h>
#include <tss2/tss2_sys.h>
#include <tss2/tss2_tcti_device.h>

#define HAS_BASENAME 1

#define ERROR_READING_SECRET 1
#define ERROR_DESRIALIZE_SECRET 2
#define ERROR_READING_CRED 3
#define ERROR_DESRIALIZE_CRED 4
#define ERROR_SIGNING 5
#define ERROR_INIT_TPM 6
#define ERROR_INIT_TPMECDAA 7
#define KEY_HANDLE 81010000
#define OK 0
#define TPM_PATH "/dev/tpm0"


int sign(uint8_t raw_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH], uint8_t *message, int msg_len, char *basename, int basename_len, char *credential_path)
{
    int ret;
    size_t size;
    uint8_t buffer[1024];
    uint8_t tcti_buffer[256];
    TPM2_HANDLE handle;

    read_key_handle(&handle, HANDLE_FILE_PATH);

    ECP_FP256BN public_key;
    TSS2_TCTI_CONTEXT *tcti_ctx;

    struct ecdaa_tpm_context ecdaa_ctx;
    struct ecdaa_member_secret_key_FP256BN sk;
    struct ecdaa_credential_FP256BN cred;

    ecdaa_init(&ecdaa_ctx, handle, tcti_buffer, sizeof(tcti_buffer));

    if (0 != ret)
    {
        fprintf(stderr ,"Error: ecdaa_tpm_context_init failed: 0x%x\n", ret);
        return ERROR_INIT_TPMECDAA;
    }

    // Read member credential from disk
    ret =  ecdaa_credential_FP256BN_deserialize_file(&cred, credential_path);

    if (0 != ret)
    {
        fprintf(stderr, "Error deserializing member credential (%d)\n", ret);
        return ERROR_DESRIALIZE_CRED;
    }

    // Create signature
    struct ecdaa_signature_FP256BN sig;

    ret = ecdaa_signature_TPM_FP256BN_sign(&sig, message, msg_len, basename, basename_len, &cred, examples_rand, &ecdaa_ctx);
    if (ret != 0)
    {
        message[msg_len] = 0;
        fprintf(stderr, "Error signing message: \"%s\"(%d)\n", (char *)message, ret);
        return ERROR_SIGNING;
    }

    ecdaa_signature_FP256BN_serialize(raw_sig, &sig, HAS_BASENAME);

    return OK;
}
