# Digital Ocean

https://duckdns.org

Install Go:

~~~
curl -L -O https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
tar -x -f go1.20.3.linux-amd64.tar.gz
~~~

Install server:

~~~
curl -L -O https://github.com/USER/REPO/archive/refs/heads/main.tar.gz
tar -x -f main.tar.gz
~~~

https://digitalocean.com

build first. if you try `go run`, it builds and runs from a temp folder, so
harder to kill:

~~~
go/bin/go build umber-main/hello.go
~~~

then run like this to survive closing the console:

~~~
nohup ./hello
~~~

then find it later with:

~~~
ps -d -f
~~~

- https://pubs.opengroup.org/onlinepubs/009695299/utilities/ps.html
- https://pubs.opengroup.org/onlinepubs/9699919799/utilities/nohup.html
