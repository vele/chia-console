package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
	"github.com/vele/chia-console/chia"
)

var (
	done = make(chan struct{})
	wg   sync.WaitGroup

	mu  sync.Mutex // protects ctr
	ctr = 0
)

func drawEligablePlotsGraph(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)
		ok := chia.ParseLogs(600)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].Plots))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(11))
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("main")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, graph)
			return nil
		})
	}
}
func drawProcessingTimesGraph(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)
		ok := chia.ParseLogs(600)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].ParseTime))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(11))
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("totalPlots")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, graph)
			return nil
		})
	}
}
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
		fmt.Fprintf(v, "\033[32mCurrent blockchain space: %v \033[0m \n", res.BlockchainState.Space)

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

	if v, err := g.SetView("plots", maxX/2+1, maxY/3+1, maxX-20, int(float32(maxY)/2)); err != nil {
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
func leftTop(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("plots", 0, 0, maxX/4, int(float32(maxY)/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Plots Details"
		v.Frame = true
		if err != nil {
			log.Println(err)
		}
		getStorageUsage, err := chia.PrintUsage("/storage_5")
		getStorageUsage1, err := chia.PrintUsage("/storage_4")
		getStorageUsage2, err := chia.PrintUsage("/storage")
		if err != nil {
			fmt.Fprintln(v, err)
		}
		fmt.Fprint(v, getStorageUsage)
		fmt.Fprint(v, getStorageUsage1)
		fmt.Fprint(v, getStorageUsage2)

	}

	return nil
}
func middleTop(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//int(float32(maxY) / 2)
	if v, err := g.SetView("totalPlots", maxX/4+1, 0, maxX/2, int(float32(maxY)/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = false
		v.Wrap = false
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		v.Autoscroll = false
		ok := chia.ParseLogs(60)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].ParseTime))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(11))
		fmt.Fprintln(v, graph)

	}
	return nil
}
func mainLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//int(float32(maxY) / 2)
	if v, err := g.SetView("main", int(float32(maxX)/2+1), 0, maxX-20, maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = false
		v.Wrap = false
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		v.Autoscroll = false
		v.Title = "Chia plots elected , last 10 minutes l<r"
		ok := chia.ParseLogs(60)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].Plots))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(11))
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
	if err := middleTop(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("main"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
