# BackBlaze

## Remote

1. Click App Keys
2. Click Generate New Master Application Key

## Local

1. `rclone config`
2. select New remote
3. name: remote
4. Storage: b2
5. account: `[keyId]`
6. key: `[applicationKey]`
7. `hard_delete`: false
8. Edit advanced config: No
9. Remote config: Yes this is OK
10. Current remotes: Quit config

## Download

~~~ps1
# "-P" is progress
rclone sync remote:<bucket> . -P
~~~

## Upload

~~~
rclone sync . remote:minerals -P
~~~

## References

- <https://backblaze.com>
- <https://rclone.org/b2>
