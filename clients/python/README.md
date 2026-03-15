# HelixDB
> `python` package

This is the official Python client for the HelixDB project.

## Table of Contents
- [Usage Examples](#usage-examples)

## Usage Examples

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