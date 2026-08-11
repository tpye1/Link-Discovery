# Hello here is Link Discovery Layer for Linux (inspired by LDWin from Chris Hall)

Here is Chris Hall's github thanks to him for this kind of idea [https://github.com/chall32](click here)

## Hi so this kind of works so here is what is happening.
- There is no stacking switches just yet. 
- There is however a tui (Bubbletea).
- pkexec works for Linux.

### **IMPORTANT Regarding root priliviages:**
- There is now a helper function where pkexec is called for this, this should work for Linux operating systems and now the whole program doesn't need root.
- On Microsoft Windows running in this in the admin cmd/powershell is best since IPC for Windows is subideal.

# How to build and run:
### On Windows:
- Go into the root of the repository and run ``$ go build .`` 
- Then run in admin terminal (cmd or powershell): ``$ .\ldlinux.exe ``

### On Linux:
- Go into the root of the repository and run ``$ go build .``
- Then in the helper folder run ``$ go build .``
- Then go back to the root of the repository and run ``$ ./ldlinux``


### Contributors
- Contributors would be nice particularly those who know the Go language.

##### Regarding Ai use:
- Ai was used since I used this project to learn Go (which is probally not for me let's be honest).
- For the tui (Bubbletea).
- But just to clear things up a lot of the main.go wasn't Ai (only some) mainly just learning and these libaries are quite complicated i.e. the gopacket libaries.

# **IMPORTANT if there is a bug** 
- If there is a bug please or even hopefully not a security hole or a problem please or feedback please email me ASAP but please but respectful and polite: 
- My email: *tomspye1@gmail.com*
- I possibly won't see the issues on gitlab/github, email is better.

This is Bubbletea project written in Go.
