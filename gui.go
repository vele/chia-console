package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
	"github.com/kataras/tablewriter"
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
				fmt.Fprintln(v, err)
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
				fmt.Fprintln(v, err)
				return err
			}
			v.Clear()
			//fmt.Fprintf(v, "\033[32mCurrent blockchain space: %v \033[0m \n", spaceCalc)
			fmt.Fprintf(v, "\033[32mCurrent chia price ( XCH ): %f \033[0m \n", ok.ChiaPrice)
			isPositive1h := math.Signbit(ok.PercentChange1H)
			if isPositive1h {
				fmt.Fprintf(v, "Current chia price change 1 h( XCH ):\033[31m%f\033[0m \n", ok.PercentChange1H)
			} else {
				fmt.Fprintf(v, "Current chia price change 1 h( XCH ):\033[32m%f\033[0m \n", ok.PercentChange1H)
			}
			isPositive24h := math.Signbit(ok.PercentChange24h)
			if isPositive24h {
				fmt.Fprintf(v, "Current chia price change 24 h( XCH ):\033[31m%f\033[0m \n", ok.PercentChange24h)
			} else {
				fmt.Fprintf(v, "Current chia price change 24 h( XCH ):\033[32m%f\033[0m \n", ok.PercentChange24h)
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
		diskInfoS5, _ := chia.PrintUsage("/storage_5")
		diskInfoS4, _ := chia.PrintUsage("/storage_4")
		diskInfoS2, _ := chia.PrintUsage("/storage_2")
		diskInfoS1, _ := chia.PrintUsage("/storage")

		data := [][]string{
			[]string{"/storage_5", diskInfoS5.TotalDiskSpace, diskInfoS5.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS5.TotalPercent)},
			[]string{"/storage_4", diskInfoS4.TotalDiskSpace, diskInfoS4.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS4.TotalPercent)},
			[]string{"/storage_2", diskInfoS2.TotalDiskSpace, diskInfoS2.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS2.TotalPercent)},
			[]string{"/storage", diskInfoS1.TotalDiskSpace, diskInfoS1.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS1.TotalPercent)},
		}
		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		totalDiskSpace := diskInfoS5.TotalDiskSpaceBytes + diskInfoS4.TotalDiskSpaceBytes + diskInfoS2.TotalDiskSpaceBytes + diskInfoS1.TotalDiskSpaceBytes
		totalFreeSpace := diskInfoS5.TotalFreeSpaceBytes + diskInfoS4.TotalFreeSpaceBytes + diskInfoS2.TotalFreeSpaceBytes + diskInfoS1.TotalFreeSpaceBytes
		table.SetHeader([]string{"Part", "Total Disk Space", "Total Free Space", "Util %"})
		table.SetFooter([]string{"Tot", humanize.Bytes(totalDiskSpace), humanize.Bytes(totalFreeSpace), ""})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("diskspace")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, tableString.String())

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
	if v, err := g.SetView("blockchain_details", 0, maxY/4+1, maxX/4, int(float32(maxY)/3)); err != nil {

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

	if v, err := g.SetView("wallet", maxX/4+1, maxY/4+1, maxX/2, int(float32(maxY)/3)); err != nil {
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

	if v, err := g.SetView("plots", maxX/2+1, maxY/4+1, int(float32(maxX)/1.5), int(float32(maxY)/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Plot Details last 10 minutes"
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
func priceLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("chia_price", int(float32(maxX)/1.5)+1, maxY/4+1, maxX-1, int(float32(maxY)/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Price Details last 60 seconds"
		v.Frame = true
		if err != nil {
			log.Println(err)
		}
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
		tableString := &strings.Builder{}

		diskInfoS5, _ := chia.PrintUsage("/storage_5")
		diskInfoS4, _ := chia.PrintUsage("/storage_4")
		diskInfoS2, _ := chia.PrintUsage("/storage_2")
		diskInfoS1, _ := chia.PrintUsage("/storage")

		data := [][]string{
			[]string{"/storage_5", diskInfoS5.TotalDiskSpace, diskInfoS5.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS5.TotalPercent)},
			[]string{"/storage_4", diskInfoS4.TotalDiskSpace, diskInfoS4.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS4.TotalPercent)},
			[]string{"/storage_2", diskInfoS2.TotalDiskSpace, diskInfoS2.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS2.TotalPercent)},
			[]string{"/storage", diskInfoS1.TotalDiskSpace, diskInfoS1.TotalFreeSpace, fmt.Sprintf("%0.2f%%", diskInfoS1.TotalPercent)},
		}
		table := tablewriter.NewWriter(tableString)
		totalDiskSpace := diskInfoS5.TotalDiskSpaceBytes + diskInfoS4.TotalDiskSpaceBytes + diskInfoS2.TotalDiskSpaceBytes + diskInfoS1.TotalDiskSpaceBytes
		totalFreeSpace := diskInfoS5.TotalFreeSpaceBytes + diskInfoS4.TotalFreeSpaceBytes + diskInfoS2.TotalFreeSpaceBytes + diskInfoS1.TotalFreeSpaceBytes
		table.SetHeader([]string{"Part", "Total Disk Space", "Total Free Space", "Util %"})
		table.SetFooter([]string{"Tot", humanize.Bytes(totalDiskSpace), humanize.Bytes(totalFreeSpace), ""})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
		fmt.Fprintln(v, tableString.String())

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
		v.Frame = true
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
	if v, err := g.SetView("main", maxX/2+1, 0, maxX-1, maxY/4); err != nil {
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
		graph := asciigraph.Plot(data, asciigraph.Height(11))
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
	if err := plotsLayout(g); err != nil {
		return err
	}
	if err := middleTop(g); err != nil {
		return err
	}
	if err := priceLayout(g); err != nil {
		return err
	}

	if _, err := g.SetCurrentView("main"); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
