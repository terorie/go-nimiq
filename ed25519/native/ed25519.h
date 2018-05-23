#ifndef ED25519_H
#define ED25519_H

#include <stddef.h>

/* Common functions */
void ed25519_private_key_decompress(unsigned char *az, const unsigned char *private_key);
int ed25519_verify(const unsigned char *signature, const unsigned char *message, size_t message_len, const unsigned char *public_key);

/* Single signature functions */
void ed25519_public_key_derive(unsigned char *out_public_key, const unsigned char *private_key);
void ed25519_sign(unsigned char *signature, const unsigned char *message, size_t message_len, const unsigned char *public_key, const unsigned char *private_key);

/* Common multisig functions */
int ed25519_create_commitment(unsigned char *secret_r, unsigned char *commitment_R, const unsigned char *randomness);
void ed25519_aggregate_commitments(unsigned char *aggregate_commitment, const unsigned char *commitments, const size_t num_commitments);
void ed25519_add_scalars(unsigned char *scalar_AB, const unsigned char *scalar_A, const unsigned char *scalar_B);

/* Delinearized multisig functions */
void ed25519_hash_public_keys(unsigned char *hash, const unsigned char *public_keys, const size_t num_public_keys);
void ed25519_delinearize_public_key(unsigned char *delinearized_public_key, const unsigned char *public_keys_hash, const unsigned char *public_key);
void ed25519_aggregate_delinearized_public_keys(unsigned char *aggregate_public_key, const unsigned char *public_keys_hash, const unsigned char *public_keys, const size_t num_public_keys);
void ed25519_derive_delinearized_private_key(unsigned char *multisig_private_key, const unsigned char *public_keys_hash, const unsigned char *public_key, const unsigned char *private_key);
void ed25519_delinearized_partial_sign(unsigned char *partial_signature, const unsigned char *message, size_t message_len, const unsigned char* commitment_R, const unsigned char *secret_r, const unsigned char *public_keys, size_t num_cosigners, const unsigned char *public_key, const unsigned char *private_key);

#endif
