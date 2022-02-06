CC=gcc
CFLAGS=-O0 -I /usr/local
VICTIM_CPU=0
ATTACK_CPU=1

TARGETS=attacker victim

all: $(TARGETS)

%.o: %.c
	$(CC) $(CFLAGS)  -c $<

$(TARGETS): %:%.o
	$(CC) $(CFLAGS) $^ -o $@

run_victim:
	taskset -c $(VICTIM_CPU) ./victim

run_attacker: attacker
	taskset -c $(ATTACK_CPU) ./attacker

.PHONY:	clean

clean:
	$(RM) *.o $(HELPERS) $(TARGETS) 
