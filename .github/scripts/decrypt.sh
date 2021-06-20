#!/bin/sh

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
--output $HOME/secrets/saltgram-c0751de619fa.json secrets/saltgram-c0751de619fa.json.gpg

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
--output $HOME/secrets/saltgram-service-key.json secrets/saltgram-service-key.json.gpg