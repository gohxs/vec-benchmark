#include <stdio.h>
#include <immintrin.h>
#include <xmmintrin.h>


//8 floats 32bits
void VecMulf32x8(int sz, float const *vec1,float const *vec2, float *out){
	int max = sz/8; // sz = nfloats, sz / 8 floats
	for (int i = 0; i < max; i++) {
		int s = i * 8; // 32 bytes
		__m256 y0 = _mm256_loadu_ps(&vec1[s]);
		y0 = _mm256_mul_ps(*(__m256*)&vec2[s], y0);
		_mm256_storeu_ps(&out[s],y0);
	}
	// Remainder
	for (int i = max * 8; i<sz; i++) {
		out[i] = vec1[i] * vec2[i];
	}

}


//4 floats 32bits
void VecMulf32x4(int sz, float const *vec1, float const *vec2, float *out)  {
	int max = sz/4; 
	for (int i = 0; i < max; i++) {
		int s = i * 4; // 32 bytes
		__m128 y0 = _mm_loadu_ps(&vec1[s]);
		y0 = _mm_mul_ps(*(__m128*)&vec2[s], y0);
		_mm_storeu_ps(&out[s],y0);
	}
	// Remainder
	for (int i = max * 4; i< sz; i++) {
		out[i] = vec1[i] * vec2[i];
	}
}


