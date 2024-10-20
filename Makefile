.SUFFIXES: .o .go

GCCGO = gccgo

lotor: main.o
	$(GCCGO) $(GCCGOFLAGS) -o lotor main.o

.go.o:
	$(GCCGO) $(GCCGOFLAGS) -c -o $@ $<
