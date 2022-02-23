#include <string.h>

#include "common.h"

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

    int status = sign(sig, message, sizeof(message), basename, sizeof(basename));

    u_int8_t k[ECP_FP256BN_LENGTH];
    memset(k, 0x00, sizeof(k));

    status = verify(sig, message, sizeof(message), basename, sizeof(basename));
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

    status = verify(sig, "hello", sizeof("hello"), basename, sizeof(basename));
    printf("non valid status: %d\n", status);
}
