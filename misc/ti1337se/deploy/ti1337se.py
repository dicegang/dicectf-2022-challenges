#!/usr/bin/env python3
import dis
import sys

banned = ["MAKE_FUNCTION", "CALL_FUNCTION", "CALL_FUNCTION_KW", "CALL_FUNCTION_EX"]

used_gift = False

def gift(target, name, value):
	global used_gift
	if used_gift: sys.exit(1)
	used_gift = True
	setattr(target, name, value)

print("Welcome to the TI-1337 Silver Edition. Enter your calculations below:")

math = input("> ")
if len(math) > 1337:
	print("Nobody needs that much math!")
	sys.exit(1)
code = compile(math, "<math>", "exec")

bytecode = list(code.co_code)
instructions = list(dis.get_instructions(code))
for i, inst in enumerate(instructions):
	if inst.is_jump_target:
		print("Math doesn't need control flow!")
		sys.exit(1)
	nextoffset = instructions[i+1].offset if i+1 < len(instructions) else len(bytecode)
	if inst.opname in banned:
		bytecode[inst.offset:instructions[i+1].offset] = [-1]*(instructions[i+1].offset-inst.offset)

names = list(code.co_names)
for i, name in enumerate(code.co_names):
	if "__" in name: names[i] = "$INVALID$"

code = code.replace(co_code=bytes(b for b in bytecode if b >= 0), co_names=tuple(names), co_stacksize=2**20)
v = {}
exec(code, {"__builtins__": {"gift": gift}}, v)
if v: print("\n".join(f"{name} = {val}" for name, val in v.items()))
else: print("No results stored.")
