#!/usr/bin/env python3
import pickle
import sys
import io

class SafeUnpickler(pickle.Unpickler):
	def find_class(self, module, name):
		raise pickle.UnpicklingError(f"HACKING DETECTED")

data = sys.stdin.buffer.readline()+b"dice{buh2Qdj0219}"
SafeUnpickler(io.BytesIO(data)).load()
