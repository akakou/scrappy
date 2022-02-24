#pragma once

#include "common.h"

int read_public_key_from_files(uint8_t *public_key,
                               TPM2_HANDLE *key_handle,
                               const char *pub_key_filename,
                               const char *handle_filename);

int ecdaa_init(struct ecdaa_tpm_context *ecdaa_ctx, TPM2_HANDLE key_handle, uint8_t tcti_buffer[256], size_t tcti_buffer_size)
{
    memset(tcti_buffer, 0, tcti_buffer_size);

    int ret = 0;

    TSS2_TCTI_CONTEXT *tcti_ctx = (TSS2_TCTI_CONTEXT *)tcti_buffer;

    size_t size;
    ret = Tss2_Tcti_Device_Init(NULL, &size, TPM_PATH);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Failed to get allocation size for tcti context\n");
        return 2;
    }
    if (size > tcti_buffer_size)
    {
        printf("Error: device TCTI context size larger than pre-allocated buffer\n");
        return 3;
    }
    ret = Tss2_Tcti_Device_Init(tcti_ctx, &size, TPM_PATH);
    if (TSS2_RC_SUCCESS != ret)
    {
        printf("Error: Unable to initialize device TCTI context\n");
        return 4;
    }

    ret = ecdaa_tpm_context_init(ecdaa_ctx, key_handle, NULL, 0, tcti_ctx);
    if (0 != ret)
    {
        printf("Error: ecdaa_tpm_context_init failed: 0x%x\n", ret);
        return 5;
    }

    return 0;
}

int read_public_key(uint8_t *public_key, const char *pub_key_filename){
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

    return ret;
}

int read_key_handle(TPM2_HANDLE *key_handle,
                    const char *handle_filename)
{
    int ret = 0;

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
