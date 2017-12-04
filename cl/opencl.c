#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <CL/cl.h>

int platformIdCount() {
	cl_uint platformIdCount = 0;
	clGetPlatformIDs (0, NULL, &platformIdCount);
	return (int)platformIdCount;
}

char* const kernel_source= 
"__kernel void vector_mul(__global const float *A, __global const float *B, __global float *C) {\n"
"	int i = get_global_id(0);\n"
"	C[i] = A[i] * B[i];\n"
"}\n"
;

const char *mErrorString(cl_int error);

#define ERRCHECK(r) if ( r != 0) {fprintf(stderr,"Err: %s\n", mErrorString(r)); return;}

// GLOBALS
cl_kernel kernel = NULL;
cl_platform_id platform_id = NULL;
cl_device_id device_id = NULL;   
cl_context context = NULL;
cl_program program = NULL;

void init() {
	char *source_str;
	size_t source_size;
	source_str = kernel_source;
	source_size = strlen(source_str);

	// Get platform and device information
	cl_uint ret_num_devices;
	cl_uint ret_num_platforms;

	cl_int ret;

	ret = clGetPlatformIDs(1, &platform_id, &ret_num_platforms);
	ERRCHECK(ret)

	ret = clGetDeviceIDs( platform_id, CL_DEVICE_TYPE_DEFAULT, 1, &device_id, &ret_num_devices);
	ERRCHECK(ret)

		// Create an OpenCL context
	context = clCreateContext( NULL, 1, &device_id, NULL, NULL, &ret);
	ERRCHECK(ret)
	
	// Create a program from the kernel source
	program = clCreateProgramWithSource(context, 1, (const char **)&source_str, (const size_t *)&source_size, &ret);
	ERRCHECK(ret)

	// Build the program
	ret = clBuildProgram(program, 1, &device_id, NULL, NULL, NULL);
	ERRCHECK(ret)

	// Create the OpenCL kernel
	kernel = clCreateKernel(program, "vector_mul", &ret);
	ERRCHECK(ret)

}

