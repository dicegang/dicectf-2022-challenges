; vim: ft=nasm

global prog

section .data

prog:
		incbin "build/obj/prog.bin"
