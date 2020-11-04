#pragma once

#include "types.h"
#include "user_options.h"
#include "hashcat.h"
#include "potfile.h"
#include "status.h"
#include "thread.h"

typedef struct
{
    hashcat_ctx_t ctx;
    void *gowrapper;
} gocat_ctx_t;

void callback(u32 id, hashcat_ctx_t *hashcat_ctx, void *wrapper, void *buf, size_t len);
void event(const u32 id, hashcat_ctx_t *hashcat_ctx, const void *buf, const size_t len);
void freeargv(int argc, char **argv);