void VecMulf32(int sz, float const *vec1, float const *vec2, float *out) {
	if (kernel == NULL) {
		init();
	}

	cl_int ret;
	// Create a command queue
	//cl_command_queue command_queue = clCreateCommandQueue(context, device_id, 0, &ret);

	cl_command_queue command_queue = clCreateCommandQueueWithProperties(context, device_id, 0, &ret);
	ERRCHECK(ret)
	// Create memory buffers on the device for each vector 
	cl_mem a_mem_obj = clCreateBuffer(context, CL_MEM_READ_ONLY, sz * sizeof(float), NULL, &ret);
	cl_mem b_mem_obj = clCreateBuffer(context, CL_MEM_READ_ONLY, sz * sizeof(float), NULL, &ret);
	cl_mem c_mem_obj = clCreateBuffer(context, CL_MEM_WRITE_ONLY, sz * sizeof(float), NULL, &ret);

	cl_event e;
	// Copy the lists A and B to their respective memory buffers
	ret = clEnqueueWriteBuffer(command_queue, a_mem_obj, CL_TRUE, 0, sz * sizeof(float), vec1, 0, NULL, NULL);

	ret = clEnqueueWriteBuffer(command_queue, b_mem_obj, CL_TRUE, 0, sz * sizeof(float), vec2, 0, NULL, NULL);

	// Here?
	// Set the arguments of the kernel
	ret = clSetKernelArg(kernel, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	ret = clSetKernelArg(kernel, 1, sizeof(cl_mem), (void *)&b_mem_obj);
	ret = clSetKernelArg(kernel, 2, sizeof(cl_mem), (void *)&c_mem_obj);

	// Execute the OpenCL kernel on the list
	size_t global_item_size = sz; // Process the entire lists
	// Find out this
	size_t local_item_size = 1; // Divide work items into groups of 64

	ret = clEnqueueNDRangeKernel(command_queue, kernel, 1, NULL, &global_item_size, &local_item_size, 0, NULL, NULL);
	// Read the memory buffer C on the device to the local variable C
	//float *C = (float*)malloc(sizeof(float)*LIST_SIZE);
	
	ret = clEnqueueReadBuffer(command_queue, c_mem_obj, CL_TRUE, 0, sz * sizeof(float), out, 0, NULL, NULL);
	
	ret = clFinish(command_queue);
	ERRCHECK(ret)
	// Display the result to the screen
	//for(int i = 0; i < sz; i++)
	//	printf("%.2f * %.2f = %.2f\n", vec1[i], vec2[i], out[i]);

	// Clean up
	ERRCHECK(ret)
	//ret = clReleaseKernel(kernel);
	//ret = clReleaseProgram(program);
	ret = clReleaseMemObject(a_mem_obj);
	ret = clReleaseMemObject(b_mem_obj);
	ret = clReleaseMemObject(c_mem_obj);
	
	ret = clReleaseCommandQueue(command_queue);

}


const char *mErrorString(cl_int error) {
switch(error){
    // run-time and JIT compiler errors
    case 0: return "CL_SUCCESS";
    case -1: return "CL_DEVICE_NOT_FOUND";
    case -2: return "CL_DEVICE_NOT_AVAILABLE";
    case -3: return "CL_COMPILER_NOT_AVAILABLE";
    case -4: return "CL_MEM_OBJECT_ALLOCATION_FAILURE";
    case -5: return "CL_OUT_OF_RESOURCES";
    case -6: return "CL_OUT_OF_HOST_MEMORY";
    case -7: return "CL_PROFILING_INFO_NOT_AVAILABLE";
    case -8: return "CL_MEM_COPY_OVERLAP";
    case -9: return "CL_IMAGE_FORMAT_MISMATCH";
    case -10: return "CL_IMAGE_FORMAT_NOT_SUPPORTED";
    case -11: return "CL_BUILD_PROGRAM_FAILURE";
    case -12: return "CL_MAP_FAILURE";
    case -13: return "CL_MISALIGNED_SUB_BUFFER_OFFSET";
    case -14: return "CL_EXEC_STATUS_ERROR_FOR_EVENTS_IN_WAIT_LIST";
    case -15: return "CL_COMPILE_PROGRAM_FAILURE";
    case -16: return "CL_LINKER_NOT_AVAILABLE";
    case -17: return "CL_LINK_PROGRAM_FAILURE";
    case -18: return "CL_DEVICE_PARTITION_FAILED";
    case -19: return "CL_KERNEL_ARG_INFO_NOT_AVAILABLE";

    // compile-time errors
    case -30: return "CL_INVALID_VALUE";
    case -31: return "CL_INVALID_DEVICE_TYPE";
    case -32: return "CL_INVALID_PLATFORM";
    case -33: return "CL_INVALID_DEVICE";
    case -34: return "CL_INVALID_CONTEXT";
    case -35: return "CL_INVALID_QUEUE_PROPERTIES";
    case -36: return "CL_INVALID_COMMAND_QUEUE";
    case -37: return "CL_INVALID_HOST_PTR";
    case -38: return "CL_INVALID_MEM_OBJECT";
    case -39: return "CL_INVALID_IMAGE_FORMAT_DESCRIPTOR";
    case -40: return "CL_INVALID_IMAGE_SIZE";
    case -41: return "CL_INVALID_SAMPLER";
    case -42: return "CL_INVALID_BINARY";
    case -43: return "CL_INVALID_BUILD_OPTIONS";
    case -44: return "CL_INVALID_PROGRAM";
    case -45: return "CL_INVALID_PROGRAM_EXECUTABLE";
    case -46: return "CL_INVALID_KERNEL_NAME";
    case -47: return "CL_INVALID_KERNEL_DEFINITION";
    case -48: return "CL_INVALID_KERNEL";
    case -49: return "CL_INVALID_ARG_INDEX";
    case -50: return "CL_INVALID_ARG_VALUE";
    case -51: return "CL_INVALID_ARG_SIZE";
    case -52: return "CL_INVALID_KERNEL_ARGS";
    case -53: return "CL_INVALID_WORK_DIMENSION";
    case -54: return "CL_INVALID_WORK_GROUP_SIZE";
    case -55: return "CL_INVALID_WORK_ITEM_SIZE";
    case -56: return "CL_INVALID_GLOBAL_OFFSET";
    case -57: return "CL_INVALID_EVENT_WAIT_LIST";
    case -58: return "CL_INVALID_EVENT";
    case -59: return "CL_INVALID_OPERATION";
    case -60: return "CL_INVALID_GL_OBJECT";
    case -61: return "CL_INVALID_BUFFER_SIZE";
    case -62: return "CL_INVALID_MIP_LEVEL";
    case -63: return "CL_INVALID_GLOBAL_WORK_SIZE";
    case -64: return "CL_INVALID_PROPERTY";
    case -65: return "CL_INVALID_IMAGE_DESCRIPTOR";
    case -66: return "CL_INVALID_COMPILER_OPTIONS";
    case -67: return "CL_INVALID_LINKER_OPTIONS";
    case -68: return "CL_INVALID_DEVICE_PARTITION_COUNT";

    // extension errors
    case -1000: return "CL_INVALID_GL_SHAREGROUP_REFERENCE_KHR";
    case -1001: return "CL_PLATFORM_NOT_FOUND_KHR";
    case -1002: return "CL_INVALID_D3D10_DEVICE_KHR";
    case -1003: return "CL_INVALID_D3D10_RESOURCE_KHR";
    case -1004: return "CL_D3D10_RESOURCE_ALREADY_ACQUIRED_KHR";
    case -1005: return "CL_D3D10_RESOURCE_NOT_ACQUIRED_KHR";
    default: return "Unknown OpenCL error";
    }
}
