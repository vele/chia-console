package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "instances" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("instances")
	return err
}
func autoscroll(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = true
	return nil
}

func movable(v *gocui.View, nextY int) (ok bool, yLimit int) {
	switch v.Name() {
	case "instances":
		yLimit = 10
		if yLimit < nextY {
			return false, yLimit
		}
		return true, yLimit

	default:
		return true, 0
	}
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		ok, _ := movable(v, oy+cy+dy)

		if !ok {
			return nil
		}
		if err := v.SetCursor(cx, cy+dy); err != nil {
			if err := v.SetOrigin(ox, oy+dy); err != nil {
				return err
			}
		}

	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("instances", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			onMovingCursorRedrawView(g, v)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("instances", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			onMovingCursorRedrawView(g, v)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("instances", gocui.KeyCtrlF, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func detailsLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("details", 0, maxY/3+1, maxX/4, int(float32(maxY)/1.4)); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Disk Details"
		v.Frame = true
		v.Wrap = false
	}
	return nil
}
func networkLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("network", maxX/4+1, maxY/3+1, maxX/2, int(float32(maxY)/1.4)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Network Details"
		v.Frame = true
		fmt.Fprintln(v, maxX, maxY, maxX/2, maxX/3, maxX/4)

	}
	return nil
}
func redrawDetail(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView("details"); err != nil {
		return err
	}
	if err := g.DeleteView("network"); err != nil {
		return err
	}
	//_, cy := v.Cursor()
	//l, err := v.Line(cy)

	//instance_id := strings.Split(l, "|")

	//if len(strings.TrimSpace(instance_id[0])) == 0 {
	//	instance_detailed = GetInstanceDetails("eu-north-1", strings.TrimSpace(first_instance_id))
	//}
	//instance_detailed = GetInstanceDetails("eu-north-1", strings.TrimSpace(instance_id[0]))

	if err := networkLayout(g); err != nil {
		return err
	}
	if err := detailsLayout(g); err != nil {
		return err
	}
	return nil
	//return detailsLayout(g)
}
func onMovingCursorRedrawView(g *gocui.Gui, v *gocui.View) error {

	switch v.Name() {
	case "instances":
		if err := redrawDetail(g, v); err != nil {
			return err
		}
	}
	return nil
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("instances", 0, 0, maxX-10, maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		a, b, c, d, e := g.ViewPosition("instances")
		fmt.Fprintln(v, a, b, c, d, e)
		v.Highlight = true
		v.Wrap = false
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.Frame = true
		v.Autoscroll = true
		//instanceList := GetInstanceList("eu-north-1")
		//instance_count = len(instanceList.Items)
		//first_instance_id = instanceList.Items[0].InstanceId

		//for i := range instanceList.Items {
		//	fmt.Fprintf(v, "%-20s|%-50s|%-20s|%-30s|%-20s|%-10s\n",
		//		instanceList.Items[i].InstanceId,
		//		instanceList.Items[i].InstanceName,
		//		instanceList.Items[i].InstanceType,
		//		instanceList.Items[i].LaunchTime,
		//		instanceList.Items[i].PublicAddress,
		//		instanceList.Items[i].SecurityGroups)
		//}
	}
	if err := detailsLayout(g); err != nil {
		return err
	}

	if err := networkLayout(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("instances"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
