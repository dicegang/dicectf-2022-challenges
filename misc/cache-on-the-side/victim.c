#include "fake_flag.h"
#include <sys/mman.h>
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdbool.h>

#define WAYS_MAX 16
#define LOOP_MAX 1<<21

char * get_huge_page() {
	void *ptr = NULL ;
	int hugepage_size = 1 << 21;
	if ( posix_memalign(&ptr , hugepage_size , hugepage_size ) ) {
		perror("posix_memalign") ;
	}
	madvise(ptr, hugepage_size, MADV_HUGEPAGE);
	return (char *)ptr;
}


int main(int argc, char **argv)
{
	char * addr_base = get_huge_page();

	char c;
	char * ptr;
	int i = 0;
	int loops = 0;
	int way = 0;

	while (true) {

		for (i = 0; i < flag_len; i ++) {
			c = flag[i];
			printf("%c\n", c);

			ptr = addr_base + c*64;

			for (loops = 0; loops < LOOP_MAX; loops++) {
				for (way = 0; way<WAYS_MAX; way++)
					(*(ptr+ (way*1024*64)))++;
			}
		}
	}
}
