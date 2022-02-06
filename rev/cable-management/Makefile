CC ?= gcc
CFLAGS := -std=c11 -Wall -Wextra -O3 $(CFLAGS)
LDFLAGS := -Wl,-s
# CFLAGS = -Wall -Wextra -O3
# LDFLAGS = -Wl,-s 
BUILD=build
BINS=bin
OBJ=obj
SRC=src
PROG=prog

VERBOSITY=RELEASE
# RELEASE, DEBUG, DEBUGV, DEBUGVV
CFLAGS += -D$(VERBOSITY)

default: mkdirs prog.bin chall

mkdirs:
	mkdir -p $(BUILD)/$(BINS)
	mkdir -p $(BUILD)/$(OBJ)


chall: chall.o bin.o
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(BUILD)/$(BINS)/$@ $(addprefix $(BUILD)/$(OBJ)/,$^)

%.o: $(SRC)/%.c $(HEADERS)
	$(CC) $(CFLAGS) -c -o $(BUILD)/$(OBJ)/$@ $<

#prog.bin: $(SRC)/prog.s
#	nasm $(NASMFLAGS) -o $(BUILD)/$(OBJ)/$@ $<
prog.bin: $(SRC)/$(PROG)/prog.py
	python3 $< $(BUILD)/$(OBJ)/$@

%.o: $(SRC)/%.s prog.bin
	nasm $(NASMFLAGS) -f elf64 -o $(BUILD)/$(OBJ)/$@ $<

