#ifndef __STACK_H__
#define __STACK_H__

#include <stdint.h>

#define STACK_MIN_CAP 8

#define stack_top(stack) ((stack)->len - 1)
#define stack_bottom(stack) ((uint32_t)(0))
#define stack_value(stack, it) ((stack)->data[it])

#define STACK_DECL(name, type)\
struct name {\
    uint32_t len;\
    uint32_t cap;\
    type *data;\
};\
void name##_free(struct name *stack);\
uint32_t name##_push(struct name *stack);\
void name##_pop(struct name *stack);

#define STACK_IMPL(name, type)\
void name##_free(struct name *stack) {\
    if (stack == NULL)\
        return;\
    if (stack->cap > 0)\
        free(stack->data);\
    memset(stack, 0, sizeof(*stack));\
}\
uint32_t name##_push(struct name *stack) {\
    if (stack->cap == stack->len)\
        stack->data = realloc(stack->data, (stack->cap += STACK_MIN_CAP) * sizeof(type));\
    uint32_t it = stack->len;\
    stack->len++;\
    return it;\
}\
void name##_pop(struct name *stack) {\
    if (stack->len > 0)\
        stack->len--;\
}

#endif
