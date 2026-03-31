# fingerpaintfun

**fingerpaintfun** is a **finger painting app for kids (ages 3 to 5)** for GNOME, GTK 4, and Libadwaita optimized for a mobile user-experience.

## Other Files

* The **user guide** for **fingerpaintfun** is at: [GUIDE.md](GUIDE.md)
* THe **developer guide** for **fingerpaintfun** is at: [HACKING.md](HACKING.md)

## Build

### Requirements

- Go >= 1.21
- A C compiler (e.g., GCC) — required for CGo
- GTK 4 >= 4.10 (development headers)
- Libadwaita >= 1.4 (development headers)
- GLib 2.0 (development headers)
- GObject Introspection (development headers)

On Fedora:
```bash
sudo dnf install gcc gtk4-devel libadwaita-devel glib2-devel gobject-introspection-devel
```

On Debian/Ubuntu:
```bash
sudo apt install gcc libgtk-4-dev libadwaita-1-dev libglib2.0-dev libgirepository1.0-dev
```

### Development Build

```bash
go build
```

Or with vendored dependencies:
```bash
go build -mod=vendor
```

### Run

```bash
./fingerpaintfun
```

### Flatpak

```bash
flatpak install --user org.gnome.{Platform,Sdk}//47
flatpak install --user org.freedesktop.Sdk.Extension.golang//24.08
flatpak-builder --user --force-clean --install build build-aux/flatpak/link.reiver.fingerpaintfun.json
flatpak run link.reiver.fingerpaintfun
```

## Author

Software **fingerpaintfun** was written by [Charles Iliya Krempeaux](http://reiver.link)
