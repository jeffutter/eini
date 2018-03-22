# eini

[![Build Status](https://travis-ci.org/jeffutter/eini.svg?branch=master)](https://travis-ci.org/jeffutter/eini)

`eini` is a utility for managing a collection of secrets in source control. The
secrets are encrypted using [public
key](http://en.wikipedia.org/wiki/Public-key_cryptography), [elliptic
curve](http://en.wikipedia.org/wiki/Elliptic_curve_cryptography) cryptography
([NaCl](http://nacl.cr.yp.to/) [Box](http://nacl.cr.yp.to/box.html):
[Curve25519](http://en.wikipedia.org/wiki/Curve25519) +
[Salsa20](http://en.wikipedia.org/wiki/Salsa20) +
[Poly1305-AES](http://en.wikipedia.org/wiki/Poly1305-AES)). Secrets are
collected in an ini file, in which all the string values are encrypted by default. Public
keys are embedded in the file, and the decrypter looks up the corresponding
private key from its local filesystem or stdin.

`eini` is based on [ejson](http://github.com/shopify/ejson) using the same encryption code, while offering a more human-readable ini format and optional unencrypted values.

The main benefits provided by `eini` are:

* Secrets can be safely stored in a git repo.
* Changes to secrets are auditable on a line-by-line basis with `git blame`.
* Anyone with git commit access has access to write new secrets.
* Decryption access can easily be locked down to production servers only.
* Secrets change synchronously with application source (as opposed to secrets
  provisioned by Configuration Management).
* Simple, well-tested, easily-auditable source.

## Installation

You can download binaries from [Github Releases](https://github.com/jeffutter/eini/releases)

## Workflow

### 1: Create the Keydir

By default, EJSON looks for keys in `/opt/ejson/keys`. You can change this by
setting `EJSON_KEYDIR` or passing the `-keydir` option.

```
$ mkdir -p /opt/ejson/keys
```

### 2: Generate a keypair

When called with `-w`, `eini keygen` will write the keypair into the `keydir`
and print the public key. Without `-w`, it will print both keys to stdout. This
is useful if you have to distribute the key to multiple servers via
configuration management, etc.

```
$ eini keygen
Public Key:
63ccf05a9492e68e12eeb1c705888aebdcc0080af7e594fc402beb24cce9d14f
Private Key:
75b80b4a693156eb435f4ed2fe397e583f461f09fd99ec2bd1bdef0a56cf6e64
```

```
$ ./eini keygen -w
53393332c6c7c474af603c078f5696c8fe16677a09a711bba299a6c1c1676a59
$ cat /opt/ejson/keys/5339*
888a4291bef9135729357b8c70e5a62b0bbe104a679d829cdbe56d46a4481aaf
```

### 3: Create an `eini` file

The format is described in more detail [later on](#format). For now, create a
file that looks something like this. Fill in the `<key>` with whatever you got
back in step 2.

Create this file as `test.eini`:

```ini
_public_key = <key>
database_password = 1234password
```

### 4: Encrypt the file

Running `eini encrypt test.eini` will encrypt any new plaintext keys in the
file, and leave any existing encrypted keys untouched:

```ini
_public_key = 63ccf05a9492e68e12eeb1c705888aebdcc0080af7e594fc402beb24cce9d14f
database_password = EJ[1:WGj2t4znULHT1IRveMEdvvNXqZzNBNMsJ5iZVy6Dvxs=:kA6ekF8ViYR5ZLeSmMXWsdLfWr7wn9qS:fcHQtdt6nqcNOXa97/M278RX6w==]
```

Try adding another plaintext secret to the file and run `eini encrypt
test.eini` again. The `database_password` field will not be changed, but the
new secret will be encrypted.

### 5: Decrypt the file

To decrypt the file, you must have a file present in the `keydir` whose name is
the 64-byte hex-encoded public key exactly as embedded in the `eini` document.
The contents of that file must be the similarly-encoded private key. If you used
`eini keygen -w`, you've already got this covered.

Unlike `eini encrypt`, which overwrites the specified files, `eini decrypt`
only takes one file parameter, and prints the output to `stdout`:

```
$ eini decrypt foo.eini
_public_key = 63ccf05a9492e68e12eeb1c705888aebdcc0080af7e594fc402beb24cce9d14f
database_password = 1234password
```

## Format

The `eini` document format is simple, but there are a few points to be aware
of:

1. It's just an ini file.
2. There *must* be a key at the top level named `_public_key`, whose value is a
   32-byte hex-encoded (i.e. 64 ASCII byte) public key as generated by `eini
   keygen`.
3. Any string literal that isn't an object key will be encrypted by default (ie.
   in `a = b`, `"b"` will be encrypted, but `"a"` will not.
4. If a key has a `#decrypted` comment at the end of the line or on the preceding
   line, its corresponding value will not be encrypted.
