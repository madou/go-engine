package opengl

import (
	"image"
	"image/draw"
	"log"
	"runtime"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/shaders"
	"github.com/walesey/go-engine/util"
)

const (
	maxLights int = 8
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// OpenglRenderer - opengl implementation
type OpenglRenderer struct {
	onInit, onUpdate, onRender func()
	WindowWidth, WindowHeight  int
	FullScreen                 bool
	WindowTitle                string
	Window                     *glfw.Window
	transformStack             *matstack.TransformStack
	program                    uint32
	envMapId                   uint32
	envMapLOD1Id               uint32
	envMapLOD2Id               uint32
	envMapLOD3Id               uint32
	illuminanceMapId           uint32
	modelUniform               int32
	nbLights                   int32
	nbDirectionalLights        int32
	lights                     []float32
	directionalLights          []float32
	cameraLocation             mgl32.Vec3
	cameraUp                   mgl32.Vec3
	cameraLookAt               mgl32.Vec3
	cameraAngle                float32
	cameraNear                 float32
	cameraFar                  float32
	cameraOrtho                bool
	tlv, trv, blv, brv         mgl32.Vec3
	postEffectVbo              uint32
	postEffects                []postEffect
}

// NewOpenglRenderer - create new renderer
func NewOpenglRenderer(WindowTitle string, WindowWidth, WindowHeight int, FullScreen bool) *OpenglRenderer {
	return &OpenglRenderer{
		WindowTitle:  WindowTitle,
		WindowWidth:  WindowWidth,
		WindowHeight: WindowHeight,
		FullScreen:   FullScreen,
	}
}

//Init -
func (glRenderer *OpenglRenderer) Init(callback func()) {
	glRenderer.onInit = callback
}

//Update -
func (glRenderer *OpenglRenderer) Update(callback func()) {
	glRenderer.onUpdate = callback
}

//Render -
func (glRenderer *OpenglRenderer) Render(callback func()) {
	glRenderer.onRender = callback
}

//Start -
func (glRenderer *OpenglRenderer) Start() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if glRenderer.FullScreen {
		glRenderer.WindowWidth = glfw.GetPrimaryMonitor().GetVideoMode().Width
	} else if glRenderer.WindowWidth == 0 {
		glRenderer.WindowWidth = glfw.GetPrimaryMonitor().GetVideoMode().Width * 95 / 100
	}
	if glRenderer.FullScreen {
		glRenderer.WindowHeight = glfw.GetPrimaryMonitor().GetVideoMode().Height
	} else if glRenderer.WindowHeight == 0 {
		glRenderer.WindowHeight = glfw.GetPrimaryMonitor().GetVideoMode().Height * 95 / 100
	}

	var monitor *glfw.Monitor
	if glRenderer.FullScreen {
		monitor = glfw.GetPrimaryMonitor()
	}
	window, err := glfw.CreateWindow(glRenderer.WindowWidth, glRenderer.WindowHeight, glRenderer.WindowTitle, monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glRenderer.Window = window

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shaders
	vertShader, _ := compileShader(shaders.MainVert, gl.VERTEX_SHADER)
	fragShader, _ := compileShader(shaders.MainFrag, gl.FRAGMENT_SHADER)
	program, _ := newProgram(vertShader, fragShader)
	gl.UseProgram(program)
	glRenderer.program = program

	//set shader uniforms
	glRenderer.transformStack = matstack.NewTransformStack()
	glRenderer.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	glRenderer.PushTransform(mgl32.Ident4())

	textureUniform := gl.GetUniformLocation(program, gl.Str("diffuse\x00"))
	gl.Uniform1i(textureUniform, 0)
	textureUniform = gl.GetUniformLocation(program, gl.Str("normal\x00"))
	gl.Uniform1i(textureUniform, 1)
	textureUniform = gl.GetUniformLocation(program, gl.Str("specular\x00"))
	gl.Uniform1i(textureUniform, 2)
	textureUniform = gl.GetUniformLocation(program, gl.Str("roughness\x00"))
	gl.Uniform1i(textureUniform, 3)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMap\x00"))
	gl.Uniform1i(textureUniform, 4)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD1\x00"))
	gl.Uniform1i(textureUniform, 5)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD2\x00"))
	gl.Uniform1i(textureUniform, 6)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD3\x00"))
	gl.Uniform1i(textureUniform, 7)
	textureUniform = gl.GetUniformLocation(program, gl.Str("illuminanceMap\x00"))
	gl.Uniform1i(textureUniform, 8)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//setup Lights
	glRenderer.lights = make([]float32, maxLights*16, maxLights*16)
	glRenderer.directionalLights = make([]float32, maxLights*16, maxLights*16)

	glRenderer.initPostEffects()

	glRenderer.onInit()

	window.SetRefreshCallback(func(w *glfw.Window) {
		glRenderer.mainLoop()
		window.SwapBuffers()
	})

	//Main loop
	for !window.ShouldClose() {
		glRenderer.mainLoop()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (glRenderer *OpenglRenderer) mainLoop() {
	glRenderer.onUpdate()
	if len(glRenderer.postEffects) == 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.onRender()
	} else {
		//Render to the first post effect buffer
		gl.UseProgram(glRenderer.program)
		gl.BindFramebuffer(gl.FRAMEBUFFER, glRenderer.postEffects[0].fboId)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.onRender()
		//Render Post effects
		for i := 0; i < len(glRenderer.postEffects)-1; i = i + 1 {
			gl.BindFramebuffer(gl.FRAMEBUFFER, glRenderer.postEffects[i+1].fboId)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			glRenderer.renderPostEffect(glRenderer.postEffects[i])
		}
		//Render final post effect to the frame buffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.renderPostEffect(glRenderer.postEffects[len(glRenderer.postEffects)-1])
	}
}

// BackGroundColor - set background color for the scene
func (glRenderer *OpenglRenderer) BackGroundColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (glRenderer *OpenglRenderer) WindowDimensions() mgl32.Vec2 {
	return mgl32.Vec2{float32(glRenderer.WindowWidth), float32(glRenderer.WindowHeight)}
}

// Ortho - set orthogonal rendering mode
func (glRenderer *OpenglRenderer) Ortho() {
	projection := mgl32.Ortho2D(0, float32(glRenderer.WindowWidth), float32(glRenderer.WindowHeight), 0)
	projectionUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	camera := mgl32.Ident4()
	cameraUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	glRenderer.cameraOrtho = true
}

// Camera - camera settings
func (glRenderer *OpenglRenderer) Perspective(location, lookat, up mgl32.Vec3, angle, near, far float32) {
	projection := mgl32.Perspective(mgl32.DegToRad(angle), float32(glRenderer.WindowWidth)/float32(glRenderer.WindowHeight), near, far)
	projectionUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	camera := mgl32.LookAtV(location, lookat, up)
	cameraUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	glRenderer.cameraAngle = angle
	glRenderer.cameraNear = near
	glRenderer.cameraFar = far
	glRenderer.cameraLocation = location
	glRenderer.cameraLookAt = lookat
	glRenderer.cameraUp = up
	glRenderer.cameraOrtho = false
	glRenderer.updateCameraVectors()
}

func (glRenderer *OpenglRenderer) updateCameraVectors() {
	window := glRenderer.WindowDimensions()
	modelview := mgl32.LookAtV(glRenderer.cameraLocation, glRenderer.cameraLookAt, glRenderer.cameraUp)
	projection := mgl32.Perspective(mgl32.DegToRad(glRenderer.cameraAngle), window.X()/window.Y(), glRenderer.cameraNear, glRenderer.cameraFar)
	v, _ := mgl32.UnProject(mgl32.Vec3{0, 0, 0.5}, modelview, projection, 0, 0, int(window.X()), int(window.Y()))
	glRenderer.tlv = v.Sub(glRenderer.cameraLocation).Normalize()
	v, _ = mgl32.UnProject(mgl32.Vec3{window.X(), 0, 0.5}, modelview, projection, 0, 0, int(window.X()), int(window.Y()))
	glRenderer.trv = v.Sub(glRenderer.cameraLocation).Normalize()
	v, _ = mgl32.UnProject(mgl32.Vec3{0, window.Y(), 0.5}, modelview, projection, 0, 0, int(window.X()), int(window.Y()))
	glRenderer.blv = v.Sub(glRenderer.cameraLocation).Normalize()
	v, _ = mgl32.UnProject(mgl32.Vec3{window.X(), window.Y(), 0.5}, modelview, projection, 0, 0, int(window.X()), int(window.Y()))
	glRenderer.brv = v.Sub(glRenderer.cameraLocation).Normalize()
}

// FrustrumContainsSphere - calculates if a sphere is at least partially within the camera frustrum
func (glRenderer *OpenglRenderer) FrustrumContainsSphere(radius float32) bool {
	if glRenderer.cameraOrtho {
		return true
	}
	model := glRenderer.transformStack.Peek()
	point := mgl32.TransformCoordinate(mgl32.Vec3{}, model)
	r := mgl32.TransformCoordinate(mgl32.Vec3{radius, 0, 0}, model).Sub(point).Len()
	delta := point.Sub(glRenderer.cameraLocation)
	if point.ApproxEqual(glRenderer.cameraLocation) {
		delta = mgl32.Vec3{1, 0, 0}
	}
	return delta.Dot(glRenderer.tlv.Cross(glRenderer.trv)) <= r &&
		delta.Dot(glRenderer.trv.Cross(glRenderer.brv)) <= r &&
		delta.Dot(glRenderer.brv.Cross(glRenderer.blv)) <= r &&
		delta.Dot(glRenderer.blv.Cross(glRenderer.tlv)) <= r &&
		util.Vec3LenSq(delta) < (glRenderer.cameraFar+r)*(glRenderer.cameraFar+r)
	return true
}

// CameraLocation - get location of the camera
func (glRenderer *OpenglRenderer) CameraLocation() mgl32.Vec3 {
	return glRenderer.cameraLocation
}

func (glRenderer *OpenglRenderer) PushTransform(transform mgl32.Mat4) {
	glRenderer.transformStack.Push(transform)
	model := glRenderer.transformStack.Peek()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func (glRenderer *OpenglRenderer) PopTransform() {
	glRenderer.transformStack.Pop()
	model := glRenderer.transformStack.Peek()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func (glRenderer *OpenglRenderer) EnableDepthTest(depthTest bool) {
	if depthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (glRenderer *OpenglRenderer) EnableDepthMask(depthMast bool) {
	gl.DepthMask(depthMast)
}

// CreateGeometry - add geometry to the renderer
func (glRenderer *OpenglRenderer) CreateGeometry(geometry *renderer.Geometry) {

	// Configure the vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
	geometry.VboId = vbo

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
	geometry.IboId = ibo
}

//
func (glRenderer *OpenglRenderer) DestroyGeometry(geometry *renderer.Geometry) {
	gl.DeleteBuffers(1, &geometry.VboId)
	gl.DeleteBuffers(1, &geometry.IboId)
}

//setup Texture
func (glRenderer *OpenglRenderer) CreateMaterial(material *renderer.Material) {
	if material.Diffuse != nil {
		material.DiffuseId = glRenderer.newTexture(material.Diffuse, gl.TEXTURE0)
	}
	if material.Normal != nil {
		material.NormalId = glRenderer.newTexture(material.Normal, gl.TEXTURE1)
	}
	if material.Specular != nil {
		material.SpecularId = glRenderer.newTexture(material.Specular, gl.TEXTURE2)
	}
	if material.Roughness != nil {
		material.RoughnessId = glRenderer.newTexture(material.Roughness, gl.TEXTURE3)
	}
}

//
func (glRenderer *OpenglRenderer) DestroyMaterial(material *renderer.Material) {
	gl.DeleteTextures(1, &material.DiffuseId)
	gl.DeleteTextures(1, &material.NormalId)
	gl.DeleteTextures(1, &material.SpecularId)
	gl.DeleteTextures(1, &material.RoughnessId)
}

//setup Texture
func (glRenderer *OpenglRenderer) newTexture(img image.Image, textureUnit uint32) uint32 {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texId
}

func (glRenderer *OpenglRenderer) ReflectionMap(cm *renderer.CubeMap) {
	cm.Resize(256)
	glRenderer.envMapId = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE4)
	cm.Resize(80)
	glRenderer.envMapLOD1Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE5)
	cm.Resize(30)
	glRenderer.envMapLOD2Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE6)
	cm.Resize(12)
	glRenderer.envMapLOD3Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE7)
	cm.Resize(6)
	glRenderer.illuminanceMapId = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE8)
}

