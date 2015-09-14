package examples

import (
	"image/color"

	"github.com/Walesey/goEngine/actor"
	"github.com/Walesey/goEngine/assets"
	"github.com/Walesey/goEngine/controller"
	"github.com/Walesey/goEngine/effects"
	"github.com/Walesey/goEngine/renderer"
	"github.com/Walesey/goEngine/vectorMath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
)

//
func Particles(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 6000

	glRenderer := &renderer.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}

	assetLib, err := assets.LoadAssetLibrary("TestAssets/particles.asset")
	if err != nil {
		panic(err)
	}

	//setup scenegraph

	geom := assetLib.GetGeometry("skybox")
	skyboxMat := assetLib.GetMaterial("skyboxMat")
	geom.Material = &skyboxMat
	geom.Material.LightingMode = renderer.MODE_UNLIT
	geom.CullBackface = false
	skyNode := renderer.CreateNode()
	skyNode.Add(&geom)
	skyNode.SetRotation(1.57, vectorMath.Vector3{0, 1, 0})
	skyNode.SetScale(vectorMath.Vector3{5000, 5000, 5000})

	geomsphere := assetLib.GetGeometry("sphere")
	sphereMat := assetLib.GetMaterial("sphereMat")
	geomsphere.Material = &sphereMat
	sphereNode := renderer.CreateNode()
	sphereNode.Add(&geomsphere)
	sphereNode.SetTranslation(vectorMath.Vector3{1, 1, 3})

	//particle effects
	explosionMat := assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
	explosionMat.LightingMode = renderer.MODE_UNLIT
	explosionParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        4,
		ParticleEmitRate:    2,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            explosionMat,
		TotalFrames:         36,
		FramesX:             6,
		FramesY:             6,
		FaceCamera:          true,
		MaxLife:             1.0,
		MinLife:             2.0,
		StartSize:           vectorMath.Vector3{0.4, 0.4, 0.4},
		EndSize:             vectorMath.Vector3{2.4, 2.4, 2.4},
		StartColor:          color.NRGBA{254, 254, 254, 254},
		EndColor:            color.NRGBA{254, 254, 254, 254},
		MinTranslation:      vectorMath.Vector3{-0.1, -0.1, -0.1},
		MaxTranslation:      vectorMath.Vector3{0.1, 0.1, 0.1},
		MaxStartVelocity:    vectorMath.Vector3{0.2, 1.8, 0.2},
		MinStartVelocity:    vectorMath.Vector3{-0.2, 2.5, -0.2},
		Acceleration:        vectorMath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectorMath.IdentityQuaternion(),
		MinAngularVelocity:  vectorMath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})

	fireMat := assets.CreateMaterial(assetLib.GetImage("fire"), nil, nil, nil)
	fireMat.LightingMode = renderer.MODE_UNLIT
	fireParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        10,
		ParticleEmitRate:    2,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            fireMat,
		TotalFrames:         36,
		FramesX:             6,
		FramesY:             6,
		FaceCamera:          true,
		MaxLife:             1.0,
		MinLife:             1.3,
		StartSize:           vectorMath.Vector3{1.0, 1.0, 1.0},
		EndSize:             vectorMath.Vector3{1.7, 1.7, 1.7},
		StartColor:          color.NRGBA{254, 54, 0, 200},
		EndColor:            color.NRGBA{254, 100, 20, 50},
		MinTranslation:      vectorMath.Vector3{-0.1, 0.1, -0.1},
		MaxTranslation:      vectorMath.Vector3{0.1, 0.3, 0.1},
		MaxStartVelocity:    vectorMath.Vector3{0.02, 0.02, 0.02},
		MinStartVelocity:    vectorMath.Vector3{-0.02, -0.02, -0.02},
		Acceleration:        vectorMath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectorMath.IdentityQuaternion(),
		MinAngularVelocity:  vectorMath.IdentityQuaternion(),
		MaxRotationVelocity: 0.3,
		MinRotationVelocity: -0.3,
	})

	smokeMat := assets.CreateMaterial(assetLib.GetImage("smoke"), nil, nil, nil)
	smokeMat.LightingMode = renderer.MODE_UNLIT
	smokeParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        300,
		ParticleEmitRate:    600,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            smokeMat,
		TotalFrames:         64,
		FramesX:             8,
		FramesY:             8,
		FaceCamera:          true,
		MaxLife:             0.5,
		MinLife:             0.3,
		StartSize:           vectorMath.Vector3{0.4, 0.4, 0.4},
		EndSize:             vectorMath.Vector3{2.4, 2.4, 2.4},
		StartColor:          color.NRGBA{254, 254, 254, 50},
		EndColor:            color.NRGBA{254, 254, 254, 0},
		MinTranslation:      vectorMath.Vector3{-0.2, -0.2, -0.2},
		MaxTranslation:      vectorMath.Vector3{0.2, 0.2, 0.2},
		MaxStartVelocity:    vectorMath.Vector3{0.2, 13.3, 0.2},
		MinStartVelocity:    vectorMath.Vector3{-0.2, 13.5, -0.2},
		Acceleration:        vectorMath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectorMath.IdentityQuaternion(),
		MinAngularVelocity:  vectorMath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})

	birdMat := assets.CreateMaterial(assetLib.GetImage("bird"), nil, nil, nil)
	birdMat.LightingMode = renderer.MODE_UNLIT
	birdSprite := effects.CreateSprite(21, 5, 5, &birdMat)
	birdSprite.SetTranslation(vectorMath.Vector3{-2, 0, -1})

	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(&skyNode)
	sceneGraph.Add(&sphereNode)
	sceneGraph.Add(&fireParticles)
	sceneGraph.Add(&smokeParticles)
	sceneGraph.Add(&explosionParticles)
	sceneGraph.Add(&birdSprite)

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	freeMoveActor.MoveSpeed = 3.0

	glRenderer.Init = func() {
		//setup reflection map
		cubeMap := renderer.CreateCubemap(assetLib.GetMaterial("skyboxMat").Diffuse)
		glRenderer.ReflectionMap(*cubeMap)

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		//camera free move actor
		mainController := controller.NewBasicMovementController(freeMoveActor)
		controllerManager.AddController(mainController)

		//test the portabitity of the actor / entity interfaces
		mainController.BindAction(func() { freeMoveActor.Entity = camera }, glfw.KeyQ, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &sphereNode }, glfw.KeyW, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &explosionParticles }, glfw.KeyE, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &birdSprite }, glfw.KeyR, glfw.Press)
	}

	glRenderer.Update = func() {
		fps.UpdateFPSMeter()

		//update things that need updating
		explosionParticles.Update(0.018, glRenderer)
		fireParticles.Update(0.018, glRenderer)
		smokeParticles.Update(0.018, glRenderer)

		birdSprite.NextFrame()

		freeMoveActor.Update(0.018)
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}