#include <stdio.h>
#include <ast.h>

#include <analyzer.h>
#include <analyzer/context.h>

#include <codegen.h>

int main()
{
    struct analyzer_context ctx = {0};

    struct ast *parsed = parse();
    struct ast *transformed = prepare(&ctx, parsed);

    struct str code = emit_c(transformed);

    printf("%s", code.str);
    return 0;
}
