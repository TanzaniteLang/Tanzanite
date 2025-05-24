#ifndef __QUEUE_H__
#define __QUEUE_H__

#define QUEUE_DECL(name, type)\
struct name##_node {\
    struct name##_node *next;\
    type val;\
};\
\
struct name {\
    struct name##_node *head;\
    struct name##_node *tail;\
};\
void name##_free(struct name *queue);\
void name##_push(struct name *queue, type val);\
type name##_pop(struct name *queue);

#define QUEUE_IMPL(name, type)\
void name##_free(struct name *queue) {\
    if (queue == NULL)\
        return;\
    struct name##_node *iter = queue->head;\
    while (iter != NULL) {\
        struct name##_node *prev = iter;\
        iter = iter->next;\
        free(prev);\
    }\
    memset(queue, 0, sizeof(*queue));\
}\
void name##_push(struct name *queue, type val) {\
    struct name##_node *node = calloc(1, sizeof(*node));\
    if (node == NULL)\
        return;\
    node->val = val;\
    if (queue->tail == NULL)\
        queue->head = node;\
    else\
        queue->tail->next = node;\
    queue->tail = node;\
}\
type name##_pop(struct name *queue) {\
    struct name##_node *node = queue->head;\
    if (node == NULL)\
        return NULL;\
    queue->head = node->next;\
    if (queue->head == NULL)\
        queue->tail = NULL;\
    type res = node->val;\
    free(node);\
    return res;\
}

#endif
