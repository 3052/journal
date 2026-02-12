# play.movistar.com.ec

~~~
url = https://play.movistar.com.ec/details/movie/18547017
monetization = FLATRATE
country = Ecuador
~~~

site is just a piece of shit. first it has a broken redirect:

~~~
> curl -v -m 9 https://play.movistar.com.ec/details/movie/18547017
*   Trying 213.140.60.178:443...
* Connection timed out after 9015 milliseconds
~~~

to here:

https://movistar.com.ec

which returns:

> Warning: Potential Security Risk Ahead

if you continue you get multiple errors, even with new Firefox:

~~~
TypeError: Error.captureStackTrace is not a function

A client-side exception has occurred, see here for more info:
https://nextjs.org/docs/messages/client-side-exception-occurred

Uncaught (in promise) TypeError: t is undefined

Uncaught TypeError: window.location.referrer is undefined
~~~
