# Python

If you want to MITM a Python program, you might need to set one or more of
these:

~~~ps1
$env:HTTPS_PROXY = 'http://127.0.0.1:8080'
$env:REQUESTS_CA_BUNDLE = 'C:\Users\steven\.mitmproxy\mitmproxy-ca.pem'
$env:SSL_CERT_FILE = 'C:\Users\steven\.mitmproxy\mitmproxy-ca.pem'
~~~

## install

- https://bootstrap.pypa.io/get-pip.py
- https://packaging.python.org/en/latest/tutorials/installing-packages
- https://python.org/downloads/release/python-3119
- https://stackoverflow.com/questions/42666121/pip-with-embedded-python

delete:

~~~
python39._pth
~~~
