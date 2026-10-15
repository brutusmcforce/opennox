// Code generated from the go-gl v3.3-core signatures of the functions OpenNox uses; see gl.go.

package gl

import "unsafe"

var procs struct {
	ActiveTexture           uintptr
	AttachShader            uintptr
	BindAttribLocation      uintptr
	BindBuffer              uintptr
	BindFragDataLocation    uintptr
	BindFramebuffer         uintptr
	BindTexture             uintptr
	BindVertexArray         uintptr
	BlendEquation           uintptr
	BlendFunc               uintptr
	BlitFramebuffer         uintptr
	BufferData              uintptr
	CheckFramebufferStatus  uintptr
	Clear                   uintptr
	CompileShader           uintptr
	CreateProgram           uintptr
	CreateShader            uintptr
	DeleteBuffers           uintptr
	DeleteFramebuffers      uintptr
	DeleteProgram           uintptr
	DeleteShader            uintptr
	DeleteTextures          uintptr
	DeleteVertexArrays      uintptr
	Disable                 uintptr
	DrawArrays              uintptr
	DrawArraysInstanced     uintptr
	DrawElements            uintptr
	Enable                  uintptr
	EnableVertexAttribArray uintptr
	Finish                  uintptr
	FramebufferTexture2D    uintptr
	GenBuffers              uintptr
	GenFramebuffers         uintptr
	GenTextures             uintptr
	GenVertexArrays         uintptr
	GetAttribLocation       uintptr
	GetError                uintptr
	GetIntegerv             uintptr
	GetProgramInfoLog       uintptr
	GetProgramiv            uintptr
	GetShaderInfoLog        uintptr
	GetShaderiv             uintptr
	GetString               uintptr
	GetTexImage             uintptr
	GetUniformLocation      uintptr
	LinkProgram             uintptr
	PixelStorei             uintptr
	ReadPixels              uintptr
	Scissor                 uintptr
	ShaderSource            uintptr
	TexImage2D              uintptr
	TexParameteri           uintptr
	TexSubImage2D           uintptr
	Uniform1f               uintptr
	Uniform1i               uintptr
	Uniform2f               uintptr
	Uniform2i               uintptr
	Uniform3fv              uintptr
	Uniform4f               uintptr
	UseProgram              uintptr
	VertexAttribDivisor     uintptr
	VertexAttribPointer     uintptr
	Viewport                uintptr
}

