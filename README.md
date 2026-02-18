# **HelixDB**
### *A lightweight, local‑first database engine built for modern applications.*

> **NOTICE:**<br>
> HelixDB is a side project and may not get updated very often since I have other projects that have a higher priority and this is closer to the bottom of the list however I will try to release updates at least every few months and contributions are welcome as long as they help fix a known issue or improve performance in some way. — **Nyxen**

HelixDB is an open‑source, file‑backed JSON database designed to be simple enough for local development yet powerful enough to support high‑traffic, production‑grade applications. It emphasizes reliability, corruption resistance, built‑in backup and recovery, and seamless integration with both Node.js and Python.

HelixDB runs as a single binary with zero configuration required — but offers deep customization through an optional `helixdb.config.json` file.

---

## **✨ Key Features**

- **Local‑first architecture** — runs anywhere with a single binary  
- **High‑traffic capable** — optimized write‑ahead logging and safe concurrency  
- **Corruption‑resistant** — checksums, WAL integrity, and safe commits  
- **Built‑in backups** — snapshot and incremental backup modes  
- **Automatic recovery** — WAL replay, checksum verification, and repair  
- **HTTP/JSON API** — language‑agnostic and easy to integrate  
- **Official Node.js & Python clients**  
- **Optional configuration file** — `helixdb.config.json`  
- **Open source & community‑driven**

---

# **📘 Table of Contents**

- [Project Goals](#project-goals)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [HTTP API Reference](#http-api-reference)
- [Client Libraries](#client-libraries)
- [Backup & Recovery](#backup--recovery)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

# **Project Goals**

HelixDB aims to bridge the gap between lightweight embedded databases and heavy enterprise systems.

### **Primary Objectives**
- Provide a **simple, local‑first** database that requires no external services.
- Deliver **enterprise‑grade reliability** through WAL, checksums, and safe writes.
- Offer **built‑in backup and recovery** without third‑party tools.
- Support **high‑traffic applications** with predictable performance.
- Maintain a **clean, approachable API** for developers.
- Remain **fully open source** and welcoming to contributors.

---

# **Installation**

### **Download the binary**
(Placeholder — replace with actual releases)

```
helixdb serve
```

HelixDB will start with default settings and create a `helix.db` file in the current directory.

---

# **Quick Start**

<details>
<summary><strong>Start the server</strong></summary>

```
helixdb serve
```

Default behavior:

- Listens on `http://127.0.0.1:7777`
- Stores data in `./helix.db`
- Uses WAL in `./wal/`
- Auto‑recovers on startup
</details>

<details>
<summary><strong>Create a document</strong></summary>

```
POST /collections/users
```

Body:

```json
{
  "data": {
    "username": "alice",
    "email": "alice@example.com"
  }
}
```
</details>

<details>
<summary><strong>Query documents</strong></summary>

```
POST /collections/users/query
```

```json
{
  "filter": { "username": "alice" },
  "limit": 10
}
```
</details>

---

# **Configuration**

HelixDB supports an optional `helixdb.config.json` file in your project root.

<details>
<summary><strong>Example helixdb.config.json</strong></summary>

```json
{
  "server": {
    "port": 7777,
    "host": "127.0.0.1"
  },
  "storage": {
    "dataFile": "./data/helix.db",
    "walDirectory": "./data/wal",
    "autoCompact": true,
    "compactThresholdMB": 128
  },
  "backup": {
    "enabled": false,
    "mode": "incremental",
    "directory": "./backups",
    "intervalMinutes": 30
  },
  "recovery": {
    "autoRecover": true,
    "verifyChecksums": true
  },
  "logging": {
    "level": "info",
    "file": "./logs/helixdb.log"
  },
  "security": {
    "requireAuth": false,
    "token": ""
  }
}
```
</details>

---

# **HTTP API Reference**

## **Create Document**
```
POST /collections/:name
```

## **Get Document**
```
GET /collections/:name/:id
```

## **Query Collection**
```
POST /collections/:name/query
```

## **Delete Document**
```
DELETE /collections/:name/:id
```

<details>
<summary><strong>Query Example</strong></summary>

```json
{
  "filter": { "status": "active" },
  "sort": [{ "field": "createdAt", "direction": "desc" }],
  "limit": 20
}
```
</details>

---

# **Client Libraries**

## **Node.js Example**

<details>
<summary><strong>Node.js Usage</strong></summary>

```js
import { HelixDB } from "@helixdb/client";

const db = new HelixDB("http://localhost:7777");

await db.collection("users").insert({
  username: "alice",
  email: "alice@example.com"
});

const users = await db.collection("users").query({
  filter: { username: "alice" }
});
```
</details>

## **Python Example**

<details>
<summary><strong>Python Usage</strong></summary>

```python
from helixdb import Client

db = Client("http://localhost:7777")

db.collection("users").insert({
    "username": "alice",
    "email": "alice@example.com"
})

users = db.collection("users").query({
    "filter": {"username": "alice"}
})
```
</details>

---

# **Backup & Recovery**

HelixDB includes built‑in backup and recovery mechanisms.

<details>
<summary><strong>Create a snapshot backup</strong></summary>

```
helixdb backup --to=./backups/snapshot-1
```
</details>

<details>
<summary><strong>Incremental backup</strong></summary>

```
helixdb backup --incremental --to=./backups/inc/
```
</details>

<details>
<summary><strong>Recover from backup</strong></summary>

```
helixdb recover --from=./backups/snapshot-1
```
</details>

---

# **Roadmap**

- [ ] v0.1 — Core engine, WAL, basic CRUD, HTTP API  
- [ ] v0.2 — Indexing, config file support  
- [ ] v0.3 — Backups & recovery  
- [ ] v0.4 — Node.js & Python clients  
- [ ] v1.0 — Production‑ready release  
- [ ] v2.0 — Read replicas, clustering, binary protocol  

---

# **Contributing**

HelixDB is open to contributors of all experience levels.

- Read `CONTRIBUTING.md`  
- Open issues for bugs or feature requests  
- Submit PRs with clear descriptions  
- Join discussions in GitHub Issues  

---

# **License**

MIT License — free for personal and commercial use.