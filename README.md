# Hello here is Link Discovery Layer for Linux (inspired by LDWin from Chris Hall)

Here is Chris Hall's github thanks to him for this kind of idea [https://github.com/chall32](click here)

## Hi so this kind of works so here is what is happening.
- There is no stacking switches just yet. 
- There is however a gui.

### **IMPORTANT Regarding root priliviages:**
- There is no way to get this to give a sudo/admin/root prompt just yet.
- For security only use this on an Admin device/user
- On Linux you can use the command of:
```shell
sudo setcap cap_net_raw+ep your-helper
```
- Or you can use
```shell
sudo ./ldlinux
```
- If anyone has a better way of doing this please let me know I have something in mind but yeah
- On Windows I think you can Run as admin
- The issue is there is only one Go function that needs admin priliviages and the proposals I have wrriten are not ideal and not the most secure since the whole app at the moment will have root.
- I am working on an pkexec for now when this will get installed I will work on an RPM and Debian based installer at some point


### Contributors would be nice:
- Is anyone is good a frontend and actually likes that cool.

#### Windows users:
- If you are using Windows
- You need a dependency of Npcap which is annoying (a way to install this is to have wireshark or you can just have it)
- A Windows binary makes things a little harder but I will hopefully get there



##### Regarding Ai use:
- Ai was used since I used this project to learn Go and for the Gui in typescript, yeah I don't like typescript or frontend.
- The Json nonsense in dto.go is Ai so yeah if bugs let me know as well as the frontend.


# **IMPORTANT if there is a bug** 
- If there is a bug please or even hopefully not a security hole or a problem please or feedback pleas email me ASAP but please but respectful and polite: 
- My email: *tomspye1@gmail.com*
- I won't see the issues on gitlab/github, email is better

This is a wails project for the frontend, written in Go.
