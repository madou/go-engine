package editor

import (
	"bytes"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	"github.com/walesey/go-engine/util"
)

type Editor struct {
	renderer              renderer.Renderer
	gameEngine            engine.Engine
	currentMap            *editorModels.MapModel
	rootMapNode           *renderer.Node
	customController      controller.Controller
	controllerManager     *glfwController.ControllerManager
	uiAssets              ui.HtmlAssets
	mainMenu              *ui.Window
	mainMenuController    glfwController.Controller
	overviewMenu          *Overview
	mainMenuOpen          bool
	progressBar           *ui.Window
	fileBrowser           *FileBrowser
	fileBrowserOpen       bool
	fileBrowserController glfwController.Controller
	mouseMode             string
	selectSprite          renderer.Entity
}

func New() *Editor {
	return &Editor{
		uiAssets:    ui.NewHtmlAssets(),
		rootMapNode: renderer.CreateNode(),
		currentMap: &editorModels.MapModel{
			Name: "default",
			Root: editorModels.NewNodeModel("root"),
		},
		mouseMode: "scale",
	}
}

func (e *Editor) Start() {
	glRenderer := &opengl.OpenglRenderer{
		WindowTitle: "GoEngine editor",
	}
	e.renderer = glRenderer
	e.gameEngine = engine.NewEngine(e.renderer)

	e.gameEngine.Start(func() {
		//lighting
		e.renderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, mgl32.Vec3{0.5, -1, 0.3}, 0)

		//Sky
		skyImg, err := assets.ImportImage("TestAssets/Files/skybox/cubemap.png")
		if err == nil {
			e.gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		//root node
		e.gameEngine.AddSpatial(e.rootMapNode)

		//input/controller manager
		e.controllerManager = glfwController.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := e.gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 20.0
		freeMoveActor.LookSpeed = 0.002
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		e.controllerManager.AddController(mainController.(glfwController.Controller))
		e.gameEngine.AddUpdatable(freeMoveActor)

		e.initSelectSprite()
		e.gameEngine.AddUpdatable(engine.UpdatableFunc(e.updateSelectSprite))

		//editor controller
		e.controllerManager.AddController(NewEditorController(e).(glfwController.Controller))

		//custom controller
		e.customController = controller.CreateController()
		e.controllerManager.AddController(e.customController.(glfwController.Controller))

		e.setupUI()
	})
}

func (e *Editor) initSelectSprite() {
	img, _ := assets.DecodeImage(bytes.NewBuffer(util.Base64ToBytes(GeometryIconData)))
	mat := assets.CreateMaterial(img, nil, nil, nil)
	mat.LightingMode = renderer.MODE_UNLIT
	mat.Transparency = renderer.TRANSPARENCY_EMISSIVE
	mat.DepthTest = false
	selectSprite := effects.CreateSprite(1, 1, 1, mat)
	e.selectSprite = selectSprite
	e.gameEngine.AddSpatialTransparent(selectSprite)
}

func (e *Editor) updateSelectSprite(dt float64) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		size := selectedModel.Translation.Sub(e.gameEngine.Camera().GetTranslation()).Len() * 0.02
		translation, err := e.rootMapNode.RelativePosition(selectedModel.GetNode())
		if err == nil {
			e.selectSprite.SetScale(mgl32.Vec3{size, size, size})
			e.selectSprite.SetTranslation(translation)
		} else {
			e.selectSprite.SetScale(mgl32.Vec3{})
		}
	}
}
