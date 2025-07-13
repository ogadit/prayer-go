# 🕌 NextPrayer

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Status](https://img.shields.io/badge/status-active-brightgreen)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blue)

> A fast and minimal command-line utility that shows Islamic prayer times and Hijri dates — right from your terminal.

---

## ✨ Features

- 🔁 Caches API response to reduce repeated requests
- 🕰️ Parses and displays only important prayer times
- 📅 Shows Hijri date (Islamic calendar)
- 💬 Smart edge-case handling (e.g. post-midnight, wraparound)
- ⚡ Built for speed — displays on terminal startup
- 🔒 Works offline (with cache)

---

## 📦 Installation

1. **Clone this repo**
   ```bash
   git clone https://github.com/yourusername/nextprayer.git
   cd nextprayer
   ```

2. **Build the binary**
   ```bash
   go build -o prayer-times
   ```

3. **Create your config**

   Save a file named `config.json` in the same folder:

   ```json
   {
     "city": "Karachi",
     "country": "Pakistan",
     "method": 2,
     "school": 1
   }
   ```

   - 🔗 [Method codes](https://aladhan.com/calculation-methods)
   - `school`: 0 = Shafi, 1 = Hanafi

---

## 🧪 Usage

Run it directly:

```bash
./prayer-times
```

**Sample Output:**

```
Asr Time
Maghrib in 2h 15m at 7:24 PM
16 Dhul-Hijjah 1446 A.H.
```

### 📌 Optional flag

```bash
./prayer-times -a
```

Shows only the next prayer and duration, for scripting or minimal output:

```
Maghrib in 2h 15m
```

---

## 🛠️ Terminal Integration

To show full prayer info on every terminal startup:

```zsh
# In ~/.zshrc

alias pt="~/path/to/prayer-times -a"

# Show on shell startup
(
  cd ~/path/to/nextprayer && ./prayer-times
)
```

---

## 🧪 Running Tests

Tests are written using Go’s `testing` package to cover edge cases:

- Before Fajr (e.g. 3 AM)
- Exactly at prayer times
- After Isha
- Wrap-around at midnight

```bash
go test
```

---

## 🔮 Roadmap

- ⏱️ Optional countdown timer mode
- 🌐 Offline fallback support
- 🧭 Location auto-detection (via IP or GPS)
- 📦 Package as a binary release

---

## 📜 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## ✍️ Author

Made with ❤️ and ☪️ by [Osama Anees](https://github.com/ogadit)
