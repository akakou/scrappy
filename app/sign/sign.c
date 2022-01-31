#include "examples_rand.h"
#include "file_utils.h"

#include <ecdaa.h>

#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#define MAX_MESSAGE_SIZE 1024
#define HAS_BASENAME 1

#define ERROR_READING_SECRET 1
#define ERROR_DESRIALIZE_SECRET 2
#define ERROR_READING_CRED 3
#define ERROR_DESRIALIZE_CRED 4
#define ERROR_SIGNING 5
#define OK 0

int attest(uint8_t buffer[1024], uint8_t *message, int msg_len, char *basename, int basename_len, char *secret_key_path, char *credential_path)
{
    struct ecdaa_member_secret_key_FP256BN sk;
    struct ecdaa_credential_FP256BN cred;

    // Read secret key from disk
    if (ECDAA_MEMBER_SECRET_KEY_FP256BN_LENGTH != read_file_into_buffer(buffer, ECDAA_MEMBER_SECRET_KEY_FP256BN_LENGTH, secret_key_path))
    {
        fprintf(stderr, "Error reading member secret key file: \"%s\"\n", secret_key_path);
        return ERROR_READING_SECRET;
    }
    if (0 != ecdaa_member_secret_key_FP256BN_deserialize(&sk, buffer))
    {
        fputs("Error deserializing member secret key\n", stderr);
        return ERROR_DESRIALIZE_SECRET;
    }

    // Read member credential from disk
    if (ECDAA_CREDENTIAL_FP256BN_LENGTH != read_file_into_buffer(buffer, ECDAA_CREDENTIAL_FP256BN_LENGTH, credential_path))
    {
        fprintf(stderr, "Error reading member credential file: \"%s\"\n", credential_path);
        return ERROR_READING_CRED;
    }
    if (0 != ecdaa_credential_FP256BN_deserialize(&cred, buffer))
    {
        fputs("Error deserializing member credential\n", stderr);
        return ERROR_DESRIALIZE_CRED;
    }

    // Create signature
    struct ecdaa_signature_FP256BN sig;
    if (0 != ecdaa_signature_FP256BN_sign(&sig, message, msg_len, basename, basename_len, &sk, &cred, examples_rand))
    {
        message[msg_len] = 0;
        fprintf(stderr, "Error signing message: \"%s\"\n", (char *)message);
        return ERROR_SIGNING;
    }

    ecdaa_signature_FP256BN_serialize(buffer, &sig, HAS_BASENAME);

    return OK;
}

void print(uint8_t *buffer, size_t len)
{
    for (int i = 0; i < len; i++)
        printf("[%d] %02x\n", i, buffer[i]);
}

int main(int argc, uint8_t *argv[])
{
    uint8_t buffer[1024];
    memset(buffer, 0x00, sizeof(buffer));

    uint8_t message[] = "hogehoge";
    char basename[] = "hogehoge";
    char secret_key_path[] = "./ignored_workspace/member_private.bin";
    char credential_path[] = "./ignored_workspace/member_credential.bin";

    int status = attest(buffer, message, sizeof(message), basename, sizeof(basename), secret_key_path, credential_path);
    if (status) 
        fprintf(stderr, "Error: status %d", status);
    else
        print(buffer, sizeof(buffer));
}

