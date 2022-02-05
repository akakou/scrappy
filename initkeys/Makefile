WORK_SPACE = /attestation/ignored_workspace
ISSUER_PUBLIC = $(WORK_SPACE)/issuer_public.bin
ISSUER_PRIVATE = $(WORK_SPACE)/issuer_private.bin
GROUP_PUBLIC = $(WORK_SPACE)/group_public.bin
MEMBER_PUBLIC = $(WORK_SPACE)/member_public.bin
MEMBER_CRED = $(WORK_SPACE)/member_credential.bin
CRED_SIG = $(WORK_SPACE)/cred_sig.bin

ECDAA = ecdaa

$(MEMBER_CRED): $(ISSUER_PUBLIC) $(ISSUER_PRIVATE) $(MEMBER_PUBLIC)
	$(ECDAA) issuer issuecredential -p $(MEMBER_PUBLIC) -s $(ISSUER_PRIVATE) -c $(MEMBER_CRED) -r $(CRED_SIG)

$(ISSUER_PUBLIC), $(ISSUER_PRIVATE):
	$(ECDAA) issuer genkeys -p $(ISSUER_PUBLIC) -s $(ISSUER_PRIVATE) \
		&& $(ECDAA) extractgpk -p $(ISSUER_PUBLIC) -g $(GROUP_PUBLIC)

setup_issuer: $(ISSUER_PUBLIC), $(ISSUER_PRIVATE)
	echo 'ok'

clean:
	rm $(WORK_SPACE)/*