import hashlib
import hmac


def verify_webhook(payload: bytes, signature: str, secret: str) -> bool:
    expected = "sha256=" + hmac.new(secret.encode(), payload, hashlib.sha256).hexdigest()
    return hmac.compare_digest(expected, signature)


if __name__ == "__main__":
    ok = verify_webhook(b'{"event":"payout.settled"}', "sha256=...", "your-secret")
    print("valid:", ok)
