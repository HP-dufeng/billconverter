package converter

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fengdu/billconverter/util"
)

// BillBaseInfo bill base info from src file
type BillBaseInfo struct {
	AccountNo          string
	StatementDateStart time.Time
	StatementDateEnd   time.Time
	BillDate           time.Time
}

// GetBillBaseInfo parse bill base info segment
func GetBillBaseInfo(content string) (BillBaseInfo, error) {
	var result BillBaseInfo

	header, err := readHeadSegment(content)
	if err != nil {
		return result, fmt.Errorf("ERROR: readHeadSegment: %v", err)
	}

	result.AccountNo = header["Account No"]
	if b, ok := header["Bill Date"]; ok {
		if t, err := time.Parse("2006-01-02", b); err == nil {
			result.BillDate = t
		}
	}

	if statementDate, ok := header["Statement Date"]; ok {
		if strings.Contains(statementDate, "to") {
			ss := strings.Split(statementDate, "to")
			if t, err := time.Parse("2006-01-02", strings.TrimSpace(ss[0])); err == nil {
				result.StatementDateStart = t
			}
			if t, err := time.Parse("2006-01-02", strings.TrimSpace(ss[1])); err == nil {
				result.StatementDateEnd = t
			}
		}
	}

	return result, nil
}

// GetBalances parse Financial Situation segment
func GetBalances(bill BillBaseInfo, content string) [][]string {
	s := util.ParseSegment(content, "Financial Situation")

	set := make(map[string][]string)
	for _, r := range readSegment(s) {
		set[strings.TrimSpace(r[0])] = r
	}

	result := [][]string{
		{
			"Account", "Currency", "BalanceBf", "Deposit", "Withdrawal",
			"OptionPremium", "DeliveryProceed", "RealisedPL", "Commission", "Interest",
			"Others", "BalanceCf", "UnrealisedPL", "Equity", "NetOptionValue",
			"EligCollateral", "as-of-date mm/dd/yyyy",
		},
	}

	Account, Currency, BalanceBf, Deposit, Withdrawal,
		OptionPremium, DeliveryProceed, RealisedPL, Commission, Interest,
		Others, BalanceCf, UnrealisedPL, Equity, NetOptionValue,
		EligCollateral, asofdate :=
		bill.AccountNo, "CNY", "Opening", "", "",
		"", "", "Opening - Closing", "", "",
		"", "Closing", "", "", "",
		"", bill.BillDate.Format("01/02/2006")

	if v, ok := set["Deposit/Withdrawal"]; ok {
		if len(v) >= 3 {
			o := strings.Replace(strings.TrimSpace(v[2]), ",", "", -1)
			if f, err := strconv.ParseFloat(o, 32); err == nil {
				if f > 0 {
					Deposit = strconv.FormatFloat(f, 'f', 2, 32)
				} else {
					Withdrawal = strconv.FormatFloat(f, 'f', 2, 32)
				}
			}
		}
	}

	if v, ok := set["Commissions"]; ok {
		if len(v) >= 3 {
			Commission = strings.Replace(strings.TrimSpace(v[2]), ",", "", -1)
		}
	}

	if v, ok := set["Unrealized"]; ok {
		if len(v) >= 3 {
			UnrealisedPL = strings.Replace(strings.TrimSpace(v[2]), ",", "", -1)
		}
	}

	if v, ok := set["Equity"]; ok {
		if len(v) >= 3 {
			Equity = strings.Replace(strings.TrimSpace(v[2]), ",", "", -1)
		}
	}

	result = append(result, []string{
		Account, Currency, BalanceBf, Deposit, Withdrawal,
		OptionPremium, DeliveryProceed, RealisedPL, Commission, Interest,
		Others, BalanceCf, UnrealisedPL, Equity, NetOptionValue,
		EligCollateral, asofdate,
	})

	return result
}

