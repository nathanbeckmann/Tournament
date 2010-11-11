TARGETS = main

all : $(TARGETS)

clean :
	rm -rfv *~ *.6 $(TARGETS)

tournament.go : player.6

main : player.6 tournament.6

% : %.6
	6l -o $@ $<

%.6 : %.go
	6g $^
