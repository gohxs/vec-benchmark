#include <stdio.h>
#include <immintrin.h>
#include <xmmintrin.h>


//8 floats 32bits
void VecMul256(int sz, float const *vec1,float const *vec2, float *out) {

	int max = sz/8; // sz = nfloats, sz / 8 floats

	__m256 v1,v2, r;
	for (int i = 0; i <= max; i++) {
		int s = i * 8; // 32 bytes
		v1 = _mm256_load_ps(&vec1[s]);
		v2 = _mm256_load_ps(&vec2[s]);
		r = _mm256_mul_ps(v1,v2);
		_mm256_store_ps(&out[s],r);
	}
}


//4 floats 32bits
void VecMul128(int sz, float const *vec1, float const *vec2, float *out) {
	// 4 floats at a time
	int max = sz/4; // sz = nfloats, sz / 8 floats
	__m128 v1,v2, r;
	for (int i = 0; i <= max; i++) {
		int s = i * 4; // 32 bytes
		v1 = _mm_load_ps(&vec1[s]);
		v2 = _mm_load_ps(&vec2[s]);
		r = _mm_mul_ps(v1,v2);
		_mm_store_ps(&out[s],r);
	}
}

//4 floats only 32bits
void Mul128(int sz, float const *vec1, float const *vec2, float *out) {
	// 4 floats at a time
	__m128 v1,v2, r;
	v1 = _mm_load_ps(&vec1[0]);
	v2 = _mm_load_ps(&vec2[0]);
	r = _mm_mul_ps(v1,v2);
	_mm_store_ps(&out[0],r);

}

