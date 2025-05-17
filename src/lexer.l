%{
#include <ast.h>
#include <str.h>
#include <parser.h>
%}

%%
"if"                return IF_TOK;
"elsif"             return ELSIF_TOK;
"else"              return ELSE_TOK;
"then"              return THEN_TOK;
"def"               return DEF_TOK;
"end"               return END_TOK;
[A-z]+[A-z0-9]*     { yylval.str = str_init(yytext, yyleng); return IDENTIFIER_TOK; }
[0-9]+              { yylval.num = atol(yytext); return INT_TOK; }
[0-9]+\.[0-9]+      { yylval.dec = strtod(yytext, NULL); }
[ \t\n]             ;
","                 return ',';
";"                 return ';';
"="                 return '=';
":"                 return ':';
"{"                 return '{';
"}"                 return '}';
"("                 return '(';
")"                 return ')';
"+"                 return PLUS_TOK;
"-"                 return MINUS_TOK;
"*"                 return STAR_TOK;
"/"                 return SLASH_TOK;

.                   { fprintf(stderr, "???: %s\n", yytext); abort(); }

%%

int yywrap(void) {
    return 1;
}
