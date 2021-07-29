# Cryptographic usages

The goal of this document is to explain the usage of cryptographic functions in
plain English, and avoid the use of too many technical terms.

We will first explain which functions are used, and after that we will describe
our exact use-cases for them.

<!-- markdown-toc start - Don't edit this section. -->
**Table of Contents**

- [Cryptographic terms](#cryptographic-terms)
    - [Hashing/Hash](#hashinghash)
    - [Encryption](#encryption)
- [Cryptographic functions](#cryptographic-functions)
    - [JWT](#jwt)
    - [HMAC](#hmac)
    - [SHA-512](#sha-512)
    - [Chacha20poly1305](#chacha20poly1305)
    - [Chacha20](#chacha20)
    - [Poly1305](#poly1305)
- [Cryptographic usage](#cryptographic-usage)
    - [User authentication](#user-authentication)
    - [Register](#register)
    - [Password reset](#password-reset)
    - [Stats collection](#stats-collection)
    - [OAuth](#oauth)
- [Questions](#questions)

<!-- markdown-toc end -->


## Cryptographic terms

### Hashing/Hash

Hashing is the process of a (cryptographic) hashing function that takes an
input, and returns a long stream of characters (string) which cannot be reversed
or linked back to the original input. This output is called a "hash".

### Encryption

Encryption is a term used to describe a process where a given input and a given
key is turned into a unique string. However, unlike the hashing process, this
unique string can be linked/reversed back to the original input if you have the
given key where this input was encrypted.


## Cryptographic functions

### JWT

JSON Web Token("JWT") is our primary way to authenticate a user with the server
and have confidence that it isn't a malicious attack.

The confidence is guaranteed by the built-in [HMAC](#HMAC) signature which is
mandatory. The HMAC is made out of a secret key (which only the server operators
know), the payload (which has information regarding the user like the UserID),
and the header (which contains information about the payload).

The header + payload are base64'ed (separated from each other). And with the
secret key thrown into the HMAC function, which after that is being hashed by
[SHA- 512](#SHA-512) function:

```dart
HMAC_SHA512(
  secret,
  base64urlEncoding(header) + '.' +
  base64urlEncoding(payload)
)
```

### HMAC

Hash-based message authentication code("HMAC") is a type of message authentication code("MAC"), which involves a cryptographic hash function and a secret cryptographic key.

Hereby the server can create HMAC with the specific secret and later-on verify with the same secret that the message is authentic. Thanks to this, we have the guarantee it's authentic and can safely be assumed it's safe.

We currently only use the cryptographic hash function [SHA-512](#SHA-512) as hash function for the HMAC generation.

### SHA-512

Secure Hash Algorithm("SHA") is a family of secure hash functions made by National Security Agency("NSA"). In particular, SHA-512 is a member of the SHA-2 series.

The SHA-2 series is a well known standard and seen in real-life application for hashing.

### Chacha20poly1305

Chacha20poly1305 is an authenticated encryption based of the cryptographic encryption function [chacha20](#chacha20) and the message authentication code [poly1305](#poly1305).

We use a slight derivation of Chacha20poly1305, the Chacha20poly1305x function.
It allows us to add more "nonce" into the string. The nonce is forged with the output.
Due to this you can still have the same input, but with a different nonce it will be a completely different string.
It will make it harder to brute-force the end result.

### Chacha20

ChaCha is a family of stream ciphers. Whereby ChaCha20 the secured member is.

Chacha20 is faster than the well-know AES standard, which is one of the reasons is was chosen.
Hereby it also comes in the [Chacha20poly1305](#chacha20poly1305) form, which has built-in message authentication code to authenticate the data.
It takes 20-rounds of encrypting the data, unlike the 2 other members which do it in less rounds.

### Poly1305

Poly1305 is a message authentication code("MAC").

The function takes an one-time key and a message as its input to produce a 16-characters long tag. 
The tag is then used to authenticate the given message.


## Cryptographic usage

### User authentication

A "cookie" is stored on the user's device which contains a [JWT](#JWT).
This is used on the server to know if this device is authenticated with a specific user account.

### Register

When someone wants to register a new account (not via OAuth).
An email is send to the user, it contains a [Chacha20poly1305](#chacha20poly1305) encrypted [JWT](#JWT).
The JWT contains simple information:
- username.
- password.
- email.

The JWT is signed by an unique key which is only used for this action and the reset password action.
Thereby the Chacha20poly1305 is given 24 length nonce, 
and that nonce is after the encryption scrambled into the output. That provides another layer of security.

### Password reset

It's somewhat similar to the [register](#register) process.
An email is send to the user, it contains a [Chacha20poly1305](#chacha20poly1305) encrypted [JWT](#JWT).
The JWT contains simple information:
- email.

The JWT is signed by an unique key which is only used for this action and the register action.
Thereby the Chacha20poly1305 is given 24 length nonce, 
and that nonce is after the encryption scrambled into the output. That provides another layer of security..

### Stats collection

We collect minimal stats for userstyle installs and views.
To ensure this privacy-sensitive data isn't stored in plain text and still usable,
we use the [HMAC](#HMAC) function to get a hash that is then stored in the database.

The process is along the lines:
Get a combined output of the User's IP and the given Style's ID

```dart
record = "IP of the user" + " " + "ID of the style"
HMAC_SHA512(
    secret,
    record
)
```

The secret used for this is only used for this action, and is a long pseudo-random string.

### OAuth

We allow third-party's to register OAuth's so they can do certain actions based on the user's behalf.
The third-party's receive a [JWT](#JWT) token with minimal information.

Information whereby it's linked to a specific style:
- specific style's ID.
- user's ID.

Information whereby it's linked to a user:
- scopes.
- user's ID.

The JWT is signed by a long pseudo-random string which is only used for this action.


## Questions

If you still have any questions after reading this document, please don't hesitate to contact [security@userstyles.world](mailto:security@userstyles.world).
