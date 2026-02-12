# drm

https://xdaforums.com/t/dev-tools-simg2img-for-windows.3156459

1. download this zip: https://samfw.com/firmware/SM-N920C/XSG/N920CXXU5CVG2

2. open it in winrar and open the first large file within that file
   (`AP_N920CXXU5CVG2_CL11762721_QB54457305_REV00_user_low_ship_meta.tar.md5`)

3. extract the system.img file

4. run simg2img from here: <https://github.com/KinglyWayne/simg2img_win> like
   this: `simg2img system.img system.img.raw`

5. run ext2explore.exe and import the .raw file

6. go to /etc/security/.drm where the files are
