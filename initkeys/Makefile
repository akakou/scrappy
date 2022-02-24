WORK_SPACE = /attestation/ignored_workspace
# TPM_WORK_SPACE = /attestation/thirdparty/ecdaa/build/testBin/
ISSUER_PUBLIC = $(WORK_SPACE)/issuer_public.bin
ISSUER_PRIVATE = $(WORK_SPACE)/issuer_private.bin
GROUP_PUBLIC = $(WORK_SPACE)/group_public.bin
MEMBER_PUBLIC = $(WORK_SPACE)/member_public.bin
MEMBER_CRED = $(WORK_SPACE)/member_credential.bin
CRED_SIG = $(WORK_SPACE)/cred_sig.bin
TPM_PARAM = $(WORK_SPACE)/pub_key.txt $(WORK_SPACE)/handle.txt 
INIT_TPM = /attestation/thirdparty/ecdaa/build/testBin/ecdaa-create_tpm_key-util
ECDAA = ecdaa

FLAGS = `pkg-config ecdaa --libs --cflags` `pkg-config ecdaa-tpm --libs --cflags` `pkg-config tss2-tcti-device --libs` `pkg-config tss2-sys --libs --cflags`
SRC_FOLDER = ./tpm-member-cred
COMMON_SRC_FILES = /attestation/common/examples_rand.h /attestation/common/examples_rand.c


SRC_FILES = $(SRC_FOLDER)/*.c $(COMMON_SRC_FILES)
BIN = ./bin
GEN_CRED_BIN = $(BIN)/gen_cred

all: $(ISSUER_PUBLIC) $(ISSUER_PRIVATE) $(GEN_CRED_BIN) $(TPM_PARAM)
	$(GEN_CRED_BIN)

setup_tpm_param: $(TPM_PARAM)
	echo 'ok'

setup_member_cred: $(GEN_CRED_BIN)
	echo 'ok'

$(MEMBER_CRED): $(ISSUER_PUBLIC) $(ISSUER_PRIVATE) $(MEMBER_PUBLIC)
	$(ECDAA) issuer issuecredential -p $(MEMBER_PUBLIC) -s $(ISSUER_PRIVATE) -c $(MEMBER_CRED) -r $(CRED_SIG)

$(ISSUER_PUBLIC) $(ISSUER_PRIVATE):
	$(ECDAA) issuer genkeys -p $(ISSUER_PUBLIC) -s $(ISSUER_PRIVATE) \
		&& $(ECDAA) extractgpk -p $(ISSUER_PUBLIC) -g $(GROUP_PUBLIC)

$(GEN_CRED_BIN):
	mkdir -p $(BIN) && gcc $(SRC_FILES) -o $(GEN_CRED_BIN) $(FLAGS) 

$(TPM_PARAM):
	$(INIT_TPM) $(TPM_PARAM)

setup_issuer: $(ISSUER_PUBLIC), $(ISSUER_PRIVATE)
	echo 'ok'

clean:
	rm -rf $(WORK_SPACE)/* $(BIN)