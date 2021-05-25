package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
	"github.com/vele/chia-console/chia"
)

var ()

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
	if v, err := g.SetView("blockchain_details", 0, maxY/3+1, maxX/4, int(float32(maxY)/2)); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "BlockChain Details"
		v.Frame = true
		v.Wrap = false

		blockChainClient := chia.NewClient(os.Getenv("CHIA_FULL_NODE_CRT"), os.Getenv("CHIA_FULL_NODE_KEY"), os.Getenv("CHIA_CA_CRT"))
		res, err := blockChainClient.GetChiaBlockchainState(os.Getenv("CHIA_SERVER_URL"))
		if err != nil {
			return err
		}
		fmt.Fprintf(v, "Current blockchain difficulty: %v \n", res.BlockchainState.Difficulty)
		fmt.Fprintf(v, "Current blockchain mempool: %v \n", res.BlockchainState.MempoolSize)
		spaceCalc := chia.ByteCountSI(res.BlockchainState.Space)
		fmt.Fprintf(v, "\033[32mCurrent blockchain space: %v \033[0m \n", spaceCalc)

	}
	return nil
}
func walletLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("wallet", maxX/4+1, maxY/3+1, maxX/2, int(float32(maxY)/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Wallet Details"
		v.Frame = true
		blockChainClient := chia.NewClient(os.Getenv("CHIA_WALLET_CRT"), os.Getenv("CHIA_WALLET_KEY"), os.Getenv("CHIA_CA_CRT"))
		res, err := blockChainClient.GetChiaWallet(os.Getenv("CHIA_WALLET_URL"))
		if err != nil {
			return err
		}
		fmt.Fprintf(v, "Current wallet balance: \033[32m%v\033[0m \n", res.WalletBalance.ConfirmedWalletBalance)
		fmt.Fprintf(v, "Pending wallet balance: \033[32m%v\033[0m \n", res.WalletBalance.PendingChange)
		fmt.Fprintf(v, "Spendable wallet balance: \033[32m%v\033[0m \n", res.WalletBalance.SpendableBalance)
		fmt.Fprintf(v, "Unconfirmed wallet balance: \033[32m%v\033[0m \n", res.WalletBalance.UnconfirmedWalletBalance)

	}
	return nil
}
func plotsLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("plots", maxX/2+1, maxY/3+1, maxX-10, int(float32(maxY)/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Plots Details"
		v.Frame = true
		blockChainClient := chia.NewClient(os.Getenv("CHIA_HARVESTER_CRT"), os.Getenv("CHIA_HARVESTER_KEY"), os.Getenv("CHIA_CA_CRT"))
		res, err := blockChainClient.GetChiaPlots(os.Getenv("CHIA_HARVESTER_URL"))
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(v, len(res.Plots))
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

	if err := walletLayout(g); err != nil {
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
	case "transactions":
		if err := redrawDetail(g, v); err != nil {
			return err
		}
	}
	return nil
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("transactions", 0, 0, maxX/2, int(float32(maxY)/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = false
		v.Wrap = false
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		v.Autoscroll = false
		ok := chia.ParseLogs()
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].ParseTime))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(20))
		fmt.Fprintln(v, graph)
	}
	if err := detailsLayout(g); err != nil {
		return err
	}
	if err := walletLayout(g); err != nil {
		return err
	}
	if err := plotsLayout(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("transactions"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
