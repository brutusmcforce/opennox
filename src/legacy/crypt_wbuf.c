#include <string.h>

#include "common__crypt.h"

nox_crypt_wbuf_t nox_crypt_wbuf = {0};

size_t nox_xxx_fileReadWrite_426AC0_buffered(unsigned char* a1, size_t a2, char* fname, int line) {
	nox_crypt_wbuf_t* b = &nox_crypt_wbuf;
	if (b->on && a2 <= b->cap - b->len && b->nsizes < b->capsizes) {
		memcpy(b->data + b->len, a1, a2);
		b->len += a2;
		b->sizes[b->nsizes++] = a2;
		return 1;
	}
	return nox_xxx_fileReadWrite_426AC0_file3_fread_impl(a1, a2, fname, line);
}
