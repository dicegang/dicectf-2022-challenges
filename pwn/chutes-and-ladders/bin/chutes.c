#include <stdio.h>
#include <time.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <unistd.h>

typedef struct Slot {
    char * markers;
    int playersHere;       // bitwise of who is here
    char idx;
} Slot;

typedef struct Player {
    char * name;
    char idx;
    char boardpos;
} Player;


Slot* BOARD[100];
Player* PLAYERS[10];

int AMOUNT_PLAYERS = -1;
int CURR_TURN = 0;

int LADDERS[5][2] = {
    {27,35},
    {1,5},
    {64,98},
    {15,26},
    {30, 33},
};

int CHUTES[5][2] = {
    {93,90},
    {24,14},
    {68,3},
    {38,36},
    {6, 2}
};

void genTable(){
    for(int i = 0; i < 100; ++i){
        BOARD[i] = malloc(sizeof(Slot));

        BOARD[i]->idx = i;
        BOARD[i]->playersHere = 0;
        BOARD[i]->markers = NULL;
    }
}


void printSlot(Slot * s){
    char num = s->idx + 1;

    printf("%s|%10d|\e[0m", num & 1 ? "\e[0;37m" : "\e[0;33m", num);       // \e[0m
}

void printMarkers(Slot * s){
    char num = s->idx + 1;

    printf("%s|%10s|\e[0m", num & 1 ? "\e[0;37m" : "\e[0;33m", s->playersHere == 0 ? "" : s->markers);       // \e[0m
}


void viewTable(){
    for (int row = 4; row >= 0; --row){
        // ----- row 1
        for (int i = 0; i < 10; ++i){
            printf("%s------------\e[0m", i & 1 ? "\e[0;37m" : "\e[0;33m", i);
        }
        puts("");

        for (int i = 0; i < 10; ++i){
            printf("%s|          |\e[0m", i  & 1 ? "\e[0;37m" : "\e[0;33m", i);
        }
        puts("");

        for (int j = 19; j >= 10; --j){
            printSlot(BOARD[row*20 + j]);
        }
        puts("");

        for (int j = 19; j >= 10; --j){
            printMarkers(BOARD[row*20 + j]);
        }
        puts("");

        for (int i = 0; i < 10; ++i){
            printf("%s------------\e[0m", i & 1 ? "\e[0;37m" : "\e[0;33m", i);
        }
        puts("");


        // ----- row 2
        for (int i = 0; i < 10; ++i){
            printf("%s------------\e[0m", i & 1 ? "\e[0;33m" : "\e[0;37m", i);
        }
        puts("");

        for (int i = 0; i < 10; ++i){
            printf("%s|          |\e[0m", i & 1 ? "\e[0;33m" : "\e[0;37m", i);
        }
        puts("");

        for (int j = 0; j <= 9; ++j){
            printSlot(BOARD[row*20 + j]);
        }
        puts("");

        for (int j = 0; j <= 9; ++j){
            printMarkers(BOARD[row*20 + j]);
        }

        puts("");

        for (int i = 0; i < 10; ++i){
            printf("%s------------\e[0m", i & 1 ? "\e[0;33m" : "\e[0;37m", i);
        }
        puts("");
    }
}


int winner() {
    for(int i = 0; i < 10; ++i){
        if (PLAYERS[i] == NULL){
            continue;
        }
        if (PLAYERS[i]->boardpos == 99){
            Player * p  = PLAYERS[i];

            BOARD[p->boardpos]->playersHere ^= 1 << p->idx;
            memset(BOARD[p->boardpos]->markers + p->idx, 0x20, 1);

            p->boardpos = 0;

            BOARD[p->boardpos]->playersHere |= 1 << p->idx;
            memcpy(BOARD[p->boardpos]->markers + p->idx, p->name, strlen(p->name));

            return i;
        }

        if (PLAYERS[i]->boardpos >= 100){
            // sent back to brazil

            Player * p  = PLAYERS[i];

            BOARD[p->boardpos]->playersHere ^= 1 << p->idx;
            memset(BOARD[p->boardpos]->markers + p->idx, 0x20, 1);

            p->boardpos = 0;

            BOARD[p->boardpos]->playersHere |= 1 << p->idx;
            memcpy(BOARD[p->boardpos]->markers + p->idx, p->name, strlen(p->name));
        }
    }

    return -1;
}

