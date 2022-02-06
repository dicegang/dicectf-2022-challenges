#include<stdio.h>
#include<stdint.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>

#define u32 uint32_t
#define u16 uint16_t
#define u8 uint8_t
#define BUFFERSIZE (2324)
#define PRINT_SIZE (84)

extern u8 prog[];

u8 input_buffer[] = {2, 2, 2, 2, 2, 2, 2, 2};

typedef enum {
		O_ELEC = 0xec,
		O_COND = 0xcd,
		O_TAIL = 0xea,
		O_INP = 0x11,
		O_ACC = 0x01,
		O_HAL = 0x80,
} CELLS;

#if defined(DEBUG) || defined(DEBUGV) || defined(DEBUGVV)
#define BLK "\e[0;30m\e[40m"
#define RED "\e[0;31m\e[41m"
#define GRN "\e[0;32m\e[42m"
#define YEL "\e[0;33m\e[43m"
#define BLU "\e[0;34m\e[44m"
#define MAG "\e[0;35m\e[45m"
#define CYN "\e[0;36m\e[46m"
#define CRESET "\e[0m"

void pretty_print(u8* buffer, int size) {
				for (int spec = 0; spec < size; spec++) {
						for (int j = 0; j < size; j++) {
								switch (buffer[spec * BUFFERSIZE + j]) {
										case O_ELEC:
												printf(RED);
												printf("..");
												break;
										case O_TAIL:
												printf(BLU);
												printf("..");
												break;
										case O_COND:
												printf(YEL);
												printf("..");
												break;
										case O_INP:
												printf(MAG);
												printf("..");
												break;
										case O_ACC:
												printf(GRN);
												printf("..");
												break;
										case O_HAL:
												printf(CYN);
												printf("..");
												break;
										case 0:
												printf(BLK);
												printf("..");
												break;
												
								}
						}

						printf("\n");
				}
				printf(CRESET);
				printf("\n");
				usleep(100000);
}
#endif


__attribute__((always_inline)) static inline
int sine(int value) {
		return value * (value - 4) * (value - 8);
}

__attribute__((always_inline)) static inline
int cosine(int value) {
		return -(value - 2) * (value - 6) * (value - 10);
}

__attribute__((always_inline)) static inline
int fix_zero(int v) {
		return v / (abs(v) + 1/(v*v + 1));
}


bool moore(int position, const u8* buf) {
		int electron_count = 0;
		int pos_x = position % BUFFERSIZE;
		int pos_y = position / BUFFERSIZE;
#ifdef DEBUGVV
		printf("MOORE: ind - %d x - %d y - %d\n", position, pos_x, pos_y);
#endif
		for (int i = 0; i < 8; i++) {

				int y = pos_y + fix_zero(sine(i));
				int x = pos_x + fix_zero(cosine(i));

#ifdef DEBUGVV
				printf("MOORE: %d: (%d, %d) -> (%d, %d) {%d}\n", i, fix_zero(sine(i)), fix_zero(cosine(i)), x, y, electron_count);
				
#endif
				if (!(x < 0 || y < 0 || x >= BUFFERSIZE || y >= BUFFERSIZE)) {
#ifdef DEBUGVV
						printf("MOORE: !! %02x\n", buf[BUFFERSIZE * y + x]);
#endif
						if (buf[BUFFERSIZE * y + x] == O_ELEC) {
								electron_count += 1;
						}
				}
		}
		if (electron_count == 1 || electron_count == 2) {
				return true;
		}
		return false;
}


typedef struct {
		bool (*const input_handler)();
} io_t;


int vm(io_t* io) {
		u8 prog_dup[BUFFERSIZE * BUFFERSIZE];
		memcpy(prog_dup, prog, BUFFERSIZE * BUFFERSIZE);
		u8* new = prog;
		u8* buffer = prog_dup;
		for (;;) {

#if defined(DEBUG) || defined(DEBUGV) || defined(DEBUGVV)
				pretty_print(buffer, PRINT_SIZE);
#endif
				for (int i = 0; i < BUFFERSIZE * BUFFERSIZE; i++) {
						switch (buffer[i]) {

								case O_ELEC:
										new[i] = O_TAIL;
										break;

								case O_TAIL:
										new[i] = O_COND;
										break;

								case O_COND:
										if (moore(i, buffer)) {
												new[i] = O_ELEC;
										} else {
												new[i] = O_COND;
										}

										break;

								case O_INP:
										bool input_char = io->input_handler();
										if (input_char) {
												new[i] = O_ELEC;
										} else {
												new[i] = O_COND;
										}
#if defined(DEBUGV) || defined(DEBUGVV)
										printf("INPUT AFTER: [ ");
										for (int c = 0; c < 8; c++) {
												printf("%d ", input_buffer[c]);
										}
										printf("]\n");
#endif
										break;

								case O_ACC:
										if (moore(i, buffer)) {
												return 1;
										}
										break;

								case O_HAL:
										if (moore(i, buffer)) {
												return 0;
										}
										break;

								default:
										new[i] = 0;
										break;

						}
				}
				u8* tmp = buffer;
				buffer = new;
				new = tmp;
		}
}

bool input() {
		if (input_buffer[7] == 2) {
				u8 character = getchar();
				for (int i = 0; i < 8; i++) {
						input_buffer[i] = (bool)((character >> (7 - i)) & 0x1);
				}

#if defined(DEBUGV) || defined(DEBUGVV)
				printf("INPUT BEFORE: [ ");
				for (int i = 0; i < 8; i++) {
						printf("%d ", input_buffer[i]);
				}
				printf("]\n");
#endif
				bool ret = (bool)input_buffer[0];
				input_buffer[0] = 2;
				return ret;
		}
		else {
				for (int i = 0; i < 8; i++) {
						if (input_buffer[i] != 2) {
								bool ret = (bool)input_buffer[i];
								input_buffer[i] = 2;
								return ret;
						}
				}
		}
		return false;
}


int main() {
		setvbuf(stdout, NULL, _IOLBF, BUFSIZ);
		io_t io = {
				.input_handler = input,
		};

		int ret = vm(&io);
		if (ret == 0x1) {
				printf(":)\n");
		} else {
				printf(":(\n");
		}
}