func procNames() map[string]*uintptr {
	return map[string]*uintptr{
		"glActiveTexture":           &procs.ActiveTexture,
		"glAttachShader":            &procs.AttachShader,
		"glBindAttribLocation":      &procs.BindAttribLocation,
		"glBindBuffer":              &procs.BindBuffer,
		"glBindFragDataLocation":    &procs.BindFragDataLocation,
		"glBindFramebuffer":         &procs.BindFramebuffer,
		"glBindTexture":             &procs.BindTexture,
		"glBindVertexArray":         &procs.BindVertexArray,
		"glBlendEquation":           &procs.BlendEquation,
		"glBlendFunc":               &procs.BlendFunc,
		"glBlitFramebuffer":         &procs.BlitFramebuffer,
		"glBufferData":              &procs.BufferData,
		"glCheckFramebufferStatus":  &procs.CheckFramebufferStatus,
		"glClear":                   &procs.Clear,
		"glCompileShader":           &procs.CompileShader,
		"glCreateProgram":           &procs.CreateProgram,
		"glCreateShader":            &procs.CreateShader,
		"glDeleteBuffers":           &procs.DeleteBuffers,
		"glDeleteFramebuffers":      &procs.DeleteFramebuffers,
		"glDeleteProgram":           &procs.DeleteProgram,
		"glDeleteShader":            &procs.DeleteShader,
		"glDeleteTextures":          &procs.DeleteTextures,
		"glDeleteVertexArrays":      &procs.DeleteVertexArrays,
		"glDisable":                 &procs.Disable,
		"glDrawArrays":              &procs.DrawArrays,
		"glDrawArraysInstanced":     &procs.DrawArraysInstanced,
		"glDrawElements":            &procs.DrawElements,
		"glEnable":                  &procs.Enable,
		"glEnableVertexAttribArray": &procs.EnableVertexAttribArray,
		"glFinish":                  &procs.Finish,
		"glFramebufferTexture2D":    &procs.FramebufferTexture2D,
		"glGenBuffers":              &procs.GenBuffers,
		"glGenFramebuffers":         &procs.GenFramebuffers,
		"glGenTextures":             &procs.GenTextures,
		"glGenVertexArrays":         &procs.GenVertexArrays,
		"glGetAttribLocation":       &procs.GetAttribLocation,
		"glGetError":                &procs.GetError,
		"glGetIntegerv":             &procs.GetIntegerv,
		"glGetProgramInfoLog":       &procs.GetProgramInfoLog,
		"glGetProgramiv":            &procs.GetProgramiv,
		"glGetShaderInfoLog":        &procs.GetShaderInfoLog,
		"glGetShaderiv":             &procs.GetShaderiv,
		"glGetString":               &procs.GetString,
		"glGetTexImage":             &procs.GetTexImage,
		"glGetUniformLocation":      &procs.GetUniformLocation,
		"glLinkProgram":             &procs.LinkProgram,
		"glPixelStorei":             &procs.PixelStorei,
		"glReadPixels":              &procs.ReadPixels,
		"glScissor":                 &procs.Scissor,
		"glShaderSource":            &procs.ShaderSource,
		"glTexImage2D":              &procs.TexImage2D,
		"glTexParameteri":           &procs.TexParameteri,
		"glTexSubImage2D":           &procs.TexSubImage2D,
		"glUniform1fv":              &procs.Uniform1f,
		"glUniform1i":               &procs.Uniform1i,
		"glUniform2fv":              &procs.Uniform2f,
		"glUniform2i":               &procs.Uniform2i,
		"glUniform3fv":              &procs.Uniform3fv,
		"glUniform4fv":              &procs.Uniform4f,
		"glUseProgram":              &procs.UseProgram,
		"glVertexAttribDivisor":     &procs.VertexAttribDivisor,
		"glVertexAttribPointer":     &procs.VertexAttribPointer,
		"glViewport":                &procs.Viewport,
	}
}

func ActiveTexture(texture uint32) {
	call(procs.ActiveTexture, uintptr(texture))
}

func AttachShader(program uint32, shader uint32) {
	call(procs.AttachShader, uintptr(program), uintptr(shader))
}

func BindAttribLocation(program uint32, index uint32, name *uint8) {
	call(procs.BindAttribLocation, uintptr(program), uintptr(index), uintptr(unsafe.Pointer(name)))
}

func BindBuffer(target uint32, buffer uint32) {
	call(procs.BindBuffer, uintptr(target), uintptr(buffer))
}

func BindFragDataLocation(program uint32, color uint32, name *uint8) {
	call(procs.BindFragDataLocation, uintptr(program), uintptr(color), uintptr(unsafe.Pointer(name)))
}

func BindFramebuffer(target uint32, framebuffer uint32) {
	call(procs.BindFramebuffer, uintptr(target), uintptr(framebuffer))
}

func BindTexture(target uint32, texture uint32) {
	call(procs.BindTexture, uintptr(target), uintptr(texture))
}

func BindVertexArray(array uint32) {
	call(procs.BindVertexArray, uintptr(array))
}

func BlendEquation(mode uint32) {
	call(procs.BlendEquation, uintptr(mode))
}

func BlendFunc(sfactor uint32, dfactor uint32) {
	call(procs.BlendFunc, uintptr(sfactor), uintptr(dfactor))
}

func BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	call(procs.BlitFramebuffer, uintptr(srcX0), uintptr(srcY0), uintptr(srcX1), uintptr(srcY1), uintptr(dstX0), uintptr(dstY0), uintptr(dstX1), uintptr(dstY1), uintptr(mask), uintptr(filter))
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	call(procs.BufferData, uintptr(target), uintptr(size), uintptr(data), uintptr(usage))
}

