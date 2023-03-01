# Raven
Faily undetectable Linux Spyware
![img](https://github.com/nemyzyxt/Raven/)

## DISCLAIMER :
I will not be responsible for any damage that may arise out of unethical use of this project . Have fun :)

## Capabilities :
- Establish persistence
- Capture screenshot of victim machine
- Record the screen of the victim machine
- Take a picture through the webcam
- Record a video through the webcam
- Log keystrokes on victim machine

## Set-Up :
1. Clone this repository
2. Navigate to the project directory
3. Feel free to modify C2 Address in the raven.go file to point to your C2
4. Run the command ```go build -ldflags="-s -w" delta.go``` to create the executable(ELF) file