void changeMarker(){
    int i = CURR_TURN;
    printf("New marker for player %d: ", i+1);

    scanf("%c", PLAYERS[i]->name);
    getchar();

    printf("New marker '%c' set.\n", *PLAYERS[i]->name);
}

void togglePlayerMark(Player *p, int type){
    BOARD[p->boardpos]->playersHere ^= 1 << p->idx;

    if (type == 1){
        if (BOARD[p->boardpos]->markers == NULL){
            BOARD[p->boardpos]->markers = malloc(10);

            memset(BOARD[p->boardpos]->markers, 0x20, 10);
        }

        memcpy(BOARD[p->boardpos]->markers + p->idx, p->name, strlen(p->name));
    } else {
        memset(BOARD[p->boardpos]->markers + p->idx, 0x20, 1);

        if (BOARD[p->boardpos]->playersHere == 0){
            free(BOARD[p->boardpos]->markers);
        }
    }
}

void move(){
    printf("----- Player %d's turn! -----\n", CURR_TURN+1);

    printf("Would you like to change your marker? (y/n): ");

    char choice = getchar(); getchar();

    switch(choice){
        case 'y':
            changeMarker();
            break;
        default:
            break;
    }

    int moveBy = 0;

    printf("What did you spin? (1-6): ");

    scanf("%d", &moveBy);
    getchar();

    if (moveBy < 0 || moveBy > 6){
        puts("You have a strange spinner.");
        exit(1);
    }

    if (moveBy == 0){
        puts("Passing a turn... coward!");
        CURR_TURN += 1;
        CURR_TURN %= AMOUNT_PLAYERS;
        return;
    }

    printf("Player %d moves forward by %d!\n", CURR_TURN+1, moveBy);

    Player * p  = PLAYERS[CURR_TURN];

    char prevpos = p->boardpos;

    togglePlayerMark(p, 0);

    // move piece
    p->boardpos += moveBy;

    // togglePlayerMark(p, 1);

    // check for ladders
    bool landedOnLadder = false;

    for(int i = 0; i < 5; ++i){
        if (LADDERS[i][0] == p->boardpos){
            // togglePlayerMark(p, 0);

            p->boardpos = LADDERS[i][1];

            printf("Player %d hit a ladder! They are now on square %d.\n", CURR_TURN+1, p->boardpos+1);

            togglePlayerMark(p, 1);
            landedOnLadder = true;
        }
    }

    // check for chutes
    bool landedOnChute = false;

    for(int i = 0; i < 5; ++i){
        if (CHUTES[i][0] == p->boardpos){
            // togglePlayerMark(p, 0);

            p->boardpos = CHUTES[i][1];

            printf("Player %d hit a chute! They are now on square %d.\n", CURR_TURN+1, p->boardpos+1);
            
            togglePlayerMark(p, 1);
            landedOnChute = true;
        }
    }

    if (!(landedOnChute || landedOnLadder)){
        togglePlayerMark(p, 1);
    }

    if (BOARD[prevpos]->playersHere == 0){
        BOARD[prevpos]->markers = NULL;
    }
    

    CURR_TURN += 1;
    CURR_TURN %= AMOUNT_PLAYERS;
}



void displayChutesAndLadders(){
    puts("This board's chutes are:");
    for (int i = 0; i < 5; ++i){
        printf("%d -> %d\n", CHUTES[i][0], CHUTES[i][1]);
    }

    puts("This board's ladders are:");
    for (int i = 0; i < 5; ++i){
        printf("%d -> %d\n", LADDERS[i][0], LADDERS[i][1]);
    }
}



