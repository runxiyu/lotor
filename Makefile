.SUFFIXES: .o .go

ifeq ($(GCCGO),)
GCCGO = gccgo
endif

lotor: main.o
	$(GCCGO) $(GCCGOFLAGS) -o lotor main.o

.go.o:
	$(GCCGO) $(GCCGOFLAGS) -c -o $@ $<
