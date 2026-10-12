#ifndef NOX_COMMON_CRYPT_H
#define NOX_COMMON_CRYPT_H

#include <stdio.h>

int nox_crypt_IsReadOnly();
int nox_xxx_cryptOpen_426910(char* a1, int a2, int a3);
void nox_xxx_cryptClose_4269F0();
void nox_xxx_cryptSetTypeMB_426A50(int a1);
int nox_xxx_cryptFlush_4268E0();
FILE* nox_xxx_mapgenGetSomeFile_426A60();
int nox_xxx_cryptSeekCur_40E0A0(int a1);
void sub_4268F0(int a1);
size_t nox_xxx_fileReadWrite_426AC0_file3_fread_impl(unsigned char* a1, size_t a2, char* fname, int line);

// Saving a map writes tens of thousands of small fields; calling into Go for each of them takes about 20 ms.
// While the crypt file is open for writing, the fields are collected here instead. Go writes them to the file
// before any other operation on it (see legacy/crypt_wbuf.go), so the result is the same.
typedef struct {
	unsigned char* data;
	unsigned int len;
	unsigned int cap;
	unsigned int* sizes; // size of each field: the file checksum is computed field by field
	unsigned int nsizes;
	unsigned int capsizes;
	int on;
} nox_crypt_wbuf_t;
extern nox_crypt_wbuf_t nox_crypt_wbuf;

size_t nox_xxx_fileReadWrite_426AC0_buffered(unsigned char* a1, size_t a2, char* fname, int line);
#define nox_xxx_fileReadWrite_426AC0_file3_fread(a1, a2) nox_xxx_fileReadWrite_426AC0_buffered((unsigned char*)(a1), a2, __FILE__, __LINE__)
void nox_xxx_fileCryptReadCrcMB_426C20(unsigned char* a1, size_t a2);
void nox_xxx_crypt_426C90();
void nox_xxx_crypt_426D40();

#endif // NOX_COMMON_CRYPT_H
