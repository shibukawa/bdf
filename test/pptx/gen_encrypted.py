"""Generates the encrypted PowerPoint test decks in
converter/internal/offcrypto/testdata with msoffcrypto-tool.

Usage: python3 test/pptx/gen_encrypted.py   (pip install msoffcrypto-tool)

Both decks are converter/pptx/testdata/basic.pptx encrypted with PASSWORD:
agile.pptx with Agile encryption (AES-256, SHA-512, as Office 2010 and later
write it) and standard.pptx with Standard encryption (AES-128, SHA-1, as
Office 2007 writes it). msoffcrypto-tool writes Agile files itself; the
Standard file is put together here from its key derivation and compound file
writer. Each file is checked by decrypting it again with msoffcrypto-tool.
"""
import io
import os
import struct
from hashlib import sha1

import msoffcrypto
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from msoffcrypto.format.ooxml import OOXMLFile
from msoffcrypto.method.container.ecma376_encrypted import ECMA376Encrypted
from msoffcrypto.method.ecma376_standard import ECMA376Standard

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
SRC = os.path.join(ROOT, "converter", "pptx", "testdata", "basic.pptx")
OUT = os.path.join(ROOT, "converter", "internal", "offcrypto", "testdata")

# Non-ASCII, with a character outside the BMP: the key derivations hash the
# password as UTF-16LE, surrogate pairs included.
PASSWORD = "パスワード🔑bdf"


def check(data, plain):
    f = msoffcrypto.OfficeFile(io.BytesIO(data))
    f.load_key(password=PASSWORD)
    out = io.BytesIO()
    f.decrypt(out, verify_integrity=True) if f.type == "agile" else f.decrypt(out)
    assert out.getvalue() == plain, "decrypted package differs"


def agile(plain):
    out = io.BytesIO()
    OOXMLFile(io.BytesIO(plain)).encrypt(PASSWORD, out)
    return out.getvalue()


def standard(plain):
    alg_id, alg_id_hash, provider, key_bits = 0x660E, 0x8004, 0x18, 128  # AES-128, SHA-1, PROV_RSA_AES
    salt = os.urandom(16)
    key = ECMA376Standard.makekey_from_password(PASSWORD, alg_id, alg_id_hash, provider, key_bits, 16, salt)

    def ecb(data):
        e = Cipher(algorithms.AES(key), modes.ECB()).encryptor()
        return e.update(data) + e.finalize()

    flags = 0x24  # fCryptoAPI | fAES
    csp = "Microsoft Enhanced RSA and AES Cryptographic Provider\0".encode("utf-16le")
    header = struct.pack("<8I", flags, 0, alg_id, alg_id_hash, key_bits, provider, 0, 0) + csp
    verifier = os.urandom(16)
    verifier_hash = sha1(verifier).digest() + b"\0" * 12
    info = (
        struct.pack("<HHII", 3, 2, flags, len(header))
        + header
        + struct.pack("<I", 16) + salt + ecb(verifier)
        + struct.pack("<I", 20) + ecb(verifier_hash)
    )
    padded = plain + b"\0" * (-len(plain) % 16)
    package = struct.pack("<Q", len(plain)) + ecb(padded)
    out = io.BytesIO()
    ECMA376Encrypted(package, info).write_to(out)
    return out.getvalue()


def main():
    plain = open(SRC, "rb").read()
    os.makedirs(OUT, exist_ok=True)
    for name, encrypt in (("agile.pptx", agile), ("standard.pptx", standard)):
        data = encrypt(plain)
        check(data, plain)
        with open(os.path.join(OUT, name), "wb") as f:
            f.write(data)
        print(name, len(data), "bytes")


if __name__ == "__main__":
    main()