func CheckFramebufferStatus(target uint32) uint32 {
	return uint32(call(procs.CheckFramebufferStatus, uintptr(target)))
}

func Clear(mask uint32) {
	call(procs.Clear, uintptr(mask))
}

func CompileShader(shader uint32) {
	call(procs.CompileShader, uintptr(shader))
}

func CreateProgram() uint32 {
	return uint32(call(procs.CreateProgram))
}

func CreateShader(xtype uint32) uint32 {
	return uint32(call(procs.CreateShader, uintptr(xtype)))
}

func DeleteBuffers(n int32, buffers *uint32) {
	call(procs.DeleteBuffers, uintptr(n), uintptr(unsafe.Pointer(buffers)))
}

func DeleteFramebuffers(n int32, framebuffers *uint32) {
	call(procs.DeleteFramebuffers, uintptr(n), uintptr(unsafe.Pointer(framebuffers)))
}

func DeleteProgram(program uint32) {
	call(procs.DeleteProgram, uintptr(program))
}

func DeleteShader(shader uint32) {
	call(procs.DeleteShader, uintptr(shader))
}

func DeleteTextures(n int32, textures *uint32) {
	call(procs.DeleteTextures, uintptr(n), uintptr(unsafe.Pointer(textures)))
}

func DeleteVertexArrays(n int32, arrays *uint32) {
	call(procs.DeleteVertexArrays, uintptr(n), uintptr(unsafe.Pointer(arrays)))
}

func Disable(cap uint32) {
	call(procs.Disable, uintptr(cap))
}

func DrawArrays(mode uint32, first int32, count int32) {
	call(procs.DrawArrays, uintptr(mode), uintptr(first), uintptr(count))
}

func DrawArraysInstanced(mode uint32, first int32, count int32, instancecount int32) {
	call(procs.DrawArraysInstanced, uintptr(mode), uintptr(first), uintptr(count), uintptr(instancecount))
}

func DrawElements(mode uint32, count int32, xtype uint32, indices unsafe.Pointer) {
	call(procs.DrawElements, uintptr(mode), uintptr(count), uintptr(xtype), uintptr(indices))
}

func Enable(cap uint32) {
	call(procs.Enable, uintptr(cap))
}

func EnableVertexAttribArray(index uint32) {
	call(procs.EnableVertexAttribArray, uintptr(index))
}

func Finish() {
	call(procs.Finish)
}

func FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	call(procs.FramebufferTexture2D, uintptr(target), uintptr(attachment), uintptr(textarget), uintptr(texture), uintptr(level))
}

func GenBuffers(n int32, buffers *uint32) {
	call(procs.GenBuffers, uintptr(n), uintptr(unsafe.Pointer(buffers)))
}

func GenFramebuffers(n int32, framebuffers *uint32) {
	call(procs.GenFramebuffers, uintptr(n), uintptr(unsafe.Pointer(framebuffers)))
}

func GenTextures(n int32, textures *uint32) {
	call(procs.GenTextures, uintptr(n), uintptr(unsafe.Pointer(textures)))
}

func GenVertexArrays(n int32, arrays *uint32) {
	call(procs.GenVertexArrays, uintptr(n), uintptr(unsafe.Pointer(arrays)))
}

func GetAttribLocation(program uint32, name *uint8) int32 {
	return int32(call(procs.GetAttribLocation, uintptr(program), uintptr(unsafe.Pointer(name))))
}

func GetError() uint32 {
	return uint32(call(procs.GetError))
}

func GetIntegerv(pname uint32, data *int32) {
	call(procs.GetIntegerv, uintptr(pname), uintptr(unsafe.Pointer(data)))
}

func GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	call(procs.GetProgramInfoLog, uintptr(program), uintptr(bufSize), uintptr(unsafe.Pointer(length)), uintptr(unsafe.Pointer(infoLog)))
}

