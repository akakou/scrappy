#include <string.h>

#include "examples_rand.h"
#include "file_utils.h"

#include <ecdaa.h>
#include "verify.h"
#include "sign.h"

void print(uint8_t *buffer, size_t len)
{
    for (int i = 0; i < len; i++)
        printf("[%d] %02x\n", i, buffer[i]);
}

int main(int argc, char *argv[])
{
    uint8_t sig[ECDAA_SIGNATURE_FP256BN_WITH_NYM_LENGTH];
    memset(sig, 0x00, sizeof(sig));

    uint8_t message[] = "hogehoge";

    char basename[] = "hogehoge";
    char credential_path[] = "/attestation/ignored_workspace/member_credential.bin";
    char gpk_path[] = "/attestation/ignored_workspace/group_public.bin";

    int status = sign(sig, message, sizeof(message), basename, sizeof(basename), credential_path);

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
        // print(k, ECP_FP256BN_LENGTH);

    status = verify(sig, "hello", sizeof("hello"), basename, sizeof(basename), gpk_path);
    printf("non valid status: %d\n", status);
}