// GetPos parse Gathered Open Positions segment
func GetPos(bill BillBaseInfo, content string) [][]string {
	s := util.ParseSegment(content, "Gathered Open Positions")
	result := [][]string{
		{
			"Account", "Tradedate", "Long", "Short", "FutOpt",
			"Exchange", "Contract", "ContractMonth", "Contractyear", "StrikePrice",
			"Price", "SettPrice", "Currency", "UnrealisedPL", "TradeNo",
			"BUY/Sell 1=BUY 0=SELL", "SubType P=Put C=Call", "Commodity", "Commission", "Option Delta",
			"Firm/Office", "as-of-date (mm/dd/yyyy)",
		},
	}

	if ss := readSegment(s); len(ss) >= 2 {
		for _, r := range ss[1 : len(ss)-1] {
			buy := strings.TrimSpace(r[3])
			sale := strings.TrimSpace(r[4])
			market := strings.TrimSpace(r[0])
			contract := strings.TrimSpace(r[2])
			contractMonth, contractYear, _ := util.ParseMonthAndYear(contract)
			matchPrice := strings.TrimSpace(r[5])
			settlementPrice := strings.TrimSpace(r[6])
			currency := strings.TrimSpace(r[10])
			positionProfit := strings.TrimSpace(r[7])
			product := strings.TrimSpace(r[1])

			result = append(result, []string{
				bill.AccountNo, bill.StatementDateEnd.Format("2006-01-02"), buy, sale, "F",
				market, contract, contractMonth, contractYear, "",
				matchPrice, settlementPrice, currency, positionProfit, "",
				"", "", product, "", "",
				"Shanghai Bunge", bill.StatementDateEnd.Format("01/02/2006"),
			})
		}
	}

	return result
}

// GetTrades parse Trade Confirmation segment
func GetTrades(bill BillBaseInfo, content string) [][]string {
	s := util.ParseSegment(content, "Trade Confirmation")
	result := [][]string{
		{
			"Account", "Tradedate", "Long", "Short", "FutOpt",
			"Exchange", "Contract", "ContractMonth", "Contractyear", "StrikePrice",
			"Price", "SettPrice", "Currency", "UnrealisedPL", "TradeNo",
			"BUY/Sell 1=BUY 0=SELL", "SubType P=Put C=Call", "Commodity", "Commission", "Option Delta",
			"Firm/Office", "as-of-date (mm/dd/yyyy)",
		},
	}

	if ss := readSegment(s); len(ss) >= 2 {
		for _, r := range ss[1 : len(ss)-1] {
			date := strings.TrimSpace(r[0])
			matchQty := strings.TrimSpace(r[7])
			market := strings.TrimSpace(r[1])
			contract := strings.TrimSpace(r[3])
			contractMonth, contractYear, _ := util.ParseMonthAndYear(contract)
			matchPrice := strings.TrimSpace(r[8])
			currency := strings.TrimSpace(r[11])
			buySale := strings.TrimSpace(r[6])
			product := strings.TrimSpace(r[2])
			fee := strings.TrimSpace(r[10])

			long, short := "0", "0"
			if buySale == "Sale" {
				short = matchQty
			} else if buySale == "Buy" {
				long = matchQty
			}

			var asOfDate string
			if t, err := time.Parse("2006-01-02", date); err == nil {
				asOfDate = t.Format("01/02/2006")
			}

			result = append(result, []string{
				bill.AccountNo, date, long, short, "F",
				market, contract, contractMonth, contractYear, "",
				matchPrice, "", currency, "", "",
				buySale, "", product, fee, "",
				"Shanghai Bunge", asOfDate,
			})
		}
	}

	return result
}

func readHeadSegment(content string) (map[string]string, error) {
	s := util.ParseSegment(content, "Account No")
	s = strings.Replace(s, "：", ":", -1)

	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "|")
		if regexp.MustCompile(`^-+`).MatchString(line) {
			continue
		}
		p := regexp.MustCompile(`\s{2,}`).Split(line, -1)
		for _, f := range p {
			kv := strings.Split(f, ":")
			if len(kv) != 2 {
				return nil, errors.New("Parse bill base info errors")
			}
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	return result, nil
}

// readSegment parse table style segment to array
func readSegment(segment string) [][]string {
	result := [][]string{}
	scanner := bufio.NewScanner(strings.NewReader(segment))
	for scanner.Scan() {
		line := strings.Trim(strings.TrimSpace(scanner.Text()), "|")
		if !strings.Contains(line, "|") {
			continue
		}
		result = append(result, strings.Split(line, "|"))
	}

	return result
}
