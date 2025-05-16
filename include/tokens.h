#ifndef __TOKENS_H__
#define __TOKENS_H__

extern const char *tokens[];

enum token_type {
    EOF_TOK,
    UNKNOWN_TOK,

    /* Identifier and literals */
    IDENTIFIER_TOK,
    STRING_TOK,
    CHAR_TOK,
    INT_TOK,
    FLOAT_TOK,
    BOOL_TOK,

    /* Statements */
    ASSING_TOK,
    IF_TOK,
    UNLESS_TOK,
    ELSE_TOK,
    FOR_TOK,
    WHILE_TOK,
    UNTIL_TOK,
    BREAK_TOK,
    NEXT_TOK,
    SWITCH_TOK,
    CASE_TOK,
    BEGIN_TOK,

    /* Operators */

    /* Delimiters */
    QUESTION_MARK_TOK,
    DOT_TOK,
    COMMA_TOK,
    COLON_TOK,
    SEMICOLON_TOK,
    LBR_TOK,
    RBR_TOK,
    LSQUAREBR_TOK,
    RSQUAREBR_TOK,
    LSQUIGLYBR_TOK,
    RSQUIGLYBR_TOK,

    /* Reserved keywords */
    AUTO_TOK,
    FN_TOK,
    RETURN_TOK,
};

#endif
