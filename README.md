# AutoSQL
AutoSQL is a mitmproxy framework written in Go language, with penatration testing abilities.

## Usage
After setting it as your proxy, you can just hang around on the Internet. AutoSQL will
automatically attack every URL you've accessed in the background, by submitting them to
a standalone SQLmap backend.

## Important Features
* Simple, efficient, easy to extend, with the power of Go libraries.
* Allow HTTPS decryption with fake CA certificate.
* Packed with some SQLmap APIs wrapped in Go language.

## Pending Improvements
* Support multiple SQLmap backends to attack concurrently.
* Parse HTML tags to attack URLs which user haven't accessed.
