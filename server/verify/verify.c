/******************************************************************************
 *
 * Copyright 2017 Xaptum, Inc.
 * 
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 * 
 *        http://www.apache.org/licenses/LICENSE-2.0
 * 
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License
 *
 *****************************************************************************/


#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <ecdaa.h>
#define DEBUG 1

#ifdef DEBUG
#include "../../client/app/sign/sign.c"
#else
#include "file_utils.h"
#endif


#define MAX_SIZE 1024
#define HAS_BASENAME 1

#define OK 0
#define WRONG 1
#define ERROR_DESRIALIZE_SIG 2
#define ERROR_READING_GPK 3
#define ERROR_DESRIALIZE_GPK 4
#define ERROR_SIGNING 5

#define ECP_FP256BN_LENGTH (2 * MODBYTES_256_56 + 1)

#include "examples_rand.h"

int verify(u_int8_t raw_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH], uint8_t *message, int msg_len, uint8_t *basename, int basename_len, char *gpk_path)
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
    if (ECDAA_GROUP_PUBLIC_KEY_FP256BN_LENGTH != read_file_into_buffer(buffer, ECDAA_GROUP_PUBLIC_KEY_FP256BN_LENGTH, gpk_path))
    {
        fprintf(stderr, "Error reading group public key file: \"%s\"\n", gpk_path);
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
        printf("test:%d\n", ecdaa_signature_FP256BN_verify(&sig, &gpk, &revocations, message, msg_len, basename, basename_len));
        return WRONG;
    }

    return OK;
}

// void print(uint8_t *buffer, size_t len)
// {
//     for (int i = 0; i < len; i++)
//         printf("[%d] %02x\n", i, buffer[i]);
// }


int parse_k(u_int8_t k[ECP_FP256BN_LENGTH], u_int8_t raw_sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH])
{
    int padding = 3 * MODBYTES_256_56 + 4 * ECP_FP256BN_LENGTH;
    memcpy(k, raw_sig + padding, ECP_FP256BN_LENGTH);

    return OK;
}

int main(int argc, char *argv[])
{
    uint8_t sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH];
    memset(sig, 0x00, sizeof(sig));

    uint8_t message[] = "hogehoge";

    char basename[] = "hogehoge";
    char secret_key_path[] = "/attestation/ignored_workspace/member_private.bin";
    char credential_path[] = "/attestation/ignored_workspace/member_credential.bin";
    char gpk_path[] = "/attestation/ignored_workspace/group_public.bin";

    int status = sign(sig, message, sizeof(message), basename, sizeof(basename), secret_key_path, credential_path);

    u_int8_t k[ECP_FP256BN_LENGTH];
    memset(k, 0x00, sizeof(k));

    status = verify(sig, message, sizeof(message), basename, sizeof(basename), gpk_path);
    printf("status: %d\n", status);

    if (status)
        fprintf(stderr, "Error: status on verify (%d)", status);
    else
        printf("result: %d\n", status);

    status = parse_k(k, sig);

    if (status)
        fprintf(stderr, "Error: status on parase k (%d)", status);
    // else
    //     print(k, ECP_FP256BN_LENGTH);

    status = verify(sig, "hello", sizeof("hello"), basename, sizeof(basename), gpk_path);
    printf("failed status: %d\n", status);
}

