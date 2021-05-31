package main

import (
	"fmt"
	"log"
	"math"
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
				//fmt.Fprintln(v, err)
				return err
			}
			v.Clear()
			fmt.Fprintln(v, graph)
			return nil
		})
	}
}
func updateChiaPriceDB(g *gocui.Gui) error {
	defer wg.Done()
	for {
		time.Sleep(60 * time.Second)
		result, err := chia.FetchCoinData("XCH")
		if err != nil {
			log.Panicln(err)
		}
		err = chia.FetchChiaPrice(result)
		if err != nil {
			log.Panicln(err)
		}
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("banner")
			if err != nil {
				//fmt.Fprintln(v, err)
				return err
			}
			v.Clear()
			fmt.Fprintf(v, "Last update: %s", time.Now())
			return nil
		})
	}

}

func updateChiaPriceGUI(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)
		ok := chia.FetchChiaPriceDB()

		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("chia_price")
			if err != nil {
				//fmt.Fprintln(v, err)
				return err
			}
			v.Clear()
			v.Title = fmt.Sprintf("Chia price \u2B50 %f \u2B50", ok.ChiaPrice)
			fmt.Fprintf(v, "Current chia price ( XCH ): \u2B50 \033[34mUSD%f\033[0m \u2B50 \n", ok.ChiaPrice)
			isPositive1h := math.Signbit(ok.PercentChange1H)
			if isPositive1h {
				fmt.Fprintf(v, "\nCurrent chia price change 1 h( XCH ):\033[31m%0.2f%%\033[0m \n", ok.PercentChange1H)
			} else {
				fmt.Fprintf(v, "Current chia price change 1 h( XCH ):\033[32m%0.2f%%\033[0m \n", ok.PercentChange1H)
			}
			isPositive24h := math.Signbit(ok.PercentChange24h)
			if isPositive24h {
				fmt.Fprintf(v, "Current chia price change 24 h( XCH ):\033[31m%0.2f%%\033[0m \n", ok.PercentChange24h)
			} else {
				fmt.Fprintf(v, "Current chia price change 24 h( XCH ):\033[32m%0.2f%%\033[0m \n", ok.PercentChange24h)
			}
			fmt.Fprintf(v, "Total chia( XCH ):\033[32m%f\033[0m \n", ok.TotalSupply)
			return nil
		})
	}
}
func drawProcessingTimesGraph(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)
		ok := chia.ParseLogs(300)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].ParseTime))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(11))
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("totalPlots")
			if err != nil {
				fmt.Fprintln(v, err)
				return err
			}
			v.Clear()
			fmt.Fprintln(v, graph)
			return nil
		})
	}
}
func drawFreeSpaceTable(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)

		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("diskspace")
			if err != nil {
				return err
			}
			v.Clear()
			ok := chia.ParseLogs(10)
			var data []float64
			for item := range ok {
				data = append(data, float64(ok[item].ParseTime))
			}
			var plots []float64
			for item := range ok {
				plots = append(data, float64(ok[item].Plots))
			}
			blockChainClient := chia.NewClient(os.Getenv("CHIA_HARVESTER_CRT"), os.Getenv("CHIA_HARVESTER_KEY"), os.Getenv("CHIA_CA_CRT"))
			res, err := blockChainClient.GetChiaPlots(os.Getenv("CHIA_HARVESTER_URL"))
			fmt.Fprintf(v, "\u2705\t Total space utilized by plots: %d TB \n", len(res.Plots)*108/1024)
			fmt.Fprintf(v, "\u2705\t Total plots: %d  \n", len(res.Plots))
			if len(data) == 0 {
				data = append(data, 0)
			}
			if len(plots) == 0 {
				plots = append(plots, 0)
			}
			if data[0] >= 1.00 {
				fmt.Fprintf(v, "\u2705\t Last transaction took  \033[31m\u25BC\033[0m: \033[31m%0.2f\033[0m sec\n", data[0])
			} else {
				fmt.Fprintf(v, "\u2705\t Last transaction took  \033[32m\u25BC\033[0m: \033[32m%0.2f\033[0m sec\n", data[0])
			}
			if plots[0] <= 10.00 {
				fmt.Fprintf(v, "\u2705\t Last eligable plots  \033[31m\u25BC\033[0m: \033[31m%0.2f\033[0m sec\n", plots[0])
			} else {
				fmt.Fprintf(v, "\u2705\t Last transaction took  \033[32m\u25BC\033[0m: \033[32m%0.2f\033[0m sec\n", plots[0])
			}

			return nil
		})
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func detailsLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("blockchain_details", 0, maxY/4+1, maxX/4, int(float32(maxY)/2)); err != nil {

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

	if v, err := g.SetView("wallet", maxX/4+1, maxY/4+1, maxX/2, int(float32(maxY)/2)); err != nil {
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

	if v, err := g.SetView("plots", maxX/2+1, maxY/4+1, int(float32(maxX)/1.8), int(float32(maxY)/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Plot Details last 10 minutes"
		v.Frame = false
		blockChainClient := chia.NewClient(os.Getenv("CHIA_HARVESTER_CRT"), os.Getenv("CHIA_HARVESTER_KEY"), os.Getenv("CHIA_CA_CRT"))
		res, err := blockChainClient.GetChiaPlots(os.Getenv("CHIA_HARVESTER_URL"))
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(v, len(res.Plots))
	}

	return nil
}
func priceLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("chia_price", int(float32(maxX)/2)+1, int(float32(maxY)/4)+1, int(float32(maxX)-1.5), int(float32(maxY)/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//v.Title = "Price Details last 60 seconds"
		v.Frame = false
	}

	return nil
}
func banner(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("banner", 0, int(float32(maxY)/2)+1, int(float32(maxX)-1.5), int(float32(maxY)/1.5)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Updates"
		v.Frame = false

	}
	return nil
}
func leftTop(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("diskspace", 0, 0, maxX/4, int(float32(maxY)/4)); err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatal("POOP")
			return err
		}
		v.Title = "Disk details"
		v.Frame = false

	}

	return nil
}
func middleTop(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//int(float32(maxY) / 2)
	if v, err := g.SetView("totalPlots", maxX/4+1, 0, maxX/2, int(float32(maxY)/4)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.ColorCyan
		v.Frame = false
		v.Title = "Chia plots processing speed , last 10 minutes l<r"
		ok := chia.ParseLogs(600)
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
	if v, err := g.SetView("main", maxX/2+1, 0, int(float32(maxX)-1), maxY/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.ColorYellow
		v.Frame = true
		v.Autoscroll = false
		v.Title = "Chia plots elected , last 10 minutes l<r"
		ok := chia.ParseLogs(600)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].Plots))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(12))
		fmt.Fprintln(v, graph)
	}
	if err := leftTop(g); err != nil {
		return err
	}

	if err := detailsLayout(g); err != nil {
		return err
	}
	if err := walletLayout(g); err != nil {
		return err
	}
	if err := middleTop(g); err != nil {
		return err
	}
	if err := priceLayout(g); err != nil {
		return err
	}
	if err := banner(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("main"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
