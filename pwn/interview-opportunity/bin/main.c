#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void sig_handler(int signum) {
  printf("You need to be faster to join us smh\n");
  exit(0);
}

int env_setup() {
  setvbuf(stdin, NULL, _IONBF, 0);
  setvbuf(stdout, NULL, _IONBF, 0);
  setvbuf(stderr, NULL, _IONBF, 0);
  signal(SIGALRM, sig_handler);
  alarm(15);
  return 0;
}

int main(int argc, char **argv) {
  char reason[10];

  env_setup();
  printf("Thank you for you interest in applying to DiceGang. We need great "
         "pwners like you to continue our traditions and competition against "
         "perfect blue.\n");
  printf("So tell us. Why should you join DiceGang?\n");
  read(0, reason, 70);

  puts("Hello: ");
  puts(reason);
}
