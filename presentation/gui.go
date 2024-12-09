        // presentation/gui.go
package presentation

import (
	"fmt"
	"os"
	"path/filepath"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ParkingView struct {
	Slots       [20]*fyne.Container
	Indicator   *canvas.Rectangle
	InfoLabel   *widget.Label
	HelpText    *widget.Label
	MainWindow  fyne.Window
	ImageAssets []string
}

type ViewHandler struct {
	view *ParkingView
}

func NewViewHandler(view *ParkingView) *ViewHandler {
	return &ViewHandler{view: view}
}

func (vh *ViewHandler) UpdateSlot(slotIndex int, isOccupied bool, imagePath string) {
	slotContainer := vh.view.Slots[slotIndex]
	if isOccupied {
		if len(slotContainer.Objects) <= 2 {
			carImage := canvas.NewImageFromFile(imagePath)
			carImage.Resize(fyne.NewSize(50, 50))
			carImage.FillMode = canvas.ImageFillContain
			slotContainer.Add(carImage)
		}
	} else {
		if len(slotContainer.Objects) > 2 {
			slotContainer.Objects = slotContainer.Objects[:2]
		}
	}
	slotContainer.Refresh()
}

func (vh *ViewHandler) ChangeIndicatorColor(state int) {
	switch state {
	case 1:
		vh.view.Indicator.FillColor = color.RGBA{0, 123, 255, 255}
	case -1:
		vh.view.Indicator.FillColor = color.RGBA{255, 123, 0, 255}
	default:
		vh.view.Indicator.FillColor = color.RGBA{40, 167, 69, 255}
	}
	vh.view.Indicator.Refresh()
}

func CreateParkingView() *ParkingView {
	appInstance := app.New()
	mainWin := appInstance.NewWindow("Parking Simulator")

	rootContainer := container.NewPadded()
	uiContent := container.NewVBox()

	infoLabel := widget.NewLabelWithStyle("Vehicles in queue: 0",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})
	infoContainer := container.NewHBox(layout.NewSpacer(), infoLabel, layout.NewSpacer())

	helpCard := widget.NewCard("Legend", "", nil)
	helpMessage := "🟢 Green: Free entry\n" +
		"🔵 Blue: Vehicle entering\n" +
		"🟠 Orange: Vehicle exiting\n" +
		"⚪ Gray: Free space\n"
	helpLabel := widget.NewLabelWithStyle(helpMessage,
		fyne.TextAlignLeading,
		fyne.TextStyle{Monospace: true})
	helpCard.SetContent(helpLabel)

	var parkingSlots [20]*fyne.Container
	slotsGrid := container.NewGridWithColumns(10)

	var assetImages []string
	files, _ := os.ReadDir("sprites")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			assetImages = append(assetImages, filepath.Join("sprites", file.Name()))
		}
	}

	for i := 0; i < 20; i++ {
		spaceRect := canvas.NewRectangle(color.RGBA{200, 200, 200, 255})
		spaceRect.SetMinSize(fyne.NewSize(60, 60))
		spaceRect.Resize(fyne.NewSize(60, 60))

		slotNumber := canvas.NewText(fmt.Sprintf("%d", i+1), color.Black)
		slotNumber.TextSize = 12
		slotNumber.Alignment = fyne.TextAlignCenter

		slotContainer := container.NewStack(spaceRect, slotNumber)
		parkingSlots[i] = slotContainer
		slotsGrid.Add(container.NewPadded(slotContainer))
	}

	indicatorCard := widget.NewCard("Entrance/Exit", "", nil)
	indicatorRect := canvas.NewRectangle(color.RGBA{0, 255, 0, 255})
	indicatorRect.SetMinSize(fyne.NewSize(300, 40))
	indicatorContainer := container.NewHBox(layout.NewSpacer(), indicatorRect, layout.NewSpacer())
	indicatorCard.SetContent(indicatorContainer)

	startButton := widget.NewButton("Start Simulation", nil)
	startButton.Importance = widget.HighImportance
	startButton.Resize(fyne.NewSize(200, 40))
	buttonContainer := container.NewHBox(layout.NewSpacer(), startButton, layout.NewSpacer())

	uiContent.Add(infoContainer)
	uiContent.Add(widget.NewSeparator())
	uiContent.Add(helpCard)
	uiContent.Add(widget.NewSeparator())
	uiContent.Add(slotsGrid)
	uiContent.Add(widget.NewSeparator())
	uiContent.Add(indicatorCard)
	uiContent.Add(buttonContainer)

	rootContainer.Add(uiContent)
	mainWin.SetContent(rootContainer)
	mainWin.Resize(fyne.NewSize(900, 600))

	return &ParkingView{
		Slots:       parkingSlots,
		Indicator:   indicatorRect,
		InfoLabel:   infoLabel,
		HelpText:    helpLabel,
		MainWindow:  mainWin,
		ImageAssets: assetImages,
	}
}
