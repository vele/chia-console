package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/guptarohit/asciigraph"
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
		ok := chia.ParseLogs(900)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].Plots))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Caption("Chia plots elected , last 15 minutes r>l"), asciigraph.Width(0), asciigraph.Precision(0))
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

func drawProcessingTimesGraph(g *gocui.Gui) error {
	defer wg.Done()

	for {
		time.Sleep(1 * time.Second)
		ok := chia.ParseLogs(900)
		var data []float64
		for item := range ok {
			data = append(data, float64(ok[item].ParseTime))
		}
		graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Caption("Chia plots processing speed , last 15 minutes r>l"))
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
		const precision = 12
		chia_mojo_calc, _ := new(big.Float).SetPrec(precision).SetString("1000000000000")
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
				plots = append(plots, float64(ok[item].Plots))
			}
			wallet := getWalletDetails()

			blockChainClient := chia.NewClient(os.Getenv("CHIA_HARVESTER_CRT"), os.Getenv("CHIA_HARVESTER_KEY"), os.Getenv("CHIA_CA_CRT"))
			res, err := blockChainClient.GetChiaPlots(os.Getenv("CHIA_HARVESTER_URL"))
			fmt.Fprintf(v, "\u2705 Total space utilized by plots: %d TB \n", len(res.Plots)*108/1024)
			fmt.Fprintf(v, "\u2705 Total plots: %d  \n", len(res.Plots))
			fmt.Fprintf(v, "\u2705 Total netspace: %s  \n", returnBlockChainDetails())
			chia_mojo_balance, _ := new(big.Float).SetPrec(precision).SetString(fmt.Sprintf("%d", wallet.WalletBalance.ConfirmedWalletBalance))
			formula_result := new(big.Float).Quo(chia_mojo_balance, chia_mojo_calc)
			chia_mojo_balance_spendable, _ := new(big.Float).SetPrec(precision).SetString(fmt.Sprintf("%d", wallet.WalletBalance.SpendableBalance))
			formula_result_spendable := new(big.Float).Quo(chia_mojo_balance_spendable, chia_mojo_calc)
			fmt.Fprintf(v, "\u2705 Current wallet ballance : %0.12f  \n", formula_result)
			fmt.Fprintf(v, "\u2705 Spendable wallet ballance: %0.12f  \n", formula_result_spendable)
			chia_price := returnChiaPriceDetails()
			fmt.Fprintf(v, "\u2705 Total chia( XCH ):\033[32m%0.1f\033[0m \n", chia_price.TotalSupply)
			fmt.Fprintf(v, "\u2705 Current chia price ( XCH ):\033[34mUSD %f\033[0m\n", chia_price.ChiaPrice)
			isPositive1h := math.Signbit(chia_price.PercentChange1H)
			if isPositive1h {
				fmt.Fprintf(v, "\u2705 Current chia price change 1 h( XCH ):\033[31m%0.2f%%\033[0m \n", chia_price.PercentChange1H)
			} else {
				fmt.Fprintf(v, "\u2705 Current chia price change 1 h( XCH ):\033[32m%0.2f%%\033[0m \n", chia_price.PercentChange1H)
			}
			isPositive24h := math.Signbit(chia_price.PercentChange24h)
			if isPositive24h {
				fmt.Fprintf(v, "\u2705 Current chia price change 24 h( XCH ):\033[31m%0.2f%%\033[0m \n", chia_price.PercentChange24h)
			} else {
				fmt.Fprintf(v, "\u2705 Current chia price change 24 h( XCH ):\033[32m%0.2f%%\033[0m \n", chia_price.PercentChange24h)
			}
			if len(data) == 0 {
				data = append(data, 0)
			}
			if len(plots) == 0 {
				plots = append(plots, 0)
			}
			if data[0] >= 1.00 {
				fmt.Fprintf(v, "\u2705 Last transaction took  \033[31m\u25BC\033[0m: \033[31m%0.2f\033[0m sec\n", data[0])
			} else {
				fmt.Fprintf(v, "\u2705 Last transaction took  \033[32m\u25BC\033[0m: \033[32m%0.2f\033[0m sec\n", data[0])
			}
			if plots[0] <= 10.00 {
				fmt.Fprintf(v, "\u2705 Last eligable plots  \033[31m\u25BC\033[0m: \033[31m%0.1f\033[0m plots\n", plots[0])
			} else {
				fmt.Fprintf(v, "\u2705 Last eligable plots  \033[32m\u25BC\033[0m: \033[32m%0.1f\033[0m plots\n", plots[0])
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

func returnBlockChainDetails() string {
	blockChainClient := chia.NewClient(os.Getenv("CHIA_FULL_NODE_CRT"), os.Getenv("CHIA_FULL_NODE_KEY"), os.Getenv("CHIA_CA_CRT"))
	res, _ := blockChainClient.GetChiaBlockchainState(os.Getenv("CHIA_SERVER_URL"))
	spaceCalc := chia.ByteCountSI(res.BlockchainState.Space)
	return spaceCalc
}
func returnChiaPriceDetails() *chia.ChiaTableDbResponse {
	chiaPrice := chia.FetchChiaPriceDB()
	return chiaPrice
}
func getWalletDetails() chia.WalletBallance {
	blockChainClient := chia.NewClient(os.Getenv("CHIA_WALLET_CRT"), os.Getenv("CHIA_WALLET_KEY"), os.Getenv("CHIA_CA_CRT"))
	res, _ := blockChainClient.GetChiaWallet(os.Getenv("CHIA_WALLET_URL"))
	return res
}

func leftTop(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("diskspace", 0, 0, int(0.25*float32(maxX)), int(0.3*float32(maxY)), gocui.LEFT); err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatal("POOP")
			return err
		}
		v.Title = "Disk details"
		v.Frame = true
		v.FrameColor = gocui.ColorMagenta
		v.Subtitle = "Subtitle"
		v.Wrap = true

	}

	return nil
}
func secondRowGraph(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//int(float32(maxY) / 2)
	if v, err := g.SetView("totalPlots", int(0.3*float32(maxX)), int(0.3*float32(maxY))+1, maxX-10, int(float32(maxY)/2), gocui.BOTTOM); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.ColorYellow
		v.Frame = false
	}
	return nil
}
func firstRowGraph(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", int(0.3*float32(maxX)+5), 0, maxX-1, int(0.3*float32(maxY)), gocui.LEFT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.NewRGBColor(145, 204, 227)
		v.Wrap = true
	}
	if err := leftTop(g); err != nil {
		return err
	}

	if err := secondRowGraph(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("main"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
