// THE AUTOGENERATED LICENSE. ALL THE RIGHTS ARE RESERVED BY ROBOTS.

// WARNING: This file has automatically been generated on Sun, 04 Sep 2016 02:34:29 MSK.
// By http://git.io/cgogen. DO NOT EDIT.

package portaudio

/*
#cgo pkg-config: portaudio-2.0
#include <portaudio.h>
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import "unsafe"

// GetVersion function as declared in portaudio/portaudio.h:57
func GetVersion() int32 {
	__ret := C.Pa_GetVersion()
	__v := (int32)(__ret)
	return __v
}

// GetVersionText function as declared in portaudio/portaudio.h:63
func GetVersionText() string {
	__ret := C.Pa_GetVersionText()
	__v := packPCharString(__ret)
	return __v
}

// GetErrorText function as declared in portaudio/portaudio.h:110
func GetErrorText(errorCode Error) string {
	cerrorCode, _ := (C.PaError)(errorCode), cgoAllocsUnknown
	__ret := C.Pa_GetErrorText(cerrorCode)
	__v := packPCharString(__ret)
	return __v
}

// Initialize function as declared in portaudio/portaudio.h:132
func Initialize() Error {
	__ret := C.Pa_Initialize()
	__v := (Error)(__ret)
	return __v
}

// Terminate function as declared in portaudio/portaudio.h:151
func Terminate() Error {
	__ret := C.Pa_Terminate()
	__v := (Error)(__ret)
	return __v
}

// GetHostApiCount function as declared in portaudio/portaudio.h:199
func GetHostApiCount() HostApiIndex {
	__ret := C.Pa_GetHostApiCount()
	__v := (HostApiIndex)(__ret)
	return __v
}

// GetDefaultHostApi function as declared in portaudio/portaudio.h:210
func GetDefaultHostApi() HostApiIndex {
	__ret := C.Pa_GetDefaultHostApi()
	__v := (HostApiIndex)(__ret)
	return __v
}

// GetHostApiInfo function as declared in portaudio/portaudio.h:289
func GetHostApiInfo(hostApi HostApiIndex) *HostApiInfo {
	chostApi, _ := (C.PaHostApiIndex)(hostApi), cgoAllocsUnknown
	__ret := C.Pa_GetHostApiInfo(chostApi)
	__v := NewHostApiInfoRef(unsafe.Pointer(__ret))
	return __v
}

// HostApiTypeIdToHostApiIndex function as declared in portaudio/portaudio.h:307
func HostApiTypeIdToHostApiIndex(kind HostApiTypeId) HostApiIndex {
	ckind, _ := (C.PaHostApiTypeId)(kind), cgoAllocsUnknown
	__ret := C.Pa_HostApiTypeIdToHostApiIndex(ckind)
	__v := (HostApiIndex)(__ret)
	return __v
}

// HostApiDeviceIndexToDeviceIndex function as declared in portaudio/portaudio.h:331
func HostApiDeviceIndexToDeviceIndex(hostApi HostApiIndex, hostApiDeviceIndex int32) DeviceIndex {
	chostApi, _ := (C.PaHostApiIndex)(hostApi), cgoAllocsUnknown
	chostApiDeviceIndex, _ := (C.int)(hostApiDeviceIndex), cgoAllocsUnknown
	__ret := C.Pa_HostApiDeviceIndexToDeviceIndex(chostApi, chostApiDeviceIndex)
	__v := (DeviceIndex)(__ret)
	return __v
}

// GetLastHostErrorInfo function as declared in portaudio/portaudio.h:358
func GetLastHostErrorInfo() *HostErrorInfo {
	__ret := C.Pa_GetLastHostErrorInfo()
	__v := NewHostErrorInfoRef(unsafe.Pointer(__ret))
	return __v
}

// GetDeviceCount function as declared in portaudio/portaudio.h:371
func GetDeviceCount() DeviceIndex {
	__ret := C.Pa_GetDeviceCount()
	__v := (DeviceIndex)(__ret)
	return __v
}

// GetDefaultInputDevice function as declared in portaudio/portaudio.h:380
func GetDefaultInputDevice() DeviceIndex {
	__ret := C.Pa_GetDefaultInputDevice()
	__v := (DeviceIndex)(__ret)
	return __v
}

// GetDefaultOutputDevice function as declared in portaudio/portaudio.h:398
func GetDefaultOutputDevice() DeviceIndex {
	__ret := C.Pa_GetDefaultOutputDevice()
	__v := (DeviceIndex)(__ret)
	return __v
}

// GetDeviceInfo function as declared in portaudio/portaudio.h:482
func GetDeviceInfo(device DeviceIndex) *DeviceInfo {
	cdevice, _ := (C.PaDeviceIndex)(device), cgoAllocsUnknown
	__ret := C.Pa_GetDeviceInfo(cdevice)
	__v := NewDeviceInfoRef(unsafe.Pointer(__ret))
	return __v
}

// IsFormatSupported function as declared in portaudio/portaudio.h:558
func IsFormatSupported(inputParameters *StreamParameters, outputParameters *StreamParameters, sampleRate float64) Error {
	cinputParameters, _ := inputParameters.PassRef()
	coutputParameters, _ := outputParameters.PassRef()
	csampleRate, _ := (C.double)(sampleRate), cgoAllocsUnknown
	__ret := C.Pa_IsFormatSupported(cinputParameters, coutputParameters, csampleRate)
	__v := (Error)(__ret)
	return __v
}

// OpenStream function as declared in portaudio/portaudio.h:841
func OpenStream(stream **Stream, inputParameters *StreamParameters, outputParameters *StreamParameters, sampleRate float64, framesPerBuffer uint, streamFlags StreamFlags, streamCallback StreamCallback, userData unsafe.Pointer) Error {
	cstream, _ := (*unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	cinputParameters, _ := inputParameters.PassRef()
	coutputParameters, _ := outputParameters.PassRef()
	csampleRate, _ := (C.double)(sampleRate), cgoAllocsUnknown
	cframesPerBuffer, _ := (C.ulong)(framesPerBuffer), cgoAllocsUnknown
	cstreamFlags, _ := (C.PaStreamFlags)(streamFlags), cgoAllocsUnknown
	cstreamCallback, _ := streamCallback.PassRef()
	cuserData, _ := (unsafe.Pointer)(unsafe.Pointer(userData)), cgoAllocsUnknown
	__ret := C.Pa_OpenStream(cstream, cinputParameters, coutputParameters, csampleRate, cframesPerBuffer, cstreamFlags, cstreamCallback, cuserData)
	__v := (Error)(__ret)
	return __v
}

// OpenDefaultStream function as declared in portaudio/portaudio.h:881
func OpenDefaultStream(stream **Stream, numInputChannels int32, numOutputChannels int32, sampleFormat SampleFormat, sampleRate float64, framesPerBuffer uint, streamCallback StreamCallback, userData unsafe.Pointer) Error {
	cstream, _ := (*unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	cnumInputChannels, _ := (C.int)(numInputChannels), cgoAllocsUnknown
	cnumOutputChannels, _ := (C.int)(numOutputChannels), cgoAllocsUnknown
	csampleFormat, _ := (C.PaSampleFormat)(sampleFormat), cgoAllocsUnknown
	csampleRate, _ := (C.double)(sampleRate), cgoAllocsUnknown
	cframesPerBuffer, _ := (C.ulong)(framesPerBuffer), cgoAllocsUnknown
	cstreamCallback, _ := streamCallback.PassRef()
	cuserData, _ := (unsafe.Pointer)(unsafe.Pointer(userData)), cgoAllocsUnknown
	__ret := C.Pa_OpenDefaultStream(cstream, cnumInputChannels, cnumOutputChannels, csampleFormat, csampleRate, cframesPerBuffer, cstreamCallback, cuserData)
	__v := (Error)(__ret)
	return __v
}

// CloseStream function as declared in portaudio/portaudio.h:894
func CloseStream(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_CloseStream(cstream)
	__v := (Error)(__ret)
	return __v
}

// SetStreamFinishedCallback function as declared in portaudio/portaudio.h:932
func SetStreamFinishedCallback(stream *Stream, streamFinishedCallback StreamFinishedCallback) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	cstreamFinishedCallback, _ := streamFinishedCallback.PassRef()
	__ret := C.Pa_SetStreamFinishedCallback(cstream, cstreamFinishedCallback)
	__v := (Error)(__ret)
	return __v
}

// StartStream function as declared in portaudio/portaudio.h:937
func StartStream(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_StartStream(cstream)
	__v := (Error)(__ret)
	return __v
}

// StopStream function as declared in portaudio/portaudio.h:943
func StopStream(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_StopStream(cstream)
	__v := (Error)(__ret)
	return __v
}

// AbortStream function as declared in portaudio/portaudio.h:949
func AbortStream(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_AbortStream(cstream)
	__v := (Error)(__ret)
	return __v
}

// IsStreamStopped function as declared in portaudio/portaudio.h:964
func IsStreamStopped(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_IsStreamStopped(cstream)
	__v := (Error)(__ret)
	return __v
}

// IsStreamActive function as declared in portaudio/portaudio.h:980
func IsStreamActive(stream *Stream) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_IsStreamActive(cstream)
	__v := (Error)(__ret)
	return __v
}

// GetStreamInfo function as declared in portaudio/portaudio.h:1034
func GetStreamInfo(stream *Stream) *StreamInfo {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_GetStreamInfo(cstream)
	__v := NewStreamInfoRef(unsafe.Pointer(__ret))
	return __v
}

// GetStreamTime function as declared in portaudio/portaudio.h:1052
func GetStreamTime(stream *Stream) Time {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_GetStreamTime(cstream)
	__v := (Time)(__ret)
	return __v
}

// GetStreamCpuLoad function as declared in portaudio/portaudio.h:1071
func GetStreamCpuLoad(stream *Stream) float64 {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_GetStreamCpuLoad(cstream)
	__v := (float64)(__ret)
	return __v
}

// ReadStream function as declared in portaudio/portaudio.h:1095
func ReadStream(stream *Stream, buffer unsafe.Pointer, frames uint) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	cbuffer, _ := (unsafe.Pointer)(unsafe.Pointer(buffer)), cgoAllocsUnknown
	cframes, _ := (C.ulong)(frames), cgoAllocsUnknown
	__ret := C.Pa_ReadStream(cstream, cbuffer, cframes)
	__v := (Error)(__ret)
	return __v
}

// WriteStream function as declared in portaudio/portaudio.h:1122
func WriteStream(stream *Stream, buffer unsafe.Pointer, frames uint) Error {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	cbuffer, _ := (unsafe.Pointer)(unsafe.Pointer(buffer)), cgoAllocsUnknown
	cframes, _ := (C.ulong)(frames), cgoAllocsUnknown
	__ret := C.Pa_WriteStream(cstream, cbuffer, cframes)
	__v := (Error)(__ret)
	return __v
}

// GetStreamReadAvailable function as declared in portaudio/portaudio.h:1135
func GetStreamReadAvailable(stream *Stream) int {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_GetStreamReadAvailable(cstream)
	__v := (int)(__ret)
	return __v
}

// GetStreamWriteAvailable function as declared in portaudio/portaudio.h:1146
func GetStreamWriteAvailable(stream *Stream) int {
	cstream, _ := (unsafe.Pointer)(unsafe.Pointer(stream)), cgoAllocsUnknown
	__ret := C.Pa_GetStreamWriteAvailable(cstream)
	__v := (int)(__ret)
	return __v
}

// GetSampleSize function as declared in portaudio/portaudio.h:1157
func GetSampleSize(format SampleFormat) Error {
	cformat, _ := (C.PaSampleFormat)(format), cgoAllocsUnknown
	__ret := C.Pa_GetSampleSize(cformat)
	__v := (Error)(__ret)
	return __v
}

// Sleep function as declared in portaudio/portaudio.h:1167
func Sleep(msec int) {
	cmsec, _ := (C.long)(msec), cgoAllocsUnknown
	C.Pa_Sleep(cmsec)
}
