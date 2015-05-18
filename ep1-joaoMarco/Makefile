CC = gcc
CFLAGS = -Wall -ansi -pedantic -pthread
DEBUGFLAGS = -g

OBJS = ominium.o
PROGRAM = ominium

all: ominium

ominium: $(OBJS)
	$(CC) $(CFLAGS) $(OBJS) -o $(PROGRAM)

%.o: %.c
	$(CC) $(CFLAGS) $(DEBUGFLAGS) -c $<

clean:
	-rm -f $(OBJS) $(PROGRAM) *~ core*
