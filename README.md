# cryptopass

CryptoPass util uses public/private key crypto to provide issue/revoke/validate lifecycle for tokens that could be used as passwords.

Password/token management tool based on curve25519.

Username is the `firstPart.secondPart` and contains two parts.
Token generated only from second part.

* second part - is the dot separated second part of the username, if dot not present, the whole username would be used
* days - is the number of days starting since Jan 1, 2020 00:00:00 UFC 
* separator: `*` is for issue token, '/' is for revoke token.

#### Formulas
```
    Sign(secondPart + sep + days, PrivateKey)
    Verify(secondPart + sep + days, token, PublicKey) < Now   
```

#### Deployment
```
    PublicKey on server
    PrivateKey on token/password issue tool
```

### User Interface

Cryptopass has integrated user GUI interface for easy key and end-to-end encryption operation.



Generate token key-pair:
```
cryptopass.exe gen
```

You can export two system variables to bypass future prompts:
```
export CRYPTOPASS_PRIVATE_KEY
export CRYPTOPASS_PUBLIC_KEY
```

Issue token for username:
```
cryptopass.exe issue username YYY-MM-DD
```

Revoke token:
```
cryptopass.exe revoke username YYY-MM-DD
```

Revoke tokens must be stored on server and applied on authentication after successful verification.

Encrypting file for recipient:
```
cryptopass.exe enc input_file [output_file]
It will prompt recipient's public key
```
If output_file is not defined, it would be saved under input_file.cp

Decrypting received file:
```
cryptopass.exe dec input_file [output_file]
It will prompt your private key if CRYPTOPASS_PRIVATE_KEY env not defined
```
If output_file is not defined and ending with `.cp`, it would be saved under input_file[:-3]
