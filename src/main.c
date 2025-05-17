#include <ast.h>

struct ast *parse();

int main()
{
    struct ast *parsed = parse();
    describe(parsed);
    return 0;
}
