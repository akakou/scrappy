#pragma once

#include "init.h"
#include "common.h"

#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <ecdaa.h>

#include "examples_rand.h"
#include "file_utils.h"

#define MAX_SIZE 1024

#define OK 0
#define WRONG 1
#define ERROR_DESRIALIZE_SIG 2
#define ERROR_READING_GPK 3
#define ERROR_DESRIALIZE_GPK 4
#define ERROR_SIGNING 5


int verify(u_int8_t raw_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH], uint8_t *message, int msg_len, uint8_t *basename, int basename_len)
{
    uint8_t buffer[MAX_SIZE];

    struct ecdaa_signature_FP256BN sig;
    struct ecdaa_revocations_FP256BN revocations;

    revocations.sk_list = NULL;
    revocations.sk_length = 0;
    revocations.bsn_list = NULL;
    revocations.bsn_length = 0;

    if (0 != ecdaa_signature_FP256BN_deserialize(&sig, raw_sig, HAS_BASENAME))
    {
        fputs("Error deserializing signature\n", stderr);
        return ERROR_DESRIALIZE_SIG;
    }

    // Read group public key from disk
    struct ecdaa_group_public_key_FP256BN gpk;
    if (ECDAA_GROUP_PUBLIC_KEY_FP256BN_LENGTH != read_file_into_buffer(buffer, ECDAA_GROUP_PUBLIC_KEY_FP256BN_LENGTH, GPK_PATH))
    {
        fprintf(stderr, "Error reading group public key file: \"%s\"\n", GPK_PATH);
        return ERROR_READING_GPK;
    }
    if (0 != ecdaa_group_public_key_FP256BN_deserialize(&gpk, buffer))
    {
        fputs("Error deserializing group public key\n", stderr);
        return ERROR_DESRIALIZE_GPK;
    }

    // Verify signature
    if (0 != ecdaa_signature_FP256BN_verify(&sig, &gpk, &revocations, message, msg_len, basename, basename_len))
    {
        return WRONG;
    }

    return OK;
}


int parse_k(u_int8_t k[ECP_FP256BN_LENGTH], u_int8_t raw_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH])
{
    int padding = 3 * MODBYTES_256_56 + 4 * ECP_FP256BN_LENGTH;
    memcpy(k, raw_sig + padding, ECP_FP256BN_LENGTH);

    return OK;
}
