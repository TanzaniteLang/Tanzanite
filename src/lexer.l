%{
#include <ast.h>
#include <parser.h>
%}

%%


%%

int yywrap(void) {
    return 1;
}
