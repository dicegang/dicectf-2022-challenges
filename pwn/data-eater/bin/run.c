#include <stdio.h>
#include <string.h>

char buf[32];

int main (int argc, char *argv[], char *envp[]) {
  while (*argv)
    *argv++ = 0;
  argv = 0;
  while (*envp)
    *envp++ = 0;
  envp = 0;

  char fmt[8];
  fgets(fmt, sizeof(fmt), stdin);
  scanf(fmt, buf);
  memset(buf, 0, sizeof(buf));

  *(char *)(0) = 0;
}
