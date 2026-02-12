# shield

havent gotten to the part of adding mitm cert etc yet

https://florisse.nl/shield-downgrade/

but used this to downgrade firmware

https://xdaforums.com/t/tool-bootmod-for-flashing-root-and-stock-images-automatic-apk-installation-apk-removal-file-pushing-file-pulling-and-more-all-models.4524873/

this tool option 0 then option 1

~~~
mkdir -p -m 700 /data/local/tmp/ca-copy
cp /system/etc/security/cacerts/* /data/local/tmp/ca-copy/
mount -t tmpfs tmpfs /system/etc/security/cacerts
mv /data/local/tmp/ca-copy/* /system/etc/security/cacerts/
mv /data/local/tmp/xxx.0 /system/etc/security/cacerts/
chown root:root /system/etc/security/cacerts/*
chmod 644 /system/etc/security/cacerts/*
chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*
~~~