bool validChute(int val, int type){
    if (type){
        for (int i = 0; i < 5; ++i){
            if (val == CHUTES[i][1]){
                return false;
            }
        }
        return true;
    } else {
        for (int i = 0; i < 5; ++i){
            if (val == CHUTES[i][0]){
                return false;
            }
        }
        return true;
    }
}

bool validLadder(int val, int type){
    if (type){
        for (int i = 0; i < 5; ++i){
            if (val == LADDERS[i][1]){
                return false;
            }
        }
        return true;
    } else {
        for (int i = 0; i < 5; ++i){
            if (val == LADDERS[i][0]){
                return false;
            }
        }
        return true;
    }
}




void changeChutesAndLadders(){
    int start = 0;
    int end = 0;

    int c;
    while ((c = getchar()) != '\n' && c != EOF) { }
    
    for (int i = 0; i < 5; ++i){
        printf("Enter a chute in the format [start][space][end]: ");
        scanf("%d %d", &start, &end);
        getchar();
        if (start < 0 || end < 0 || start >= 100 || end >= 100 || start <= end){
            puts("That's not a chute!");
            exit(1);
        }

        if (!validChute(start, 0) || !validChute(end, 1)){
            puts("Cannot chain chutes!");
            exit(1);
        }

        CHUTES[i][0] = start;
        CHUTES[i][1] = end;
    }
    
    for (int i = 0; i < 5; ++i){
        printf("Enter a ladder in the format [start][space][end]: ");
        scanf("%d %d", &start, &end);
        getchar();
        if (start < 0 || end < 0 || start >= 100 || end >= 100 || start >= end){
            puts("That's not a ladder!");
            exit(1);
        }

        if (!validLadder(start, 0) || !validLadder(end, 1)){
            puts("Cannot chain ladders!");
            exit(1);
        }

        LADDERS[i][0] = start;
        LADDERS[i][1] = end;
    }

}




int main() {
    srand(time(NULL));

    setbuf(stdout, NULL);
    setbuf(stdin, NULL);

    printf("Number of players (max 10): ");
    scanf("%d", &AMOUNT_PLAYERS);
    getchar();

    if (AMOUNT_PLAYERS > 10){
        puts("Too many people. Go play CS:GO.");
        exit(1);
    }

    if (AMOUNT_PLAYERS < 2){
        puts("Too little people. Go play Solitaire.");
        exit(1);
    }

    for (int i = 0; i < AMOUNT_PLAYERS; ++i){
        PLAYERS[i] = malloc(sizeof(Player));
        PLAYERS[i]->name = malloc(1);     // lol

        PLAYERS[i]->idx = i;
        PLAYERS[i]->boardpos = 0;

        printf("Player %d marker (1 character): ", i+1);
        scanf("%c", PLAYERS[i]->name);
        getchar();
    }

    genTable();

    displayChutesAndLadders();

    printf("Would you like to change the chutes and ladders? (y/n): ");

    char choice = getchar();

    if (choice == 'y'){
        changeChutesAndLadders();
        puts("");
        displayChutesAndLadders();
    }



    BOARD[0]->markers = malloc(10);

    for (int i = 0; i < 10; ++i){
        Player * p = PLAYERS[i];

        if (p == NULL){
            continue;
        }

        BOARD[0]->playersHere |= 1 << i;
        
        memcpy(BOARD[p->boardpos]->markers + p->idx, p->name, strlen(p->name));
    }

    viewTable();

    while (1) {
        move();

        int w = winner();

        if(w >= 0){
            printf("Player %d won! Here is your prize: %p\n", w+1, &puts);
        }

        printf("Would you like to take a look at the board now? (y/n): ");
        char choice = getchar(); getchar();

        if (choice == 'y'){
            viewTable();
        }
    }

    puts("\e[0;91mbruh\e[0m");
    exit(0);
}