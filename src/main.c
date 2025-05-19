#include <stdio.h>
#include <ast.h>
#include <str.h>
#include <codegen.h>

int main()
{
    printf("TZN: Parsing...\n");
    struct ast *parsed = parse();
    printf("TZN: AST codegen...\n");
    struct str code = emit_c(parsed);

    printf("Code:\n\n%s", code.str);
    return 0;
}
