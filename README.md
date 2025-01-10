# Vault

## Overview
Vault is a multi-user, offline, privacy-preserving encrypted note storage system. The only plaintext data stored is the username and the identifiers used to lookup the users encrypted Keyset and Metadata. Each user has a BaseKey that is derived from their password using Argon2. This BaseKey is then used to derive the encryption keys and identifiers used to lookup the Keyset and Metadata. Since the Keyset and Metadata IDs are derived from the BaseKey it is not possible to determine which Keyset or Metadata belongs to which user, when there is more than one user.

### Threat Model
While vault can be used to store passwords as notes, it is not a password manager and should not be treated as one when considering potential threats. In addition, Vault is not able to keep your UnlockedVault safe on a machine that has been compromised.

## Accounts
Vault is a multi-user system but there is no concept of an administrator. User's can register for an account using a unique username and as long as they don't share their password, they are the only ones who can read their data.

### Registering
To register for an account you need to provide a username and a password that is at least 16 characters long. There are no other password requirements and Unicode passwords are acceptable. If the username is already in use, you will receive an error while registering.

When you register your account you must remember the password you used. If you do not, there is no way to decrypt the data. Any data in your vault will be lost until your remember the password.

### Authenticating
To login to your Vault, you need to provide your username and password. If they are correct, the Vault will be unlocked and you will be able to read, add, and update Items in your vault.

### Password Changes
To change the password on your vault you must provide the username, old password, and new password. Vault will derive a new AuthToken and CryptKey and reencrypt your User and Keyset. In addition, it will add a new BaseKey to your Keyset and reencrypt the Metadata with the new BaseKey. Each time you login after changing your password, Vault will begin reencrypting your Items with the new BaseKey. Over time, all of the Items will be reencrypted and the old key will be purged.

## Cryptography
### Algorithms
Vault uses xChaCha20 to encrypt all data, Argon2id for slow key derivation, and Blake2b for fast key derivation. These cryptographic primitives are imported from the golang/x/crypto repository. Vault is purposely designed for cryptographic agility, making it "relatively easy" to upgrade encryption and key derivation algorithms in the future. Vault is not designed for sharing data so no public key encryption is used, which means we do not have to worry about post-quantum cryptography at this time.

### Keys
Vault uses a number of keys for encryption, some are derived from the user's password (using Argon2id) and some are derived from the user's BaseKey (using Blake2b). Each of the key types is defined below.

*BaseKey* - The BaseKey is 256 bits and is derived from the user's password using Argon2id and the recommended settings from the RFC.

*AuthKey* - This key is used to encrypt the User object that contains the identifiers for the user's Keyset and Metadata. This key is derived from the BaseKey and the phrase "This key will be used for authentication.", using Blake2b.

*CryptKey* - This key is used to encrypt the user's Keyset and is derived from the BaseKey and the phrase "This key will be used for encryption.", using Blake2b.

*Metadata CryptKey* - This key is used to encrypt the user's Metadata and is derived from a BaseKey in the Keyset and the MetadataId using Blake2b.

*Item CryptKey* - This key is used to encrypt the user's Items and is derived from a BaseKey in the Keyset and the ItemId using Blake2b.

*Keyset BaseKey* - The Keyset BaseKey is used to derive new Metadata and Item CryptKeys. The Keyset BaseKey is randomly generated when a user registers their account or when they change their password. Once a Keyset BaseKey is no longer in use, it is no longer used to derive the Metadata and Item CryptKeys, it is purged from the Keyset.

### Tokens
Vault uses randomly generated tokens as identifiers for all objects stored in the database. The tokens have a prefix that identifies the type of token it is. Each of the token types is defined below:

*UserToken* - A unique identifier for each user in the database. This token is associated with the username.

*AuthToken* - A unique identifier for the encrypted User object. The token is used as associated data when encrypting a User to strongly bind the identifier with the User it identifies.

This token is derived from the AuthKey and the UserToken. Since it is derived, it is not possible to associate a username with a User object without knowing the user's password, as long as there is more than one user in the database.

*KeysetToken* - A randomly generated unique identifier for a Keyset. The token is used as associated data when encrypting a Keyset to strongly bind the identifier with the Keyset it identifies. 

*MetadataToken* - A randomly generated unique identifier for a Metadata. The token is used as associated data when encrypting a Metadata to strongly bind the identifier with the Metadata it identifies. 

*ItemToken* - A randomly generated unique identifier for an Item. The token is used as associated data when encrypting an Item to strongly bind the identifier with the Item it identifies. 

*VersionToken* - A randomly generated unique identifier for cryptographic algorithm and BaseKey versions.


## Database Structure
Vault uses the BoltDB key/value store because it is simple, stable, and performant. Vault is designed so that many other backend databases could be used with relative ease.

### Buckets
BoltDB stores data in buckets and Vault users separate buckets for each type of data. The buckets in use are defined below:

*User* - This bucket holds the UserToken keyed on the username.
*Auth* - This bucket holds the encrypted User objects keyed on the derived AuthToken.
*Keyset* - This bucket holds the encrypted Keyset objects keyed on the KeysetId. All Keyset objects for all users are stored in this bucket.
*Metadata* - This bucket holds encrypted Metadata objects keyed on the MetadataId. All Metadata objects for all users are stored in this bucket.
*Item* - This bucket holds the encrypted Item objects keyed on the ItemId. All Items for all users are stored in this bucket.