#pragma once

#include <string.h>
#include <ecdaa.h>
#include <ecdaa-tpm.h>
#include <tss2/tss2_tcti.h>
#include <tss2/tss2_sys.h>
#include <tss2/tss2_tcti_device.h>
#include "examples_rand.h"
#include "file_utils.h"

#define HANDLE_FILE_PATH "/attestation/ignored_workspace/handle.txt"
#define PUB_KEY_PATH "/attestation/ignored_workspace/pub_key.txt"
#define ISSUER_PRIV_KEY_PATH "/attestation/ignored_workspace/issuer_private.bin"
#define CRED_PATH "/attestation/ignored_workspace/member_credential.bin"
#define GPK_PATH "/attestation/ignored_workspace/group_public.bin"
#define TPM_PATH "/dev/tpm0"

#define HAS_BASENAME 1

#define PUBKEY_NUF_LEN 65
#define ECP_FP256BN_LENGTH (2 * MODBYTES_256_56 + 1)