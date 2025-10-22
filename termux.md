# installation
- install f-droid
- install termux


# openssh
- `pkg upgrade`
- `pkg install openssh`
- `passwd`  
    current password is prem
- start ssh server 
  - `sshd`
  - find logs `logcat -s 'sshd:*'`
- find user using command `whoami`
  - current user is `u0_a195`
- find ip address `ifconfig`


# access from laptop
- ssh -p 8022 u0_a195@192.168.1.12

# access storage
- run in phone only as it will prompt for access
- `termux-setup-storage`
  - allow storage access
- check symbolic links to all storage in ~/storage
  - `cd ~/storage`
  - `ls -l `
    ```
    ~/storage $ ls -l
    total 0
    lrwxrwxrwx. 1 u0_a195 u0_a195 26 Oct 22 19:57 dcim -> /storage/emulated/0/DCIM
    lrwxrwxrwx. 1 u0_a195 u0_a195 30 Oct 22 19:57 downloads -> /storage/emulated/0/Download
    lrwxrwxrwx. 1 u0_a195 u0_a195 50 Oct 22 19:57 external-1 -> /storage/37A2-1924/Android/data/com.termux/files
    lrwxrwxrwx. 1 u0_a195 u0_a195 30 Oct 22 19:57 movies -> /storage/emulated/0/Movies
    lrwxrwxrwx. 1 u0_a195 u0_a195 30 Oct 22 19:57 music -> /storage/emulated/0/Music
    lrwxrwxrwx. 1 u0_a195 u0_a195 30 Oct 22 19:57 pictures -> /storage/emulated/0/Pictures
    lrwxrwxrwx. 1 u0_a195 u0_a195 22 Oct 22 19:57 shared -> /storage/emulated/0
    ```


# termux API
- install termux API from f-droid
- install package termux-api-package
  `pkg install termux-api`

