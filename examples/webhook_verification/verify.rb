require 'openssl'

def verify_webhook(payload, signature, secret)
  expected = 'sha256=' + OpenSSL::HMAC.hexdigest('SHA256', secret, payload)
  ActiveSupport::SecurityUtils.secure_compare(expected, signature)
end