func (glRenderer *OpenglRenderer) newCubeMap(right, left, top, bottom, back, front image.Image, textureUnit uint32) uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texId)

	for i := 0; i < 6; i++ {
		img := right
		texIndex := (uint32)(gl.TEXTURE_CUBE_MAP_POSITIVE_X)
		if i == 1 {
			img = left
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_X
		} else if i == 2 {
			img = top
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Y
		} else if i == 3 {
			img = bottom
			texIndex = gl.TEXTURE_CUBE_MAP_POSITIVE_Y
		} else if i == 4 {
			img = back
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Z
		} else if i == 5 {
			img = front
			texIndex = gl.TEXTURE_CUBE_MAP_POSITIVE_Z
		}
		img = imaging.FlipV(img)
		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			log.Fatal("unsupported stride")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		gl.TexImage2D(
			texIndex,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	return texId
}

//
func (glRenderer *OpenglRenderer) DrawGeometry(geometry *renderer.Geometry) {

	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.VboId)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.IboId)

	//update buffers
	if geometry.VboDirty && len(geometry.Verticies) > 0 && len(geometry.Indicies) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
		geometry.VboDirty = false
	}

	//set back face culling
	if geometry.CullBackface {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}

	//set depthbuffer modes
	glRenderer.EnableDepthTest(geometry.Material.DepthTest)
	glRenderer.EnableDepthMask(geometry.Material.DepthMask)

	//set lighting mode
	lightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("mode\x00"))
	gl.Uniform1i(lightsUniform, geometry.Material.LightingMode)

	//world camera position
	cam := glRenderer.CameraLocation()
	worldCamPosUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("worldCamPos\x00"))
	gl.Uniform4f(worldCamPosUniform, cam[0], cam[1], cam[2], 0)

	//transparency mode
	if geometry.Material.Transparency == renderer.TRANSPARENCY_NON_EMISSIVE {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else if geometry.Material.Transparency == renderer.TRANSPARENCY_EMISSIVE {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	} else {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	//set verticies attribute
	vertAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(0))
	//set normals attribute
	normAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(3*4))
	//set texture coord attribute
	texCoordAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(6*4))
	//vertex color attribute
	colorAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(8*4))

	//setup textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.DiffuseId)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.NormalId)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.SpecularId)
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.RoughnessId)
	gl.ActiveTexture(gl.TEXTURE4)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapId)
	gl.ActiveTexture(gl.TEXTURE5)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD1Id)
	gl.ActiveTexture(gl.TEXTURE6)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD2Id)
	gl.ActiveTexture(gl.TEXTURE7)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD3Id)
	gl.ActiveTexture(gl.TEXTURE8)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.illuminanceMapId)

	useVertexColorUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("useVertexColor\x00"))
	if geometry.Material.Diffuse == nil {
		gl.Uniform1i(useVertexColorUniform, 1)
	} else {
		gl.Uniform1i(useVertexColorUniform, 0)
	}

	gl.DrawElements(gl.TRIANGLES, (int32)(len(geometry.Indicies)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// ambient, diffuse and specular light values ( i is the light index )
func (glRenderer *OpenglRenderer) CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position mgl32.Vec3, i int) {
	lights := glRenderer.lights
	if directional {
		lights = glRenderer.directionalLights
		glRenderer.nbDirectionalLights = int32(i + 1)
	} else {
		glRenderer.nbLights = int32(i + 1)
	}

	//position
	lights[(i * 16)] = position.X()
	lights[(i*16)+1] = position.Y()
	lights[(i*16)+2] = position.Z()
	lights[(i*16)+3] = 1
	//ambient
	lights[(i*16)+4] = ar
	lights[(i*16)+5] = ag
	lights[(i*16)+6] = ab
	lights[(i*16)+7] = 1
	//diffuse
	lights[(i*16)+8] = dr
	lights[(i*16)+9] = dg
	lights[(i*16)+10] = db
	lights[(i*16)+11] = 1
	//specular
	lights[(i*16)+12] = sr
	lights[(i*16)+13] = sg
	lights[(i*16)+14] = sb
	lights[(i*16)+15] = 1

	//set uniform array
	uniformName := "lights\x00"
	if directional {
		uniformName = "directionalLights\x00"
	}
	lightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str(uniformName))
	gl.Uniform4fv(lightsUniform, (int32)(maxLights*16), &lights[0])
	nbLightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("nbLights\x00"))
	gl.Uniform1i(nbLightsUniform, glRenderer.nbLights)
	nbLightsUniform = gl.GetUniformLocation(glRenderer.program, gl.Str("nbDirectionalLights\x00"))
	gl.Uniform1i(nbLightsUniform, glRenderer.nbDirectionalLights)
}

func (glRenderer *OpenglRenderer) DestroyLight(i int) {
	//TODO
}

func (glRenderer *OpenglRenderer) LockCursor(lock bool) {
	glRenderer.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
