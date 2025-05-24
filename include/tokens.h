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
    ASSIGN_TOK,
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
    RETURN_TOK,

    /* Operators */
    EQUALS_TOK,
    PLUS_TOK,
    DPLUS_TOK,
    PLUS_ASSIGN_TOK,
    MINUS_TOK,
    DMINUS_TOK,
    MINUS_ASSIGN_TOK,
    STAR_TOK,
    STAR_ASSIGN_TOK,
    SLASH_TOK,
    SLASH_ASSIGN_TOK,
    DSLASH_TOK,
    DSLASH_ASSIGN_TOK,
    MOD_TOK,
    MOD_ASSIGN_TOK,
    BANG_TOK,
    NOT_EQL_TOK,
    TILDA_TOK,
    TILDA_ASSIGN_TOK,
    AMPERSAND_TOK,
    AMPERSAND_ASSIGN_TOK,
    AND_TOK,
    PIPE_TOK,
    PIPE_ASSIGN_TOK,
    OR_TOK,
    PIPE_FORWARD_TOK,
    CARET_TOK,
    CARET_ASSIGN_TOK,
    LESS_TOK,
    LESS_EQL_TOK,
    LEFT_SHIFT_TOK,
    LEFT_SHIFT_ASSIGN_TOK,
    GREATER_TOK,
    GREATER_EQL_TOK,
    RIGHT_SHIFT_TOK,
    RIGHT_SHIFT_ASSIGN_TOK,

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
    RESCUE_TOK,
};

#endif
