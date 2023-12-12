package simulation

import (
	"fmt"
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"slices"
	"strconv"
	"strings"
)

var idle = image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 255, A: 180})
var hover = image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 255, A: 220})
var pressed = image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 255, A: 255})

var idleBackground = image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 180})
var hoverBackground = image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 220})

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func newPauseButton(sim *Simulation) *widget.Button {
	face, _ := loadFont(20)

	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		}),

		widget.ButtonOpts.Text("Pause", face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			sim.Paused = !sim.Paused
		}),
	)

	return button
}

func newDebugButton(sim *Simulation) *widget.Button {
	face, _ := loadFont(20)

	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		}),

		widget.ButtonOpts.Text("Debug", face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			sim.Debug = !sim.Debug
		}),
	)

	return button
}

func newTPSSlider(sim *Simulation) *widget.Slider {
	slider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		widget.SliderOpts.MinMax(3, 50),

		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(200, 6),
		),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  idleBackground,
				Hover: hoverBackground,
			},
			&widget.ButtonImage{
				Idle:    idle,
				Hover:   hover,
				Pressed: pressed,
			},
		),
		widget.SliderOpts.FixedHandleSize(6),
		widget.SliderOpts.TrackOffset(0),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			sim.TPS = args.Current
		}),
	)

	return slider
}

type AgentChoice struct {
	id string
}

func makeAgentList(sim *Simulation) []any {
	var agents []any
	sim.Agents.Range(func(_, ag interface{}) bool {
		if ag == nil {
			return true
		}
		agent := ag.(*typing.Agent)
		agents = append(agents, AgentChoice{id: fmt.Sprintf("%s", agent.ID)})
		return true
	})

	slices.SortFunc(agents, func(i, j any) int {
		iId, _ := strconv.Atoi(strings.Split(i.(AgentChoice).id, "-")[1])
		jId, _ := strconv.Atoi(strings.Split(j.(AgentChoice).id, "-")[1])
		if iId < jId {
			return -1
		} else if iId > jId {
			return 1
		} else {
			return 0
		}
	})

	return agents
}

func newAgentChoice(sim *Simulation) *widget.List {
	face, _ := loadFont(20)

	agents := makeAgentList(sim)

	list := widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
			}),
		)),
		widget.ListOpts.Entries(agents),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     idleBackground,
				Disabled: idleBackground,
				Mask:     idleBackground,
			}),
		),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(
				&widget.SliderTrackImage{
					Idle:  idleBackground,
					Hover: hoverBackground,
				},
				&widget.ButtonImage{
					Idle:    idle,
					Hover:   hover,
					Pressed: pressed,
				},
			),
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryFontFace(face),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Unselected:                 color.RGBA{R: 255, G: 255, B: 255, A: 180},
			SelectedBackground:         color.RGBA{R: 0, G: 0, B: 255, A: 220},
			SelectedFocusedBackground:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
			FocusedBackground:          color.RGBA{R: 0, G: 0, B: 255, A: 220},
			DisabledUnselected:         color.RGBA{R: 255, G: 255, B: 255, A: 180},
			DisabledSelected:           color.RGBA{R: 255, G: 255, B: 255, A: 220},
			DisabledSelectedBackground: color.RGBA{R: 0, G: 0, B: 255, A: 200},
		}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(AgentChoice).id
		}),
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			sim.SelectedAgent = args.Entry.(AgentChoice).id
		}),
	)

	return list
}

func makeAgentDetails() *widget.TextArea {
	face, _ := loadFont(14)

	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  300,
					MaxHeight: 100,
				}),
				widget.WidgetOpts.MinSize(200, 200),
			),
		),
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontColor(color.White),
		widget.TextAreaOpts.FontFace(face),
		widget.TextAreaOpts.Text("Select an agent to see it's statistics"),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(10)),
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: idleBackground,
				Mask: idleBackground,
			}),
		),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				&widget.SliderTrackImage{
					Idle:  idleBackground,
					Hover: hoverBackground,
				},
				&widget.ButtonImage{
					Idle:    idle,
					Hover:   hover,
					Pressed: pressed,
				},
			),
		),
	)

	return textarea
}

func BuildUI(sim *Simulation) (*ebitenui.UI, *widget.List, *widget.TextArea) {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(10, 10),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}))),
	)

	leftContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(10, 10),
			widget.GridLayoutOpts.Stretch([]bool{false}, []bool{true, false}))),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)

	pauseButton := newPauseButton(sim)
	buttonContainer.AddChild(pauseButton)

	debugButton := newDebugButton(sim)
	buttonContainer.AddChild(debugButton)

	leftContainer.AddChild(buttonContainer)

	tpsSlider := newTPSSlider(sim)
	tpsSlider.Current = sim.TPS
	leftContainer.AddChild(tpsSlider)

	rightContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(10, 10),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, false}))),
	)

	agentChoice := newAgentChoice(sim)
	rightContainer.AddChild(agentChoice)

	agentDetails := makeAgentDetails()
	rightContainer.AddChild(agentDetails)

	rootContainer.AddChild(leftContainer)
	rootContainer.AddChild(rightContainer)

	rootContainer.BackgroundImage = image.NewNineSliceColor(color.RGBA{})

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui, agentChoice, agentDetails
}
