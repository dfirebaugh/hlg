# ggez

ggez is a simple game library built on top of SDL2.

(this project is a WIP -- i'm playing around with dynamically linking system libraries)


## Setup

1. install ggez
```bash
go get github.com/dfirebaugh/ggez
```

2. Install SDL2 (see install instructions below)
3. build your project

### Examples
check the `./examples` dir for some basic examples

### Installing SDL2 System Libraries

#### **Windows**:

1. **Using vcpkg (Recommended)**:
   - If you haven't installed `vcpkg`, follow the instructions on [their GitHub page](https://github.com/microsoft/vcpkg).
   - Once `vcpkg` is installed:
     ```bash
     vcpkg install sdl2:x64-windows
     ```

2. **Manual installation**:
   - Go to the [SDL2 download page](https://libsdl.org/download-2.0.php) and download the SDL2 development library for Visual C++.
   - Extract the downloaded ZIP file.
   - Copy the contents of the `lib/x64` (or `lib/x86` for 32-bit applications) directory to your project's directory.
   - Copy `SDL2.dll` from the extracted `x64` (or `x86`) directory to your project's executable directory.

#### **macOS**:

1. **Using Homebrew (Recommended)**:
```bash
brew install sdl2
```

#### **Linux (Ubuntu/Debian)**:

1. **Using APT (Recommended)**:
```bash
sudo apt update
sudo apt install libsdl2-dev
```

#### **Linux (Fedora)**:

1. **Using DNF**:
```bash
sudo dnf install SDL2-devel
```

#### **Linux (Arch Linux)**:

1. **Using Pacman**:
```bash
sudo pacman -S sdl2
```

