#include <stdio.h>
#include <ast.h>

#include <analyzer.h>
#include <analyzer/context.h>

int main()
{
    struct analyzer_context ctx = {0};

    printf("TZN: Parsing...\n");
    struct ast *parsed = parse();
    printf("TZN: AST transformation...\n");
    struct ast *transformed = prepare(&ctx, parsed);

    describe(transformed);
    return 0;
}
