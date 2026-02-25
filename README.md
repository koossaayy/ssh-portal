# 🌀 ssh-portal

A glamorous, interactive SSH portal built with [Charm](https://charm.land)'s libraries.
Connect and explore: `ssh ssh.koossaayy.tn -p 2222`

```
  ██╗  ██╗ ██████╗  ██████╗ ███████╗███████╗ █████╗  █████╗ ██╗   ██╗██╗   ██╗
  ██║ ██╔╝██╔═══██╗██╔═══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗╚██╗ ██╔╝╚██╗ ██╔╝
  █████╔╝ ██║   ██║██║   ██║███████╗███████╗███████║███████║ ╚████╔╝  ╚████╔╝ 
  ██╔═██╗ ██║   ██║██║   ██║╚════██║╚════██║██╔══██║██╔══██║  ╚██╔╝    ╚██╔╝  
  ██║  ██╗╚██████╔╝╚██████╔╝███████║███████║██║  ██║██║  ██║   ██║      ██║   
  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝      ╚═╝  
```

## Features

- 👋 **Welcome / About** — Gorgeous banner + personal blurb
- 🚀 **Portfolio** — Browse your projects with tech stack tags
- 🖧  **Server Directory** — Wishlist-style SSH server menu
- 🐍 **Snake Game** — Full playable Snake with high score tracking

Built with:
- [`wish`](https://github.com/charmbracelet/wish) — SSH server framework  
- [`bubbletea`](https://github.com/charmbracelet/bubbletea) — TUI framework  
- [`lipgloss`](https://github.com/charmbracelet/lipgloss) — Terminal styling  
- [`bubbles`](https://github.com/charmbracelet/bubbles) — UI components  

---

## 🚀 Deploying to Coolify (Step by Step)

### Step 1 — Push to GitHub

```bash
git init
git add .
git commit -m "feat: initial ssh portal"
git remote add origin https://github.com/YOUR_USER/ssh-portal.git
git push -u origin main
```

### Step 2 — Create a new Resource in Coolify

1. Open your Coolify dashboard
2. Click **New Resource** → **Application**
3. Select **GitHub** (or public git URL)
4. Pick your `ssh-portal` repo

### Step 3 — Build settings

| Setting | Value |
|---|---|
| Build Pack | **Dockerfile** *(recommended)* or Nixpacks |
| Dockerfile path | `./Dockerfile` |
| Port | `2222` |
| Start command | *(leave empty, handled by entrypoint)* |

> If using **Nixpacks**, the `nixpacks.toml` is already configured. Set port to `2222`.

### Step 4 — Persistent Volume (IMPORTANT!)

The SSH host key must persist between deployments, otherwise every deploy will
rotate the key and users will get scary "host key changed" warnings.

In Coolify, go to **Storage** tab of your app:

| Host path | Container path |
|---|---|
| `/data/ssh-portal` | `/app/data` |

This ensures `/app/data/.ssh/id_ed25519` survives redeployments.

### Step 5 — Domain / Port Routing

SSH traffic is **not HTTP**, so you cannot use Coolify's normal reverse proxy (Traefik) for this.

**Option A — Direct port mapping (simplest):**
- In Coolify → **Network** tab, add port mapping: `2222:2222`
- Users connect with: `ssh ssh.koossaayy.tn -p 2222`

**Option B — Map to standard port 22:**
- Map `22:2222` on the host
- But make sure your server's own SSH isn't on port 22 (change it first!)
- Users connect with: `ssh ssh.koossaayy.tn` (no port needed)

### Step 6 — Deploy!

Click **Deploy**. Watch the build logs. Once green:

```bash
ssh ssh.koossaayy.tn -p 2222
```

🎉 You're in!

---

## ✏️ Customizing

### Change your about info
Edit `internal/ui/model.go` → find the `homeView()` function → update the about blurb text.

### Add portfolio projects
Edit `internal/portfolio/portfolio.go` → update the `projects` slice at the top.

### Add servers to the directory
Edit `internal/servers/servers.go` → update the `serverList` slice at the top.

### Change colors
All colors use Dracula palette by default. Edit the color variables at the top of each file — they're all `lipgloss.Color("#XXXXXX")` values.

---

## Running locally

```bash
go run .
# then in another terminal:
ssh localhost -p 2222
```

---

## License

MIT — do whatever you want with it. Star it if you like it! ⭐
