src/compile_args.h:
	echo '#pragma once' > $@
	echo '#define CONFIGURE_ARGS "$(CONFIGURE_PATH) $(CONFIG_ARGS)"' >> $@

include/parser.h: src/parser.h
	mkdir -p include
	rm -f $@
	ln $< $@

bin_PROGRAMS = Tanzanite
Tanzanite_SOURCES = ./src/main.c ./src/tokens.c ./src/str.c ./src/ast.c ./src/lexer.l ./src/parser.y ./src/str_builder.c ./src/codegen.c

AM_YFLAGS = -d -Wcounterexamples
AM_LFLAGS =
AM_CPPFLAGS = -Wall -Wextra -I$(srcdir)/include -I$(builddir)/include

BUILT_SOURCES = ./src/compile_args.h ./include/parser.h 
CLEANFILES = ./src/compile_args.h ./include/parser.h 