func GetProgramiv(program uint32, pname uint32, params *int32) {
	call(procs.GetProgramiv, uintptr(program), uintptr(pname), uintptr(unsafe.Pointer(params)))
}

func GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	call(procs.GetShaderInfoLog, uintptr(shader), uintptr(bufSize), uintptr(unsafe.Pointer(length)), uintptr(unsafe.Pointer(infoLog)))
}

func GetShaderiv(shader uint32, pname uint32, params *int32) {
	call(procs.GetShaderiv, uintptr(shader), uintptr(pname), uintptr(unsafe.Pointer(params)))
}

func GetString(name uint32) *uint8 {
	return (*uint8)(unsafe.Pointer(call(procs.GetString, uintptr(name)))) //nolint:govet // native string
}

func GetTexImage(target uint32, level int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	call(procs.GetTexImage, uintptr(target), uintptr(level), uintptr(format), uintptr(xtype), uintptr(pixels))
}

func GetUniformLocation(program uint32, name *uint8) int32 {
	return int32(call(procs.GetUniformLocation, uintptr(program), uintptr(unsafe.Pointer(name))))
}

func LinkProgram(program uint32) {
	call(procs.LinkProgram, uintptr(program))
}

func PixelStorei(pname uint32, param int32) {
	call(procs.PixelStorei, uintptr(pname), uintptr(param))
}

func ReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	call(procs.ReadPixels, uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(pixels))
}

func Scissor(x int32, y int32, width int32, height int32) {
	call(procs.Scissor, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}

func ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	call(procs.ShaderSource, uintptr(shader), uintptr(count), uintptr(unsafe.Pointer(xstring)), uintptr(unsafe.Pointer(length)))
}

func TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	call(procs.TexImage2D, uintptr(target), uintptr(level), uintptr(internalformat), uintptr(width), uintptr(height), uintptr(border), uintptr(format), uintptr(xtype), uintptr(pixels))
}

func TexParameteri(target uint32, pname uint32, param int32) {
	call(procs.TexParameteri, uintptr(target), uintptr(pname), uintptr(param))
}

func TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	call(procs.TexSubImage2D, uintptr(target), uintptr(level), uintptr(xoffset), uintptr(yoffset), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(pixels))
}

func Uniform1f(location int32, v0 float32) {
	v := [1]float32{v0}
	call(procs.Uniform1f, uintptr(location), 1, uintptr(unsafe.Pointer(&v)))
}

func Uniform1i(location int32, v0 int32) {
	call(procs.Uniform1i, uintptr(location), uintptr(v0))
}

func Uniform2f(location int32, v0 float32, v1 float32) {
	v := [2]float32{v0, v1}
	call(procs.Uniform2f, uintptr(location), 1, uintptr(unsafe.Pointer(&v)))
}

func Uniform2i(location int32, v0 int32, v1 int32) {
	call(procs.Uniform2i, uintptr(location), uintptr(v0), uintptr(v1))
}

func Uniform3fv(location int32, count int32, value *float32) {
	call(procs.Uniform3fv, uintptr(location), uintptr(count), uintptr(unsafe.Pointer(value)))
}

func Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	v := [4]float32{v0, v1, v2, v3}
	call(procs.Uniform4f, uintptr(location), 1, uintptr(unsafe.Pointer(&v)))
}

func UseProgram(program uint32) {
	call(procs.UseProgram, uintptr(program))
}

func VertexAttribDivisor(index uint32, divisor uint32) {
	call(procs.VertexAttribDivisor, uintptr(index), uintptr(divisor))
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	call(procs.VertexAttribPointer, uintptr(index), uintptr(size), uintptr(xtype), cbool(normalized), uintptr(stride), uintptr(pointer))
}

func VertexAttribPointerWithOffset(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uintptr) {
	call(procs.VertexAttribPointer, uintptr(index), uintptr(size), uintptr(xtype), cbool(normalized), uintptr(stride), offset)
}

func Viewport(x int32, y int32, width int32, height int32) {
	call(procs.Viewport, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}
