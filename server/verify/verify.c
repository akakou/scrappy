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

#include "file_utils.h"

#include <ecdaa.h>

#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#define MAX_SIZE 512
#define HAS_BASENAME 0

#define ERROR_DESRIALIZE_SIG 1
#define ERROR_READING_GPK 2
#define ERROR_DESRIALIZE_GPK 3
#define ERROR_SIGNING 4
#define OK 0

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
        printf("Signature successfully verified!\n");

    return OK;
}

int main(int argc, char *argv[])
{
    uint8_t sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH];

    fread(sig, sizeof(uint8_t), ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH, stdin);
    
    uint8_t message[] = "hogehoge";
    char basename[] = "hogehoge";
    char gpk_path[] = "/attestation/ignored_workspace/group_public.bin";

    int status = verify(sig, message, sizeof(message), basename, sizeof(basename), gpk_path);

    if (status)
        fprintf(stderr, "Error: status %d", status);
    else
        printf("result: %d\n", status);
}

