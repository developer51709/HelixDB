# HelixDB
> `node.js` package

This is the official Node.js client for the HelixDB project.

## Table of Contents
- [Usage Examples](#usage-examples)

## Usage Examples

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