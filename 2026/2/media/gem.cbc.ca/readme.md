# gem.cbc.ca

~~~
url = https://gem.cbc.ca/joan/s01e01
monetization = ADS
count = 1
country = Canada
~~~

this works:

~~~
GET /media/validation/v2/?appCode=gem&idMedia=975359&output=json&tech=hls HTTP/1.1
Host: services.radio-canada.ca
x-claims-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJIYXNIRCI6IkZhbHNlIiwiVGllciI6Ik1lbWJlciIsIkhhc0FkcyI6IlRydWUiLCJSY0lkIjoiZmE4N2M1NDAtMDRlNy00NDYxLTk3N2EtYzhiZDQ3MGQxNDBhIiwiTWF4aW11bU51bWJlck9mU3RyZWFtcyI6IjMiLCJQYXltZW50VHlwZSI6IiIsIlBwaWQiOiI0ZWQyZTIxZGEyNGY5ODVjODY4MmIyMmVhNDAyMjQ1MDI4NDc4MTQyMDM1ZmFmNWM2OTYzNGU1YWJkZTVkMjNhIiwiZXhwIjoxNzI5NDUzNDk2fQ.q2O-ScQZfcGg5iN1uJrVXdfiaZbv4-OV8rvMNQS7bJw
~~~

but `dash` and `mpd` fail. also no DRM.
