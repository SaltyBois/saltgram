#!/bin/sh

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
--output $HOME/secrets/saltgram-c0751de619fa.json saltgram-c0751de619fa.json.gpg

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
--output $HOME/secrets/saltgram-service-key.json saltgram-service-key.json.gpg